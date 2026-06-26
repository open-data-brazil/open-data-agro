//go:build integration

package mdic

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFetchComexSnapshotLive(t *testing.T) {
	if os.Getenv("MDIC_INTEGRATION") != "1" {
		t.Skip("set MDIC_INTEGRATION=1 to run live Comex Stat API test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:   catalog.MustParseDatasetID("mdic.comex-exportacao-ncm-mes"),
		ComexFlow:   "export",
		ComexNCMs:   []string{"12019000", "10059000"},
		PeriodStart: 2024,
		PeriodEnd:   2024,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	client := NewClient()
	raw, sourceURL, err := client.FetchComexSnapshot(ctx, entry, "2024-01-01")
	if err != nil {
		t.Fatalf("FetchComexSnapshot: %v", err)
	}
	if sourceURL == "" {
		t.Fatal("expected sourceURL")
	}

	var rows []comexRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(rows) < 2 {
		t.Fatalf("rows: got %d want >= 2", len(rows))
	}

	headers, flat, err := FlattenComex(entry, raw)
	if err != nil {
		t.Fatalf("FlattenComex: %v", err)
	}
	if len(flat) != len(rows) {
		t.Fatalf("flatten rows: got %d want %d", len(flat), len(rows))
	}
	t.Logf("live fetch ok: %d rows, source=%s, columns=%d", len(rows), sourceURL, len(headers))
}
