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
	Register(ctx context.Context, request model.RegisterRequest) error
	Login(ctx context.Context, request model.LoginRequest) (string, error)
}

type authUseCase struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthUseCase(cfg *config.Config, ur repository.UserRepository) authUseCase {
	return authUseCase{
		cfg:      cfg,
		userRepo: ur,
	}
}

func (u authUseCase) Register(ctx context.Context, request model.RegisterRequest) error {
	ctxWT, cancel := context.WithTimeout(ctx, time.Duration(u.cfg.Context.Timeout*int(time.Second)))
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

func (u authUseCase) Login(ctx context.Context, request model.LoginRequest) (string, error) {
	ctxWT, cancel := context.WithTimeout(ctx, time.Duration(u.cfg.Context.Timeout*int(time.Second)))
	defer cancel()

	var token string

	user, _ := u.userRepo.SelectWhere(ctxWT, "email", request.Email)
	if err := user.VerifyPassword(request.Password); err != nil {
		return token, err
	}

	token, err := user.GenerateToken(u.cfg.Token.SecretKey)
	if err != nil {
		return token, err
	}

	return token, nil
}
