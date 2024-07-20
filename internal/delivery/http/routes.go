package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/delivery/http/handler"
	"github.com/prawirdani/golang-restapi/internal/delivery/http/helper"
	"github.com/prawirdani/golang-restapi/internal/delivery/http/middleware"
)

func registerAuthRoutes(r chi.Router, h *handler.AuthHandler, mw *middleware.Collection) {
	r.Post("/auth/register", helper.HandlerFn(h.HandleRegister))
	r.Post("/auth/login", helper.HandlerFn(h.HandleLogin))
	r.With(mw.AuthAccessToken).Get("/auth/current", helper.HandlerFn(h.CurrentUser))
	r.With(mw.AuthRefreshToken).Get("/auth/refresh", helper.HandlerFn(h.RefreshToken))
}
