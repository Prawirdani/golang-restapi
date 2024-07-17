package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/middleware"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

var handlerFn = httputil.HandlerWrapper

func MapAuthRoutes(r chi.Router, h AuthHandler, mw middleware.MiddlewareManager) {
	r.Post("/auth/register", handlerFn(h.HandleRegister))
	r.Post("/auth/login", handlerFn(h.HandleLogin))
	r.With(mw.Auth.AccessToken).Get("/auth/current", handlerFn(h.CurrentUser))
	r.With(mw.Auth.RefreshToken).Get("/auth/refresh", handlerFn(h.RefreshToken))
}
