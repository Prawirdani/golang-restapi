package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/transport/http/handler"
)

var fn = handler.Handler

type authMiddleware = func(next http.Handler) http.Handler

func RegisterAuthRoutes(r chi.Router, h *handler.AuthHandler, authMw authMiddleware) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", fn(h.LoginHandler))
		r.Post("/register", fn(h.RegisterHandler))
		r.Delete("/logout", fn(h.LogoutHandler))
		r.Post("/password/forgot", fn(h.ForgotPasswordHandler))
		r.Get("/password/reset/{token}", fn(h.GetResetPasswordTokenHandler))
		r.Post("/password/reset", fn(h.ResetPasswordHandler))

		r.Get("/refresh", fn(h.RefreshTokenHandler))
		r.With(authMw).Group(func(r chi.Router) {
			r.Get("/current", fn(h.CurrentUserHandler))
			r.Post("/password/change", fn(h.ChangePasswordHandler))
		})
	})
}

func RegisterUserRoutes(r chi.Router, h *handler.UserHandler, authMw authMiddleware) {
	r.With(authMw).Route("/users", func(r chi.Router) {
		r.Post("/profile/upload", fn(h.ChangeProfilePictureHandler))
	})
}
