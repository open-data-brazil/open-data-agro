package japan

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the MAFF Japan portal URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	switch entry.DatasetID.String() {
	case "japan.maff-ag-trade":
		if u := strings.TrimSpace(entry.SourceURL); u != "" {
			return u, nil
		}
		return "https://www.maff.go.jp/e/data/stat/export/index.html", nil
	default:
		return "", fmt.Errorf("unsupported japan dataset %s", entry.DatasetID)
	}
}
