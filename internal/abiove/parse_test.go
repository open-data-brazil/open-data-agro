package abiove

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenBalancoComplexoSoja(t *testing.T) {
	t.Parallel()
	headers, rows := flattenFixture(t, "balanco_complexo_soja.sample.xlsx", "abiove.balanco-complexo-soja", "Rel_Exp2026")
	if len(rows) < 12 {
		t.Fatalf("rows: got %d want >= 12", len(rows))
	}
	if headers[0] != "section" {
		t.Fatalf("headers: %v", headers)
	}
}

func TestFlattenExportacoesComplexoSoja(t *testing.T) {
	t.Parallel()
	_, rows := flattenFixture(t, "exportacoes_complexo_soja.sample.xlsx", "abiove.exportacoes-complexo-soja", "materia-prima_anual")
	if len(rows) < 5 {
		t.Fatalf("rows: got %d want >= 5", len(rows))
	}
}

func TestFlattenCapacidadeEsmagamento(t *testing.T) {
	t.Parallel()
	_, rows := flattenFixture(t, "capacidade_esmagamento.sample.xlsx", "abiove.capacidade-instalada-esmagamento", "projecoes_mensais")
	if len(rows) < 6 {
		t.Fatalf("rows: got %d want >= 6", len(rows))
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("abiove.balanco-complexo-soja"),
		SourceURL: "https://abiove.org.br/abiove_content/Abiove/exp_202605.xlsx",
	}
	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" {
		t.Fatal("expected url")
	}
}

func flattenFixture(t *testing.T, fileName, datasetID, sheet string) ([]string, [][]string) {
	t.Helper()
	raw := readTestdata(t, fileName)
	book, err := OpenWorkbook(raw)
	if err != nil {
		t.Fatalf("OpenWorkbook: %v", err)
	}
	defer func() { _ = book.Close() }()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID(datasetID),
		XLSXSheet: sheet,
	}
	headers, rows, err := ConvertWorkbook(entry, book)
	if err != nil {
		t.Fatalf("ConvertWorkbook: %v", err)
	}
	return headers, rows
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
