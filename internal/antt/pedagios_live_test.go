//go:build integration

package antt

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestDownloadPracasPedagioLive(t *testing.T) {
	if os.Getenv("ANTT_INTEGRATION") != "1" {
		t.Skip("set ANTT_INTEGRATION=1 to run live ANTT CKAN download test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:          catalog.MustParseDatasetID("antt.pracas-pedagio"),
		CKANPackageID:      "praca-de-pedagio",
		CKANResourceFormat: "CSV",
	}

	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	client := NewClient()
	result, err := client.Download(ctx, url)
	if err != nil {
		t.Fatalf("Download: %v", err)
	}
	if len(result.Body) < 100 {
		t.Fatalf("body too small: %d bytes", len(result.Body))
	}
	t.Logf("live download ok: %d bytes from %s", len(result.Body), url)
}
