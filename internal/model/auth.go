package model

import "github.com/prawirdani/golang-restapi/pkg/utils"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Syntatic Validator
func (r RegisterRequest) ValidateRequest() error {
	return utils.Validate.Struct(r)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Syntatic Validator
func (r LoginRequest) ValidateRequest() error {
	return utils.Validate.Struct(r)
}

type TokenResponse struct {
	Token string `json:"token"`
}

type TokenInfoResponse struct {
	TokenInfo interface{} `json:"tokenInfo"`
}
