package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/transport/http/handler"
	"github.com/prawirdani/golang-restapi/internal/transport/http/middleware"
	res "github.com/prawirdani/golang-restapi/internal/transport/http/response"
)

// HandlerFn is a custom handler wrapper function type.
// It is used to wrap handler function to make it easier handling errors from handler function.
// Also help handler code to be more readable without to many early return statement.
// Every handler function should use this function signature.
type HandlerFn func(w http.ResponseWriter, r *http.Request) error

// fn is an adapter to use custom HandlerFn with std http.Handler
func fn(fn HandlerFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			res.HandleError(w, err)
			return
		}
	}
}

func RegisterAuthRoutes(r chi.Router, h *handler.AuthHandler, mw *middleware.Collection) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", fn(h.LoginHandler))
		r.Post("/register", fn(h.RegisterHandler))
		r.Delete("/logout", fn(h.LogoutHandler))
		r.Post("/password/forgot", fn(h.ForgotPasswordHandler))
		r.Get("/password/reset/{token}", fn(h.GetResetPasswordTokenHandler))
		r.Post("/password/reset", fn(h.ResetPasswordHandler))

		r.With(mw.Auth).Group(func(r chi.Router) {
			r.Get("/refresh", fn(h.RefreshTokenHandler))
			r.Get("/current", fn(h.CurrentUserHandler))
			r.Post("/password/change", fn(h.ChangePasswordHandler))
		})
	})
}

func RegisterUserRoutes(r chi.Router, h *handler.UserHandler, mw *middleware.Collection) {
	r.With(mw.Auth).Route("/users", func(r chi.Router) {
		r.Post("/profile/upload", fn(h.ChangeProfilePictureHandler))
	})
}
