package otlp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const operation = "trace"

func NewTracer(ctx context.Context, endpoint, name string) (*Tracer, error) {
	exp, err := newOTLPExporter(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create otlp exporter: %w", operation, err)
	}

	tp, err := newTraceProvider(exp)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create trace provider: %w", operation, err)
	}

	otel.SetTracerProvider(tp)
	tracer := tp.Tracer(name)

	return &Tracer{
		tracer: tracer,
	}, nil
}

type Tracer struct {
	tracer trace.Tracer
}

func newOTLPExporter(ctx context.Context, endpoint string) (trace.SpanExporter, error) {
	insecureOpt := otlptracehttp.WithInsecure()
	endpointOpt := otlptracehttp.WithEndpoint(endpoint)
	return otlptracehttp.New(ctx, insecureOpt, endpointOpt)
}

func newTraceProvider(exp trace.SpanExporter) (*trace.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("myapp"),
		),
	)
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(r),
	), nil
}
