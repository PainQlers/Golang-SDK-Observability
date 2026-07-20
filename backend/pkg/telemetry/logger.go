package telemetry

import (
	"context"

	otlploghttp "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
)

// NewLoggerProvider creates a new logger provider with OTLP exporter.
func NewLoggerProvider(
	ctx context.Context,
	alloyEndpoint string,
	insecure bool,
	res *sdkresource.Resource,
) (*sdklog.LoggerProvider, error) {

	opts := []otlploghttp.Option{
		otlploghttp.WithEndpoint(alloyEndpoint),
		otlploghttp.WithURLPath("/v1/logs"),
	}

	if insecure {
		opts = append(opts, otlploghttp.WithInsecure())
	}

	logExporter, err := otlploghttp.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)

	return lp, nil
}
