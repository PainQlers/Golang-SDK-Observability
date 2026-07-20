package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	DatabaseQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "Database queries.",
		},
		[]string{"operation"},
	)
)
