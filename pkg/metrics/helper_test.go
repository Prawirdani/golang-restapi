package metrics

import "github.com/prometheus/client_golang/prometheus"

// newTestMetrics builds the metric vectors without touching the global
// Prometheus registry, so tests can run independently and repeatedly.
func newTestMetrics() *Metrics {
	return &Metrics{
		ReqDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "app",
				Name:      "request_duration",
				Buckets:   prometheus.DefBuckets,
			}, []string{"path", "method", "status_code"}),
		ReqCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "app",
				Name:      "request_total",
			}, []string{"path", "method", "status_code"}),
	}
}
