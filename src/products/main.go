package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type server struct {
	version  string
	catalog  *catalog
	verifier *tokenVerifier
}

func main() {
	port := envOrDefault("PORT", "8080")
	version := envOrDefault("VERSION", "dev")
	dataPath := envOrDefault("PRODUCTS_DATA_PATH", "data/products.json")
	pubKeyPath := envOrDefault("PUB_KEY_PATH", os.Getenv("JWT_PUBLIC_KEY_PATH"))

	if pubKeyPath == "" {
		log.Fatal("PUB_KEY_PATH or JWT_PUBLIC_KEY_PATH must be set")
	}

	catalog, err := loadCatalog(dataPath)
	if err != nil {
		log.Fatalf("load catalog: %v", err)
	}
	setCatalogProductsLoaded(len(catalog.products))

	verifier, err := newTokenVerifier(pubKeyPath)
	if err != nil {
		log.Fatalf("init auth: %v", err)
	}

	srv := &server{
		version:  version,
		catalog:  catalog,
		verifier: verifier,
	}

	addr := ":" + port
	log.Printf("product-catalog listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.handler()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func (s *server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writePlain(w, http.StatusOK, "ok")
}

func (s *server) handleReady(w http.ResponseWriter, _ *http.Request) {
	writePlain(w, http.StatusOK, "ok")
}

func (s *server) handleVersion(w http.ResponseWriter, _ *http.Request) {
	writePlain(w, http.StatusOK, s.version)
}

func (s *server) handleProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writePlain(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if _, err := s.verifier.verify(r.Header.Get("Authorization")); err != nil {
		writePlain(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	productType := strings.TrimSpace(r.URL.Query().Get("type"))
	products := s.catalog.list(productType)
	recordProductsRequest(productType)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]Product{"products": products}); err != nil {
		log.Printf("encode products: %v", err)
	}
}

func writePlain(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
