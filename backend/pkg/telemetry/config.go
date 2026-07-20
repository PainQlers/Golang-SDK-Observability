package telemetry

// Config holds the configuration for the telemetry SDK.
type Config struct {
	// AlloyEndpoint is the endpoint of the Alloy collector (for OTLP).
	AlloyEndpoint string

	// ServiceName is the name of the service for telemetry.
	ServiceName string

	MetricsPort string

	Insecure bool

	EnableTrace bool

	EnableMetrics bool

	EnableLogs bool
}
