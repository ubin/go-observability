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

// InitSentry initializes Sentry for error monitoring
func New(cfg config.Config) (*Sentry, error) {
	s := Sentry{
		cfg: cfg,
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.GetCollectorEndpoint(),
		EnableTracing:    cfg.IsInsecure(),
		Environment:      cfg.GetEnvironment(),
		TracesSampleRate: 1.0, // Adjust for performance (percentage of transactions to capture)
		Debug:            cfg.IsDebugMode(),
	})
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
