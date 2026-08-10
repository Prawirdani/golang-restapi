package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// InstrumentHandler is Prometheus metrics instrumentation middleware
func (m *Metrics) InstrumentHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// M6/L6: chi's wrapper preserves optional interfaces (Flusher, Hijacker,
		// ReaderFrom) that a naive struct embedding drops, and exposes
		// Status()/BytesWritten() for the metrics.
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			// Use the chi route template (e.g. "/users/{id}") rather than the raw
			// URL path to keep label cardinality bounded. An unmatched route yields
			// an empty pattern, which we bucket as "unknown".
			path := chi.RouteContext(r.Context()).RoutePattern()
			if path == "" {
				path = "unknown"
			}
			duration := time.Since(start).Seconds()
			status := ww.Status()
			if status == 0 { // handler never wrote a header
				status = http.StatusOK
			}
			m.ReqDuration.WithLabelValues(path, r.Method, strconv.Itoa(status)).Observe(duration)
			m.ReqCounter.WithLabelValues(path, r.Method, strconv.Itoa(status)).Inc()
		}()
		next.ServeHTTP(ww, r)
	})
}
