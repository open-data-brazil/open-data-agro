package ingest

import (
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestAbioveBalancoComplexoSojaGoldenVector(t *testing.T) {
	t.Parallel()
	runAbioveGolden(t, "balanco_complexo_soja.sample.xlsx", "abiove.balanco-complexo-soja", "Rel_Exp2026", 27)
}

func TestAbioveExportacoesComplexoSojaGoldenVector(t *testing.T) {
	t.Parallel()
	runAbioveGolden(t, "exportacoes_complexo_soja.sample.xlsx", "abiove.exportacoes-complexo-soja", "materia-prima_anual", 50)
}

func TestAbioveCapacidadeEsmagamentoGoldenVector(t *testing.T) {
	t.Parallel()
	runAbioveGolden(t, "capacidade_esmagamento.sample.xlsx", "abiove.capacidade-instalada-esmagamento", "projecoes_mensais", 6)
}

func TestB3FuturoCafeGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readB3IngestTestdata(t, "futuro_cafe.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("b3.futuro-cafe"),
		Format:    catalog.FormatJSON,
	}
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func TestB3FuturoAcucarGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readB3IngestTestdata(t, "futuro_acucar.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("b3.futuro-acucar"),
		Format:    catalog.FormatJSON,
	}
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func runAbioveGolden(t *testing.T, fileName, datasetID, sheet string, minRows int) {
	t.Helper()
	path := filepath.Join("..", "abiove", "testdata", fileName)
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID(datasetID),
		Format:    catalog.FormatXLSX,
		XLSXSheet: sheet,
	}
	_, rowCount, err := ConvertToParquetFromFile(entry, path)
	if err != nil {
		t.Fatalf("ConvertToParquetFromFile: %v", err)
	}
	if rowCount < minRows {
		t.Fatalf("rowCount: got %d want >= %d", rowCount, minRows)
	}
}
