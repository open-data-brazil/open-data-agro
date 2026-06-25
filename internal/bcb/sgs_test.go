package bcb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenSGSIPCA(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "sgs_ipca.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("bcb.sgs-ipca"),
		SGSCode:   433,
	}
	headers, rows, err := FlattenSGS(entry, raw)
	if err != nil {
		t.Fatalf("FlattenSGS: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["sgs_codigo"]]; got != "433" {
		t.Fatalf("sgs_codigo: got %q want 433", got)
	}
	if got := rows[0][idx["data"]]; got != "2020-01-01" {
		t.Fatalf("data: got %q want 2020-01-01", got)
	}
	if got := rows[0][idx["valor"]]; got != "0.21" {
		t.Fatalf("valor: got %q want 0.21", got)
	}
}

func TestChunkDateRangeTenYears(t *testing.T) {
	t.Parallel()

	start, err := parseBCBDate("01/01/1995")
	if err != nil {
		t.Fatal(err)
	}
	end, err := parseBCBDate("31/12/2024")
	if err != nil {
		t.Fatal(err)
	}
	chunks := chunkDateRange(start, end, 10)
	if len(chunks) < 3 {
		t.Fatalf("chunks: got %d want >= 3", len(chunks))
	}
	if chunks[0].from != "01/01/1995" {
		t.Fatalf("first chunk from: %q", chunks[0].from)
	}
}

func TestBuildSGSURL(t *testing.T) {
	t.Parallel()

	got := buildSGSURL(433, "01/01/2020", "31/12/2020")
	if got == "" || !contains(got, "433") || !contains(got, "dataInicial") {
		t.Fatalf("url: %q", got)
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("bcb.sgs-ipca"),
		SourceURL: "https://api.bcb.gov.br/dados/serie/bcdata.sgs.433/dados",
	}
	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" {
		t.Fatal("expected non-empty url")
	}
}

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func contains(value, part string) bool {
	return len(value) >= len(part) && indexOf(value, part) >= 0
}

func indexOf(value, part string) int {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}
