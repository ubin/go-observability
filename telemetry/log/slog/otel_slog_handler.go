package slog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ubin/go-observability/telemetry/config"
	"github.com/ubin/go-observability/telemetry/provider/sentry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type OtelHandler struct {
	wrapped slog.Handler
	// otelCfg config.Tracing
	otelProvider string
}

func NewOtelHandler(wrapped slog.Handler) *OtelHandler {
	return &OtelHandler{
		wrapped:      wrapped,
		otelProvider: "sentry",
	}
}

func (h *OtelHandler) Handle(ctx context.Context, r slog.Record) error {
	// Extract tracing information from context
	span := trace.SpanFromContext(ctx)
	if span != nil {
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		// Add tracing metadata to the log
		r.AddAttrs(
			slog.String("trace_id", traceID),
			slog.String("span_id", spanID),
		)
	}

	// if h.otelCfg.ExporterType == config.ExporterTypeSentry {
	if h.otelProvider == string(config.ExporterTypeSentry) {
		// Send logs as messages to Sentry
		sentry.CaptureLogMessage(r)
	}

	// Delegate actual logging to the wrapped handler
	return h.wrapped.Handle(ctx, r)
}

func (h *OtelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &OtelHandler{wrapped: h.wrapped.WithAttrs(attrs)}
}

func (h *OtelHandler) WithGroup(name string) slog.Handler {
	return &OtelHandler{wrapped: h.wrapped.WithGroup(name)}
}

func (h *OtelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Delegate to the wrapped handler's Enabled method
	return h.wrapped.Enabled(ctx, level)
}

func AddLogToSpan(ctx context.Context, level slog.Level, msg string, keyvals ...interface{}) {
	span := trace.SpanFromContext(ctx)
	// traceID := span.SpanContext().TraceID().String()
	// spanID := span.SpanContext().SpanID().String()

	if span.IsRecording() {
		attrs := make([]attribute.KeyValue, 0, len(keyvals)/2+2)
		attrs = append(attrs, attribute.String("log.level", level.String()))
		attrs = append(attrs, attribute.String("log.message", msg))
		for i := 0; i < len(keyvals); i += 2 {
			key, ok := keyvals[i].(string)
			if !ok {
				continue
			}
			value := keyvals[i+1]
			// attrs = append(attrs, attribute.Any(key, value))
			switch v := value.(type) {
			case string:
				attrs = append(attrs, attribute.String(key, v))
			case int:
				attrs = append(attrs, attribute.Int(key, v))
			case bool:
				attrs = append(attrs, attribute.Bool(key, v))
			case float64:
				attrs = append(attrs, attribute.Float64(key, v))
			default:
				attrs = append(attrs, attribute.String(key, fmt.Sprintf("%v", v)))
			}
		}
		span.AddEvent("log", trace.WithAttributes(attrs...))
	}
}
