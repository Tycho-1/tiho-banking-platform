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
