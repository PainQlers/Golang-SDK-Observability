package instrument

import (
	"time"

	appmetrics "github.com/PainQlers/backend/pkg/telemetry/metrics"
)

func RecordAIRequest() {
	appmetrics.AIRequestsTotal.Inc()
}

func RecordAIRequestSucceeded() {
	appmetrics.AIRequestsSucceededTotal.Inc()
}

func RecordAIRequestFailed() {
	appmetrics.AIRequestsFailedTotal.Inc()
}

func RecordAIRequestDuration(start time.Time) {
	appmetrics.AIRequestDuration.Observe(
		time.Since(start).Seconds(),
	)
}
