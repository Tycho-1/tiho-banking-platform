package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func otlpEndpoint() string {
	return strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
}

func tracingEnabled() bool {
	if otlpEndpoint() == "" {
		return false
	}
	if v := strings.TrimSpace(os.Getenv("ENABLE_TRACING")); v != "" && !strings.EqualFold(v, "true") {
		return false
	}
	return true
}

func serviceName() string {
	return envOrDefault("OTEL_SERVICE_NAME", "product-catalog")
}

// initTracer starts an OTLP HTTP tracer when OTEL_EXPORTER_OTLP_ENDPOINT is set.
// Returns a shutdown func (no-op if tracing is off).
func initTracer() func() {
	if !tracingEnabled() {
		log.Printf("tracing disabled (set OTEL_EXPORTER_OTLP_ENDPOINT to enable OTLP)")
		return func() {}
	}

	ctx := context.Background()
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		log.Printf("tracing: otlp exporter: %v (continuing without traces)", err)
		return func() {}
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(serviceName())),
	)
	if err != nil {
		log.Printf("tracing: resource: %v (continuing without traces)", err)
		return func() {}
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	log.Printf("tracing enabled (otlp) service=%s endpoint=%s", serviceName(), otlpEndpoint())

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := provider.Shutdown(shutdownCtx); err != nil {
			log.Printf("tracing shutdown: %v", err)
		}
	}
}

func skipTracePath(path string) bool {
	switch path {
	case "/metrics", "/health", "/ready":
		return true
	default:
		return false
	}
}

func wrapTracing(next http.Handler) http.Handler {
	if !tracingEnabled() {
		return next
	}
	return otelhttp.NewHandler(next, "product-catalog",
		otelhttp.WithFilter(func(r *http.Request) bool {
			return !skipTracePath(r.URL.Path)
		}),
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}
