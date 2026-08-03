package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func Cors(origins []string, allowCredentials, debug bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return cors.Handler(
			cors.Options{
				AllowedOrigins:   origins,
				AllowCredentials: allowCredentials,
				// AllowedHeaders:   []string{"*"},
				AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
				ExposedHeaders: []string{
					"X-RateLimit-Limit",
					"X-RateLimit-Remaining",
					"X-RateLimit-Reset",
					"Retry-After",
					"X-Request-Id",
				},
				AllowedMethods: []string{"OPTIONS", "HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
				MaxAge:         600, // 600s/10m
				Debug:          debug,
			},
		)(next)
	}
}
