package inpe

import (
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BuildProbeURL returns a WFS GetCapabilities URL for health checks.
func BuildProbeURL(entry catalog.RegistryEntry) (string, error) {
	raw, err := ResolveURL(entry)
	if err != nil {
		return "", err
	}
	lower := strings.ToLower(raw)
	if strings.Contains(lower, "/wfs") && !strings.Contains(lower, "request=") {
		sep := "?"
		if strings.Contains(raw, "?") {
			sep = "&"
		}
		return raw + sep + "service=WFS&request=GetCapabilities", nil
	}
	return raw, nil
}
