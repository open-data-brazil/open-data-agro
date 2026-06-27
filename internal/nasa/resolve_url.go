package nasa

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the NASA POWER API URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	switch entry.DatasetID.String() {
	case "nasa.power-agroclimatology":
		if u := strings.TrimSpace(entry.SourceURL); u != "" {
			return u, nil
		}
		return defaultPOWERURL, nil
	default:
		return "", fmt.Errorf("unsupported nasa dataset %s", entry.DatasetID)
	}
}
