package embrapa

import (
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BuildProbeURL returns a probe URL for Embrapa catalog entries.
func BuildProbeURL(entry catalog.RegistryEntry) (string, error) {
	if entry.DatasetID.String() == "embrapa.agroapi-agrofit" {
		return "https://www.agroapi.cnptia.embrapa.br/store/apis/info?name=AGROFIT&provider=agroapi&version=v1", nil
	}
	url, err := ResolveURL(entry)
	if err != nil {
		return "", err
	}
	return url, nil
}

// ProbeHeaders returns HTTP headers for Embrapa probes.
func ProbeHeaders(entry catalog.RegistryEntry) map[string]string {
	_ = entry
	return map[string]string{
		"Accept": "text/html,application/xhtml+xml,application/json,*/*",
	}
}
