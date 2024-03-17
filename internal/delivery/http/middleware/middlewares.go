package middleware

import "github.com/prawirdani/golang-restapi/pkg/utils"

type Collection struct {
	jwtAuth *utils.JWTProvider
}

func New(jp *utils.JWTProvider) *Collection {
	return &Collection{
		jwtAuth: jp,
	}
}
