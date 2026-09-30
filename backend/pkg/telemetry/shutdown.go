package telemetry

import (
	"context"
	"errors"

	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// ShutdownTelemetry shuts down telemetry providers.
func ShutdownTelemetry(
	ctx context.Context,
	tp *sdktrace.TracerProvider,
	mp *sdkmetric.MeterProvider,
	lp *sdklog.LoggerProvider,
) error {
	var err error

	if tp != nil {
		if shutdownErr := tp.Shutdown(ctx); shutdownErr != nil {
			err = errors.Join(err, shutdownErr)
		}
	}

	if mp != nil {
		if shutdownErr := mp.Shutdown(ctx); shutdownErr != nil {
			err = errors.Join(err, shutdownErr)
		}
	}

	if lp != nil {
		if shutdownErr := lp.Shutdown(ctx); shutdownErr != nil {
			err = errors.Join(err, shutdownErr)
		}
	}

	return err
}
