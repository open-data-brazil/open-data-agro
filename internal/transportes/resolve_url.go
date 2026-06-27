package transportes

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/dnit"
)

// ResolveURL returns the DNIT SNV jurisdiction CSV used as MTR BIT road metadata.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	return dnit.ResolveURL(entry)
}

// ResolveURLString validates a resolved download URL for tests.
func ResolveURLString(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url", entry.DatasetID)
	}
	if !strings.Contains(strings.ToLower(raw), "gov.br") {
		return "", fmt.Errorf("source_url for %s must be on gov.br", entry.DatasetID)
	}
	return raw, nil
}
