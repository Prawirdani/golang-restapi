package service

import (
	"context"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/repository"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
	timeout  time.Duration
	logger   logging.Logger
}

func NewAuthService(
	cfg *config.Config,
	l logging.Logger,
	ur *repository.UserRepository,
) *AuthService {
	return &AuthService{
		cfg:      cfg,
		logger:   l,
		userRepo: ur,
		timeout:  time.Duration(5 * int(time.Second)),
	}
}

func (u *AuthService) Register(ctx context.Context, request model.RegisterRequest) error {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	newUser, err := entity.NewUser(request)
	if err != nil {
		u.logger.Error(logging.Service, "AuthService.Register", err.Error())
		return err
	}
	if err := u.userRepo.InsertUser(ctxWT, newUser); err != nil {
		return err
	}
	return nil
}

// Login is a method to authenticate the user, returning access token, refresh token, and error if any.
func (u *AuthService) Login(
	ctx context.Context,
	request model.LoginRequest,
) (string, string, error) {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, _ := u.userRepo.SelectWhere(ctxWT, "email", request.Email)
	if err := user.VerifyPassword(request.Password); err != nil {
		return "", "", err
	}

	accessToken, err := user.GenerateAccessToken(
		u.cfg.Token.SecretKey,
		u.cfg.Token.AccessTokenExpiry,
	)
	if err != nil {
		u.logger.Error(logging.Service, "AuthService.Login.GenerateAccessToken", err.Error())
		return "", "", err
	}

	sess, err := auth.NewSession(
		user.ID,
		request.UserAgent,
		u.cfg.Token.RefreshTokenExpiry,
	)
	if err != nil {
		u.logger.Error(logging.Service, "AuthService.Login.NewSession", err.Error())
		return "", "", err
	}

	if err = u.userRepo.InsertSession(ctxWT, sess); err != nil {
		return "", "", err
	}

	return accessToken, sess.RefreshToken, nil
}

// TODO: Should also refreshing the refresh token, maybe by checking the exp time, if its nearly N to expire, then refresh it.
func (u *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	sess, err := u.userRepo.SelectSession(ctxWT, "refresh_token", refreshToken)
	if err != nil {
		return "", err
	}

	if sess.IsExpired() {
		_ = u.userRepo.DeleteSession(ctxWT, "id", sess.ID)
		return "", auth.ErrSessionExpired
	}

	user, err := u.userRepo.SelectWhere(ctxWT, "id", sess.UserID)
	if err != nil {
		return "", err
	}

	accessToken, err := user.GenerateAccessToken(
		u.cfg.Token.SecretKey,
		u.cfg.Token.AccessTokenExpiry,
	)
	if err != nil {
		u.logger.Error(
			logging.Service,
			"AuthService.RefreshAccessToken.GenerateAccessToken",
			err.Error(),
		)
		return "", err
	}

	return accessToken, nil
}

func (u *AuthService) Logout(ctx context.Context, refreshToken string) error {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.userRepo.DeleteSession(ctxWT, "refresh_token", refreshToken)
}

func (u *AuthService) IdentifyUser(ctx context.Context) (entity.User, error) {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	tokenPayload, err := auth.GetContext(ctxWT)
	if err != nil {
		return entity.User{}, err
	}

	user, err := u.userRepo.SelectWhere(ctxWT, "id", tokenPayload["id"])
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
