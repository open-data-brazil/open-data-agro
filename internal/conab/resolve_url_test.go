package conab_test

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/conab"
)

func TestResolveURLValidatesPortalHost(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		SourceURL: "https://portaldeinformacoes.conab.gov.br/downloads/arquivos/LevantamentoGraos.txt",
	}
	url, err := conab.ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url != entry.SourceURL {
		t.Fatalf("got %q want %q", url, entry.SourceURL)
	}
}

func TestResolveURLRejectsExternalHost(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		SourceURL: "https://example.com/file.txt",
	}
	if _, err := conab.ResolveURL(entry); err == nil {
		t.Fatal("expected error for external host")
	}
}
