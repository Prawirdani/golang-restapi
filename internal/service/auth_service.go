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

func (u *AuthService) Login(
	ctx context.Context,
	request model.LoginRequest,
) (accessToken string, refreshToken string, err error) {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, _ := u.userRepo.SelectWhere(ctxWT, "email", request.Email)
	if err = user.VerifyPassword(request.Password); err != nil {
		return
	}

	accessToken, refreshToken, err = user.GenerateTokenPair(
		u.cfg.Token.SecretKey,
		u.cfg.Token.AccessTokenExpiry,
		u.cfg.Token.RefreshTokenExpiry,
	)
	if err != nil {
		u.logger.Error(logging.Service, "AuthService.Login.GenerateTokenPair", err.Error())
		return
	}

	return
}

// TODO: Should also refreshing the refresh token, maybe by checking the exp time, if its nearly N to expire, then refresh it.
func (u *AuthService) RefreshToken(ctx context.Context) (string, error) {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	tokenPayload, err := auth.GetContext(ctxWT)
	if err != nil {
		return "", err
	}

	user, err := u.userRepo.SelectWhere(ctxWT, "id", tokenPayload["id"])
	if err != nil {
		return "", err
	}
	accessToken, err := user.GenerateToken(
		auth.AccessToken,
		u.cfg.Token.SecretKey,
		u.cfg.Token.AccessTokenExpiry,
	)
	if err != nil {
		u.logger.Error(logging.Service, "AuthService.RefreshToken.GenerateAccessToken", err.Error())
		return "", err
	}
	return accessToken, nil
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
