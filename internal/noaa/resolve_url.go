package noaa

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns the NOAA source URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	if url := strings.TrimSpace(entry.NOAAIndexURL); url != "" {
		return url, nil
	}
	switch entry.DatasetID.String() {
	case "noaa.enso-indices":
		return defaultONIURL, nil
	case "noaa.global-temp-anomaly":
		start := entry.PeriodStart
		if start == 0 {
			start = 2010
		}
		end := entry.PeriodEnd
		if end == 0 {
			end = start
		}
		return fmt.Sprintf(defaultGlobalTempURL, start, end), nil
	case "noaa.gpcc-precipitation":
		raw := strings.TrimSpace(entry.SourceURL)
		if raw == "" {
			return "https://opendata.dwd.de/climate_environment/GPCC/Monitoring/", nil
		}
		return raw, nil
	default:
		return "", fmt.Errorf("unsupported noaa dataset %s", entry.DatasetID)
	}
}
