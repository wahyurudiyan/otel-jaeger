package telemetry

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// enumeration constant for which protocol used
const (
	HTTP uint8 = iota
	GRPC
)

// setup client to connect web-service with Jaeger
func SetupTraceClient(ctx context.Context, protocol uint8, endpoint string) otlptrace.Client {
	switch protocol {
	case GRPC:
		return otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure(), otlptracegrpc.WithCompressor("gzip"))
	default:
		return otlptracehttp.NewClient(otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithInsecure(), otlptracehttp.WithCompression(otlptracehttp.NoCompression))
	}
}

func setupTraceProvider(ctx context.Context, traceClient otlptrace.Client) (*trace.TracerProvider, error) {
	// set resource
	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, err
	}

	// init trace exporter
	traceExporter, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, err
	}

	// init trace exporter
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			traceExporter,
			trace.WithBatchTimeout(time.Duration(time.Second*3)),
		),
		trace.WithResource(res),
	)

	return traceProvider, nil
}

func SetupTelemetrySDK(ctx context.Context, traceClient otlptrace.Client) (func(context.Context) error, error) {
	var err error
	var shutdownFuncs []func(context.Context) error

	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	traceProvider, err := setupTraceProvider(ctx, traceClient)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}

	shutdownFuncs = append(shutdownFuncs, traceProvider.Shutdown)
	otel.SetTracerProvider(traceProvider)

	return shutdown, nil
}
