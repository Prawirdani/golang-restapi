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
				Debug:            debug,
			},
		)(next)
	}
}
