package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/token"
	"github.com/prawirdani/golang-restapi/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorWrongCredentials = httputil.ErrUnauthorized("Check your credentials")
)

type User struct {
	ID        uuid.UUID `validate:"required,uuid"`
	Name      string    `validate:"required"`
	Email     string    `validate:"required,email"`
	Password  string    `validate:"required,min=6"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Create new user from request payload
func NewUser(request model.RegisterRequest) (User, error) {
	u := User{
		ID:       uuid.New(),
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	if err := utils.Validate(u); err != nil {
		return User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	u.Password = string(hashedPassword)

	return u, nil
}

// Verify / Decrypt user password
func (u User) VerifyPassword(plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	if err != nil {
		return ErrorWrongCredentials
	}
	return nil
}

func (u User) GenerateAccessToken(cfg *config.Config) (token.JWT, error) {
	payload := map[string]interface{}{
		"id":   u.ID.String(),
		"name": u.Name,
	}
	return token.GenerateJWT(cfg, payload, token.Access)
}

func (u User) GenerateRefreshToken(cfg *config.Config) (token.JWT, error) {
	payload := map[string]interface{}{
		"id": u.ID.String(),
	}
	return token.GenerateJWT(cfg, payload, token.Refresh)
}

func (u User) GenerateTokenPair(cfg *config.Config) ([]token.JWT, error) {
	accessToken, err := u.GenerateAccessToken(cfg)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.GenerateRefreshToken(cfg)
	if err != nil {
		return nil, err
	}

	return []token.JWT{accessToken, refreshToken}, nil
}
