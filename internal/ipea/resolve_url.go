package ipea

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL validates the catalog base URL for an IPEA dataset.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url", entry.DatasetID)
	}
	if !strings.Contains(strings.ToLower(raw), "ipeadata.gov.br") {
		return "", fmt.Errorf("source_url for %s must be on ipeadata.gov.br", entry.DatasetID)
	}
	return raw, nil
}
