package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Product struct {
	ID                  string `json:"id"`
	Type                string `json:"type"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	APY                 string `json:"apy,omitempty"`
	MinimumBalance      string `json:"minimumBalance,omitempty"`
	Term                string `json:"term,omitempty"`
	Rate                string `json:"rate,omitempty"`
	EarlyWithdrawalNote string `json:"earlyWithdrawalNote,omitempty"`
	RiskLevel           string `json:"riskLevel,omitempty"`
	Disclaimer          string `json:"disclaimer,omitempty"`
}

type catalog struct {
	products []Product
}

func loadCatalog(path string) (*catalog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read catalog: %w", err)
	}

	var payload struct {
		Products []Product `json:"products"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("parse catalog: %w", err)
	}

	return &catalog{products: payload.Products}, nil
}

func (c *catalog) list(productType string) []Product {
	if productType == "" {
		out := make([]Product, len(c.products))
		copy(out, c.products)
		return out
	}

	filter := strings.ToLower(strings.TrimSpace(productType))
	out := make([]Product, 0)
	for _, product := range c.products {
		if strings.EqualFold(product.Type, filter) {
			out = append(out, product)
		}
	}
	return out
}
