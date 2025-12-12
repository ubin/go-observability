package telemetry

import (
	"context"
	"fmt"

	"github.com/ubin/go-observability/telemetry/config"
	"github.com/ubin/go-observability/telemetry/provider/sentry"

	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// InitTracer initializes OpenTelemetry tracing with configurable exporters.
func InitTracer(cfg config.Config) (*trace.TracerProvider, error) {
	ctx := context.Background()

	var exporter trace.SpanExporter
	var err error

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	switch cfg.GetExporterType() {
	case config.ExporterTypeSentry:
		//initialize sentry sdk
		s, err := sentry.New(cfg)
		if err != nil {
			return nil, fmt.Errorf("sentry initialization failed: %w", err)
		}
		tracerProvider := s.TracerProvider()
		otel.SetTracerProvider(tracerProvider)
		return tracerProvider, nil

	case "signoz":
		// endpoint := os.Getenv("OTEL_SIGNOZ_ENDPOINT")
		// exporter, err = otlptrace.New(ctx, otlphttp.NewClient(otlphttp.WithEndpoint(endpoint)))

	case config.ExporterTypeHTTP:
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(cfg.GetCollectorEndpoint()),
		}
		if cfg.IsInsecure() {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		exporter, err = otlptracehttp.New(
			ctx,
			opts...,
		)

	case config.ExporterTypeGRPC:
		secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		if cfg.IsInsecure() {
			secureOption = otlptracegrpc.WithInsecure()
		}
		exporter, err = otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithEndpoint(cfg.GetCollectorEndpoint()),
			secureOption,
		)
	case config.ExporterTypeStdout:
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)

	default:
		return nil, fmt.Errorf("unknown exporter type: %s", cfg.GetExporterType())
	}

	if err != nil {
		return nil, err
	}

	resource, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.GetServiceName()),
			attribute.String("library.language", "go"),
			attribute.String("environment", cfg.GetEnvironment()),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil
}
