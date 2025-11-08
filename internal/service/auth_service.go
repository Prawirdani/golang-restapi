package service

import (
	"context"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/internal/domain/user"
	"github.com/prawirdani/golang-restapi/internal/infra/repository"
	"github.com/prawirdani/golang-restapi/internal/messages"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

type AuthService struct {
	cfg         *config.Config
	tx          repository.Transactor
	authRepo    auth.Repository
	userRepo    user.Repository
	userService *UserService
	publisher   auth.MessagePublisher
}

func NewAuthService(
	cfg *config.Config,
	transactor repository.Transactor,
	userRepo user.Repository,
	authRepo auth.Repository,
	userService *UserService,
	publisher auth.MessagePublisher,
) *AuthService {
	return &AuthService{
		cfg:         cfg,
		tx:          transactor,
		userRepo:    userRepo,
		authRepo:    authRepo,
		userService: userService,
		publisher:   publisher,
	}
}

func (s *AuthService) Register(ctx context.Context, payload model.CreateUserInput) error {
	userExists, err := s.userRepo.GetUserBy(ctx, "email", payload.Email)
	if err != nil {
		return err
	}

	if userExists != nil {
		return user.ErrEmailExist
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	newUser, err := user.New(payload, string(hashedPassword))
	if err != nil {
		log.ErrorCtx(ctx, "Failed to create new user", err)
		return err
	}

	return s.userRepo.Insert(ctx, newUser)
}

// Login is a method to authenticate the user, returning access token, refresh token, and error if any.
func (s *AuthService) Login(
	ctx context.Context,
	payload model.LoginInput,
) (accessToken string, refreshToken string, err error) {
	u, err := s.userRepo.GetUserBy(ctx, "email", payload.Email)
	if err != nil {
		if err == user.ErrUserNotFound {
			return accessToken, refreshToken, auth.ErrWrongCredentials
		}
		return accessToken, refreshToken, err
	}

	if err := auth.VerifyPassword(payload.Password, u.Password); err != nil {
		return accessToken, refreshToken, err
	}

	accessToken, err = s.generateAccessToken(map[string]any{
		"id":   u.ID.String(),
		"name": u.Name,
	})
	if err != nil {
		log.ErrorCtx(ctx, "Failed to generate access token", err)
		return accessToken, refreshToken, err
	}

	sess, err := auth.NewSession(
		u.ID,
		payload.UserAgent,
		s.cfg.Token.RefreshTokenExpiry,
	)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to create new session", err)
		return accessToken, refreshToken, err
	}

	if err = s.authRepo.InsertSession(ctx, sess); err != nil {
		return accessToken, refreshToken, err
	}

	return accessToken, sess.RefreshToken, nil
}

// TODO: Should also refreshing the refresh token, maybe by checking the exp time, if its nearly N to expire, then refresh it.
func (s *AuthService) RefreshAccessToken(
	ctx context.Context,
	refreshToken string,
) (string, error) {
	sess, err := s.authRepo.GetUserSessionBy(ctx, "refresh_token", refreshToken)
	if err != nil {
		return "", err
	}

	if sess.IsExpired() {
		_ = s.authRepo.DeleteSession(ctx, "id", sess.ID)
		return "", auth.ErrSessionExpired
	}

	user, err := s.userRepo.GetUserBy(ctx, "id", sess.UserID)
	if err != nil {
		return "", err
	}

	newAccessToken, err := s.generateAccessToken(
		map[string]any{"id": user.ID.String(), "name": user.Name},
	)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to generate new access token", err)
		return "", err
	}

	return newAccessToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.authRepo.DeleteSession(ctx, "refresh_token", refreshToken)
}

func (s *AuthService) IdentifyUser(ctx context.Context) (*user.User, error) {
	tokenPayload, err := auth.GetAuthCtx(ctx)
	if err != nil {
		return nil, err
	}

	u, err := s.userService.GetUserByID(ctx, tokenPayload["id"].(string))
	if err != nil {
		return nil, err
	}

	return u, nil
}

// ForgotPassword initiates the password reset process by sending a reset link or token to the user's email.
func (s *AuthService) ForgotPassword(ctx context.Context, i model.ForgotPasswordInput) error {
	return s.tx.Transact(ctx, func(ctx context.Context) error {
		u, err := s.userRepo.GetUserBy(ctx, "email", i.Email)
		if err != nil {
			if err == user.ErrUserNotFound {
				return user.ErrEmailNotVerified
			}
			return err
		}

		expiresAt := time.Now().Add(s.cfg.Token.ResetPasswordTokenExpiry)
		token, err := auth.NewResetPasswordToken(u.ID, expiresAt)
		if err != nil {
			log.ErrorCtx(ctx, "Failed to create reset password token", err)
			return err
		}

		// Save token to db
		if err := s.authRepo.InsertResetPasswordToken(ctx, token); err != nil {
			return err
		}

		// Publish email job to message queue
		msg := messages.ResetPasswordEmail{
			To:       u.Email,
			Name:     u.Name,
			ResetURL: s.cfg.Token.ResetPasswordFormEndpoint + "?token=" + token.Value,
			Expiry:   s.cfg.Token.ResetPasswordTokenExpiry,
		}

		return s.publisher.SendResetPasswordEmail(ctx, msg)
	})
}

func (s *AuthService) GetResetPasswordToken(
	ctx context.Context,
	token string,
) (*auth.ResetPasswordToken, error) {
	return s.authRepo.GetResetPasswordTokenObj(ctx, token)
}

// ResetPassword resets a user's password using a valid reset password token from email.
func (s *AuthService) ResetPassword(ctx context.Context, i model.ResetPasswordInput) error {
	return s.tx.Transact(ctx, func(ctx context.Context) error {
		token, err := s.authRepo.GetResetPasswordTokenObj(ctx, i.Token)
		if err != nil {
			return err
		}

		if token.Expired() || token.Used() {
			return auth.ErrResetPasswordTokenInvalid
		}

		user, err := s.userRepo.GetUserBy(ctx, "id", token.UserId)
		if err != nil {
			return err
		}

		newHashedPassword, err := auth.HashPassword(i.NewPassword)
		if err != nil {
			log.ErrorCtx(ctx, "Failed to hash new password", err)
			return err
		}
		user.Password = string(newHashedPassword)

		if err := s.authRepo.InvalidateResetPasswordToken(ctx, token); err != nil {
			return err
		}

		return s.userRepo.UpdateUser(ctx, user)
	})
}

// ChangePassword updates the authenticated user's password after verifying the current password.
func (s *AuthService) ChangePassword(ctx context.Context, i model.ChangePasswordInput) error {
	tokenPayload, err := auth.GetAuthCtx(ctx)
	if err != nil {
		return err
	}

	u, err := s.userRepo.GetUserBy(ctx, "id", tokenPayload["id"])
	if err != nil {
		return err
	}

	// Verify old password
	if err := auth.VerifyPassword(i.Password, u.Password); err != nil {
		return err
	}

	// Hash new password
	newHashedPassword, err := auth.HashPassword(i.NewPassword)
	if err != nil {
		return err
	}

	u.Password = string(newHashedPassword)

	return s.userRepo.UpdateUser(ctx, u)
}

func (s *AuthService) generateAccessToken(payload map[string]any) (string, error) {
	return auth.GenerateJWT(
		s.cfg.Token.SecretKey,
		s.cfg.Token.AccessTokenExpiry,
		&payload,
	)
}
