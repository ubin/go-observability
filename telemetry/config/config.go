// config.go
package config

// Pioneer shipment statuses
type ExporterType string

const (
	ExporterTypeHTTP   ExporterType = "http"
	ExporterTypeGRPC   ExporterType = "grpc"
	ExporterTypeStdout ExporterType = "stdout"
	ExporterTypeSentry ExporterType = "sentry"
)

// Config is the interface for configuration
type Config interface {
	GetServiceName() string
	GetEnvironment() string
	IsEnabled() bool
	GetExporterType() ExporterType
	GetCollectorEndpoint() string
	IsInsecure() bool
	IsDebugMode() bool
	GetTracesSampleRate() float64
	GetRelease() string
	IsLogsEnabled() bool
}

// TracingConfig implements the Config interface
type TracingConfig struct {
	ServiceName       string       `koanf:"service_name"`
	Environment       string       `koanf:"environment"`
	Enabled           bool         `koanf:"enabled"`
	ExporterType      ExporterType `koanf:"exporter_type"`
	CollectorEndpoint string       `koanf:"collector_endpoint"`
	Insecure          bool         `koanf:"insecure"`
	DebugMode         bool         `koanf:"debug_mode"`
	TracesSampleRate  float64      `koanf:"traces_sample_rate"` // 0.0 to 1.0 (0.1 = 10%, 1.0 = 100%)
	Release           string       `koanf:"release"`            // Release version (e.g., "v1.0.0", git commit hash)
	EnableLogs        bool         `koanf:"enable_logs"`        // Send logs to Sentry
}

// GetServiceName returns the service name
func (c *TracingConfig) GetServiceName() string {
	return c.ServiceName
}

// GetEnvironment returns the environment
func (c *TracingConfig) GetEnvironment() string {
	return c.Environment
}

// IsEnabled returns if the service is enabled
func (c *TracingConfig) IsEnabled() bool {
	return c.Enabled
}

// GetExporterType returns the exporter type
func (c *TracingConfig) GetExporterType() ExporterType {
	return c.ExporterType
}

// GetExporterEndpoint returns the exporter endpoint
func (c *TracingConfig) GetCollectorEndpoint() string {
	return c.CollectorEndpoint
}

// IsInsecure returns if the connection is insecure
func (c *TracingConfig) IsInsecure() bool {
	return c.Insecure
}

// IsDebugMode returns if the debug mode is enabled
func (c *TracingConfig) IsDebugMode() bool {
	return c.DebugMode
}

// GetTracesSampleRate returns the traces sample rate (0.0 to 1.0)
func (c *TracingConfig) GetTracesSampleRate() float64 {
	if c.TracesSampleRate <= 0 {
		return 1.0 // Default to 100% if not set
	}
	return c.TracesSampleRate
}

// GetRelease returns the release version
func (c *TracingConfig) GetRelease() string {
	return c.Release
}

// IsLogsEnabled returns if logs should be sent to Sentry
func (c *TracingConfig) IsLogsEnabled() bool {
	return c.EnableLogs
}
