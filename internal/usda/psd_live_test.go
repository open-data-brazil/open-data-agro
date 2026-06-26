//go:build integration

package usda

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFetchPSDSojaLive(t *testing.T) {
	if os.Getenv("USDA_INTEGRATION") != "1" {
		t.Skip("set USDA_INTEGRATION=1 to run live USDA PSD SOAP test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:        catalog.MustParseDatasetID("usda.psd-soja"),
		PSDCommodityCode: "2222000",
		PSDCommoditySlug: "soja",
		PeriodStart:      2024,
		PeriodEnd:        2024,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	body, sourceURL, err := NewClient().FetchPSDSnapshot(ctx, entry, "2024")
	if err != nil {
		t.Fatalf("FetchPSDSnapshot: %v", err)
	}
	if len(body) < 100 {
		t.Fatalf("body too small: %d", len(body))
	}
	headers, rows, err := FlattenPSD(entry, body)
	if err != nil {
		t.Fatalf("FlattenPSD: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("no rows flattened")
	}
	t.Logf("live ok: %d rows, headers=%v, source=%s", len(rows), headers, sourceURL)
}
