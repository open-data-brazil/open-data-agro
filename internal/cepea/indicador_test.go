package cepea

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestParseIndicatorHTMLParanagua(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "soja_paranagua.sample.html")
	rows, err := ParseIndicatorHTML(raw, "Paranaguá")
	if err != nil {
		t.Fatalf("ParseIndicatorHTML: %v", err)
	}
	if len(rows) < 5 {
		t.Fatalf("rows: got %d want >= 5", len(rows))
	}
	if rows[0].Data == "" || rows[0].PrecoRsSc == "" {
		t.Fatalf("first row missing fields: %+v", rows[0])
	}
}

func TestFlattenIndicadorHistorico(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "soja_paranagua_historico.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("cepea.soja-paranagua"),
	}
	headers, rows, err := FlattenIndicador(entry, raw)
	if err != nil {
		t.Fatalf("FlattenIndicador: %v", err)
	}
	if len(rows) != 5 {
		t.Fatalf("rows: got %d want 5", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["data"]]; got != "2010-01-04" {
		t.Fatalf("data: got %q", got)
	}
	if got := rows[0][idx["preco_rs_sc"]]; got != "52.30" {
		t.Fatalf("preco_rs_sc: got %q", got)
	}
	if got := rows[0][idx["praca"]]; got != "Paranaguá" {
		t.Fatalf("praca: got %q", got)
	}
}

func TestFilterFromDate2010(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "soja_paranagua_historico.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("cepea.soja-paranagua"),
		StartDate: "2010-01-01",
	}
	start, err := resolveFromDate(entry, "2010-01-01")
	if err != nil {
		t.Fatal(err)
	}

	var rows []Observation
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	filtered := filterFromDate(rows, start)
	if len(filtered) != 5 {
		t.Fatalf("filtered: got %d want 5", len(filtered))
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:       catalog.MustParseDatasetID("cepea.soja-paranagua"),
		SourcePortalURL: "https://www.cepea.org.br/br/indicador/soja.aspx",
	}
	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" {
		t.Fatal("expected non-empty url")
	}
}

func TestMirrorURL(t *testing.T) {
	t.Parallel()

	url, err := MirrorURL("cepea.soja-paranagua")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(url, "noticiasagricolas.com.br") {
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
