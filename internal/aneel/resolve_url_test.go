package aneel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestResolveURLCKANAcionamento(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ckanPackageResponse{
			Success: true,
			Result: struct {
				Resources []ckanResource `json:"resources"`
			}{
				Resources: []ckanResource{
					{
						Name:         "Bandeira Tarifária - Adicional CSV",
						Format:       "CSV",
						URL:          "https://dadosabertos.aneel.gov.br/resource/adicional.csv",
						LastModified: "2026-06-01T00:00:00",
					},
					{
						Name:         "Bandeira Tarifária - Acionamento CSV",
						Format:       "CSV",
						URL:          "https://dadosabertos.aneel.gov.br/resource/acionamento.csv",
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
		DatasetID:                catalog.MustParseDatasetID("aneel.tarifas-energia"),
		CKANPackageID:            "bandeiras-tarifarias",
		CKANResourceFormat:       "CSV",
		CKANResourceNameContains: "Acionamento",
	}

	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url != "https://dadosabertos.aneel.gov.br/resource/acionamento.csv" {
		t.Fatalf("url: got %q", url)
	}
}
