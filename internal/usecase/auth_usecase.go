package usecase

import (
	"context"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/repository"
)

type AuthUseCase struct {
	userRepo    repository.UserRepository
	tokenConfig config.TokenConfig
}

func NewAuthUseCase(tokenCfg config.TokenConfig, ur repository.UserRepository) AuthUseCase {
	return AuthUseCase{
		tokenConfig: tokenCfg,
		userRepo:    ur,
	}
}

func (u AuthUseCase) CreateNewUser(ctx context.Context, request model.RegisterRequestPayload) error {
	newUser := entity.NewUser(request)

	if err := newUser.Validate(); err != nil {
		return err
	}

	if err := newUser.EncryptPassword(); err != nil {
		return err
	}

	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return err
	}
	return nil
}

func (u AuthUseCase) Login(ctx context.Context, request model.LoginRequestPayload) (string, error) {
	var token string

	// Query user from database by request email
	user, err := u.userRepo.SelectByEmail(ctx, request.Email)
	if err != nil {
		return token, err
	}

	if err := user.VerifyPassword(request.Password); err != nil {
		return token, err
	}

	token, err = user.GenerateToken(u.tokenConfig.SecretKey)
	if err != nil {
		return token, err
	}

	return token, nil
}
