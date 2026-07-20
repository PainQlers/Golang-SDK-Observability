package instrument

import (
	appmetrics "github.com/PainQlers/backend/pkg/telemetry/metrics"
)

func RecordHTTPRequest(
	method string,
	route string,
	status string,
	duration float64,
) {

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
