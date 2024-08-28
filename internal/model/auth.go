package model

import (
	"github.com/prawirdani/golang-restapi/pkg/common"
	"github.com/prawirdani/golang-restapi/pkg/sanitizer"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

type RegisterRequest struct {
	Name           string `json:"name" validate:"required,min=3"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6"`
	RepeatPassword string `json:"repeatPassword" validate:"required,eqfield=Password,min=6"`
}

func (r *RegisterRequest) Validate() error {
	return validator.Struct(r)
}

func (r *RegisterRequest) Sanitize() error {
	r.Email = sanitizer.TrimSpaces(r.Email)
	r.Name = sanitizer.TrimSpacesConcat(r.Name)
	return nil
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (l *LoginRequest) Validate() error {
	return validator.Struct(l)
}

func (l *LoginRequest) Sanitize() error {
	return nil
}

type AccessTokenPayload struct {
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Type common.TokenType `json:"type"`
}

type RefreshTokenPayload struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	Type common.TokenType `json:"type"`
}
