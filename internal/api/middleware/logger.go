package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (m *Collection) Logger(next http.Handler) http.Handler {
	return middleware.Logger(next)
}
