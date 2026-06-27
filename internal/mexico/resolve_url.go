package mexico

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the SIAP portal URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	switch entry.DatasetID.String() {
	case "mexico.siap-produccion-agricola":
		if u := strings.TrimSpace(entry.SourceURL); u != "" {
			return u, nil
		}
		return "https://www.gob.mx/siap", nil
	default:
		return "", fmt.Errorf("unsupported mexico dataset %s", entry.DatasetID)
	}
}
