// Package user provides the domain model and business logic for managing users in system.
package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/internal/domain"
	"github.com/prawirdani/golang-restapi/pkg/nullable"
)

var (
	ErrRequiredName     = domain.ErrValidation("Name is required")
	ErrRequiredEmail    = domain.ErrValidation("Email is required")
	ErrRequiredPassword = domain.ErrValidation("Password is required")
	ErrEmailExists      = domain.ErrDuplicate("Email already exists")
	ErrNotFound         = domain.ErrNotFound("User not found")
	ErrEmailNotVerified = domain.ErrForbidden("Email is not registered or not verified")
)

type User struct {
	ID           uuid.UUID                 `db:"id"            json:"id"`
	Name         string                    `db:"name"          json:"name"`
	Email        string                    `db:"email"         json:"email"`
	Password     string                    `db:"password"      json:"-"`
	Phone        nullable.Nullable[string] `db:"phone"         json:"phone"`
	ProfileImage nullable.Nullable[string] `db:"profile_image" json:"profile_image"`
	CreatedAt    time.Time                 `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time                 `db:"updated_at"    json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrRequiredName
	}
	if u.Email == "" {
		return ErrRequiredEmail
	}
	if u.Password == "" {
		return ErrRequiredPassword
	}

	return nil
}

// New creates new user, returns an error if validation fails.
func New(name, email, phone, hashedPassword string) (*User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	u := User{
		ID:        id,
		Name:      name,
		Email:     email,
		Phone:     nullable.New(phone, false),
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return &u, nil
}
