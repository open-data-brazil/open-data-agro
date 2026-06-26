package ibge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFlattenPEVSGolden(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "pevs_producao_vegetal.sample.json")
	headers, rows, err := FlattenPEVS("ibge.pevs-producao-vegetal", raw)
	if err != nil {
		t.Fatalf("FlattenPEVS: %v", err)
	}
	if len(rows) != 14 {
		t.Fatalf("rows: got %d want 14", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["sidra_tabela"]]; got != "289" {
		t.Fatalf("sidra_tabela: got %q want 289", got)
	}
	if got := rows[0][idx["codigo_uf"]]; got == "" {
		t.Fatalf("codigo_uf should not be empty")
	}
}

func TestBuildPEVSURL(t *testing.T) {
	t.Parallel()

	got := buildPEVSURL("289", "2022,2023", "144,145")
	want := "https://apisidra.ibge.gov.br/values/t/289/n3/all/p/2022,2023/v/144,145"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func readIBGETestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func TestFlattenPEVSParsesSampleJSON(t *testing.T) {
	t.Parallel()

	var rows []map[string]any
	raw := readIBGETestdata(t, "pevs_producao_vegetal.sample.json")
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatalf("unmarshal sample: %v", err)
	}
	if len(rows) < 2 {
		t.Fatalf("sample too small")
	}
}
