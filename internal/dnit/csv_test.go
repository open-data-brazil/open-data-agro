package dnit

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripMetadataRows(t *testing.T) {
	t.Parallel()

	raw := readDNITTestdata(t, "snv_rodovias_federais.sample.csv")
	stripped := StripMetadataRows(raw)
	first := firstNonEmptyLine(stripped)
	if !strings.HasPrefix(strings.TrimSpace(strings.Split(first, ";")[0]), "BR") {
		t.Fatalf("expected header row starting with BR after strip, got %q", first)
	}
	if bytes.Contains(stripped, []byte("CONTATO:")) {
		t.Fatalf("metadata rows should be removed")
	}
}

func TestStripMetadataRowsAlreadyClean(t *testing.T) {
	t.Parallel()

	raw := readDNITTestdata(t, "jurisdicao_vias_trimmed.sample.csv")
	first := firstNonEmptyLine(raw)
	if !strings.EqualFold(strings.TrimSpace(strings.Split(first, ";")[0]), "BR") {
		t.Skip("trimmed fixture missing BR header")
	}
	stripped := StripMetadataRows(raw)
	if !bytes.Equal(stripped, raw) {
		t.Fatalf("clean csv should remain unchanged")
	}
}

func readDNITTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func firstNonEmptyLine(raw []byte) string {
	for _, line := range bytes.Split(raw, []byte("\n")) {
		trimmed := strings.TrimSpace(string(line))
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
