package sagis

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the SAGIS portal URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	switch entry.DatasetID.String() {
	case "sagis.grain-supply-statistics":
		if u := strings.TrimSpace(entry.SourceURL); u != "" {
			return u, nil
		}
		return "https://www.sagis.org.za/", nil
	default:
		return "", fmt.Errorf("unsupported sagis dataset %s", entry.DatasetID)
	}
}
