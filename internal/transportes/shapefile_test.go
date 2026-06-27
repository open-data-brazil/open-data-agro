package transportes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDBFAttributesFromBaseFerroZIP(t *testing.T) {
	t.Parallel()

	zipPath := os.Getenv("TRANSPORTES_SHAPEFILE_BULK_PATH")
	if zipPath == "" {
		zipPath = "/tmp/BaseFerro.zip"
	}
	raw, err := os.ReadFile(zipPath)
	if err != nil {
		t.Skipf("shapefile zip not available at %s: %v", zipPath, err)
	}

	csvBytes, member, err := extractDBFFromZIP(raw)
	if err != nil {
		t.Fatalf("extractDBFFromZIP: %v", err)
	}
	if member == "" {
		t.Fatal("empty member name")
	}
	if len(csvBytes) < 100 {
		t.Fatalf("csv too short: %d", len(csvBytes))
	}
}

func TestParseDBFAttributesSample(t *testing.T) {
	t.Parallel()

	path := filepath.Join("testdata", "baseferro_attributes.sample.csv")
	sample, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}
	if len(sample) < 50 {
		t.Fatal("sample too short")
	}
}
