package bndes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestResolveURLCKANCNAE(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ckanPackageResponse{
			Success: true,
			Result: struct {
				Resources []ckanResource `json:"resources"`
			}{
				Resources: []ckanResource{
					{
						Name:         "Por UF - Desembolsos CSV",
						Format:       "CSV",
						URL:          "https://dadosabertos.bndes.gov.br/resource/uf.csv",
						LastModified: "2026-06-01T00:00:00",
					},
					{
						Name:         "Por setor CNAE - Desembolsos CSV",
						Format:       "CSV",
						URL:          "https://dadosabertos.bndes.gov.br/resource/cnae.csv",
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
		DatasetID:                catalog.MustParseDatasetID("bndes.financiamento-agro"),
		CKANPackageID:            "desembolsos",
		CKANResourceFormat:       "CSV",
		CKANResourceNameContains: "Por setor CNAE",
	}

	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url != "https://dadosabertos.bndes.gov.br/resource/cnae.csv" {
		t.Fatalf("url: got %q", url)
	}
}
