package b3

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const pesquisaPregaoPortalURL = "https://www.b3.com.br/pt_br/market-data-e-indices/servicos-de-dados/market-data/historico/boletins-diarios/pesquisa-por-pregao/"

// ResolveURL validates the catalog portal URL for a B3 futures dataset.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourcePortalURL)
	if raw == "" {
		raw = strings.TrimSpace(entry.SourceURL)
	}
	if raw == "" {
		return pesquisaPregaoPortalURL, nil
	}
	if !strings.Contains(strings.ToLower(raw), "b3.com.br") {
		return "", fmt.Errorf("source_url for %s must be on b3.com.br", entry.DatasetID)
	}
	return raw, nil
}
