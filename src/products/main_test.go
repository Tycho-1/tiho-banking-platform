package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestCatalogListFiltersByType(t *testing.T) {
	catalog := &catalog{
		products: []Product{
			{ID: "1", Type: "savings", Name: "A"},
			{ID: "2", Type: "cd", Name: "B"},
			{ID: "3", Type: "investment", Name: "C"},
		},
	}

	if got := len(catalog.list("")); got != 3 {
		t.Fatalf("expected 3 products, got %d", got)
	}
	if got := len(catalog.list("savings")); got != 1 {
		t.Fatalf("expected 1 savings product, got %d", got)
	}
	if got := len(catalog.list("cd")); got != 1 {
		t.Fatalf("expected 1 cd product, got %d", got)
	}
}

func TestHandleProductsRequiresJWT(t *testing.T) {
	dir := t.TempDir()
	dataPath := filepath.Join(dir, "products.json")
	if err := os.WriteFile(dataPath, []byte(`{"products":[{"id":"x","type":"savings","name":"Demo"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}

	keyPath, token := writeTestJWT(t)
	catalog, err := loadCatalog(dataPath)
	if err != nil {
		t.Fatal(err)
	}
	verifier, err := newTokenVerifier(keyPath)
	if err != nil {
		t.Fatal(err)
	}

	srv := &server{
		version:  "test",
		catalog:  catalog,
		verifier: verifier,
	}

	unauthReq := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	unauthRec := httptest.NewRecorder()
	srv.handleProducts(unauthRec, unauthReq)
	if unauthRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", unauthRec.Code)
	}

	authReq := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	authReq.Header.Set("Authorization", "Bearer "+token)
	authRec := httptest.NewRecorder()
	srv.handleProducts(authRec, authReq)
	if authRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", authRec.Code)
	}

	var payload struct {
		Products []Product `json:"products"`
	}
	if err := json.Unmarshal(authRec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Products) != 1 || payload.Products[0].ID != "x" {
		t.Fatalf("unexpected payload: %+v", payload.Products)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	srv := testServer(t)
	handler := srv.handler()

	// Prime middleware so counters exist in exposition output.
	prime := httptest.NewRequest(http.MethodGet, "/health", nil)
	handler.ServeHTTP(httptest.NewRecorder(), prime)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "product_catalog_http_requests_total") {
		t.Fatalf("expected product_catalog_http_requests_total in metrics output")
	}
	if !strings.Contains(body, "product_catalog_http_request_duration_seconds") {
		t.Fatalf("expected product_catalog_http_request_duration_seconds in metrics output")
	}
	if !strings.Contains(body, "go_goroutines") {
		t.Fatalf("expected go_goroutines from Go collector in metrics output")
	}
}

func TestMetricsMiddlewareRecordsRequests(t *testing.T) {
	srv := testServer(t)
	handler := srv.handler()

	healthReq := httptest.NewRequest(http.MethodGet, "/health", nil)
	healthRec := httptest.NewRecorder()
	handler.ServeHTTP(healthRec, healthReq)
	if healthRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", healthRec.Code)
	}

	metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsRec := httptest.NewRecorder()
	handler.ServeHTTP(metricsRec, metricsReq)
	body := metricsRec.Body.String()
	if !strings.Contains(body, `product_catalog_http_requests_total{method="GET",path="/health",status="200"}`) {
		t.Fatalf("expected /health request in metrics, got:\n%s", body)
	}
	if strings.Contains(body, `path="/metrics"`) {
		t.Fatalf("did not expect /metrics scrape in http_requests_total, got:\n%s", body)
	}
}

func TestProductsRequestByTypeMetric(t *testing.T) {
	srv, token := testServerWithToken(t)
	handler := srv.handler()

	for _, q := range []struct {
		query string
	}{
		{""},
		{"?type=savings"},
	} {
		req := httptest.NewRequest(http.MethodGet, "/api/products"+q.query, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("query %q: expected 200, got %d", q.query, rec.Code)
		}
	}

	metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsRec := httptest.NewRecorder()
	handler.ServeHTTP(metricsRec, metricsReq)
	body := metricsRec.Body.String()

	if !strings.Contains(body, `product_catalog_products_requests_total{type="all"}`) {
		t.Fatalf("expected all-type product request metric, got:\n%s", body)
	}
	if !strings.Contains(body, `product_catalog_products_requests_total{type="savings"}`) {
		t.Fatalf("expected savings product request metric, got:\n%s", body)
	}
}

func testServer(t *testing.T) *server {
	srv, _ := testServerWithToken(t)
	return srv
}

func testServerWithToken(t *testing.T) (*server, string) {
	t.Helper()

	dir := t.TempDir()
	dataPath := filepath.Join(dir, "products.json")
	if err := os.WriteFile(dataPath, []byte(`{"products":[{"id":"x","type":"savings","name":"Demo"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}

	keyPath, token := writeTestJWT(t)
	catalog, err := loadCatalog(dataPath)
	if err != nil {
		t.Fatal(err)
	}
	verifier, err := newTokenVerifier(keyPath)
	if err != nil {
		t.Fatal(err)
	}

	return &server{
		version:  "test",
		catalog:  catalog,
		verifier: verifier,
	}, token
}

func writeTestJWT(t *testing.T) (string, string) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	pubDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubDER,
	})

	keyPath := filepath.Join(t.TempDir(), "publickey")
	if err := os.WriteFile(keyPath, pubPEM, 0o644); err != nil {
		t.Fatal(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user": "testuser",
		"acct": "1234567890",
		"name": "Test User",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	return keyPath, signed
}
