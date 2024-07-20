package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func (c *Collection) Cors(next http.Handler) http.Handler {
	return cors.Handler(
		cors.Options{
			AllowedOrigins:   c.cfg.Cors.Origins,
			AllowCredentials: c.cfg.Cors.Credentials,
			Debug:            c.cfg.IsProduction(),
		},
	)(next)
}
