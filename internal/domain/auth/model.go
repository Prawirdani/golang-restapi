// Package auth provides authentication and authorization functionality.
// This package handles user authentication through sessions, access tokens, and
// password management including secure hashing and password reset flows. It manages
// the complete authentication lifecycle from login through logout, including token
// generation, validation, and session management.
package auth

import (
	"time"

	"github.com/prawirdani/golang-restapi/pkg/strings"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

type RegisterInput struct {
	Name           string `json:"name"            validate:"required"`
	Email          string `json:"email"           validate:"required,email"`
	Phone          string `json:"phone"`
	Password       string `json:"password"        validate:"required,min=8"`
	RepeatPassword string `json:"repeat_password" validate:"required,eqfield=Password,min=8"`
}

// Validate implements [handler.JSONRequestBody]
func (r *RegisterInput) Validate() error {
	return validator.Struct(r)
}

// Sanitize implements [handler.JSONRequestBody]
func (r *RegisterInput) Sanitize() error {
	r.Email = strings.TrimSpaces(r.Email)
	r.Name = strings.TrimSpacesConcat(r.Name)
	r.Phone = strings.TrimSpaces(r.Phone)
	return nil
}

type LoginInput struct {
	Email     string `json:"email"    validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	UserAgent string
}

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordInput struct {
	Token             string `json:"token"               validate:"required"`
	NewPassword       string `json:"new_password"        validate:"required,min=8"`
	RepeatNewPassword string `json:"repeat_new_password" validate:"required,eqfield=NewPassword"`
}

type ChangePasswordInput struct {
	Password          string
	NewPassword       string `json:"new_password"        validate:"required,min=8"`
	RepeatNewPassword string `json:"repeat_new_password" validate:"required,eqfield=NewPassword"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ResetPasswordEmailMessage struct {
	To       string        `json:"to"`         // Recipient's email address
	Name     string        `json:"name"`       // Recipient's name
	ResetURL string        `json:"reset_url"`  // Link for resetting the password
	Expiry   time.Duration `json:"expiry_min"` // Expiration time of the reset token in minutes
}
