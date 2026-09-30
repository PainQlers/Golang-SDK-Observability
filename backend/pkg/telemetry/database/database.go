package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/PainQlers/backend/pkg/telemetry/instrument"
	appmetrics "github.com/PainQlers/backend/pkg/telemetry/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Operation describes the kind of database operation being instrumented.
type Operation string

const (
	// Select represents a SELECT query.
	Select Operation = "SELECT"
	// Insert represents an INSERT query.
	Insert Operation = "INSERT"
	// Update represents an UPDATE query.
	Update Operation = "UPDATE"
	// Delete represents a DELETE query.
	Delete Operation = "DELETE"
	// Begin represents a transaction BEGIN statement.
	Begin Operation = "BEGIN"
	// Commit represents a transaction COMMIT statement.
	Commit Operation = "COMMIT"
	// Rollback represents a transaction ROLLBACK statement.
	Rollback Operation = "ROLLBACK"
)

// Observe executes a database operation while recording metrics, tracing, and logs.
// The context should usually come from operation.Start so spans remain linked to the business operation.
func Observe(
	ctx context.Context,
	operation Operation,
	statement string,
	fn func(context.Context) error,
) error {
	if fn == nil {
		return fmt.Errorf("database instrumentation: nil function")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	dbCtx, span := instrument.StartSpan(ctx, "database."+string(operation))
	defer span.End()

	start := time.Now()
	metricsLabel := string(operation)

	appmetrics.DatabaseQueriesTotal.WithLabelValues(metricsLabel).Inc()

	span.SetAttributes(
		attribute.String("db.operation", metricsLabel),
		attribute.String("db.query", statement),
	)

	logger := slog.With(
		"db.operation", metricsLabel,
		"db.statement", statement,
	)
	logger.InfoContext(dbCtx, "database operation started")

	err := fn(dbCtx)
	duration := time.Since(start)

	appmetrics.DatabaseQueryDuration.WithLabelValues(metricsLabel).Observe(duration.Seconds())

	if err != nil {
		appmetrics.DatabaseQueryFailuresTotal.WithLabelValues(metricsLabel).Inc()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.ErrorContext(dbCtx, "database operation failed", slog.Any("error", err), "db.duration_ms", duration.Milliseconds())
		return fmt.Errorf("database %s failed: %w", metricsLabel, err)
	}

	span.SetStatus(codes.Ok, "")
	logger.InfoContext(dbCtx, "database operation completed", "db.duration_ms", duration.Milliseconds())
	return nil
}

// ObserveResult executes a database operation and returns its result while recording metrics, tracing, and logs.
func ObserveResult[T any](ctx context.Context, operation Operation, query string, fn func(context.Context) (T, error)) (T, error) {
	var zero T
	if fn == nil {
		return zero, fmt.Errorf("database instrumentation: nil function")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	dbCtx, span := instrument.StartSpan(ctx, "database."+string(operation))
	defer span.End()

	start := time.Now()
	metricsLabel := string(operation)

	appmetrics.DatabaseQueriesTotal.WithLabelValues(metricsLabel).Inc()

	span.SetAttributes(
		attribute.String("db.operation", metricsLabel),
		attribute.String("db.query", query),
	)

	logger := slog.With("db.operation", metricsLabel, "db.query", query)
	logger.InfoContext(dbCtx, "database operation started")

	result, err := fn(dbCtx)
	duration := time.Since(start)

	appmetrics.DatabaseQueryDuration.WithLabelValues(metricsLabel).Observe(duration.Seconds())

	if err != nil {
		appmetrics.DatabaseQueryFailuresTotal.WithLabelValues(metricsLabel).Inc()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.ErrorContext(dbCtx, "database operation failed", slog.Any("error", err), "db.duration_ms", duration.Milliseconds())
		return zero, fmt.Errorf("database %s failed: %w", metricsLabel, err)
	}

	span.SetStatus(codes.Ok, "")
	logger.InfoContext(dbCtx, "database operation completed", "db.duration_ms", duration.Milliseconds())
	return result, nil
}
