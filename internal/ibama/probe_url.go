package ibama

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BuildProbeURL returns a minimal URL aligned with ingest fetch semantics.
func BuildProbeURL(entry catalog.RegistryEntry) (string, error) {
	id := entry.DatasetID.String()
	switch id {
	case "ibama.licencas-ambientais":
		return fmt.Sprintf("%s?id=%s", defaultCKANPackageShowURL, licencasPackageID), nil
	case "ibama.autos-infracao":
		return fmt.Sprintf("%s?id=%s", defaultCKANPackageShowURL, autosPackageID), nil
	}
	return ResolveURL(entry)
}

// ProbeHeaders returns HTTP headers for IBAMA probes.
func ProbeHeaders(entry catalog.RegistryEntry) map[string]string {
	headers := map[string]string{"Accept": "*/*"}
	if entry.DatasetID.String() == "ibama.licencas-ambientais" {
		headers["Accept"] = "application/json"
	}
	if entry.DatasetID.String() == "ibama.autos-infracao" {
		headers["Accept"] = "application/json"
	}
	if strings.HasPrefix(entry.DatasetID.String(), "ibama.") {
		return headers
	}
	return headers
}
