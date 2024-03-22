package usecase

import (
	"context"

	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/repository"
)

type AuthUseCase struct {
	userRepo *repository.UserRepository
}

func NewAuthUseCase(ur *repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo: ur,
	}
}

func (u AuthUseCase) CreateNewUser(ctx context.Context, request model.RegisterRequestPayload) error {
	newUser := entity.NewUser(request)

	// Validate Struct
	if err := newUser.Validate(); err != nil {
		return err
	}

	// Encrypt password
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

	//TODO validate LoginRequest

	// Query user from database by request email
	user, err := u.userRepo.SelectByEmail(ctx, request.Email)
	if err != nil {
		return token, err
	}

	if err := user.VerifyPassword(request.Password); err != nil {
		return token, err
	}

	token, err = user.GenerateToken("secret")
	if err != nil {
		return token, err
	}

	return token, nil
}
