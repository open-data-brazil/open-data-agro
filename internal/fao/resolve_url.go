package fao

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the FAOSTAT bulk download URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	url := strings.TrimSpace(entry.FAOBulkURL)
	if url == "" {
		url = defaultPricesBulkURL
	}
	if strings.TrimSpace(entry.DatasetID.String()) == "" {
		return "", fmt.Errorf("empty dataset id")
	}
	return url, nil
}
