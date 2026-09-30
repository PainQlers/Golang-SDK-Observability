package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	DatabaseQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "Database queries.",
		},
		[]string{
			"operation",
		},
	)

	DatabaseQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{
			"operation",
		},
	)

	DatabaseQueryFailuresTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_query_failures_total",
			Help: "Failed database queries.",
		},
		[]string{
			"operation",
		},
	)
)
