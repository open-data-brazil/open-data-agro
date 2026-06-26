package fao

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the official download URL for a FAO catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	if entry.DatasetID.String() == "fao.food-price-index" {
		url := strings.TrimSpace(entry.SourceURL)
		if url == "" {
			url = defaultFFPICSVURL
		}
		return url, nil
	}
	url := strings.TrimSpace(entry.FAOBulkURL)
	if url == "" {
		url = defaultPricesBulkURL
	}
	if strings.TrimSpace(entry.DatasetID.String()) == "" {
		return "", fmt.Errorf("empty dataset id")
	}
	return url, nil
}
