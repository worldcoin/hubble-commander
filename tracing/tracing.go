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

func Initialize(cfg *config.TracingConfig) (func(), error) {
	log.Infof("Initializing tracing with endpoint %s", cfg.Endpoint)

	ctx := context.Background()
	//nolint
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(3*time.Second),
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock(), grpc.WithTimeout(3*time.Second)))

	if err != nil {
		return nil, fmt.Errorf("error creating trace exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(newResource(cfg.ServiceName, cfg.Version, cfg.Env)),
	)

	otel.SetTracerProvider(tp)

	shutdownFunc := func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Errorf("error shutting down tracer: %s", err)
		}
	}

	return shutdownFunc, nil
}

func newResource(service, version, env string) *resource.Resource {
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
