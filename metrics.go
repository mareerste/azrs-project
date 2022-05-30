package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Init counter value at
	currentCount = 0

	// Prometheus metric being exposed/available
	total_hits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "total_hits",
			Help: "Total number of http hits.",
		},
	)

	create_config_hits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "create_config_hits",
			Help: "Total number of create config hits.",
		},
	)

	// Add all metrics that will be resisted
	metricsList = []prometheus.Collector{
		total_hits,
		create_config_hits,
	}

	// Prometheus Registry to register metrics.
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	// Register metrics that will be exposed.
	prometheusRegistry.MustRegister(metricsList...)
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func count(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		total_hits.Inc()
		create_config_hits.Inc()
		f(w, r) // original function call
	}
}
