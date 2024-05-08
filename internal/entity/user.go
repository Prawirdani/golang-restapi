package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
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
func NewUser(request model.RegisterRequest) User {
	return User{
		ID:       uuid.New(),
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
}

// Semantic Validation
func (u User) Validate() error {
	return utils.Validate.Struct(u)
}

// Encrypt user password
func (u *User) EncryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Verify / Decrypt user password
func (u User) VerifyPassword(plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	if err != nil {
		return ErrorWrongCredentials
	}
	return nil
}

// Generate JWT Token
func (u User) GenerateToken(secret string, expiryHour int) (string, error) {
	return utils.GenerateToken(u.ID.String(), u.Name, secret, time.Duration(expiryHour)*time.Hour)
}
