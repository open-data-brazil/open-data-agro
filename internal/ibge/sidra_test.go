package ibge

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestBuildSIDRAURL(t *testing.T) {
	t.Parallel()

	got := buildSIDRAURL("1612", []string{"51", "43"}, 2015, "109,216,214", "81", "2713")
	want := "https://apisidra.ibge.gov.br/values/t/1612/n6/in%20n3%2051,43/p/2015/v/109,216,214/c81/2713"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestResolveCropCodesSingle(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.pam-area-quantidade"),
		SidraCrops: map[string]int{
			"soja":  2713,
			"milho": 2711,
		},
	}

	crops, err := resolveCropCodes(entry, "soja")
	if err != nil {
		t.Fatalf("resolveCropCodes: %v", err)
	}
	if len(crops) != 1 || crops["soja"] != "2713" {
		t.Fatalf("crops: %#v", crops)
	}
}

func TestResolveCropCodesAll(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.pam-area-quantidade"),
		SidraCrops: map[string]int{
			"soja":  2713,
			"milho": 2711,
		},
	}

	crops, err := resolveCropCodes(entry, "all")
	if err != nil {
		t.Fatalf("resolveCropCodes: %v", err)
	}
	if len(crops) != 2 {
		t.Fatalf("crops: %#v", crops)
	}
}

func TestResolveYearRangeDefaults(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		PeriodStart: 2010,
		PeriodEnd:   2020,
	}

	from, to, err := resolveYearRange(entry, 0, 0)
	if err != nil {
		t.Fatalf("resolveYearRange: %v", err)
	}
	if from != 2010 || to != 2020 {
		t.Fatalf("range: %d-%d", from, to)
	}
}

func TestResolvePAMURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.pam-area-quantidade"),
		SourceURL: "https://apisidra.ibge.gov.br/values/t/1612",
	}
	url, err := ResolvePAMURL(entry)
	if err != nil {
		t.Fatalf("ResolvePAMURL: %v", err)
	}
	if url == "" {
		t.Fatal("expected non-empty url")
	}
}

func TestChunkUFsRequested(t *testing.T) {
	t.Parallel()

	chunks := chunkUFs([]string{"51", "43"})
	if len(chunks) != 1 || len(chunks[0]) != 2 {
		t.Fatalf("chunks: %#v", chunks)
	}
}
