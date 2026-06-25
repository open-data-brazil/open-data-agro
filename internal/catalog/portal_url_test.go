package catalog

import "testing"

func TestSourcePortalURLCONAB(t *testing.T) {
	t.Parallel()

	id := MustParseDatasetID("conab.estimativa-graos")
	if got := SourcePortalURL(id); got != CONABSourcePortalURL {
		t.Fatalf("got %q want %q", got, CONABSourcePortalURL)
	}
}

func TestSourcePortalURLIBGE(t *testing.T) {
	t.Parallel()

	id := MustParseDatasetID("ibge.localidades-municipios")
	if got := SourcePortalURL(id); got != IBGELocalidadesDocsURL {
		t.Fatalf("got %q want %q", got, IBGELocalidadesDocsURL)
	}
}

func TestSourcePortalURLUnknownAgency(t *testing.T) {
	t.Parallel()

	id := MustParseDatasetID("foo.bar-baz")
	if got := SourcePortalURL(id); got != "" {
		t.Fatalf("got %q want empty", got)
	}
}
