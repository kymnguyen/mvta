package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
)

type Telemetry struct {
	MeterProvider  metric.MeterProvider
	TracerProvider *trace.TracerProvider
	logger         *zap.Logger
}

func InitTelemetry(serviceName, jaegerEndpoint string, logger *zap.Logger) (*Telemetry, error) {
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Setup tracing with noop exporter (can be replaced with OTLP later)
	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
	)
	otel.SetTracerProvider(tracerProvider)

	// Setup metrics
	meterProvider, err := setupMetrics(res)
	if err != nil {
		return nil, err
	}
	otel.SetMeterProvider(meterProvider)

	return &Telemetry{
		MeterProvider:  meterProvider,
		TracerProvider: tracerProvider,
		logger:         logger,
	}, nil
}

func setupMetrics(res *resource.Resource) (metric.MeterProvider, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(res),
	)

	return mp, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := t.TracerProvider.Shutdown(ctx); err != nil {
		t.logger.Error("Failed to shutdown tracer provider", zap.Error(err))
		return err
	}

	return nil
}
