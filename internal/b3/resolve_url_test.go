package b3

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestResolveURLPortal(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:       catalog.MustParseDatasetID("b3.futuro-soja"),
		SourcePortalURL: "https://www.b3.com.br/pt_br/market-data-e-indices/servicos-de-dados/market-data/historico/boletins-diarios/pesquisa-por-pregao/",
	}

	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" || !contains(url, "b3.com.br") {
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
