package middleware

import (
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func (m *Collection) Gzip(next http.Handler) http.Handler {
	return chiMiddleware.Compress(6)(next)
}
