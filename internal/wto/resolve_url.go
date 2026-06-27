package wto

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the WTO Stats API URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	switch entry.DatasetID.String() {
	case "wto.its-trade-statistics":
		if u := strings.TrimSpace(entry.SourceURL); u != "" {
			return u, nil
		}
		return defaultWTOAPIURL, nil
	default:
		return "", fmt.Errorf("unsupported wto dataset %s", entry.DatasetID)
	}
}
