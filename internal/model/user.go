package model

import (
	"github.com/prawirdani/golang-restapi/pkg/strings"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

type CreateUserInput struct {
	Name           string `json:"name"            validate:"required"`
	Email          string `json:"email"           validate:"required,email"`
	Phone          string `json:"phone"`
	Password       string `json:"password"        validate:"required,min=8"`
	RepeatPassword string `json:"repeat_password" validate:"required,eqfield=Password,min=8"`
}

// Implements http.RequestBody
func (r *CreateUserInput) Validate() error {
	return validator.Struct(r)
}

// Implements http.RequestBody
func (r *CreateUserInput) Sanitize() error {
	r.Email = strings.TrimSpaces(r.Email)
	r.Name = strings.TrimSpacesConcat(r.Name)
	r.Phone = strings.TrimSpaces(r.Phone)
	return nil
}
