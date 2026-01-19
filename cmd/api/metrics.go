package main

import "github.com/prometheus/client_golang/prometheus"

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nailit",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests received",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "nailit",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request latency",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func registerMetrics() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
	)
}
