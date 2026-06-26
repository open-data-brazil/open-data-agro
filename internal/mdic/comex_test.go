package mdic

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenComexExport(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "comex_export.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mdic.comex-exportacao-ncm-mes"),
		ComexFlow: "export",
		ComexNCMs: []string{"12019000", "10059000"},
	}

	headers, rows, err := FlattenComex(entry, raw)
	if err != nil {
		t.Fatalf("FlattenComex: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["co_ncm"]]; got != "12019000" {
		t.Fatalf("co_ncm: got %q want 12019000", got)
	}
	if got := rows[0][idx["data"]]; got != "2024-01-01" {
		t.Fatalf("data: got %q want 2024-01-01", got)
	}
	if got := rows[0][idx["produto_slug"]]; got != "soja" {
		t.Fatalf("produto_slug: got %q want soja", got)
	}
	if got := rows[0][idx["valor_fob_usd"]]; got != "1454912473" {
		t.Fatalf("valor_fob_usd: got %q", got)
	}
}

func TestMonthToDate(t *testing.T) {
	t.Parallel()

	got, err := monthToDate("2024", "3")
	if err != nil {
		t.Fatal(err)
	}
	if got != "2024-03-01" {
		t.Fatalf("got %q want 2024-03-01", got)
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mdic.comex-exportacao-ncm-mes"),
		SourceURL: "https://api-comexstat.mdic.gov.br/general",
	}
	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" {
		t.Fatal("expected non-empty url")
	}
}

func TestResolveComexRangeFromDate(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:   catalog.MustParseDatasetID("mdic.comex-exportacao-ncm-mes"),
		PeriodStart: 2015,
	}
	start, end, err := resolveComexRange(entry, "2020-06-15")
	if err != nil {
		t.Fatalf("resolveComexRange: %v", err)
	}
	if start.Format("2006-01-02") != "2020-06-15" {
		t.Fatalf("start: got %s", start.Format("2006-01-02"))
	}
	if end.Before(start) {
		t.Fatal("end before start")
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
