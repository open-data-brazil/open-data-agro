package transportes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStripMetadataRowsGolden(t *testing.T) {
	t.Parallel()

	raw := readTransportesTestdata(t, "mtr_bit_malha_rodoviaria.sample.csv")
	stripped := StripMetadataRows(raw)
	if len(stripped) == 0 {
		t.Fatal("expected stripped csv bytes")
	}
	if string(stripped[:2]) != "BR" && string(stripped[:3]) != "\"BR" {
		t.Fatalf("expected BR header row, got %q", string(stripped[:20]))
	}
}

func readTransportesTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
