package jrc

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the official JRC MARS yield forecast CSV URL.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	switch entry.DatasetID.String() {
	case "jrc.mars-crop-yield":
		if u := strings.TrimSpace(entry.SourceURL); u != "" {
			return u, nil
		}
		return defaultMARSYieldURL, nil
	default:
		return "", fmt.Errorf("unsupported jrc dataset %s", entry.DatasetID)
	}
}
