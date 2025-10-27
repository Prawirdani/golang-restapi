package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

const (
	HeaderXRequestID = "X-Request-ID"
)

// RequestID attaches or generates an X-Request-ID header per request.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx := log.WithContext(r.Context(), "request_id", reqID)
		w.Header().Set(HeaderXRequestID, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
