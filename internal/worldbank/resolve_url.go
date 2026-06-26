package worldbank

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the World Bank Pink Sheet download URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	url := strings.TrimSpace(entry.WorldBankPinkSheetURL)
	if url == "" {
		url = defaultPinkSheetURL
	}
	if strings.TrimSpace(entry.DatasetID.String()) == "" {
		return "", fmt.Errorf("empty dataset id")
	}
	return url, nil
}
