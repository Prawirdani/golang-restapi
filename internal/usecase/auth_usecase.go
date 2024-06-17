package usecase

import (
	"context"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/repository"
	"github.com/prawirdani/golang-restapi/pkg/utils"
)

type AuthUseCase interface {
	Register(ctx context.Context, request model.RegisterRequest) error
	Login(ctx context.Context, request model.LoginRequest) ([]utils.JWT, error)
	RefreshToken(ctx context.Context, userID string) (utils.JWT, error)
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

func (u authUseCase) Login(ctx context.Context, request model.LoginRequest) ([]utils.JWT, error) {
	ctxWT, cancel := context.WithTimeout(ctx, time.Duration(u.cfg.Context.Timeout*int(time.Second)))
	defer cancel()

	user, _ := u.userRepo.SelectWhere(ctxWT, "email", request.Email)
	if err := user.VerifyPassword(request.Password); err != nil {
		return nil, err
	}

	return user.GenerateTokenPair(u.cfg)
}

func (u authUseCase) RefreshToken(ctx context.Context, userID string) (utils.JWT, error) {
	ctxWT, cancel := context.WithTimeout(ctx, time.Duration(u.cfg.Context.Timeout*int(time.Second)))
	defer cancel()

	user, _ := u.userRepo.SelectWhere(ctxWT, "id", userID)
	return user.GenerateRefreshToken(u.cfg)
}
