package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

func (c *Collection) RateLimit(next http.Handler) http.Handler {
	return httprate.Limit(10, 10*time.Second, httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint))(next)
}
