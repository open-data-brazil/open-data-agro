package ibama

import (
	"context"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const (
	sisfogoDirectURL        = "https://dadosabertos.ibama.gov.br/dados/SISFOGO/ROI.csv"
	licencasDirectURL       = "https://stibamadadosabertosprd.blob.core.windows.net/dados-abertos/dados/SISLIC/sislic-licencas.csv"
	autosPackageID          = "fiscalizacao-auto-de-infracao"
	autosResourceName       = "Autos de infração"
	licencasPackageID       = "licencas-ambientais-de-atividades-e-empreendimentos-licenciados-pelo-ibama"
)

// ResolveURL returns the download URL for an IBAMA catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	id := entry.DatasetID.String()
	switch id {
	case "ibama.sisfogo-incendios":
		if raw := strings.TrimSpace(entry.SourceURL); raw != "" {
			return raw, nil
		}
		return sisfogoDirectURL, nil
	case "ibama.licencas-ambientais":
		if raw := strings.TrimSpace(entry.SourceURL); raw != "" && strings.Contains(strings.ToLower(raw), "blob.core.windows.net") {
			return raw, nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		url, err := resolveCKANResourceURL(ctx, licencasPackageID, "sislic-licencas", "CSV")
		if err == nil {
			return url, nil
		}
		return licencasDirectURL, nil
	case "ibama.autos-infracao":
		if raw := strings.TrimSpace(entry.SourceURL); raw != "" {
			return raw, nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		return resolveCKANResourceURL(ctx, autosPackageID, autosResourceName, "CSV")
	default:
		return resolveDirectURL(entry)
	}
}
