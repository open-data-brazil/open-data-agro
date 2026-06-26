//go:build integration

package noaa

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFetchENSOLive(t *testing.T) {
	if os.Getenv("NOAA_INTEGRATION") != "1" {
		t.Skip("set NOAA_INTEGRATION=1 to run live NOAA ONI download test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:   catalog.MustParseDatasetID("noaa.enso-indices"),
		PeriodStart: 2020,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	body, sourceURL, err := NewClient().FetchENSOSnapshot(ctx, entry, "2020")
	if err != nil {
		t.Fatalf("FetchENSOSnapshot: %v", err)
	}
	if len(body) < 50 {
		t.Fatalf("body too small: %d", len(body))
	}
	_, rows, err := FlattenENSO(entry, body)
	if err != nil {
		t.Fatalf("FlattenENSO: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("no rows flattened")
	}
	t.Logf("live ONI ok: %d rows, source=%s", len(rows), sourceURL)
}

func TestFetchGlobalTempLive(t *testing.T) {
	if os.Getenv("NOAA_INTEGRATION") != "1" {
		t.Skip("set NOAA_INTEGRATION=1 to run live NOAA global temp download test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:   catalog.MustParseDatasetID("noaa.global-temp-anomaly"),
		PeriodStart: 2020,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	body, sourceURL, err := NewClient().FetchGlobalTempSnapshot(ctx, entry, "2020-01")
	if err != nil {
		t.Fatalf("FetchGlobalTempSnapshot: %v", err)
	}
	if len(body) < 50 {
		t.Fatalf("body too small: %d", len(body))
	}
	_, rows, err := FlattenGlobalTemp(entry, body)
	if err != nil {
		t.Fatalf("FlattenGlobalTemp: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("no rows flattened")
	}
	t.Logf("live global temp ok: %d rows, source=%s", len(rows), sourceURL)
}
