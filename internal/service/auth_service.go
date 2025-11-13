package service

import (
	"context"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/internal/entity/user"
	"github.com/prawirdani/golang-restapi/internal/infra/repository"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

type AuthService struct {
	cfg       config.Auth
	tx        repository.Transactor
	authRepo  auth.Repository
	userRepo  user.Repository
	publisher auth.MessagePublisher
}

func NewAuthService(
	cfg config.Auth,
	transactor repository.Transactor,
	userRepo user.Repository,
	authRepo auth.Repository,
	publisher auth.MessagePublisher,
) *AuthService {
	return &AuthService{
		cfg:       cfg,
		tx:        transactor,
		userRepo:  userRepo,
		authRepo:  authRepo,
		publisher: publisher,
	}
}

func (s *AuthService) Register(ctx context.Context, inp model.CreateUserInput) error {
	userExists, err := s.userRepo.GetByEmail(ctx, inp.Email)
	if err != nil && err != user.ErrNotFound {
		return err
	}

	if userExists != nil {
		return user.ErrEmailExist
	}

	hashedPassword, err := auth.HashPassword(inp.Password)
	if err != nil {
		return err
	}

	newUser, err := user.New(inp, string(hashedPassword))
	if err != nil {
		log.ErrorCtx(ctx, "Failed to create new user", err)
		return err
	}

	return s.userRepo.Store(ctx, newUser)
}

// Login is a method to authenticate the user, returning access token, refresh token, and error if any.
func (s *AuthService) Login(
	ctx context.Context,
	inp model.LoginInput,
) (accessToken string, sessID string, err error) {
	usr, err := s.userRepo.GetByEmail(ctx, inp.Email)
	if err != nil {
		if err == user.ErrNotFound {
			return accessToken, sessID, auth.ErrWrongCredentials
		}
		return accessToken, sessID, err
	}

	if err := auth.VerifyPassword(inp.Password, usr.Password); err != nil {
		return accessToken, sessID, err
	}

	accessToken, err = s.generateAccessToken(*usr)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to generate access token", err)
		return accessToken, sessID, err
	}

	sess, err := auth.NewSession(
		usr.ID,
		inp.UserAgent,
		s.cfg.SessionTTL,
	)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to create new session", err)
		return accessToken, sessID, err
	}

	if err = s.authRepo.StoreSession(ctx, sess); err != nil {
		return accessToken, sessID, err
	}

	return accessToken, sess.ID.String(), nil
}

// TODO: Should also refreshing the refresh token, maybe by checking the exp time, if its nearly N to expire, then refresh it.
func (s *AuthService) RefreshAccessToken(
	ctx context.Context,
	sessID string,
) (string, error) {
	sess, err := s.authRepo.GetSession(ctx, sessID)
	if err != nil {
		return "", err
	}

	if sess.IsExpired() {
		return "", auth.ErrSessionExpired
	}

	usr, err := s.userRepo.GetByID(ctx, sess.UserID.String())
	if err != nil {
		return "", err
	}

	newAccessToken, err := s.generateAccessToken(*usr)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to generate new access token", err)
		return "", err
	}

	return newAccessToken, nil
}

func (s *AuthService) Logout(ctx context.Context, sessID string) error {
	// return s.authRepo.DeleteSession(ctx, sessID)
	return nil
}

// ForgotPassword initiates the password reset process by sending a reset link or token to the user's email.
func (s *AuthService) ForgotPassword(ctx context.Context, inp model.ForgotPasswordInput) error {
	return s.tx.Transact(ctx, func(ctx context.Context) error {
		usr, err := s.userRepo.GetByEmail(ctx, inp.Email)
		if err != nil {
			if err == user.ErrNotFound {
				return user.ErrEmailNotVerified
			}
			return err
		}

		token, err := auth.NewResetPasswordToken(usr.ID, s.cfg.ResetPasswordTTL)
		if err != nil {
			log.ErrorCtx(ctx, "Failed to create reset password token", err)
			return err
		}

		// Save token to db
		if err := s.authRepo.StoreResetPasswordToken(ctx, token); err != nil {
			return err
		}

		// Publish email job to message queue
		msg := model.ResetPasswordEmailMessage{
			To:       usr.Email,
			Name:     usr.Name,
			ResetURL: s.cfg.ResetPasswordFormEndpoint + "?token=" + token.Value,
			Expiry:   s.cfg.ResetPasswordTTL,
		}

		return s.publisher.SendResetPasswordEmail(ctx, msg)
	})
}

func (s *AuthService) GetResetPasswordToken(
	ctx context.Context,
	token string,
) (*auth.ResetPasswordToken, error) {
	return s.authRepo.GetResetPasswordToken(ctx, token)
}

// ResetPassword resets a user's password using a valid reset password token from email.
func (s *AuthService) ResetPassword(ctx context.Context, inp model.ResetPasswordInput) error {
	return s.tx.Transact(ctx, func(ctx context.Context) error {
		token, err := s.authRepo.GetResetPasswordToken(ctx, inp.Token)
		if err != nil {
			return err
		}

		if token.Expired() || token.Used() {
			return auth.ErrResetPasswordTokenInvalid
		}

		user, err := s.userRepo.GetByID(ctx, token.UserID.String())
		if err != nil {
			return err
		}

		newHashedPassword, err := auth.HashPassword(inp.NewPassword)
		if err != nil {
			log.ErrorCtx(ctx, "Failed to hash new password", err)
			return err
		}
		user.Password = string(newHashedPassword)

		token.Revoke()
		if err := s.authRepo.UpdateResetPasswordToken(ctx, token); err != nil {
			return err
		}

		return s.userRepo.Update(ctx, user)
	})
}

// ChangePassword updates the authenticated user's password after verifying the current password.
func (s *AuthService) ChangePassword(
	ctx context.Context,
	userID string,
	inp model.ChangePasswordInput,
) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := auth.VerifyPassword(inp.Password, u.Password); err != nil {
		return err
	}

	// Hash new password
	newHashedPassword, err := auth.HashPassword(inp.NewPassword)
	if err != nil {
		return err
	}

	u.Password = string(newHashedPassword)

	return s.userRepo.Update(ctx, u)
}

func (s *AuthService) generateAccessToken(user user.User) (string, error) {
	return auth.SignAccessToken(
		s.cfg.JwtSecret,
		auth.AccessTokenClaims{UserID: user.ID.String()},
		s.cfg.JwtTTL,
	)
}
