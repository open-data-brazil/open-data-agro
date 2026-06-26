//go:build integration

package worldbank

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFetchPinkSheetMonthlyLive(t *testing.T) {
	if os.Getenv("WORLDBANK_INTEGRATION") != "1" {
		t.Skip("set WORLDBANK_INTEGRATION=1 to run live World Bank Pink Sheet download test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:   catalog.MustParseDatasetID("worldbank.pink-sheet-monthly"),
		PeriodStart: 2024,
		PeriodEnd:   2024,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	body, sourceURL, err := NewClient().FetchPinkSheetSnapshot(ctx, entry, "2024-01")
	if err != nil {
		t.Fatalf("FetchPinkSheetSnapshot: %v", err)
	}
	if len(body) < 100 {
		t.Fatalf("body too small: %d", len(body))
	}
	headers, rows, err := FlattenPinkSheet(entry, body)
	if err != nil {
		t.Fatalf("FlattenPinkSheet: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("no rows flattened")
	}
	t.Logf("live ok: %d rows, headers=%v, source=%s", len(rows), headers, sourceURL)
}
