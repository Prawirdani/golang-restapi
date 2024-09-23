package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/app/middleware"
	"github.com/prawirdani/golang-restapi/internal/handler"
	"github.com/prawirdani/golang-restapi/pkg/response"
)

// HandlerFn is a custom handler wrapper function type.
// It is used to wrap handler function to make it easier handling errors from handler function.
// Also help handler code to be more readable without to many early return statement.
// Every handler function should use this function signature.
type HandlerFn func(w http.ResponseWriter, r *http.Request) error

// handlerFn is a helper function to handle error in CustomHandler function
func handlerFn(fn HandlerFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			response.HandleError(w, err)
			return
		}
	}
}

func registerAuthRoutes(r chi.Router, h *handler.AuthHandler, mw *middleware.Collection) {
	r.Post("/auth/register", handlerFn(h.HandleRegister))
	r.Post("/auth/login", handlerFn(h.HandleLogin))
	r.Delete("/auth/logout", handlerFn(h.HandleLogout))
	r.With(mw.Auth).Get("/auth/current", handlerFn(h.CurrentUser))
	r.Get("/auth/refresh", handlerFn(h.RefreshToken))
}
