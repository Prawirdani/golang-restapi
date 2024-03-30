package model

import "github.com/prawirdani/golang-restapi/pkg/utils"

type RegisterRequestPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Syntatic Validator
func (r RegisterRequestPayload) ValidateRequest() error {
	return utils.Validate.Struct(r)
}

type LoginRequestPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Syntatic Validator
func (r LoginRequestPayload) ValidateRequest() error {
	return utils.Validate.Struct(r)
}

type TokenResponse struct {
	Token string `json:"token"`
}

type TokenInfoResponse struct {
	TokenInfo interface{} `json:"tokenInfo"`
}
