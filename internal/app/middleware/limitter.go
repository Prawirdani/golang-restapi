package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

func (c *Collection) RateLimit(next http.Handler) http.Handler {
	// Adjust as needed
	// This will limit to 10 requests per 10 seconds per IP address on the same endpoint
	return httprate.Limit(10, 10*time.Second, httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint))(next)
}
