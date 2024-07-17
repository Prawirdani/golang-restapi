package middleware

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/utils"
)

type MiddlewareManager struct {
	cfg  *config.Config
	Auth AuthMiddleware
}

type AuthMiddleware struct {
	AccessToken  func(next http.Handler) http.Handler
	RefreshToken func(next http.Handler) http.Handler
}

func NewMiddlewareManager(cfg *config.Config) MiddlewareManager {
	mw := MiddlewareManager{
		cfg: cfg,
	}
	mw.Auth.AccessToken = mw.authorize(utils.AccessToken)
	mw.Auth.RefreshToken = mw.authorize(utils.RefreshToken)

	return mw
}
