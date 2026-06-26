//go:build integration

package fao

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFetchPricesAgroLive(t *testing.T) {
	if os.Getenv("FAO_INTEGRATION") != "1" {
		t.Skip("set FAO_INTEGRATION=1 to run live FAOSTAT bulk download test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:   catalog.MustParseDatasetID("fao.prices-agro"),
		PeriodStart: 2020,
		PeriodEnd:   2020,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	body, sourceURL, err := NewClient().FetchPricesSnapshot(ctx, entry, "2020")
	if err != nil {
		t.Fatalf("FetchPricesSnapshot: %v", err)
	}
	if len(body) < 100 {
		t.Fatalf("body too small: %d", len(body))
	}
	headers, rows, err := FlattenPrices(entry, body)
	if err != nil {
		t.Fatalf("FlattenPrices: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("no rows flattened")
	}
	t.Logf("live ok: %d rows, headers=%v, source=%s", len(rows), headers, sourceURL)
}
