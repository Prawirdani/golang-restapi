package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func (m *Collection) Cors(next http.Handler) http.Handler {
	return cors.Handler(
		cors.Options{
			AllowedOrigins:   m.cfg.Cors.Origins,
			AllowCredentials: m.cfg.Cors.Credentials,
			Debug:            m.cfg.IsProduction(),
		},
	)(next)
}
