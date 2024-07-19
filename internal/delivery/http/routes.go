package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/delivery/http/middleware"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

var handlerFn = httputil.HandlerWrapper

func RegisterAuthRoutes(r chi.Router, h *AuthHandler, mw *middleware.Collection) {
	r.Post("/auth/register", handlerFn(h.HandleRegister))
	r.Post("/auth/login", handlerFn(h.HandleLogin))
	r.With(mw.AuthAccessToken).Get("/auth/current", handlerFn(h.CurrentUser))
	r.With(mw.AuthRefreshToken).Get("/auth/refresh", handlerFn(h.RefreshToken))
}
