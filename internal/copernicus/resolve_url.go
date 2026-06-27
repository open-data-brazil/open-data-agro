package copernicus

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the Copernicus CDS portal URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	switch entry.DatasetID.String() {
	case "copernicus.era5-agroclimate":
		if u := strings.TrimSpace(entry.SourceURL); u != "" {
			return u, nil
		}
		return "https://cds.climate.copernicus.eu/", nil
	default:
		return "", fmt.Errorf("unsupported copernicus dataset %s", entry.DatasetID)
	}
}
