package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID
	Name      string `validate:"required"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=6"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(request model.RegisterRequestPayload) User {
	return User{
		ID:       uuid.New(),
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
}

func (u User) Validate() error {
	return utils.Validate.Struct(u)
}

func (u *User) EncryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u User) VerifyPassword(plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	if err != nil {
		return httputil.ErrUnauthorized("check your credentials")
	}
	return nil
}

func (u User) GenerateToken(secret string) (string, error) {
	return utils.GenerateToken(u.ID.String(), secret, 60*time.Minute)
}
