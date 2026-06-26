package antt

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestResolveLatestCKANResourceURL(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("id") != "praca-de-pedagio" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(ckanPackageResponse{
			Success: true,
			Result: struct {
				Resources []ckanResource `json:"resources"`
			}{
				Resources: []ckanResource{
					{
						Name:         "older",
						Format:       "CSV",
						URL:          "https://dados.antt.gov.br/resource/old.csv",
						LastModified: "2024-01-01T00:00:00",
					},
					{
						Name:         "latest",
						Format:       "CSV",
						URL:          "https://dados.antt.gov.br/resource/latest.csv",
						LastModified: "2026-06-01T00:00:00",
					},
				},
			},
		})
	}))
	t.Cleanup(srv.Close)

	prev := ckanPackageShowURL
	ckanPackageShowURL = srv.URL
	t.Cleanup(func() { ckanPackageShowURL = prev })

	got, err := ResolveLatestCKANResourceURL(context.Background(), "praca-de-pedagio", "CSV")
	if err != nil {
		t.Fatalf("ResolveLatestCKANResourceURL: %v", err)
	}
	if got != "https://dados.antt.gov.br/resource/latest.csv" {
		t.Fatalf("url: got %q", got)
	}
}

func TestResolveURLCKAN(t *testing.T) {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ckanPackageResponse{
			Success: true,
			Result: struct {
				Resources []ckanResource `json:"resources"`
			}{
				Resources: []ckanResource{
					{
						Format:       "CSV",
						URL:          "https://dados.antt.gov.br/resource/pracas.csv",
						LastModified: "2026-06-01T00:00:00",
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
		DatasetID:          catalog.MustParseDatasetID("antt.pracas-pedagio"),
		SourceURL:          "https://dados.antt.gov.br/dataset/praca-de-pedagio",
		CKANPackageID:      "praca-de-pedagio",
		CKANResourceFormat: "CSV",
	}

	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url != "https://dados.antt.gov.br/resource/pracas.csv" {
		t.Fatalf("url: got %q", url)
	}
}

func TestResolveURLCKANNameFilter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ckanPackageResponse{
			Success: true,
			Result: struct {
				Resources []ckanResource `json:"resources"`
			}{
				Resources: []ckanResource{
					{
						Name:         "Volume 2024 Mensal consolidado",
						Format:       "CSV",
						URL:          "https://dados.antt.gov.br/resource/volume-2024.csv",
						LastModified: "2025-01-01T00:00:00",
					},
					{
						Name:         "Volume 2026 Mensal consolidado",
						Format:       "CSV",
						URL:          "https://dados.antt.gov.br/resource/volume-2026.csv",
						LastModified: "2026-01-01T00:00:00",
					},
					{
						Name:         "Volume 2026 Diário",
						Format:       "CSV",
						URL:          "https://dados.antt.gov.br/resource/volume-daily.csv",
						LastModified: "2026-06-01T00:00:00",
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
		DatasetID:                catalog.MustParseDatasetID("antt.volume-trafego-pedagio"),
		CKANPackageID:            "volume-trafego-praca-pedagio",
		CKANResourceFormat:       "CSV",
		CKANResourceNameContains: "Mensal consolidado",
	}

	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url != "https://dados.antt.gov.br/resource/volume-2026.csv" {
		t.Fatalf("url: got %q", url)
	}
}

func TestResolveURLCKANYearInName(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ckanPackageResponse{
			Success: true,
			Result: struct {
				Resources []ckanResource `json:"resources"`
			}{
				Resources: []ckanResource{
					{
						Name:         "Receita por Praça - 2023",
						Format:       "CSV",
						URL:          "https://dados.antt.gov.br/resource/receita-2023.csv",
						LastModified: "2024-01-01T00:00:00",
					},
					{
						Name:         "Receita por Praça - 2025",
						Format:       "CSV",
						URL:          "https://dados.antt.gov.br/resource/receita-2025.csv",
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
		DatasetID:                catalog.MustParseDatasetID("antt.receita-por-praca"),
		CKANPackageID:            "receita-por-praca",
		CKANResourceFormat:       "CSV",
		CKANResourceNameContains: "Receita por Praça -",
	}

	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url != "https://dados.antt.gov.br/resource/receita-2025.csv" {
		t.Fatalf("url: got %q", url)
	}
}
