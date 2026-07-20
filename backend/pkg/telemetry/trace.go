package telemetry

import (
	"context"

	otlptracehttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	resource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewTraceProvider creates a new trace provider with OTLP exporter.
func NewTraceProvider(
	ctx context.Context,
	alloyEndpoint string,
	insecure bool,
	res *resource.Resource,
) (*sdktrace.TracerProvider, error) {

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(alloyEndpoint),
		otlptracehttp.WithURLPath("/v1/traces"),
	}

	if insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	traceExporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	return tp, nil
}
