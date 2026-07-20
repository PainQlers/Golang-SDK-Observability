package telemetry

import (
	"context"
	"net/http"

	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartMetricsServer starts an HTTP server for Prometheus metrics on port 2112.
func StartMetricsServer(ctx context.Context) {
	go func() {
		mux := http.NewServeMux()

		mux.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(":2112", nil); err != nil {
			return
		}

	}()
}
