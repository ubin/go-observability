package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// OtelHandler wraps a slog.Handler and adds OpenTelemetry tracing information
type OtelHandler struct {
	wrapped slog.Handler
}

// NewOtelHandler creates a new OtelHandler that wraps the provided handler
func NewOtelHandler(wrapped slog.Handler) *OtelHandler {
	return &OtelHandler{
		wrapped: wrapped,
	}
}

// Handle implements slog.Handler and adds trace_id and span_id to logs
func (h *OtelHandler) Handle(ctx context.Context, r slog.Record) error {
	// Extract tracing information from context
	span := trace.SpanFromContext(ctx)
	if span != nil && span.SpanContext().IsValid() {
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		// Add tracing metadata to the log
		r.AddAttrs(
			slog.String("trace_id", traceID),
			slog.String("span_id", spanID),
		)
	}

	// Delegate actual logging to the wrapped handler
	return h.wrapped.Handle(ctx, r)
}

// WithAttrs implements slog.Handler
func (h *OtelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &OtelHandler{wrapped: h.wrapped.WithAttrs(attrs)}
}

// WithGroup implements slog.Handler
func (h *OtelHandler) WithGroup(name string) slog.Handler {
	return &OtelHandler{wrapped: h.wrapped.WithGroup(name)}
}

// Enabled implements slog.Handler
func (h *OtelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.wrapped.Enabled(ctx, level)
}

// AddLogToSpan adds log information as events to the current span
func AddLogToSpan(ctx context.Context, level slog.Level, msg string, keyvals ...interface{}) {
	span := trace.SpanFromContext(ctx)

	if span == nil || !span.IsRecording() {
		return
	}

	attrs := make([]attribute.KeyValue, 0, len(keyvals)/2+2)
	attrs = append(attrs, attribute.String("log.level", level.String()))
	attrs = append(attrs, attribute.String("log.message", msg))

	var capturedError error

	// Parse keyvals pairs
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 >= len(keyvals) {
			break
		}

		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		value := keyvals[i+1]

		switch v := value.(type) {
		case string:
			attrs = append(attrs, attribute.String(key, v))
		case int:
			attrs = append(attrs, attribute.Int(key, v))
		case int64:
			attrs = append(attrs, attribute.Int64(key, v))
		case bool:
			attrs = append(attrs, attribute.Bool(key, v))
		case float64:
			attrs = append(attrs, attribute.Float64(key, v))
		case error:
			// Capture the error if present
			capturedError = v
			attrs = append(attrs, attribute.String(key, v.Error()))
		default:
			attrs = append(attrs, attribute.String(key, fmt.Sprintf("%v", v)))
		}
	}

	// Record error or add event to span
	if level == slog.LevelError {
		if capturedError != nil {
			span.RecordError(capturedError, trace.WithAttributes(attrs...))
		} else {
			span.RecordError(errors.New(msg), trace.WithAttributes(attrs...))
		}
	} else {
		span.AddEvent("log", trace.WithAttributes(attrs...))
	}
}
