package telemetry

import (
	prometheus "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
)

// NewMetricProvider creates a new meter provider with Prometheus exporter.
func NewMetricProvider(res *sdkresource.Resource) (*sdkmetric.MeterProvider, *prometheus.Exporter, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(res),
	)
	return mp, exporter, nil
}
