package catalog

import "testing"

func TestSourcePortalURLCONAB(t *testing.T) {
	t.Parallel()

	id := MustParseDatasetID("conab.estimativa-graos")
	if got := SourcePortalURL(id); got != CONABSourcePortalURL {
		t.Fatalf("got %q want %q", got, CONABSourcePortalURL)
	}
}

func TestSourcePortalURLBCB(t *testing.T) {
	t.Parallel()

	id := MustParseDatasetID("bcb.sgs-ipca")
	if got := SourcePortalURL(id); got != BCBDadosAbertosURL {
		t.Fatalf("got %q want %q", got, BCBDadosAbertosURL)
	}
}

func TestSourcePortalURLIBGEPAM(t *testing.T) {
	t.Parallel()

	id := MustParseDatasetID("ibge.pam-area-quantidade")
	if got := SourcePortalURL(id); got != IBGESIDRAPAMURL {
		t.Fatalf("got %q want %q", got, IBGESIDRAPAMURL)
	}
}

func TestSourcePortalURLINMET(t *testing.T) {
	t.Parallel()

	id := MustParseDatasetID("inmet.bdmep-diario")
	if got := SourcePortalURL(id); got != INMETBDMEPPortalURL {
		t.Fatalf("got %q want %q", got, INMETBDMEPPortalURL)
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
