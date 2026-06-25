package catalog

import "testing"

func TestSourcePortalURLCONAB(t *testing.T) {
	t.Parallel()

	id := MustParseDatasetID("conab.estimativa-graos")
	if got := SourcePortalURL(id); got != CONABSourcePortalURL {
		t.Fatalf("got %q want %q", got, CONABSourcePortalURL)
	}
}

func TestSourcePortalURLNonCONAB(t *testing.T) {
	t.Parallel()

	id := MustParseDatasetID("ibge.foo-bar")
	if got := SourcePortalURL(id); got != "" {
		t.Fatalf("got %q want empty", got)
	}
}
