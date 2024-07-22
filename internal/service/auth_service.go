package service

import (
	"context"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
	timeout  time.Duration
}

func NewAuthService(cfg *config.Config, ur *repository.UserRepository) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userRepo: ur,
		timeout:  time.Duration(5 * int(time.Second)),
	}
}

func (u *AuthService) Register(ctx context.Context, request model.RegisterRequest) error {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	newUser, err := entity.NewUser(request)
	if err != nil {
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

	user, _ := u.userRepo.SelectWhere(ctxWT, "email", request.Email)
	if err := user.VerifyPassword(request.Password); err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err = user.GenerateTokenPair(u.cfg)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *AuthService) RefreshToken(ctx context.Context, userID string) (string, error) {
	ctxWT, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	user, _ := u.userRepo.SelectWhere(ctxWT, "id", userID)
	return user.GenerateAccessToken(u.cfg)
}
