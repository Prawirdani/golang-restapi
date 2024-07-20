package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/app/middleware"
	"github.com/prawirdani/golang-restapi/internal/handler"
)

func registerAuthRoutes(r chi.Router, h *handler.AuthHandler, mw *middleware.Collection) {
	r.Post("/auth/register", handlerFn(h.HandleRegister))
	r.Post("/auth/login", handlerFn(h.HandleLogin))
	r.With(mw.AuthAccessToken).Get("/auth/current", handlerFn(h.CurrentUser))
	r.With(mw.AuthRefreshToken).Get("/auth/refresh", handlerFn(h.RefreshToken))
}
