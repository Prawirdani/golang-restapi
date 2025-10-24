package service

import (
	"context"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/internal/entity/user"
	"github.com/prawirdani/golang-restapi/internal/infra/mq"
	"github.com/prawirdani/golang-restapi/internal/infra/repository"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/contextx"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

type AuthService struct {
	logger      logging.Logger
	cfg         *config.Config
	tx          repository.Transactor
	authRepo    auth.Repository
	userRepo    user.Repository
	userService *UserService
	mqProducer  mq.MessageProducer
}

func NewAuthService(
	cfg *config.Config,
	l logging.Logger,
	transactor repository.Transactor,
	userRepo user.Repository,
	authRepo auth.Repository,
	userService *UserService,
	mqProducer mq.MessageProducer,
) *AuthService {
	return &AuthService{
		cfg:         cfg,
		logger:      l,
		tx:          transactor,
		userRepo:    userRepo,
		authRepo:    authRepo,
		userService: userService,
		mqProducer:  mqProducer,
	}
}

func (s *AuthService) Register(ctx context.Context, i model.CreateUserInput) error {
	hashedPassword, err := auth.HashPassword(i.Password)
	if err != nil {
		return err
	}

	newUser, err := user.New(i, string(hashedPassword))
	if err != nil {
		s.logger.Error(logging.Service, "AuthService.Register", err.Error())
		return err
	}
	if err := s.userRepo.Insert(ctx, newUser); err != nil {
		return err
	}
	return nil
}

// Login is a method to authenticate the user, returning access token, refresh token, and error if any.
func (s *AuthService) Login(
	ctx context.Context,
	request model.LoginInput,
) (string, string, error) {
	u, _ := s.userRepo.GetUserBy(ctx, "email", request.Email)

	if err := auth.VerifyPassword(request.Password, u.Password); err != nil {
		return "", "", err
	}

	accessToken, err := s.generateAccessToken(map[string]any{
		"id":   u.ID.String(),
		"name": u.Name,
	})
	if err != nil {
		s.logger.Error(
			logging.Service,
			"AuthService.Login.GenerateAccessToken",
			err.Error(),
		)
		return "", "", err
	}

	sess, err := auth.NewSession(
		u.ID,
		request.UserAgent,
		s.cfg.Token.RefreshTokenExpiry,
	)
	if err != nil {
		s.logger.Error(logging.Service, "AuthService.Login.NewSession", err.Error())
		return "", "", err
	}

	if err = s.authRepo.InsertSession(ctx, sess); err != nil {
		return "", "", err
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
		s.logger.Error(
			logging.Service,
			"AuthService.RefreshAccessToken.GenerateAccessToken",
			err.Error(),
		)
		return "", err
	}

	return newAccessToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.authRepo.DeleteSession(ctx, "refresh_token", refreshToken)
}

func (s *AuthService) IdentifyUser(ctx context.Context) (user.User, error) {
	tokenPayload, err := contextx.GetAuthCtx(ctx)
	if err != nil {
		return user.User{}, err
	}

	u, err := s.userService.GetUserByID(ctx, tokenPayload["id"].(string))
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// ForgotPassword initiates the password reset process by sending a reset link or token to the user's email.
func (s *AuthService) ForgotPassword(ctx context.Context, i model.ForgotPasswordInput) error {
	return s.tx.Transact(ctx, func(ctx context.Context) error {
		u, err := s.userRepo.GetUserBy(ctx, "email", i.Email)
		if err != nil {
			if err == user.ErrUserNotFound {
				return errors.Forbidden("Email is not registered or not verified")
			}
			return err
		}

		expiresAt := time.Now().Add(s.cfg.Token.ResetPasswordTokenExpiry)
		token, err := auth.NewResetPasswordToken(u.ID, expiresAt)
		if err != nil {
			s.logger.Error(
				logging.Service,
				"AuthService.ForgotPassword.NewResetPasswordToken",
				err.Error(),
			)
			return err
		}

		// Save token to db
		if err := s.authRepo.InsertResetPasswordToken(ctx, *token); err != nil {
			return err
		}

		// Publish email job to message queue
		emailJob := mq.EmailResetPasswordJob{
			Type:      "password_reset",
			To:        u.Email,
			Name:      u.Name,
			ResetURL:  s.cfg.Token.ResetPasswordFormEndpoint + "?token=" + token.Value,
			ExpiryMin: int(s.cfg.Token.ResetPasswordTokenExpiry.Minutes()),
		}

		// Non-blocking - just queue the job
		return s.mqProducer.Publish(ctx, mq.EmailResetPasswordeJobKey, emailJob)
	})
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
			s.logger.Error(
				logging.Service,
				"AuthService.ResetPassword.NewPassword",
				err.Error(),
			)
			return err
		}
		user.Password = string(newHashedPassword)

		if err := s.authRepo.UseResetPasswordToken(ctx, token); err != nil {
			return err
		}

		return s.userRepo.UpdateUser(ctx, user)
	})
}

// ChangePassword updates the authenticated user's password after verifying the current password.
func (s *AuthService) ChangePassword(ctx context.Context, i model.ChangePasswordInput) error {
	tokenPayload, err := contextx.GetAuthCtx(ctx)
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
