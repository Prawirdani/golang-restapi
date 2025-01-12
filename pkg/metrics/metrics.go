package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Port        int
	Info        *prometheus.GaugeVec
	ReqDuration *prometheus.HistogramVec
	ReqCounter  *prometheus.CounterVec
}

func Init() *Metrics {
	m := &Metrics{
		Info: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "app",
				Name:      "info",
				Help:      "Application Information",
			}, []string{"version", "environment"}),
		ReqDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "app",
				Name:      "request_duration",
				Help:      "Request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			}, []string{"path", "method", "status_code"}),
		ReqCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "app",
				Name:      "request_total",
				Help:      "Total number of requests",
			}, []string{"path", "method", "status_code"}),
	}
	prometheus.MustRegister(m.ReqDuration, m.Info, m.ReqCounter)
	return m
}

func (m *Metrics) RunServer(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

func (m *Metrics) SetAppInfo(version, env string) {
	m.Info.WithLabelValues(version, env).Set(1)
}
