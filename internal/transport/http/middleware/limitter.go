package middleware

import (
	"net/http"
	"time"

	m "github.com/go-chi/httprate"
	httperr "github.com/prawirdani/golang-restapi/internal/transport/http/error"
	"github.com/prawirdani/golang-restapi/internal/transport/http/handler"
)

func RateLimit(reqLimit int, windowLength time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.Limit(reqLimit, windowLength,
			m.WithKeyFuncs(m.KeyByIP, m.KeyByEndpoint),
			m.WithLimitHandler(handler.Handler(func(ctx *handler.Context) error {
				return httperr.New(
					http.StatusTooManyRequests,
					"too many request, try again later",
					nil,
				)
			})),
		)(next)
	}
}
