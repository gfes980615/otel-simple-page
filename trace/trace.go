package trace

import (
	"context"
	"fmt"
	"log"
	"os"
	"otel-demo/config"
	"time"

	"github.com/go-logr/stdr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	_ "google.golang.org/grpc/encoding/gzip"
)

func Setup() []trace.TracerProviderOption {
	l := log.New(os.Stdout, "", 0)
	var tracerProviders []trace.TracerProviderOption
	batchOpts := batchOptions()
	if endpoint := config.AppConfig.OtelGrpcEndpoint; endpoint != "" {
		otlpGrpcExp, err := newOTLPGrpcExporter(context.Background(), endpoint)
		if err != nil {
			l.Fatal(err)
			return nil
		}
		tracerProviders = append(tracerProviders, trace.WithBatcher(otlpGrpcExp, batchOpts...))
	}

	if endpoint := config.AppConfig.OtelHttpEndpoint; endpoint != "" {
		otlpHttpExp, err := newOTLPHttpExporter(context.Background(), endpoint)
		if err != nil {
			l.Fatal(err)
			return nil
		}
		tracerProviders = append(tracerProviders, trace.WithBatcher(otlpHttpExp, batchOpts...))
	}

	if endpoint := config.AppConfig.JaegerEndpoint; endpoint != "" {
		jaegerExp, err := newJaegerExporter(endpoint)
		if err != nil {
			l.Fatal(err)
			return nil
		}
		tracerProviders = append(tracerProviders, trace.WithBatcher(jaegerExp, batchOpts...))
	}

	if service := config.AppConfig.Service; service != "" {
		tracerProviders = append(tracerProviders, trace.WithResource(newResource(service)))
	} else {
		l.Fatal("service parameter cannot be empty in app.yaml file")
		return nil
	}
	stdr.SetVerbosity(5)

	tracerProviders = append(tracerProviders, trace.WithSampler(trace.AlwaysSample()))
	return tracerProviders
}

func batchOptions() []trace.BatchSpanProcessorOption {
	opts := []trace.BatchSpanProcessorOption{}

	opts = append(opts, trace.WithMaxQueueSize(config.AppConfig.MaxQueueSize))
	opts = append(opts, trace.WithMaxExportBatchSize(config.AppConfig.MaxExportBatchSize))
	opts = append(opts, trace.WithBatchTimeout(time.Millisecond*time.Duration(config.AppConfig.BatchTimeout)))
	return opts
}

func NewSession(tracerProviders []trace.TracerProviderOption) *trace.TracerProvider {
	tp := trace.NewTracerProvider(
		tracerProviders...,
	)
	otel.SetTracerProvider(tp)
	return tp
}

// newOTLPGrpcExporter returns am otlp exporter.
func newOTLPGrpcExporter(ctx context.Context, endpoint string, additionalOpts ...otlptracegrpc.Option) (*otlptrace.Exporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithHeaders(map[string]string{"X-Test-Routing": "routing"}),
		otlptracegrpc.WithEndpoint(endpoint),
	}

	opts = append(opts, additionalOpts...)
	client := otlptracegrpc.NewClient(opts...)
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create an otlp grpc exporter: %w", err)
	}

	return exp, nil
}

// newOTLPHttpExporter returns am otlp exporter.
func newOTLPHttpExporter(ctx context.Context, endpoint string, additionalOpts ...otlptracehttp.Option) (*otlptrace.Exporter, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
	}

	opts = append(opts, additionalOpts...)
	client := otlptracehttp.NewClient(opts...)
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create an otlp http exporter: %w", err)
	}

	return exp, nil
}

func newJaegerExporter(endpoint string) (*jaeger.Exporter, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, fmt.Errorf("failed to create a jaeger exporter: %w", err)
	}

	return exp, nil
}

// newResource returns a resource describing this application.
func newResource(service string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
