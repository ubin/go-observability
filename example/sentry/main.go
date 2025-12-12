package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/ubin/go-observability/telemetry"
	"github.com/ubin/go-observability/telemetry/config"
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

	tracer := otel.Tracer("example-tracer")
	ctx, span := tracer.Start(context.Background(), "example-span")
	defer span.End()

	// Simulate work
	time.Sleep(2 * time.Second)
	span.AddEvent("Processing completed")
	sentry.CaptureMessage("Test message")

	// Simulate an error
	err = fmt.Errorf("simulated error for demonstration purposes")
	if err != nil {
		span.RecordError(err)
		sentry.CaptureException(err)
		log.Printf("error occurred: %v", err)
	}

	fmt.Println("Trace completed")

}
