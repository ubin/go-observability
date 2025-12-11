package http

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// FiberMiddleware returns a Fiber middleware that adds OpenTelemetry tracing
func FiberMiddleware(config *Config) fiber.Handler {
	if config == nil {
		config = DefaultConfig()
	}

	return func(c *fiber.Ctx) error {
		// Skip if path is in skip list
		if config.shouldSkipPath(c.Path()) {
			return c.Next()
		}

		startTime := time.Now()

		// Generate request ID if enabled
		var requestID string
		if config.GenerateRequestID {
			requestID = c.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}
			c.Set(RequestIDHeader, requestID)
			c.Locals("request_id", requestID)
		}

		// Skip tracing if no tracer provider
		if config.TracerProvider == nil {
			// Still log the request if logger is configured
			if config.Logger != nil && !config.SkipLogging {
				defer func() {
					config.Logger.Info("HTTP request completed",
						"method", c.Method(),
						"path", c.Path(),
						"status", c.Response().StatusCode(),
						"duration_ms", time.Since(startTime).Milliseconds(),
						"request_id", requestID)
				}()
			}
			return c.Next()
		}

		// Extract trace context from incoming headers (for distributed tracing)
		ctx := c.Context()
		carrier := make(propagation.MapCarrier)
		c.Request().Header.VisitAll(func(key, value []byte) {
			carrier[string(key)] = string(value)
		})
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		// Create a span for this HTTP request
		tracer := otel.Tracer(config.ServiceName)
		ctx, span := tracer.Start(ctx, fmt.Sprintf("%s %s", c.Method(), c.Route().Path),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.path", c.Path()),
				attribute.String("http.route", c.Route().Path),
				attribute.String("http.scheme", c.Protocol()),
				attribute.String("http.target", string(c.Request().RequestURI())),
				attribute.String("http.host", c.Hostname()),
				attribute.String("http.user_agent", c.Get("User-Agent")),
				attribute.String("http.remote_addr", c.IP()),
			),
		)
		defer span.End()

		// Add request ID to span
		if requestID != "" {
			span.SetAttributes(attribute.String("http.request_id", requestID))
		}

		// Extract trace and span IDs
		spanContext := span.SpanContext()
		traceID := spanContext.TraceID().String()
		spanID := spanContext.SpanID().String()

		// Add trace context to response headers
		c.Set(TraceIDHeader, traceID)
		c.Set(SpanIDHeader, spanID)

		// Store trace info in Fiber locals
		c.Locals("trace_id", traceID)
		c.Locals("span_id", spanID)

		// Store context in Fiber context
		c.SetUserContext(ctx)

		// Log the incoming request
		if config.Logger != nil && !config.SkipLogging {
			config.Logger.InfoContext(ctx, "HTTP request received",
				"method", c.Method(),
				"path", c.Path(),
				"remote_addr", c.IP(),
				"request_id", requestID)
		}

		// Handle the request
		err := c.Next()

		// Record response details in span
		statusCode := c.Response().StatusCode()
		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.Int("http.response_size", len(c.Response().Body())),
		)

		// Set span status based on HTTP status code
		if statusCode >= 500 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
			span.SetAttributes(attribute.Bool("error", true))

			// Record error if returned
			if err != nil {
				span.RecordError(err)
			}
		} else if statusCode >= 400 {
			span.SetAttributes(attribute.Bool("error", false))
		}

		// Log the completed request
		if config.Logger != nil && !config.SkipLogging {
			duration := time.Since(startTime)
			config.Logger.InfoContext(ctx, "HTTP request completed",
				"method", c.Method(),
				"path", c.Path(),
				"status", statusCode,
				"duration_ms", duration.Milliseconds(),
				"bytes", len(c.Response().Body()),
				"request_id", requestID)
		}

		return err
	}
}
