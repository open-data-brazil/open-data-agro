package ibge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenLSPAGolden(t *testing.T) {
	t.Parallel()

	raw := readLSPATestdata(t, "lspa_area_producao.sample.json")
	headers, rows, err := FlattenLSPA("ibge.lspa-area-producao", raw)
	if err != nil {
		t.Fatalf("FlattenLSPA: %v", err)
	}
	if len(headers) != 12 {
		t.Fatalf("headers: got %d want 12", len(headers))
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}
	if rows[0][8] != "soja" {
		t.Fatalf("produto_slug: got %q", rows[0][8])
	}
	if rows[0][0] != "6588" {
		t.Fatalf("sidra_tabela: got %q", rows[0][0])
	}
}

func TestBuildLSPAURL(t *testing.T) {
	t.Parallel()

	url := buildLSPAURL("6588", []string{"43"}, "202401,202402", "109,216,35", "48", "39443")
	if !strings.Contains(url, "/t/6588/n3/in%2043/p/202401,202402") {
		t.Fatalf("url missing n3 period chunk: %q", url)
	}
	if !strings.Contains(url, "/c48/39443") {
		t.Fatalf("url missing classification: %q", url)
	}
}

func TestResolveLSPAURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:  catalog.MustParseDatasetID("ibge.lspa-area-producao"),
		SourceURL:  "https://apisidra.ibge.gov.br/values/t/6588",
	}
	url, err := ResolveLSPAURL(entry)
	if err != nil {
		t.Fatalf("ResolveLSPAURL: %v", err)
	}
	if !strings.Contains(url, "apisidra.ibge.gov.br") {
		t.Fatalf("url: got %q", url)
	}
}

func readLSPATestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
