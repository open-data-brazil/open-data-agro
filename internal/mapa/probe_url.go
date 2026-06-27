package mapa

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BuildProbeURL returns a minimal URL aligned with ingest fetch semantics.
func BuildProbeURL(entry catalog.RegistryEntry) (string, error) {
	id := entry.DatasetID.String()
	if strings.HasPrefix(id, "mapa.sipeagro-") {
		return fmt.Sprintf("%s?id=%s", defaultCKANPackageShowURL, sipeagroPackageID), nil
	}
	if id == "mapa.sisser-seguro-rural" {
		return fmt.Sprintf("%s?id=%s", defaultCKANPackageShowURL, sisserPackageID), nil
	}
	if strings.HasPrefix(id, "mapa.sigef-") {
		packageID := strings.TrimSpace(entry.CKANPackageID)
		if packageID == "" {
			packageID = "dados-referentes-ao-controle-da-producao-de-sementes-sigef"
		}
		return fmt.Sprintf("%s?id=%s", defaultCKANPackageShowURL, packageID), nil
	}
	return ResolveURL(entry)
}
