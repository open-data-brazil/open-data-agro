package antaq

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the official ANTAQ bulk export URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url", entry.DatasetID)
	}
	if !strings.Contains(strings.ToLower(raw), "gov.br") && !strings.Contains(strings.ToLower(raw), "antaq.gov.br") {
		return "", fmt.Errorf("source_url for %s must be on gov.br or antaq.gov.br", entry.DatasetID)
	}
	return raw, nil
}
