package telemetry

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartMetricsServer starts an HTTP server for Prometheus metrics.
func StartMetricsServer(ctx context.Context, port string) {
	if port == "" {
		port = "2112"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			slog.Warn("metrics server shutdown failed", "error", err)
		}
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Warn("metrics server stopped", "error", err)
		}
	}()
}
