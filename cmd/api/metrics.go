package main

import "github.com/prometheus/client_golang/prometheus"

var (
	// traffic
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nailit",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests received",
		},
		[]string{"method", "path", "status"},
	)

	// latency
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

	// error metrics
	httpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "nailit",
			Subsystem: "http",
			Name:      "errors_total",
			Help:      "Total number of HHTP error responses",
		},
		[]string{"method", "path", "status_class"},
	)

	// in-flight requests
	httpRequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "nailit",
			Subsystem: "http",
			Name:      "requets_in_flight",
			Help:      "Current number of in-flight HTTP requests",
		},
		[]string{"path"},
	)
)

func registerMetrics() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpErrorsTotal,
		httpRequestsInFlight,
	)
}
