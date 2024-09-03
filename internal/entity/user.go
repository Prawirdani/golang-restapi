package entity

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongCredentials = errors.Unauthorized("Check your credentials")
	ErrEmailExist       = errors.Conflict("Email already exists")
	ErrUserNotFound     = errors.NotFound("User not found")
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Name      string    `db:"name" json:"name" validate:"required"`
	Email     string    `db:"email" json:"email" validate:"required,email"`
	Password  string    `db:"password" json:"-" validate:"required,min=6"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
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
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	if err != nil {
		return ErrWrongCredentials
	}
	return nil
}

// Generate access token for user
func (u User) GenerateAccessToken(cfg *config.Config) (string, error) {
	// If you change the map structure, you must adjust the AccessTokenPayload struct from auth package
	p := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   u.ID.String(),
			"name": u.Name,
		},
		"type": auth.AccessToken,
	}

	return auth.TokenEncode(cfg.Token.SecretKey, p, cfg.Token.AccessTokenExpiry)
}

// Generate refresh token for user
func (u User) GenerateRefreshToken(cfg *config.Config) (string, error) {
	// Same as GenerateAccessToken, if you change the map structure, you must adjust the RefreshTokenPayload struct from auth package
	p := map[string]interface{}{
		"user": map[string]interface{}{
			"id": u.ID.String(),
		},
		"type": auth.RefreshToken,
	}

	return auth.TokenEncode(cfg.Token.SecretKey, p, cfg.Token.RefreshTokenExpiry)
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
