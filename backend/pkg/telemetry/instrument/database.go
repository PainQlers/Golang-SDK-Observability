package instrument

import (
	"time"

	appmetrics "github.com/PainQlers/backend/pkg/telemetry/metrics"
)

func RecordDatabaseQuery(operation string) {

	appmetrics.DatabaseQueriesTotal.
		WithLabelValues(operation).
		Inc()
}

func RecordDatabaseFailure(operation, query string) {
	appmetrics.DatabaseQueryFailuresTotal.
		WithLabelValues(operation).
		Inc()
}

func RecordDatabaseDuration(operation, query string, start time.Time) {
	appmetrics.DatabaseQueryDuration.
		WithLabelValues(operation).
		Observe(time.Since(start).Seconds())
}
