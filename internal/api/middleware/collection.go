package middleware

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/token"
)

type Collection struct {
	cfg  *config.Config
	Auth authMiddleware
}

type authMiddleware struct {
	AccessToken  func(next http.Handler) http.Handler
	RefreshToken func(next http.Handler) http.Handler
}

func NewCollection(cfg *config.Config) *Collection {
	mw := Collection{
		cfg: cfg,
	}
	mw.Auth.AccessToken = mw.authorize(token.Access)
	mw.Auth.RefreshToken = mw.authorize(token.Refresh)

	return &mw
}
