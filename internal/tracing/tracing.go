package tracing

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	SERVICE_NAME = "order-manager"
)

func createResource() *resource.Resource {
	r, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		resource.Default().SchemaURL(),
		semconv.ServiceNameKey.String(SERVICE_NAME),
		semconv.ServiceVersionKey.String("v0.1.0"),
	))

	if err != nil {
		log.Fatal(err.Error())
	}
	return r
}

type TraceContext interface {
	NewSpan(ctx context.Context, name string) (context.Context, trace.Span)
	NewSpanWithAttributes(ctx context.Context, name string, kvs ...attribute.KeyValue) (context.Context, trace.Span)
	SetSpanFailed(trace.Span, error)
	Close() error
}

type TraceContextNoop struct {
}

func (tcm *TraceContextNoop) NewSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return trace.NewNoopTracerProvider().Tracer("noop").Start(ctx, name)
}

func (tcm *TraceContextNoop) NewSpanWithAttributes(ctx context.Context, name string, kvs ...attribute.KeyValue) (context.Context, trace.Span) {
	return trace.NewNoopTracerProvider().Tracer("noop").Start(ctx, name)
}

func (tcm *TraceContextNoop) SetSpanFailed(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func (tcm *TraceContextNoop) Close() error { return nil }

type OtelTracingContext struct {
	Resource      *resource.Resource
	TraceProvider *sdktrace.TracerProvider
}

func (tc *OtelTracingContext) NewSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	spanCtx, span := tc.TraceProvider.Tracer(SERVICE_NAME).Start(ctx, name)
	return spanCtx, span
}

func (tc *OtelTracingContext) NewSpanWithAttributes(ctx context.Context, name string, kvs ...attribute.KeyValue) (context.Context, trace.Span) {
	spanCtx, span := tc.TraceProvider.Tracer(SERVICE_NAME).Start(ctx, name)
	span.SetAttributes(kvs...)
	return spanCtx, span
}

func (tc *OtelTracingContext) SetSpanFailed(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func (tc *OtelTracingContext) Close() error {
	if err := tc.TraceProvider.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
	return nil
}

func NewOtelTracingContext() TraceContext {
	otel_host, found := os.LookupEnv("OTEL_HOST")

	if !found {
		otel_host = "otelp"
	}

	res := createResource()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	connection, err := grpc.DialContext(ctx, otel_host+":4317", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create gRPC connection to collector: %v", err.Error())
	}

	log.Printf("connection established: %s ", connection.Target())

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(connection), otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create traceExporter: %s", err.Error())
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}))

	return &OtelTracingContext{
		Resource:      res,
		TraceProvider: tracerProvider,
	}
}
