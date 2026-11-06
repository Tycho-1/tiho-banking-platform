package main

import (
	"context"
	"io"
	"log"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// initLogs sends standard-library logs to the Collector when
// OTEL_EXPORTER_OTLP_ENDPOINT is set. Stdout logging stays in place.
// Returns a shutdown func (no-op if the endpoint is unset).
func initLogs() func() {
	if otlpEndpoint() == "" {
		return func() {}
	}

	ctx := context.Background()
	exporter, err := otlploghttp.New(ctx)
	if err != nil {
		log.Printf("logs: otlp exporter: %v (continuing with stdout only)", err)
		return func() {}
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(serviceName())),
	)
	if err != nil {
		log.Printf("logs: resource: %v (continuing with stdout only)", err)
		return func() {}
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)
	log.SetOutput(io.MultiWriter(log.Writer(), &otelLogWriter{
		logger: provider.Logger(serviceName()),
	}))
	log.Printf("OTLP logs enabled endpoint=%s service=%s", otlpEndpoint(), serviceName())

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := provider.Shutdown(shutdownCtx); err != nil {
			log.Printf("logs shutdown: %v", err)
		}
	}
}

// otelLogWriter copies each standard-log line to the OTel logger.
type otelLogWriter struct {
	logger otellog.Logger
}

func (w *otelLogWriter) Write(p []byte) (int, error) {
	var record otellog.Record
	record.SetTimestamp(time.Now())
	record.SetSeverity(otellog.SeverityInfo)
	record.SetBody(attribute.StringValue(string(p)))
	w.logger.Emit(context.Background(), record)
	return len(p), nil
}
