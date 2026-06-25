package inmet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenEstacoesAutomaticas(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "estacoes_automaticas.sample.csv")
	headers, rows, err := FlattenEstacoes("inmet.estacoes-automaticas", raw)
	if err != nil {
		t.Fatalf("FlattenEstacoes: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows: got %d want 2", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["cd_estacao"]]; got != "A901" {
		t.Fatalf("cd_estacao: got %q want A901", got)
	}
	if got := rows[0][idx["uf"]]; got != "MT" {
		t.Fatalf("uf: got %q want MT", got)
	}
}

func TestParseAnnualZIPToDailyLong(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "bdmep_2023_mt.sample.zip")
	rows, err := parseAnnualZIPToDailyLong(raw, 2023, map[string]struct{}{"MT": {}}, nil)
	if err != nil {
		t.Fatalf("parseAnnualZIPToDailyLong: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected daily rows")
	}
}

func TestIsMissingValue(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"9999", "-9999", "Null", "", "//"} {
		if !IsMissingValue(value) {
			t.Fatalf("expected missing for %q", value)
		}
	}
	if IsMissingValue("12.5") {
		t.Fatal("expected numeric value to be present")
	}
}

func TestResolveAnnualZIPURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("inmet.bdmep-diario"),
		SourceURL: "https://portal.inmet.gov.br/uploads/dadoshistoricos/{year}.zip",
	}
	url, err := ResolveAnnualZIPURL(entry, 2023)
	if err != nil {
		t.Fatalf("ResolveAnnualZIPURL: %v", err)
	}
	if url == "" || !contains(url, "2023.zip") {
		t.Fatalf("url: %q", url)
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
	return len(value) >= len(part) && (value == part || len(part) == 0 || indexOfString(value, part) >= 0)
}

func indexOfString(value, part string) int {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}
