package middleware

import (
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func (c *Collection) ReqLogger(next http.Handler) http.Handler {
	return chiMiddleware.Logger(next)
}
