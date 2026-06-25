package inmet

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSplitStationFileSample(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "bdmep_station_mt.sample.csv"))
	if err != nil {
		t.Fatal(err)
	}
	meta, dataRaw, err := splitStationFile(raw)
	if err != nil {
		t.Fatalf("splitStationFile: %v", err)
	}
	if meta.Code != "A901" || meta.State != "MT" {
		t.Fatalf("meta: %+v", meta)
	}
	rows, err := parseStationDailyLong(meta, dataRaw, 2023, nil)
	if err != nil {
		t.Fatalf("parseStationDailyLong: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected rows from sample csv")
	}
}
