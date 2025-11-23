// Package auth provides authentication and authorization functionality.
// This package handles user authentication through sessions, access tokens, and
// password management including secure hashing and password reset flows. It manages
// the complete authentication lifecycle from login through logout, including token
// generation, validation, and session management.
package auth

import (
	"context"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/domain/user"
	"github.com/prawirdani/golang-restapi/internal/infrastructure/repository"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

type Service struct {
	cfg        config.Auth
	transactor repository.Transactor
	authRepo   Repository
	userRepo   user.Repository
	publisher  MessagePublisher
}

func NewService(
	cfg config.Auth,
	transactor repository.Transactor,
	userRepo user.Repository,
	authRepo Repository,
	publisher MessagePublisher,
) *Service {
	return &Service{
		cfg:        cfg,
		transactor: transactor,
		userRepo:   userRepo,
		authRepo:   authRepo,
		publisher:  publisher,
	}
}

func (s *Service) Register(ctx context.Context, inp RegisterInput) error {
	userExists, err := s.userRepo.GetByEmail(ctx, inp.Email)
	if err != nil && err != user.ErrNotFound {
		return err
	}

	if userExists != nil {
		return user.ErrEmailExists
	}

	hashedPassword, err := HashPassword(inp.Password)
	if err != nil {
		return err
	}

	newUser, err := user.New(
		inp.Name,
		inp.Email,
		inp.Phone,
		string(hashedPassword),
	)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to create new user", err)
		return err
	}

	return s.userRepo.Store(ctx, newUser)
}

// Login is a method to authenticate the user, returning access token, refresh token, and error if any.
func (s *Service) Login(
	ctx context.Context,
	inp LoginInput,
) (accessToken string, sessID string, err error) {
	usr, err := s.userRepo.GetByEmail(ctx, inp.Email)
	if err != nil {
		if err == user.ErrNotFound {
			return accessToken, sessID, ErrWrongCredentials
		}
		return accessToken, sessID, err
	}

	if err := VerifyPassword(inp.Password, usr.Password); err != nil {
		return accessToken, sessID, err
	}

	accessToken, err = s.generateAccessToken(*usr)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to generate access token", err)
		return accessToken, sessID, err
	}

	sess, err := NewSession(
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

// TODO: Should also refreshing the refresh token, maybe by checking the exp time, if its nearly N to expire, then
// refresh it.
func (s *Service) RefreshAccessToken(
	ctx context.Context,
	sessID string,
) (string, error) {
	sess, err := s.authRepo.GetSession(ctx, sessID)
	if err != nil {
		return "", err
	}

	if sess.IsExpired() {
		return "", ErrSessionExpired
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

func (s *Service) Logout(ctx context.Context, sessID string) error {
	session, err := s.authRepo.GetSession(ctx, sessID)
	if err != nil {
		return err
	}

	if session.IsExpired() {
		return nil
	}

	session.Revoke()

	return s.authRepo.UpdateSession(ctx, session)
}

// ForgotPassword initiates the password reset process by sending a reset link or token to the user's email.
func (s *Service) ForgotPassword(ctx context.Context, inp ForgotPasswordInput) error {
	return s.transactor.Transact(ctx, func(ctx context.Context) error {
		usr, err := s.userRepo.GetByEmail(ctx, inp.Email)
		if err != nil {
			if err == user.ErrNotFound {
				return user.ErrEmailNotVerified
			}
			return err
		}

		token, err := NewResetPasswordToken(usr.ID, s.cfg.ResetPasswordTTL)
		if err != nil {
			log.ErrorCtx(ctx, "Failed to create reset password token", err)
			return err
		}

		// Save token to db
		if err := s.authRepo.StoreResetPasswordToken(ctx, token); err != nil {
			return err
		}

		// Publish email job to message queue
		msg := ResetPasswordEmailMessage{
			To:       usr.Email,
			Name:     usr.Name,
			ResetURL: s.cfg.ResetPasswordFormEndpoint + "?token=" + token.Value,
			Expiry:   s.cfg.ResetPasswordTTL,
		}

		return s.publisher.SendResetPasswordEmail(ctx, msg)
	})
}

func (s *Service) GetResetPasswordToken(
	ctx context.Context,
	token string,
) (*ResetPasswordToken, error) {
	return s.authRepo.GetResetPasswordToken(ctx, token)
}

// ResetPassword resets a user's password using a valid reset password token from email.
func (s *Service) ResetPassword(ctx context.Context, inp ResetPasswordInput) error {
	return s.transactor.Transact(ctx, func(ctx context.Context) error {
		token, err := s.authRepo.GetResetPasswordToken(ctx, inp.Token)
		if err != nil {
			return err
		}

		if token.Expired() || token.Used() {
			return ErrResetPasswordTokenInvalid
		}

		user, err := s.userRepo.GetByID(ctx, token.UserID.String())
		if err != nil {
			return err
		}

		newHashedPassword, err := HashPassword(inp.NewPassword)
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
func (s *Service) ChangePassword(
	ctx context.Context,
	userID string,
	inp ChangePasswordInput,
) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := VerifyPassword(inp.Password, u.Password); err != nil {
		return err
	}

	// Hash new password
	newHashedPassword, err := HashPassword(inp.NewPassword)
	if err != nil {
		return err
	}

	u.Password = string(newHashedPassword)

	return s.userRepo.Update(ctx, u)
}

func (s *Service) generateAccessToken(user user.User) (string, error) {
	return SignAccessToken(
		s.cfg.JwtSecret,
		AccessTokenClaims{UserID: user.ID.String()},
		s.cfg.JwtTTL,
	)
}
