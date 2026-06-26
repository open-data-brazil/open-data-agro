//go:build integration

package eia

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFetchPetroleumPricesLive(t *testing.T) {
	if os.Getenv("EIA_INTEGRATION") != "1" {
		t.Skip("set EIA_INTEGRATION=1 to run live EIA petroleum download test")
	}
	if strings.TrimSpace(os.Getenv("EIA_API_KEY")) == "" {
		t.Skip("set EIA_API_KEY for live EIA petroleum download test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:    catalog.MustParseDatasetID("eia.petroleum-prices"),
		PeriodStart:  2024,
		EIASeriesIDs: []string{"PET.RWTC.D"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	body, sourceURL, err := NewClient().FetchPetroleumSnapshot(ctx, entry, "2024-01-01")
	if err != nil {
		t.Fatalf("FetchPetroleumSnapshot: %v", err)
	}
	_, rows, err := FlattenPetroleum(entry, body)
	if err != nil {
		t.Fatalf("FlattenPetroleum: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("no rows flattened")
	}
	t.Logf("live ok: %d rows, source=%s", len(rows), sourceURL)
}
