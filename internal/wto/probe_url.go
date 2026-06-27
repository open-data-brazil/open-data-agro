package wto

import (
	"fmt"
	"os"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BuildProbeRequest returns a minimal WTO Stats API URL and headers.
func BuildProbeRequest(entry catalog.RegistryEntry) (string, map[string]string, error) {
	headers := map[string]string{"Accept": "application/json"}
	apiKey := strings.TrimSpace(os.Getenv("WTO_API_KEY"))
	if apiKey == "" {
		return "", headers, fmt.Errorf("WTO_API_KEY not set")
	}
	headers["Ocp-Apim-Subscription-Key"] = apiKey

	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = fmt.Sprintf("%s?i=HS_P_0070&r=840&p=000&ps=2023", defaultWTOAPIURL)
	} else if !strings.Contains(sourceURL, "?") {
		sourceURL = sourceURL + "?i=HS_P_0070&r=840&p=000&ps=2023"
	}

	return sourceURL, headers, nil
}
