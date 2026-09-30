package telemetry

import (
	"context"
	"log/slog"

	otelslog "go.opentelemetry.io/contrib/bridges/otelslog"
	otel "go.opentelemetry.io/otel"
	propagation "go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	metric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitSharedTelemetry initializes all telemetry components (traces, metrics, logs)
// and returns a shutdown function.
func InitSharedTelemetry(ctx context.Context, cfg Config, serviceName string) (func() error, error) {
	// 1. Create resource
	res, err := NewResource(ctx, cfg.ServiceName)
	if err != nil {
		return nil, err
	}

	var (
		tp            *sdktrace.TracerProvider
		mp            *metric.MeterProvider  // ใช้ type ที่ NewMetricProvider คืนมา
		lp            *sdklog.LoggerProvider // ใช้ type ที่ NewLoggerProvider คืนมา
		cancelMetrics func()
	)

	// 2. Set up trace provider
	if cfg.EnableTrace {
		tp, err = NewTraceProvider(
			ctx,
			cfg.AlloyEndpoint,
			cfg.Insecure,
			res,
		)
		if err != nil {
			return nil, err
		}

		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
			),
		)
	}

	// 3. Set up metric provider
	if cfg.EnableMetrics {
		mp, _, err = NewMetricProvider(res)
		if err != nil {
			if tp != nil {
				_ = tp.Shutdown(ctx)
			}
			return nil, err
		}

		otel.SetMeterProvider(mp)
		metricsCtx, metricsCancel := context.WithCancel(context.Background())
		cancelMetrics = metricsCancel
		StartMetricsServer(metricsCtx, cfg.MetricsPort)
	}

	// 4. Set up logger provider
	if cfg.EnableLogs {
		lp, err = NewLoggerProvider(
			ctx,
			cfg.AlloyEndpoint,
			cfg.Insecure,
			res,
		)
		if err != nil {
			// Clean up trace and meter providers on error
			if tp != nil {
				_ = tp.Shutdown(ctx)
			}

			if mp != nil {
				_ = mp.Shutdown(ctx)
			}
			return nil, err
		}

		// Create slog logger bridged with OpenTelemetry
		logger := otelslog.NewLogger(
			cfg.ServiceName,
			otelslog.WithLoggerProvider(lp),
		)
		slog.SetDefault(logger)
	}

	// Return shutdown function
	return func() error {
		if tp != nil {
			_ = tp.Shutdown(ctx)
		}

		if mp != nil {
			if cancelMetrics != nil {
				cancelMetrics()
			}
			_ = mp.Shutdown(ctx)
		}

		if lp != nil {
			_ = lp.Shutdown(ctx)
		}

		return nil
	}, nil
}
