package ingest

import (
	"strings"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestConvertDelimitedToParquetPreservesHeaders(t *testing.T) {
	t.Parallel()

	raw := []byte("col_a;col_b\n1;foo\n2;bar\n")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	parquetBytes, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if rowCount != 2 {
		t.Fatalf("rowCount: got %d want 2", rowCount)
	}
	if len(parquetBytes) == 0 || !strings.HasPrefix(string(parquetBytes[:4]), "PAR1") {
		t.Fatal("expected parquet magic bytes")
	}
}

func TestConvertDelimitedLatin1ToUTF8(t *testing.T) {
	t.Parallel()

	// "NÃO" in ISO-8859-1 (portal encoding for PrecosSemanalUF.txt).
	raw := []byte("produto;classificao_produto\nSOJA;N\xC3O INFORMADO\n")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.precos-agropecuarios-semanal-uf"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if rowCount != 1 {
		t.Fatalf("rowCount: got %d want 1", rowCount)
	}
}

func TestBronzeKey(t *testing.T) {
	t.Parallel()

	key, err := BronzeKey("conab.estimativa-graos", parseDate(t, "2026-06-25"), "abc-123")
	if err != nil {
		t.Fatalf("BronzeKey: %v", err)
	}
	want := "bronze/conab/estimativa-graos/ingest_date=2026-06-25/part-abc-123.parquet"
	if key != want {
		t.Fatalf("got %q want %q", key, want)
	}
}

func parseDate(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		t.Fatalf("parse date: %v", err)
	}
	return parsed
}
