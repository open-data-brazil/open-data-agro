package ingest

import (
	"strings"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestBuildProbeSpecBCB(t *testing.T) {
	t.Parallel()

	spec, err := BuildProbeSpec(catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("bcb.sgs-ptax-usd-venda"),
		SourceURL: "https://api.bcb.gov.br/dados/serie/bcdata.sgs.1/dados",
		SGSCode:   1,
	})
	if err != nil {
		t.Fatalf("BuildProbeSpec: %v", err)
	}
	if !strings.Contains(spec.URL, "formato=json") {
		t.Fatalf("url missing formato=json: %q", spec.URL)
	}
	if !strings.Contains(spec.URL, "dataInicial=") {
		t.Fatalf("url missing date window: %q", spec.URL)
	}
}

func TestBuildProbeSpecIBGESIDRA(t *testing.T) {
	t.Parallel()

	spec, err := BuildProbeSpec(catalog.RegistryEntry{
		DatasetID:             catalog.MustParseDatasetID("ibge.pam-area-quantidade"),
		SourceURL:             "https://apisidra.ibge.gov.br/values/t/1612",
		SidraTable:            "1612",
		SidraVariables:        []int{109},
		SidraClassification:   "c81",
		SidraCrops:          map[string]int{"soja": 2713},
		PeriodEnd:             2024,
	})
	if err != nil {
		t.Fatalf("BuildProbeSpec: %v", err)
	}
	if spec.URL == "https://apisidra.ibge.gov.br/values/t/1612" {
		t.Fatalf("expected full SIDRA path, got bare table url")
	}
	if !strings.Contains(spec.URL, "/p/2024/") {
		t.Fatalf("url missing year: %q", spec.URL)
	}
}

func TestBuildProbeSpecMDICPOST(t *testing.T) {
	t.Parallel()

	spec, err := BuildProbeSpec(catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mdic.comex-exportacao-ncm-mes"),
		SourceURL: "https://api-comexstat.mdic.gov.br/general",
		ComexFlow: "export",
		ComexNCMs: []string{"12019000"},
	})
	if err != nil {
		t.Fatalf("BuildProbeSpec: %v", err)
	}
	if spec.Method != "POST" {
		t.Fatalf("method: got %q want POST", spec.Method)
	}
	if spec.Body == "" {
		t.Fatal("expected POST body")
	}
}

func TestBuildProbeSpecDirectDownload(t *testing.T) {
	t.Parallel()

	direct := "https://dados.agricultura.gov.br/dataset/x/resource/y/download/agrofit.csv"
	spec, err := BuildProbeSpec(catalog.RegistryEntry{
		DatasetID:     catalog.MustParseDatasetID("mapa.agrofit-produtos-tecnicos"),
		SourceURL:     direct,
		CKANPackageID: "ignored-for-probe",
	})
	if err != nil {
		t.Fatalf("BuildProbeSpec: %v", err)
	}
	if spec.URL != direct {
		t.Fatalf("url: got %q want direct download", spec.URL)
	}
}

func TestBuildProbeSpecINMETYear(t *testing.T) {
	t.Parallel()

	spec, err := BuildProbeSpec(catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("inmet.bdmep-diario"),
		SourceURL: "https://portal.inmet.gov.br/uploads/dadoshistoricos/{year}.zip",
	})
	if err != nil {
		t.Fatalf("BuildProbeSpec: %v", err)
	}
	if strings.Contains(spec.URL, "{year}") {
		t.Fatalf("year placeholder not substituted: %q", spec.URL)
	}
	if !strings.HasSuffix(spec.URL, ".zip") {
		t.Fatalf("expected zip url: %q", spec.URL)
	}
}
