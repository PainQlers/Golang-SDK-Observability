package httpclient

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PainQlers/backend/pkg/telemetry/instrument"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Client wraps an underlying http.Client and adds telemetry for outgoing requests.
type Client struct {
	base *http.Client
}

// NewHTTPClient creates a telemetry-enabled HTTP client.
func NewHTTPClient() *Client {
	return WrapHTTPClient(http.DefaultClient)
}

// WrapHTTPClient wraps an existing http.Client with telemetry instrumentation.
func WrapHTTPClient(base *http.Client) *Client {
	if base == nil {
		base = http.DefaultClient
	}

	return &Client{base: base}
}

// Do executes the request and records telemetry for the call.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c == nil || c.base == nil {
		return nil, fmt.Errorf("http client instrumentation: nil client")
	}
	if req == nil {
		return nil, fmt.Errorf("http client instrumentation: nil request")
	}

	ctx := req.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	method := normalizeMethod(req.Method)
	host, port := requestHostAndPort(req)
	targetService := targetServiceLabel(host)

	// The span name stays generic so multiple HTTP methods can be grouped under one operation family.
	ctx, span := instrument.StartSpan(ctx, "http.client")
	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", method),
		attribute.String("server.address", host),
	)
	if port > 0 {
		span.SetAttributes(attribute.Int("server.port", port))
	}

	// We intentionally avoid recording URL or path data because those can contain identifiers or secrets.
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	otel.GetTextMapPropagator().Inject(ctx, propagationCarrier{header: req.Header})

	start := time.Now()
	logger := slog.Default().With("http.method", method, "server.address", host)
	logger.InfoContext(ctx, "outgoing http request started")

	outReq := req.WithContext(ctx)
	resp, err := c.base.Do(outReq)
	duration := time.Since(start)

	statusCode := 0
	statusClass := "unknown"
	if resp != nil {
		statusCode = resp.StatusCode
		statusClass = classifyStatusCode(statusCode)
		span.SetAttributes(attribute.Int("http.status_code", statusCode))
	}

	// Metrics stay low-cardinality by passing a service label and a status class rather than raw URLs or paths.
	instrument.RecordHTTPClientRequest(method, targetService, statusClass, duration.Seconds())

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.ErrorContext(ctx, "outgoing http request failed", slog.Any("error", err), "http.duration_ms", duration.Milliseconds(), "http.status", statusCode)
		return nil, err
	}

	span.SetStatus(determineSpanStatus(statusCode), "")
	logger.InfoContext(ctx, "outgoing http request completed", "http.duration_ms", duration.Milliseconds(), "http.status", statusCode)
	return resp, nil
}

func normalizeMethod(method string) string {
	if method == "" {
		return "GET"
	}
	return strings.ToUpper(method)
}

func requestHostAndPort(req *http.Request) (string, int) {
	host := req.URL.Hostname()
	if host == "" {
		host = req.Host
	}
	if host == "" {
		return "unknown", 0
	}

	port, err := strconv.Atoi(req.URL.Port())
	if err != nil || port < 1 || port > 65535 {
		return host, 0
	}

	return host, port
}

func determineSpanStatus(statusCode int) codes.Code {
	switch {
	case statusCode == 0:
		return codes.Error

	case statusCode >= 500:
		return codes.Error

	default:
		return codes.Ok
	}
}

func targetServiceLabel(host string) string {
	if host == "" || host == "unknown" {
		return "unknown"
	}
	return host
}

func classifyStatusCode(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "2xx"
	case statusCode >= 300 && statusCode < 400:
		return "3xx"
	case statusCode >= 400 && statusCode < 500:
		return "4xx"
	case statusCode >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}

type propagationCarrier struct {
	header http.Header
}

func (c propagationCarrier) Get(key string) string {
	return c.header.Get(key)
}

func (c propagationCarrier) Set(key, value string) {
	c.header.Set(key, value)
}

func (c propagationCarrier) Keys() []string {
	keys := make([]string, 0, len(c.header))
	for key := range c.header {
		keys = append(keys, key)
	}
	return keys
}
