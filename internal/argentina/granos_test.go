package argentina

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenGranosGoldenVector(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "magyp_producion_granos.sample.json"))
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("argentina.magyp-producion-granos"),
	}
	headers, rows, err := FlattenGranos(entry, raw)
	if err != nil {
		t.Fatalf("FlattenGranos: %v", err)
	}
	if len(headers) != 6 || len(rows) != 4 {
		t.Fatalf("headers=%d rows=%d", len(headers), len(rows))
	}
}

func TestResolveGranosURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:          catalog.MustParseDatasetID("argentina.magyp-producion-granos"),
		ArgentinaSeriesIDs: []string{"AGRO_A_Soja_0003"},
		PeriodStart:        2010,
	}
	url, err := ResolveGranosURL(entry)
	if err != nil {
		t.Fatalf("ResolveGranosURL: %v", err)
	}
	if url == "" {
		t.Fatal("empty url")
	}
}
