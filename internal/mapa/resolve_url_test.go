package mapa

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestSafraSortKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want int
	}{
		{"Tábua de Risco Safra 2024/2025", 2024},
		{"invalid", -1},
	}
	for _, tc := range tests {
		if got := safraSortKey(tc.name); got != tc.want {
			t.Fatalf("safraSortKey(%q) = %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestResolveLatestSafraResourceURL(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ckanPackageResponse{
			Success: true,
			Result: struct {
				Resources []ckanResource `json:"resources"`
			}{
				Resources: []ckanResource{
					{
						Name:         "Tábua de Risco Safra 2022/2023",
						Format:       "CSV",
						URL:          "https://dados.agricultura.gov.br/old.csv",
						LastModified: "2023-01-01T00:00:00",
					},
					{
						Name:         "Tábua de Risco Safra 2024/2025",
						Format:       "CSV",
						URL:          "https://dados.agricultura.gov.br/dados-abertos-tabua-de-risco-2024.csv",
						LastModified: "2025-01-01T00:00:00",
					},
				},
			},
		})
	}))
	t.Cleanup(srv.Close)

	prev := ckanPackageShowURL
	ckanPackageShowURL = srv.URL
	t.Cleanup(func() { ckanPackageShowURL = prev })

	got, err := resolveLatestSafraResourceURL(context.Background(), "tabua-de-risco-zoneamento-agricola-de-risco-climatico", "CSV")
	if err != nil {
		t.Fatalf("resolveLatestSafraResourceURL: %v", err)
	}
	if got != "https://dados.agricultura.gov.br/dados-abertos-tabua-de-risco-2024.csv" {
		t.Fatalf("url: got %q", got)
	}
}

func TestResolveURLLatestSafra(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ckanPackageResponse{
			Success: true,
			Result: struct {
				Resources []ckanResource `json:"resources"`
			}{
				Resources: []ckanResource{
					{
						Name:         "Tábua de Risco Safra 2024/2025",
						Format:       "CSV",
						URL:          "https://dados.agricultura.gov.br/tabua-de-risco-2024.csv",
						LastModified: "2025-01-01T00:00:00",
					},
				},
			},
		})
	}))
	t.Cleanup(srv.Close)

	prev := ckanPackageShowURL
	ckanPackageShowURL = srv.URL
	t.Cleanup(func() { ckanPackageShowURL = prev })

	entry := catalog.RegistryEntry{
		DatasetID:          catalog.MustParseDatasetID("mapa.zarc-tabua-risco"),
		SourceURL:          "https://dados.agricultura.gov.br/dataset/tabua-de-risco-zoneamento-agricola-de-risco-climatico",
		CKANPackageID:      "tabua-de-risco-zoneamento-agricola-de-risco-climatico",
		CKANResourceFormat: "CSV",
	}

	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url != "https://dados.agricultura.gov.br/tabua-de-risco-2024.csv" {
		t.Fatalf("url: got %q", url)
	}
}
