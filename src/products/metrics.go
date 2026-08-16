package main

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var registerCollectorsOnce sync.Once

func registerCollectors() {
	registerCollectorsOnce.Do(func() {
		registerCollector(collectors.NewGoCollector())
		registerCollector(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	})
}

func registerCollector(c prometheus.Collector) {
	if err := prometheus.Register(c); err != nil {
		var already prometheus.AlreadyRegisteredError
		if !errors.As(err, &already) {
			panic(err)
		}
	}
}

// Prometheus metric names use namespace product_catalog → product_catalog_http_requests_total, etc.
var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "product_catalog",
		Name:      "http_requests_total",
		Help:      "Total HTTP requests handled by product-catalog.",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "product_catalog",
		Name:      "http_request_duration_seconds",
		Help:      "HTTP request latency in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path"})

	// Business metric: successful product list requests by filter (?type= or all).
	productsRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "product_catalog",
		Name:      "products_requests_total",
		Help:      "Successful GET /api/products requests grouped by type filter.",
	}, []string{"type"})

	catalogProductsLoaded = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "product_catalog",
		Name:      "catalog_products_loaded",
		Help:      "Number of products loaded from products.json at startup.",
	})
)

func setCatalogProductsLoaded(count int) {
	catalogProductsLoaded.Set(float64(count))
}

func recordProductsRequest(productType string) {
	label := productType
	if label == "" {
		label = "all"
	}
	productsRequestsTotal.WithLabelValues(label).Inc()
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// metricsMiddleware records request count and latency for every route except /metrics
// (Prometheus scrapes should not inflate application traffic metrics).
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		path := r.URL.Path
		status := strconv.Itoa(rec.status)
		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())
	})
}

func (s *server) handler() http.Handler {
	registerCollectors()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/ready", s.handleReady)
	mux.HandleFunc("/version", s.handleVersion)
	mux.HandleFunc("/api/products", s.handleProducts)
	mux.Handle("/metrics", promhttp.Handler())
	return metricsMiddleware(mux)
}
