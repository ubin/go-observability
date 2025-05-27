package main

import (
	"context"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/goliatone/go-errors"
	logconfig "github.com/ubin/go-telemetry/example/sentry/config"
	"github.com/ubin/go-telemetry/logger"
	"github.com/ubin/go-telemetry/logger/loggerfactory"
	"github.com/ubin/go-telemetry/telemetry"
	"github.com/ubin/go-telemetry/telemetry/config"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()

	cfg := config.TracingConfig{
		ServiceName:       "example-service",
		Environment:       "test",
		Enabled:           true,
		ExporterType:      config.ExporterTypeSentry,
		CollectorEndpoint: "",
		Insecure:          true,
		DebugMode:         true,
	}
	tp, err := telemetry.InitTracer(&cfg)
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shut down tracer: %v", err)
		}
	}()

	initlogger()

	tracer := otel.Tracer("example-tracer")
	ctx, span := tracer.Start(ctx, "example-span")
	defer span.End()

	// Simulate work
	time.Sleep(2 * time.Second)
	logger.Log.InfoContext(ctx, "Processing completed")

	errSample := errors.New("simulated error for demonstration purposes", errors.CategoryNotFound)

	// Add context
	enrichedErr := errSample.
		WithMetadata(map[string]any{"ctx_id": 123}).
		WithRequestID("req-456").
		WithStackTrace().
		WithCode(404).
		WithTextCode("RESOURCE_NOT_FOUND")
	logger.Log.ErrorContext(ctx, "Custom Error: %v", enrichedErr)
	//capture the stack trace with sentry
	sentry.CaptureException(enrichedErr)

	logger.Log.InfoContext(ctx, "Trace completed")

}

func initlogger() {

	logCfg := logconfig.LoggerConfig{
		Code:         "slog",
		Level:        "info",
		Formatter:    "json",
		EnableCaller: true,
		FileEnabled:  false,
		Filename:     "app.log",
		MaxSize:      10,
		MaxBackups:   5,
		MaxAge:       30,
		Compress:     true,
		LocalTime:    true,
	}
	if err := loggerfactory.Register(logCfg, "test"); err != nil {
		log.Fatalf("failed to register logger: %v", err)
	}
}
