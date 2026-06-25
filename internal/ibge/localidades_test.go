package ibge

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenMunicipio3550308(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "municipios.sample.json")
	headers, rows, err := FlattenLocalidades("ibge.localidades-municipios", raw)
	if err != nil {
		t.Fatalf("FlattenLocalidades: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows: got %d want 1", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}

	row := rows[0]
	if got := row[idx["codigo_ibge"]]; got != "3550308" {
		t.Fatalf("codigo_ibge: got %q want 3550308", got)
	}
	if got := row[idx["sigla_uf"]]; got != "SP" {
		t.Fatalf("sigla_uf: got %q want SP", got)
	}
	if got := row[idx["nome_regiao"]]; got != "Sudeste" {
		t.Fatalf("nome_regiao: got %q want Sudeste", got)
	}
}

func TestFlattenUFsSnapshot(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "ufs.sample.json")
	_, rows, err := FlattenLocalidades("ibge.localidades-ufs", raw)
	if err != nil {
		t.Fatalf("FlattenLocalidades: %v", err)
	}
	if len(rows) != 27 {
		t.Fatalf("rows: got %d want 27", len(rows))
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.localidades-municipios"),
		SourceURL: "https://servicodados.ibge.gov.br/api/v1/localidades/municipios?orderBy=nome",
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
