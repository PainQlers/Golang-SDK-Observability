package operation

import (
	"context"
	"log/slog"
	"time"

	instrument "github.com/PainQlers/backend/pkg/telemetry/instrument"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Operation struct {
	ctx   context.Context
	span  trace.Span
	start time.Time
	name  string
}

func Start(ctx context.Context, name string) (*Operation, context.Context) {
	ctx, span := instrument.StartSpan(ctx, name)

	instrument.RecordAIRequest()

	slog.InfoContext(ctx,
		"operation started",
		"operation", name,
	)

	return &Operation{
		ctx:   ctx,
		span:  span,
		start: time.Now(),
		name:  name,
	}, ctx
}

func (o *Operation) Success() {
	instrument.RecordAIRequestSucceeded()

	o.span.AddEvent("operation succeeded")
	o.span.SetStatus(codes.Ok, "")

	slog.InfoContext(o.ctx,
		"operation succeeded",
		"operation", o.name,
	)
}

func (o *Operation) Failed(err error) {
	instrument.RecordAIRequestFailed()

	o.span.RecordError(err)
	o.span.SetStatus(codes.Error, err.Error())

	slog.ErrorContext(o.ctx,
		"operation failed",
		slog.Any("error", err),
		"operation", o.name,
	)
}

func (o *Operation) End() {
	instrument.RecordAIRequestDuration(o.start)

	o.span.End()
}

func (o *Operation) AddEvent(name string) {
	o.span.AddEvent(name)
}
