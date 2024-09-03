package entity

import (
	"log"
	"sync"
	"time"

	stderrors "errors"

	"github.com/google/uuid"
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

// GenerateToken generates jwt token for user authentication
// Returns the token string and error if any
func (u User) GenerateToken(tokenType auth.TokenType, secretKey string, expiry time.Duration) (string, error) {
	var payload map[string]interface{}

	// If you change one of the map structure, you must adjust the TokenPayload struct from auth package
	switch tokenType {
	case auth.AccessToken:
		payload = map[string]interface{}{
			"user": map[string]interface{}{
				"id":   u.ID.String(),
				"name": u.Name,
			},
			"type": tokenType,
		}
	case auth.RefreshToken:
		payload = map[string]interface{}{
			"user": map[string]interface{}{
				"id": u.ID.String(),
			},
			"type": tokenType,
		}
	default:
		return "", stderrors.New("Invalid token type")
	}

	return auth.TokenEncode(secretKey, payload, expiry)
}

// GenerateTokenPair generates access token and refresh token using goroutines
func (u User) GenerateTokenPair(
	secretKey string,
	accessExpiry time.Duration,
	refreshExpiry time.Duration,
) (
	accessToken string,
	refreshToken string,
	err error,
) {

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		token, e := u.GenerateToken(auth.AccessToken, secretKey, accessExpiry)
		if e != nil {
			errCh <- e
		}
		accessToken = token
		log.Println("Access token generated")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		token, e := u.GenerateToken(auth.RefreshToken, secretKey, refreshExpiry)
		if e != nil {
			errCh <- e
		}
		refreshToken = token
		log.Println("Refresh token generated")
	}()

	wg.Wait()
	close(errCh)

	// Check if any errors occurred
	for e := range errCh {
		if e != nil {
			err = e
			return
		}
	}

	return accessToken, refreshToken, nil
}
