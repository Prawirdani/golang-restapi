package model

import "github.com/prawirdani/golang-restapi/pkg/utils"

type RegisterRequestPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (r *RegisterRequestPayload) Validate() error {
	return utils.ValidateStruct(r)
}

type LoginRequestPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginRequestPayload) Validate() error {
	return utils.ValidateStruct(r)
}
