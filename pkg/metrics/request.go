package metrics

import (
	"net/http"
	"strconv"
	"time"
)

type writerRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (w *writerRecorder) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Prometheus metrics instrumentation middleware
func (m *Metrics) Instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &writerRecorder{w, http.StatusOK}
		defer func() {
			duration := time.Since(start).Seconds()
			m.ReqDuration.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(ww.statusCode)).Observe(duration)
			m.ReqCounter.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(ww.statusCode)).Inc()
		}()
		next.ServeHTTP(ww, r)
	})
}
