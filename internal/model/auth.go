package model

import "github.com/prawirdani/golang-restapi/pkg/common"

type RegisterRequest struct {
	Name           string `json:"name" validate:"required,min=3"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6"`
	RepeatPassword string `json:"repeatPassword" validate:"required,eqfield=Password,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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
