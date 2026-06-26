//go:build integration

package mapa

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestDownloadZARCTabuaRiscoLive(t *testing.T) {
	if os.Getenv("MAPA_INTEGRATION") != "1" {
		t.Skip("set MAPA_INTEGRATION=1 to run live MAPA CKAN download test")
	}

	entry := catalog.RegistryEntry{
		DatasetID:          catalog.MustParseDatasetID("mapa.zarc-tabua-risco"),
		CKANPackageID:      "tabua-de-risco-zoneamento-agricola-de-risco-climatico",
		CKANResourceFormat: "CSV",
	}

	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HEAD: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("HEAD status %d for %s", resp.StatusCode, url)
	}
	t.Logf("live HEAD ok: status=%d url=%s", resp.StatusCode, url)
}
