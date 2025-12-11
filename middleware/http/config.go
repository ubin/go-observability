package http

import (
	"log/slog"

	"go.opentelemetry.io/otel/sdk/trace"
)

// Config holds configuration for HTTP tracing middleware
type Config struct {
	// TracerProvider is the OpenTelemetry tracer provider
	// If nil, tracing will be skipped
	TracerProvider *trace.TracerProvider

	// Logger is used for logging HTTP requests
	// If nil, logging will be skipped
	Logger *slog.Logger

	// ServiceName is the name of the service for tracing (defaults to "http-server")
	ServiceName string

	// SkipPaths are paths to exclude from tracing (e.g., /health, /metrics)
	// Useful for reducing noise from health checks
	SkipPaths []string

	// SkipLogging disables request logging if true
	SkipLogging bool

	// GenerateRequestID enables request ID generation and X-Request-ID header
	GenerateRequestID bool
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		ServiceName:       "http-server",
		SkipPaths:         []string{},
		SkipLogging:       false,
		GenerateRequestID: true,
	}
}

// shouldSkipPath checks if a path should be excluded from tracing
func (c *Config) shouldSkipPath(path string) bool {
	for _, skipPath := range c.SkipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}
