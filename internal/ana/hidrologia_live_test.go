//go:build integration

package ana

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFetchHidrologiaSnapshotLive(t *testing.T) {
	if os.Getenv("ANA_INTEGRATION") != "1" {
		t.Skip("set ANA_INTEGRATION=1 to run live ANA HidroWeb test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:            catalog.MustParseDatasetID("ana.hidrologia-series"),
		SourceURL:            defaultServiceBaseURL,
		ANAStationCodes:      []string{"15400000"},
		ANATipoDados:         "3",
		ANANivelConsistencia: "2",
		StartDate:            "01/01/2024",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	body, sourceURL, err := NewClient().FetchHidrologiaSnapshot(ctx, entry, HidrologiaFetchOptions{
		DataInicio: "01/01/2024",
		DataFim:    "07/01/2024",
	})
	if err != nil {
		t.Fatalf("FetchHidrologiaSnapshot: %v", err)
	}
	if len(body) == 0 {
		t.Fatal("empty body")
	}
	t.Logf("live ok: bytes=%d url=%s", len(body), sourceURL)
}
