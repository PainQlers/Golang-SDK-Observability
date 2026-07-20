package instrument

import (
	appmetrics "github.com/PainQlers/backend/pkg/telemetry/metrics"
)

func RecordDatabaseQuery(operation string) {

	appmetrics.DatabaseQueriesTotal.
		WithLabelValues(operation).
		Inc()
}
