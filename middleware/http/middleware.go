package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	RequestIDHeader = "X-Request-ID"
	TraceIDHeader   = "X-Trace-ID"
	SpanIDHeader    = "X-Span-ID"
)

// Middleware returns a standard net/http middleware that adds OpenTelemetry tracing
func Middleware(config *Config) func(http.Handler) http.Handler {
	if config == nil {
		config = DefaultConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if path is in skip list
			if config.shouldSkipPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			startTime := time.Now()
			ctx := r.Context()

			// Generate request ID if enabled
			var requestID string
			if config.GenerateRequestID {
				requestID = r.Header.Get(RequestIDHeader)
				if requestID == "" {
					requestID = uuid.New().String()
				}
			}

			// Wrap response writer to capture status code
			rw := newResponseWriter(w)

			// Set request ID header on wrapped writer
			if requestID != "" {
				rw.Header().Set(RequestIDHeader, requestID)
			}

			// Skip tracing if no tracer provider
			if config.TracerProvider == nil {
				// Still log the request if logger is configured
				if config.Logger != nil && !config.SkipLogging {
					defer func() {
						logRequest(config, ctx, r, rw.Status(), time.Since(startTime), requestID, "", "")
					}()
				}
				next.ServeHTTP(rw, r)
				return
			}

			// Extract trace context from incoming headers (for distributed tracing)
			ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

			// Create a span for this HTTP request
			tracer := otel.Tracer(config.ServiceName)
			ctx, span := tracer.Start(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path),
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.path", r.URL.Path),
					attribute.String("http.route", r.URL.Path), // Can be improved with route patterns
					attribute.String("http.scheme", r.URL.Scheme),
					attribute.String("http.target", r.URL.RequestURI()),
					attribute.String("http.host", r.Host),
					attribute.String("http.user_agent", r.UserAgent()),
					attribute.String("http.remote_addr", r.RemoteAddr),
				),
			)
			defer span.End()

			// Add request ID to span if generated
			if requestID != "" {
				span.SetAttributes(attribute.String("http.request_id", requestID))
			}

			// Extract trace and span IDs for headers and logging
			spanContext := span.SpanContext()
			traceID := spanContext.TraceID().String()
			spanID := spanContext.SpanID().String()

			// Add trace context to response headers
			rw.Header().Set(TraceIDHeader, traceID)
			rw.Header().Set(SpanIDHeader, spanID)

			// Replace request context with traced context
			r = r.WithContext(ctx)

			// Log the incoming request
			if config.Logger != nil && !config.SkipLogging {
				logRequest(config, ctx, r, 0, 0, requestID, traceID, spanID)
			}

			// Handle panics
			defer func() {
				if err := recover(); err != nil {
					// Record panic in span
					span.RecordError(fmt.Errorf("panic: %v", err))
					span.SetStatus(codes.Error, "panic recovered")
					span.SetAttributes(attribute.Bool("error", true))

					// Log the panic
					if config.Logger != nil {
						config.Logger.ErrorContext(ctx, "HTTP handler panic",
							"error", err,
							"method", r.Method,
							"path", r.URL.Path)
					}

					// Re-panic to let the server handle it
					panic(err)
				}
			}()

			// Call the next handler
			next.ServeHTTP(rw, r)

			// Record response details in span
			statusCode := rw.Status()
			span.SetAttributes(
				attribute.Int("http.status_code", statusCode),
				attribute.Int("http.response_size", rw.BytesWritten()),
			)

			// Set span status based on HTTP status code
			if statusCode >= 500 {
				span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
				span.SetAttributes(attribute.Bool("error", true))

			}

			// Log the completed request
			if config.Logger != nil && !config.SkipLogging {
				duration := time.Since(startTime)
				config.Logger.InfoContext(ctx, "HTTP request completed",
					"method", r.Method,
					"path", r.URL.Path,
					"status", statusCode,
					"duration_ms", duration.Milliseconds(),
					"bytes", rw.BytesWritten(),
					"request_id", requestID)
			}
		})
	}
}

// logRequest logs the incoming HTTP request
func logRequest(config *Config, ctx context.Context, r *http.Request, status int, duration time.Duration, requestID, traceID, spanID string) {
	if config.Logger == nil {
		return
	}

	attrs := []any{
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	}

	if requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}

	if traceID != "" {
		attrs = append(attrs, "trace_id", traceID)
	}

	if spanID != "" {
		attrs = append(attrs, "span_id", spanID)
	}

	if status > 0 {
		attrs = append(attrs, "status", status)
	}

	if duration > 0 {
		attrs = append(attrs, "duration_ms", duration.Milliseconds())
	}

	config.Logger.InfoContext(ctx, "HTTP request received", attrs...)
}
