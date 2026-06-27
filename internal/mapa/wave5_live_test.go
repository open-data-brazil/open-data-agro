//go:build live

package mapa

import (
	"context"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestLiveFetchSIPEAGROEstabelecimentos(t *testing.T) {
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.sipeagro-estabelecimentos"),
	}
	body, sourceURL, err := NewClient().FetchSIPEAGROSnapshot(context.Background(), entry)
	if err != nil {
		t.Fatalf("FetchSIPEAGROSnapshot: %v", err)
	}
	if len(body) < 1000 {
		t.Fatalf("body too small: %d bytes from %s", len(body), sourceURL)
	}
	if !containsAll(string(body), "linha_produto", "numero_registro_estabelecimento") {
		t.Fatal("merged CSV missing canonical headers")
	}
}

func TestLiveFetchSIGEFProducaoSementes(t *testing.T) {
	entry := catalog.RegistryEntry{
		DatasetID:     catalog.MustParseDatasetID("mapa.sigef-producao-sementes"),
		CKANPackageID: "dados-referentes-ao-controle-da-producao-de-sementes-sigef",
		CKANResourceFormat: "CSV",
		CKANResourceNameContains: "campos de produção",
	}
	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	result, err := NewClient().Download(context.Background(), url)
	if err != nil {
		t.Fatalf("Download: %v", err)
	}
	if len(result.Body) < 500 {
		t.Fatalf("body too small: %d", len(result.Body))
	}
}

func TestLiveFetchSISSERSeguroRural(t *testing.T) {
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.sisser-seguro-rural"),
	}
	body, sourceURL, err := NewClient().FetchSISSERSnapshot(context.Background(), entry)
	if err != nil {
		t.Fatalf("FetchSISSERSnapshot: %v", err)
	}
	if len(body) < 1000 {
		t.Fatalf("body too small: %d bytes from %s", len(body), sourceURL)
	}
}
