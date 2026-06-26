package noaa

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenENSOGolden(t *testing.T) {
	t.Parallel()

	raw := readNOAATestdata(t, "enso.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("noaa.enso-indices"),
	}

	headers, rows, err := FlattenENSO(entry, raw)
	if err != nil {
		t.Fatalf("FlattenENSO: %v", err)
	}
	if len(headers) != 5 {
		t.Fatalf("headers: got %d want 5", len(headers))
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}
	if rows[0][1] != "DJF" {
		t.Fatalf("season_code: got %q", rows[0][1])
	}
}

func TestFlattenGlobalTempGolden(t *testing.T) {
	t.Parallel()

	raw := readNOAATestdata(t, "global_temp.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("noaa.global-temp-anomaly"),
	}

	headers, rows, err := FlattenGlobalTemp(entry, raw)
	if err != nil {
		t.Fatalf("FlattenGlobalTemp: %v", err)
	}
	if len(headers) != 5 {
		t.Fatalf("headers: got %d want 5", len(headers))
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}
	if rows[0][0] != "2024-01" {
		t.Fatalf("refmonth: got %q", rows[0][0])
	}
}

func TestParseONIASCIISample(t *testing.T) {
	t.Parallel()

	raw := readNOAATestdata(t, "oni.sample.txt")
	rows, err := parseONIASCII(raw, 2020, 2025)
	if err != nil {
		t.Fatalf("parseONIASCII: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}
	if rows[0].Anomaly != "1.53" {
		t.Fatalf("anomaly: got %q", rows[0].Anomaly)
	}
}

func TestParseGlobalTempCSVSample(t *testing.T) {
	t.Parallel()

	raw := []byte(`# Title: Global Land and Ocean Average Temperature Departures
# Units: Degrees Celsius
Date,Departure from Average
202401,0.74
202402,0.81
`)
	rows, err := parseGlobalTempCSV(raw, "2024-01", "2024-12")
	if err != nil {
		t.Fatalf("parseGlobalTempCSV: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows: got %d want 2", len(rows))
	}
	if rows[0].RefMonth != "2024-01" {
		t.Fatalf("refmonth: got %q", rows[0].RefMonth)
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	enso := catalog.RegistryEntry{
		DatasetID:    catalog.MustParseDatasetID("noaa.enso-indices"),
		NOAAIndexURL: defaultONIURL,
	}
	url, err := ResolveURL(enso)
	if err != nil || url == "" {
		t.Fatalf("ResolveURL enso: url=%q err=%v", url, err)
	}

	temp := catalog.RegistryEntry{
		DatasetID:   catalog.MustParseDatasetID("noaa.global-temp-anomaly"),
		PeriodStart: 2010,
		PeriodEnd:   2024,
	}
	url, err = ResolveURL(temp)
	if err != nil || url == "" {
		t.Fatalf("ResolveURL temp: url=%q err=%v", url, err)
	}
}

func readNOAATestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
