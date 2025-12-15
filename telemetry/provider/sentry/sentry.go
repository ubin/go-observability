package sentry

import (
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/ubin/go-observability/telemetry/config"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Sentry struct {
	cfg config.Config
}

// New initializes Sentry for error monitoring and distributed tracing
func New(cfg config.Config) (*Sentry, error) {
	s := Sentry{
		cfg: cfg,
	}

	clientOptions := sentry.ClientOptions{
		Dsn:              cfg.GetCollectorEndpoint(),
		Environment:      cfg.GetEnvironment(),
		Release:          cfg.GetRelease(),
		Debug:            cfg.IsDebugMode(),
		AttachStacktrace: true,
		EnableTracing:    true, // Always enable tracing when using Sentry exporter
		TracesSampleRate: cfg.GetTracesSampleRate(),
		EnableLogs:       cfg.IsLogsEnabled(), // Send logs to Sentry if enabled
	}

	// Only set ServerName if service name is provided
	if cfg.GetServiceName() != "" {
		clientOptions.ServerName = cfg.GetServiceName()
	}

	err := sentry.Init(clientOptions)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Sentry) TracerProvider() *trace.TracerProvider {
	tracerProvider := trace.NewTracerProvider(
		trace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)

	return tracerProvider
}
