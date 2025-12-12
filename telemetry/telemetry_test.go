package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ubin/go-observability/telemetry/config"
)

// MockConfig is a mock implementation of the Config interface for testing
type MockConfig struct {
	ServiceName       string
	Environment       string
	Enabled           bool
	ExporterType      config.ExporterType
	CollectorEndpoint string
	Insecure          bool
	DebugMode         bool
}

func (c *MockConfig) GetServiceName() string               { return c.ServiceName }
func (c *MockConfig) GetEnvironment() string               { return c.Environment }
func (c *MockConfig) IsEnabled() bool                      { return c.Enabled }
func (c *MockConfig) GetExporterType() config.ExporterType { return c.ExporterType }
func (c *MockConfig) GetCollectorEndpoint() string         { return c.CollectorEndpoint }
func (c *MockConfig) IsInsecure() bool                     { return c.Insecure }
func (c *MockConfig) IsDebugMode() bool                    { return c.DebugMode }

func TestInitTracer_Sentry(t *testing.T) {
	cfg := &MockConfig{
		ServiceName:       "test-service",
		Environment:       "test",
		Enabled:           true,
		ExporterType:      config.ExporterTypeSentry,
		CollectorEndpoint: "https://mock-user@test.ingest.us.sentry.io/sentry-id",
		Insecure:          true,
		DebugMode:         true,
	}

	tp, err := InitTracer(cfg)
	assert.NoError(t, err, "InitTracer should not return an error for Sentry")
	assert.NotNil(t, tp, "TracerProvider should be initialized for Sentry")
}

func TestInitTracer_HTTP(t *testing.T) {
	cfg := &MockConfig{
		ServiceName:       "test-service",
		Environment:       "test",
		Enabled:           true,
		ExporterType:      config.ExporterTypeHTTP,
		CollectorEndpoint: "localhost:4318",
		Insecure:          true,
		DebugMode:         false,
	}

	tp, err := InitTracer(cfg)
	assert.NoError(t, err, "InitTracer should not return an error for HTTP")
	assert.NotNil(t, tp, "TracerProvider should be initialized for HTTP")
}

func TestInitTracer_GRPC(t *testing.T) {
	cfg := &MockConfig{
		ServiceName:       "test-service",
		Environment:       "test",
		Enabled:           true,
		ExporterType:      config.ExporterTypeGRPC,
		CollectorEndpoint: "localhost:4317",
		Insecure:          true,
		DebugMode:         false,
	}

	tp, err := InitTracer(cfg)
	assert.NoError(t, err, "InitTracer should not return an error for GRPC")
	assert.NotNil(t, tp, "TracerProvider should be initialized for GRPC")
}

func TestInitTracer_Stdout(t *testing.T) {
	cfg := &MockConfig{
		ServiceName:       "test-service",
		Environment:       "test",
		Enabled:           true,
		ExporterType:      config.ExporterTypeStdout,
		CollectorEndpoint: "",
		Insecure:          false,
		DebugMode:         false,
	}

	tp, err := InitTracer(cfg)
	assert.NoError(t, err, "InitTracer should not return an error for Stdout")
	assert.NotNil(t, tp, "TracerProvider should be initialized for Stdout")
}

func TestInitTracer_UnknownExporter(t *testing.T) {
	cfg := &MockConfig{
		ServiceName:       "test-service",
		Environment:       "test",
		Enabled:           true,
		ExporterType:      "unknown",
		CollectorEndpoint: "",
		Insecure:          false,
		DebugMode:         false,
	}

	tp, err := InitTracer(cfg)
	assert.Error(t, err, "InitTracer should return an error for unknown exporter")
	assert.Nil(t, tp, "TracerProvider should not be initialized for unknown exporter")
}
