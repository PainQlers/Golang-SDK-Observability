package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	AIRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ai_agent_requests_total",
			Help: "Total number of AI Agent requests.",
		},
	)

	AIRequestsSucceededTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ai_agent_requests_succeeded_total",
			Help: "Total number of succeeded AI requests.",
		},
	)

	AIRequestsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ai_agent_requests_failed_total",
			Help: "Total number of failed AI requests.",
		},
	)

	AIRequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "ai_agent_request_duration_seconds",
			Help:    "AI request duration.",
			Buckets: prometheus.DefBuckets,
		},
	)
)
