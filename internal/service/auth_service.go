package service

import (
	"context"
	"time"

	"github.com/prawirdani/golang-restapi/config"
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

func NewAuthService(cfg *config.Config, l logging.Logger, ur *repository.UserRepository) *AuthService {
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

	user, err := u.userRepo.SelectWhere(ctxWT, repository.UserEmail, request.Email)
	if user != nil {
		return entity.ErrorEmailExists
	}

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

func (u *AuthService) Login(ctx context.Context, request model.LoginRequest) (accessToken string, refreshToken string, err error) {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, _ := u.userRepo.SelectWhere(ctxWT, repository.UserEmail, request.Email)
	if err := user.VerifyPassword(request.Password); err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err = user.GenerateTokenPair(u.cfg)
	if err != nil {
		u.logger.Error(logging.Service, "AuthService.Login.GenerateTokenPair", err.Error())
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *AuthService) RefreshToken(ctx context.Context, userID string) (string, error) {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, err := u.userRepo.SelectWhere(ctxWT, "id", userID)
	if err != nil {
		return "", err
	}

	accessToken, err := user.GenerateAccessToken(u.cfg)

	if err != nil {
		u.logger.Error(logging.Service, "AuthService.RefreshToken.GenerateAccessToken", err.Error())
		return "", err
	}

	return accessToken, nil
}
