package mapa

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestResolveURLLatestSafra(t *testing.T) {
	t.Parallel()

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
	if url == "" || !contains(url, "dados.agricultura.gov.br") {
		t.Fatalf("url: %q", url)
	}
	if !contains(url, "tabua-de-risco") && !contains(url, "dados-abertos-tabua-de-risco") {
		t.Fatalf("expected safra tabua resource url, got %q", url)
	}
}

func contains(value, part string) bool {
	return len(value) >= len(part) && indexOf(value, part) >= 0
}

func indexOf(value, part string) int {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}
