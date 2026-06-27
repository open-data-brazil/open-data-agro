package fred

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the FRED graph CSV URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	switch entry.DatasetID.String() {
	case "fred.commodity-indexes":
		if u := strings.TrimSpace(entry.SourceURL); u != "" {
			return u, nil
		}
		return defaultCommodityIndexURL, nil
	default:
		return "", fmt.Errorf("unsupported fred dataset %s", entry.DatasetID)
	}
}
