package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Prometheus metric being exposed/available
	total_hits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "total_hits",
			Help: "Total number of all http hits.",
		},
	)

	get_all_config_hits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_all_config_hits",
			Help: "Total number of get all config hits.",
		},
	)

	create_config_hits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "create_config_hits",
			Help: "Total number of create config hits.",
		},
	)

	get_config_hits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_config_hits",
			Help: "Total number of get config hits.",
		},
	)

	delete_config_hits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "delete_config_hits",
			Help: "Total number of delete config hits.",
		},
	)

	create_new_version = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "create_new_version",
			Help: "Total number of create new version hits.",
		},
	)

	add_config_to_group = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "add_config_to_group",
			Help: "Total number of add config to group hits.",
		},
	)

	get_by_labels = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_by_labels",
			Help: "Total number of get config with labels hits.",
		},
	)

	// Add all metrics that will be resisted
	metricsList = []prometheus.Collector{
		total_hits,
		get_all_config_hits,
		create_config_hits,
		get_config_hits,
		delete_config_hits,
		create_new_version,
		add_config_to_group,
		get_by_labels,
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

func countCreateConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		total_hits.Inc()
		create_config_hits.Inc()
		f(w, r) // original function call
	}
}

func countGetAll(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		total_hits.Inc()
		get_all_config_hits.Inc()
		f(w, r)
	}
}

func countGet(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		total_hits.Inc()
		get_config_hits.Inc()
		f(w, r)
	}
}

func countDelete(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		total_hits.Inc()
		delete_config_hits.Inc()
		f(w, r)
	}
}

func countCreateNewVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		total_hits.Inc()
		create_new_version.Inc()
		f(w, r)
	}
}

func countAddConfigToGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		total_hits.Inc()
		add_config_to_group.Inc()
		f(w, r)
	}
}

func countSearchByLabels(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		total_hits.Inc()
		get_by_labels.Inc()
		f(w, r)
	}
}
