package ingest

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/conab"
)

func TestCitationFor(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:   catalog.MustParseDatasetID("conab.estimativa-graos"),
		PortalLabel: "Estimativa Grãos",
	}
	citation := CitationFor(entry)
	if citation.FonteOficial != conab.PortalDownloadPage {
		t.Fatalf("fonteOficial: got %q", citation.FonteOficial)
	}
	if citation.Agencia != "CONAB" {
		t.Fatalf("agencia: got %q", citation.Agencia)
	}
	if citation.DatasetPortal != "Estimativa Grãos" {
		t.Fatalf("datasetPortal: got %q", citation.DatasetPortal)
	}
}
