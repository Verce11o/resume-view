package tracer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

type JaegerTracing struct {
	Exporter tracesdk.SpanExporter
	Provider *tracesdk.TracerProvider
	trace.Tracer
}

func NewJaegerExporter(ctx context.Context, endpoint string) (tracesdk.SpanExporter, error) {
	spanExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(endpoint))
	if err != nil {
		return nil, fmt.Errorf("otlptracegrpc.New: %w", err)
	}

	return spanExporter, nil
}

func NewTraceProvider(exp tracesdk.SpanExporter, serviceName string) (*tracesdk.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("resource.Merge: %w", err)
	}

	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(r),
	), nil
}

func InitTracer(ctx context.Context, serviceName, endpoint string) (*JaegerTracing, error) {
	exporter, err := NewJaegerExporter(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("NewJaegerExporter: %w", err)
	}

	tp, err := NewTraceProvider(exporter, serviceName)
	if err != nil {
		return nil, fmt.Errorf("NewTraceProvider: %w", err)
	}

	otel.SetTracerProvider(tp)

	tracer := tp.Tracer("main tracer")

	return &JaegerTracing{
		Exporter: exporter,
		Provider: tp,
		Tracer:   tracer,
	}, nil
}
