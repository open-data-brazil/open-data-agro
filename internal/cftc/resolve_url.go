package cftc

import (
	"fmt"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the official CFTC COT Socrata API URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	switch entry.DatasetID.String() {
	case "cftc.cot-agricultural-futures":
		return buildCOTURL(entry)
	default:
		return "", fmt.Errorf("unsupported cftc dataset %s", entry.DatasetID)
	}
}
