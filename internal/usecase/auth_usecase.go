package usecase

import (
	"context"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/repository"
)

type AuthUseCase interface {
	CreateNewUser(ctx context.Context, request model.RegisterRequestPayload) error
	Login(ctx context.Context, request model.LoginRequestPayload) (string, error)
}

type authUseCase struct {
	userRepo    repository.UserRepository
	tokenConfig config.TokenConfig
	ctxTimeout  time.Duration
}

func NewAuthUseCase(cfg *config.Config, ur repository.UserRepository) authUseCase {
	return authUseCase{
		tokenConfig: cfg.Token,
		userRepo:    ur,
		ctxTimeout:  time.Duration(cfg.Context.Timeout * int(time.Second)),
	}
}

func (u authUseCase) CreateNewUser(ctx context.Context, request model.RegisterRequestPayload) error {
	ctxWT, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	newUser := entity.NewUser(request)

	if err := newUser.Validate(); err != nil {
		return err
	}

	if err := newUser.EncryptPassword(); err != nil {
		return err
	}

	if err := u.userRepo.InsertUser(ctxWT, newUser); err != nil {
		return err
	}
	return nil
}

func (u authUseCase) Login(ctx context.Context, request model.LoginRequestPayload) (string, error) {
	ctxWT, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	var token string

	user, _ := u.userRepo.SelectWhere(ctxWT, "email", request.Email)
	if err := user.VerifyPassword(request.Password); err != nil {
		return token, err
	}

	token, err := user.GenerateToken(u.tokenConfig.SecretKey)
	if err != nil {
		return token, err
	}

	return token, nil
}
