//go:build integration

package b3

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFetchFuturoSojaLive(t *testing.T) {
	if os.Getenv("B3_INTEGRATION") != "1" {
		t.Skip("set B3_INTEGRATION=1 to run live B3 SPRD download test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:         catalog.MustParseDatasetID("b3.futuro-soja"),
		B3FilePrefix:      "SPRD",
		B3CommodityPrefix: "SOY",
		PeriodStart:       2025,
		PeriodEnd:         2025,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	body, sourceURL, err := NewClient().FetchFuturoSnapshot(ctx, entry, "2025-06-01")
	if err != nil {
		t.Fatalf("FetchFuturoSnapshot: %v", err)
	}
	if len(body) < 50 {
		t.Fatalf("body too small: %d", len(body))
	}
	headers, rows, err := FlattenFuturo(entry, body)
	if err != nil {
		t.Fatalf("FlattenFuturo: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("no rows flattened")
	}
	t.Logf("live ok: %d rows, headers=%v, source=%s", len(rows), headers, sourceURL)
}
