package antt

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestResolveURLCKAN(t *testing.T) {
	t.Parallel()

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
	if url == "" || !contains(url, "dados.antt.gov.br") {
		t.Fatalf("url: %q", url)
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
