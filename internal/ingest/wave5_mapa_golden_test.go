package ingest

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/mapa"
)

func TestMAPASIPEAGROEstabelecimentosGoldenVector(t *testing.T) {
	t.Parallel()

	samples := map[string][]byte{
		"Fertilizantes":                    readMAPATestdata(t, "sipeagro_fertilizantes.sample.csv"),
		"Qualidade Vegetal":                readMAPATestdata(t, "sipeagro_qualidade_vegetal.sample.csv"),
		"Produto Veterinário":              readMAPATestdata(t, "sipeagro_fertilizantes.sample.csv"),
		"Vinhos e Bebidas":                 readMAPATestdata(t, "sipeagro_qualidade_vegetal.sample.csv"),
		"Alimentação Animal":               readMAPATestdata(t, "sipeagro_fertilizantes.sample.csv"),
		"Material de Multiplicação Animal": readMAPATestdata(t, "sipeagro_fertilizantes.sample.csv"),
		"Aves de Reprodução":               readMAPATestdata(t, "sipeagro_fertilizantes.sample.csv"),
		"Aviação Agrícola - Registro":     readMAPATestdata(t, "sipeagro_fertilizantes.sample.csv"),
	}
	raw, err := mapa.MergeSIPEAGROSampleCSV("mapa.sipeagro-estabelecimentos", samples)
	if err != nil {
		t.Fatalf("MergeSIPEAGROSampleCSV: %v", err)
	}
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.sipeagro-estabelecimentos"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 24 {
		t.Fatalf("rowCount: got %d want 24", rowCount)
	}
}

func TestMAPASIPEAGROProdutosGoldenVector(t *testing.T) {
	t.Parallel()

	samples := map[string][]byte{
		"Fertilizantes":       readMAPATestdata(t, "sipeagro_fertilizantes.sample.csv"),
		"Qualidade Vegetal":   readMAPATestdata(t, "sipeagro_qualidade_vegetal.sample.csv"),
		"Produto Veterinário": readMAPATestdata(t, "sipeagro_fertilizantes.sample.csv"),
		"Vinhos e Bebidas":    readMAPATestdata(t, "sipeagro_qualidade_vegetal.sample.csv"),
		"Alimentação Animal":  readMAPATestdata(t, "sipeagro_fertilizantes.sample.csv"),
	}
	raw, err := mapa.MergeSIPEAGROSampleCSV("mapa.sipeagro-produtos", samples)
	if err != nil {
		t.Fatalf("MergeSIPEAGROSampleCSV: %v", err)
	}
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.sipeagro-produtos"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 15 {
		t.Fatalf("rowCount: got %d want 15", rowCount)
	}
}

func TestMAPASIGEFProducaoSementesGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readMAPATestdata(t, "sigef_producao_sementes.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.sigef-producao-sementes"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func TestMAPASIGEFAreasGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readMAPATestdata(t, "sigef_areas.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.sigef-areas"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func TestMAPASISSERSeguroRuralGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readMAPATestdata(t, "sisser_psr_2025.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.sisser-seguro-rural"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}
