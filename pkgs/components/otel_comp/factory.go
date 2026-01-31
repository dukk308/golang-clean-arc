package tracer

import (
	"context"
	"time"

	"github.com/dukk308/golang-clean-arch-starter/pkgs/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type TracerProviderFactory interface {
	Create(options ...TracerProviderOption) (*trace.TracerProvider, error)
}

type DefaultTracerProviderFactory struct {
	logger logger.Logger
}

func NewDefaultTracerProviderFactory(logger logger.Logger) TracerProviderFactory {
	return &DefaultTracerProviderFactory{
		logger: logger,
	}
}

func (f *DefaultTracerProviderFactory) Create(options ...TracerProviderOption) (*trace.TracerProvider, error) {

	appliedOptions := defaultTracerProviderOptions
	for _, opt := range options {
		opt(&appliedOptions)
	}

	ctx := context.Background()

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appliedOptions.Name),
		),
	)
	if err != nil {
		f.logger.Error("failed to create resource for tracer provider", err)

		return nil, err
	}

	spanProcessor, err := f.createSpanProcessor(ctx, appliedOptions)
	if err != nil {
		f.logger.Error("failed to create span processor for tracer provider", err)

		return nil, err
	}

	prefixProcessor := newSpanNamePrefixProcessor(spanProcessor)

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(prefixProcessor),
	)

	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tracerProvider, nil
}

func (f *DefaultTracerProviderFactory) createSpanProcessor(ctx context.Context, opts options) (trace.SpanProcessor, error) {
	switch opts.Exporter {
	case Stdout:
		stdoutExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}

		wrappedExporter := newSpanNamePrefixExporter(stdoutExporter)
		return trace.NewBatchSpanProcessor(wrappedExporter), nil
	case OtlpGrpc:
		f.logger.Info("Connecting to OTLP gRPC collector", "endpoint", opts.Collector)

		endpoint := opts.Collector
		otlpGrpcExporter, err := otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			f.logger.Error("failed to create otlp-grpc span exporter",
				"endpoint", endpoint,
				"error", err,
				"hint", "Make sure Jaeger is running with COLLECTOR_OTLP_ENABLED=true and listening on the specified endpoint")
			return nil, err
		}

		f.logger.Info("OTLP gRPC exporter created successfully",
			"endpoint", endpoint,
			"batch_timeout", "5s",
			"hint", "Traces will be batched and exported automatically every 5 seconds")

		wrappedExporter := newSpanNamePrefixExporter(otlpGrpcExporter)
		batchProcessor := trace.NewBatchSpanProcessor(
			wrappedExporter,
			trace.WithBatchTimeout(2*time.Second),
			trace.WithMaxExportBatchSize(512),
		)

		f.logger.Info("Batch span processor configured",
			"batch_timeout", "2s",
			"max_batch_size", 512)

		return batchProcessor, nil
	default:
		noopExporter := tracetest.NewNoopExporter()
		wrappedExporter := newSpanNamePrefixExporter(noopExporter)
		return trace.NewBatchSpanProcessor(wrappedExporter), nil
	}
}
