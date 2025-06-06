package entity

import (
	"time"

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
	ID        uuid.UUID `db:"id"         json:"id"         validate:"required,uuid"`
	Name      string    `db:"name"       json:"name"       validate:"required"`
	Email     string    `db:"email"      json:"email"      validate:"required,email"`
	Password  string    `db:"password"   json:"-"          validate:"required,min=6"`
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

func (u *User) GenerateAccessToken(secret string, exp time.Duration) (string, error) {
	payload := map[string]interface{}{
		"id":   u.ID.String(),
		"name": u.Name,
	}

	return auth.GenerateJWT(secret, exp, payload)
}
