package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ubin/go-observability/telemetry"
	"github.com/ubin/go-observability/telemetry/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

// WithOpenTracing will initialize the order status sync cron job
func WithOpenTracing(cfg config.Config) (*trace.TracerProvider, error) {
	if !cfg.IsEnabled() {
		// logger.Log.Info("tracing is disabled")
		return nil, nil
	}
	// logger.Log.Info("initializing otel exporter")
	return telemetry.InitTracer(cfg)
}

func main() {
	ctx := context.Background()

	cfg := config.TracingConfig{
		ServiceName:       "example-service",
		Environment:       "",
		Enabled:           true,
		ExporterType:      config.ExporterTypeStdout,
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
	time.Sleep(2 * time.Second) // Simulating work
	span.AddEvent("Processing completed")
	span.End()

	fmt.Println("Trace completed")

}
