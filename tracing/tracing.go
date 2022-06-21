package tracing

import (
	"context"
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc"
)

func Initialize(cfg *config.TracingConfig) error {
	log.Info("Initializing tracing")

	ctx := context.Background()
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(4*time.Second),
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock(), grpc.WithTimeout(2*time.Second)))

	if err != nil {
		return fmt.Errorf("failed to create new OTLP trace exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(newResource(cfg.ServiceName, cfg.Version, cfg.Env)),
	)

	otel.SetTracerProvider(tp)
	return nil
}

func newResource(service string, version string, env string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceVersionKey.String(version),
			attribute.String("env", env),
		),
	)
	return r
}
