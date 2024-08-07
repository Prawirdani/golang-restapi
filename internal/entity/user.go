package entity

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/common"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/token"
	"github.com/prawirdani/golang-restapi/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorWrongCredentials = errors.Unauthorized("Check your credentials")
	ErrorEmailExists      = errors.Conflict("Email already exists")
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

	if err := validator.Struct(u); err != nil {
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
func (u *User) VerifyPassword(plain string) error {
	if u == nil || u.Password == "" {
		return ErrorWrongCredentials
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	if err != nil {
		return ErrorWrongCredentials
	}
	return nil
}

func (u User) GenerateAccessToken(cfg *config.Config) (string, error) {
	p := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   u.ID.String(),
			"name": u.Name,
		},
		"type": common.AccessToken,
	}

	return token.Encode(cfg.Token.SecretKey, p, cfg.Token.AccessTokenExpiry)
}

func (u User) GenerateRefreshToken(cfg *config.Config) (string, error) {
	p := map[string]interface{}{
		"user": map[string]interface{}{
			"id": u.ID.String(),
		},
		"type": common.RefreshToken,
	}

	return token.Encode(cfg.Token.SecretKey, p, cfg.Token.RefreshTokenExpiry)
}

// GenerateTokenPair generates access token and refresh token using goroutines
func (u User) GenerateTokenPair(cfg *config.Config) (at string, rf string, err error) {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		token, e := u.GenerateAccessToken(cfg)
		if e != nil {
			errCh <- e
		}
		at = token
	}()

	go func() {
		defer wg.Done()
		token, e := u.GenerateRefreshToken(cfg)
		if e != nil {
			errCh <- e
		}
		rf = token
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for e := range errCh {
		if e != nil {
			return "", "", e
		}
	}

	return at, rf, nil
}
