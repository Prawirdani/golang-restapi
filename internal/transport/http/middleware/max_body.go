package middleware

import "net/http"

// Middleware to limit body size
func MaxBodySizeMiddleware(maxBodySize int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
			next.ServeHTTP(w, r)
		})
	}
}
