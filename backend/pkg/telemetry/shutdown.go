package telemetry

import (
	"context"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// ShutdownTelemetry shuts down telemetry providers.
func ShutdownTelemetry(ctx context.Context, tp *sdktrace.TracerProvider, mp *sdkmetric.MeterProvider, lp *sdklog.LoggerProvider) error {
	var err error

	if tp != nil {
		if errShutdown := tp.Shutdown(ctx); errShutdown != nil {
			err = errShutdown
		}
	}
	if mp != nil {
		if errShutdown := mp.Shutdown(ctx); errShutdown != nil {
			if err != nil {
				err = err
			} else {
				err = errShutdown
			}
		}
	}
	if lp != nil {
		if errShutdown := lp.Shutdown(ctx); errShutdown != nil {
			if err != nil {
				err = err
			} else {
				err = errShutdown
			}
		}
	}
	return err
}