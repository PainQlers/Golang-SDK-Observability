package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var registerOnce sync.Once

func Register() {
	registerOnce.Do(func() {
		prometheus.MustRegister(
			AIRequestsTotal,
			AIRequestsSucceededTotal,
			AIRequestsFailedTotal,
			AIRequestDuration,

			HTTPRequestsTotal,
			HTTPRequestDuration,

			DatabaseQueriesTotal,
		)
	})
}
