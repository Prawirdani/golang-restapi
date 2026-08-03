package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type writerRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (w *writerRecorder) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// InstrumentHandler is Prometheus metrics instrumentation middleware
func (m *Metrics) InstrumentHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &writerRecorder{w, http.StatusOK}
		defer func() {
			// Use the chi route template (e.g. "/users/{id}") rather than the raw
			// URL path to keep label cardinality bounded. An unmatched route yields
			// an empty pattern, which we bucket as "unknown".
			path := chi.RouteContext(r.Context()).RoutePattern()
			if path == "" {
				path = "unknown"
			}
			duration := time.Since(start).Seconds()
			status := strconv.Itoa(ww.statusCode)
			m.ReqDuration.WithLabelValues(path, r.Method, status).Observe(duration)
			m.ReqCounter.WithLabelValues(path, r.Method, status).Inc()
		}()
		next.ServeHTTP(ww, r)
	})
}
