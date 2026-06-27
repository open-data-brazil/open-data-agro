package mapa

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultSIFAbateURL = "https://dados.agricultura.gov.br/dataset/062166e3-b515-4274-8e7d-68aadd64b820/resource/341dc717-4716-42ab-b189-c8d7a9d2a1ba/download/sigsifrelatorioabates.csv"

// ResolveSIFAbateURL returns the official MAPA SIGSIF slaughter statistics CSV URL.
func ResolveSIFAbateURL(entry catalog.RegistryEntry) (string, error) {
	if entry.DatasetID.String() == "mapa.sif-abate-estatisticas" {
		raw := strings.TrimSpace(entry.SourceURL)
		if raw == "" {
			return defaultSIFAbateURL, nil
		}
		if !strings.Contains(strings.ToLower(raw), "gov.br") {
			return "", fmt.Errorf("source_url for %s must be on gov.br", entry.DatasetID)
		}
		return raw, nil
	}
	return ResolveURL(entry)
}
