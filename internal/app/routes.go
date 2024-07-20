package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/app/middleware"
	"github.com/prawirdani/golang-restapi/internal/handler"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/response"
)

// handlerFn is a helper function to handle error in CustomHandler function
func handlerFn(fn httputil.CustomHandler) http.HandlerFunc {
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
	r.With(mw.AuthAccessToken).Get("/auth/current", handlerFn(h.CurrentUser))
	r.With(mw.AuthRefreshToken).Get("/auth/refresh", handlerFn(h.RefreshToken))
}
