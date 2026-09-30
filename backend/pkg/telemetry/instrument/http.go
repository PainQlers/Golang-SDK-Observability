package instrument

import (
	appmetrics "github.com/PainQlers/backend/pkg/telemetry/metrics"
)

// RecordHTTPRequest records server-side HTTP metrics using the existing middleware labels.
func RecordHTTPRequest(method, route, status string, duration float64) {
	appmetrics.HTTPRequestsTotal.
		WithLabelValues(
			method,
			route,
			status,
		).
		Inc()

	appmetrics.HTTPRequestDuration.
		WithLabelValues(
			method,
			route,
		).
		Observe(duration)
}

// RecordHTTPClientRequest records client-side HTTP metrics using low-cardinality labels.
// The caller should provide a stable target service label and a status class instead of raw URLs or paths.
func RecordHTTPClientRequest(method, targetService, statusClass string, duration float64) {
	appmetrics.HTTPClientRequestsTotal.
		WithLabelValues(
			method,
			targetService,
			statusClass,
		).
		Inc()

	appmetrics.HTTPClientRequestDuration.
		WithLabelValues(
			method,
			targetService,
		).
		Observe(duration)
}
