package middleware

import (
	"net/http"
	"time"

	m "github.com/go-chi/httprate"
	res "github.com/prawirdani/golang-restapi/internal/transport/http/response"
	err "github.com/prawirdani/golang-restapi/pkg/errors"
)

func _handler(w http.ResponseWriter, r *http.Request) {
	res.HandleError(
		w,
		err.TooManyRequest("Too many request, try again later"),
	)
}

func (c *Collection) RateLimit(
	reqLimit int,
	windowLength time.Duration,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.Limit(reqLimit, windowLength,
			m.WithKeyFuncs(m.KeyByIP, m.KeyByEndpoint),
			m.WithLimitHandler(_handler),
		)(next)
	}
}
