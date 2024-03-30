package usecase

import (
	"context"

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
}

func NewAuthUseCase(tokenCfg config.TokenConfig, ur repository.UserRepository) authUseCase {
	return authUseCase{
		tokenConfig: tokenCfg,
		userRepo:    ur,
	}
}

func (u authUseCase) CreateNewUser(ctx context.Context, request model.RegisterRequestPayload) error {
	newUser := entity.NewUser(request)

	if err := newUser.Validate(); err != nil {
		return err
	}

	if err := newUser.EncryptPassword(); err != nil {
		return err
	}

	if err := u.userRepo.InsertUser(ctx, newUser); err != nil {
		return err
	}
	return nil
}

func (u authUseCase) Login(ctx context.Context, request model.LoginRequestPayload) (string, error) {
	var token string

	user, _ := u.userRepo.SelectByEmail(ctx, request.Email)
	if err := user.VerifyPassword(request.Password); err != nil {
		return token, err
	}

	token, err := user.GenerateToken(u.tokenConfig.SecretKey)
	if err != nil {
		return token, err
	}

	return token, nil
}
