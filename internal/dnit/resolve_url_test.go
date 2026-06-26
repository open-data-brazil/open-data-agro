package dnit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestResolveLatestCKANResourceURL(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ckanPackageResponse{
			Success: true,
			Result: struct {
				Resources []ckanResource `json:"resources"`
			}{
				Resources: []ckanResource{
					{Name: "Jurisdição 2020", Format: "CSV", URL: "https://servicos.dnit.gov.br/old.csv", LastModified: "2020-01-01"},
					{Name: "Jurisdição 2022", Format: "CSV", URL: "https://servicos.dnit.gov.br/new.csv", LastModified: "2022-04-01"},
				},
			},
		})
	}))
	defer srv.Close()

	prev := ckanPackageShowURL
	ckanPackageShowURL = srv.URL
	t.Cleanup(func() { ckanPackageShowURL = prev })

	got, err := ResolveLatestCKANResourceURL(context.Background(), "jurisdicao-de-vias", "CSV")
	if err != nil {
		t.Fatalf("ResolveLatestCKANResourceURL: %v", err)
	}
	if got != "https://servicos.dnit.gov.br/new.csv" {
		t.Fatalf("got %q want latest csv url", got)
	}
}

func TestResolveURLCKAN(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(ckanPackageResponse{
			Success: true,
			Result: struct {
				Resources []ckanResource `json:"resources"`
			}{
				Resources: []ckanResource{
					{Name: "SNV CSV", Format: "CSV", URL: "https://servicos.dnit.gov.br/snv.csv", LastModified: "2022-04-01"},
				},
			},
		})
	}))
	defer srv.Close()

	prev := ckanPackageShowURL
	ckanPackageShowURL = srv.URL
	t.Cleanup(func() { ckanPackageShowURL = prev })

	entry := catalog.RegistryEntry{
		DatasetID:       catalog.MustParseDatasetID("dnit.snv-rodovias-federais"),
		CKANPackageID:   "jurisdicao-de-vias",
		CKANResourceFormat: "CSV",
	}

	got, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if got != "https://servicos.dnit.gov.br/snv.csv" {
		t.Fatalf("got %q", got)
	}
}
