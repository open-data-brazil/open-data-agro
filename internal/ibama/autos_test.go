package ibama

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractAutosYearFromZIP(t *testing.T) {
	t.Parallel()

	zipBytes := buildTestAutosZIP(t)
	body, member, err := extractAutosYearFromZIP(zipBytes, 1977)
	if err != nil {
		t.Fatalf("extractAutosYearFromZIP: %v", err)
	}
	if member != "auto_infracao_1977.csv" {
		t.Fatalf("member: got %s", member)
	}
	if len(body) == 0 {
		t.Fatal("empty body")
	}
}

func TestNormalizeCSVStripsBOM(t *testing.T) {
	t.Parallel()

	raw := []byte{0xEF, 0xBB, 0xBF, 'a', ';', 'b'}
	out := NormalizeCSV(raw)
	if out[0] == 0xEF {
		t.Fatal("BOM not stripped")
	}
}

func buildTestAutosZIP(t *testing.T) []byte {
	t.Helper()
	samplePath := filepath.Join("testdata", "auto_infracao_1977.sample.csv")
	sample, err := os.ReadFile(samplePath)
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)
	file, err := writer.Create("auto_infracao_1977.csv")
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := file.Write(sample); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buf.Bytes()
}
