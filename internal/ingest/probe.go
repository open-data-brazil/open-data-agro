package ingest

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/cepea"
	"github.com/open-data-brazil/open-data-agro/internal/sourceprobe"
)

// ProbeCatalogEntry resolves the operational URL and probes a sample GET (ingest-compatible headers).
func ProbeCatalogEntry(ctx context.Context, client *sourceprobe.Client, entry catalog.RegistryEntry) sourceprobe.EndpointProbe {
	role := "source"
	url, resolveErr := resolveProbeURL(entry)
	if resolveErr != nil {
		return sourceprobe.EndpointProbe{
			Role:   role,
			URL:    strings.TrimSpace(entry.SourceURL),
			Status: sourceprobe.ProbeResolveError,
			Error:  resolveErr.Error(),
		}
	}

	result, err := client.ProbeURLWithHeaders(ctx, url, probeHeaders(entry))
	if err != nil {
		return sourceprobe.EndpointProbe{
			Role:   role,
			URL:    url,
			Status: sourceprobe.ProbeUnavailable,
			Error:  err.Error(),
		}
	}

	return sourceprobe.EndpointProbe{
		Role:         role,
		URL:          url,
		Status:       sourceprobe.ProbeOK,
		HTTPStatus:   result.HTTPStatus,
		SampleBytes:  result.SampleBytes,
		SampleSHA256: result.SampleSHA256,
		LastModified: result.LastModified,
		ContentType:  result.ContentType,
	}
}

func resolveProbeURL(entry catalog.RegistryEntry) (string, error) {
	agency, _, err := catalog.SplitDatasetID(entry.DatasetID.String())
	if err != nil {
		return "", err
	}

	if agency == "cepea" {
		mirror, mirrorErr := cepea.MirrorURL(entry.DatasetID.String())
		if mirrorErr == nil {
			return mirror, nil
		}
	}

	url, err := ResolveSourceURL(entry)
	if err != nil {
		return "", err
	}

	if agency == "eia" {
		key := strings.TrimSpace(os.Getenv("EIA_API_KEY"))
		if key != "" && !strings.Contains(url, "api_key=") {
			sep := "?"
			if strings.Contains(url, "?") {
				sep = "&"
			}
			url = url + sep + "api_key=" + key
		}
	}

	return url, nil
}

func probeHeaders(entry catalog.RegistryEntry) map[string]string {
	agency, _, err := catalog.SplitDatasetID(entry.DatasetID.String())
	if err != nil {
		return map[string]string{"Accept": "*/*"}
	}

	switch agency {
	case "bcb", "eia", "fred", "nasa", "usda", "worldbank", "oecd-fao", "igc", "eurostat", "argentina", "un", "wto", "jrc", "fao", "ons", "inpe", "ipea", "bndes", "aneel", "mdic", "ibge", "cftc", "copernicus", "sagis", "japan", "mexico":
		return map[string]string{"Accept": "application/json"}
	case "cepea":
		return map[string]string{"Accept": "text/html,application/xhtml+xml"}
	case "conab", "anp", "antt", "antaq", "mapa", "b3", "dnit", "suframa", "transportes", "inmet", "ana", "noaa":
		return map[string]string{"Accept": "*/*"}
	default:
		return map[string]string{"Accept": "*/*"}
	}
}

// ProbePortalURL probes the catalog portal page when distinct from the source URL.
func ProbePortalURL(ctx context.Context, client *sourceprobe.Client, entry catalog.RegistryEntry, sourceURL string) (sourceprobe.EndpointProbe, bool) {
	portalURL := strings.TrimSpace(entry.PortalURL())
	if portalURL == "" || portalURL == sourceURL {
		return sourceprobe.EndpointProbe{}, false
	}

	result, err := client.ProbeURLWithHeaders(ctx, portalURL, map[string]string{"Accept": "text/html,application/xhtml+xml,*/*"})
	if err != nil {
		return sourceprobe.EndpointProbe{
			Role:   "portal",
			URL:    portalURL,
			Status: sourceprobe.ProbeUnavailable,
			Error:  err.Error(),
		}, true
	}

	return sourceprobe.EndpointProbe{
		Role:         "portal",
		URL:          portalURL,
		Status:       sourceprobe.ProbeOK,
		HTTPStatus:   result.HTTPStatus,
		SampleBytes:  result.SampleBytes,
		SampleSHA256: result.SampleSHA256,
		LastModified: result.LastModified,
		ContentType:  result.ContentType,
	}, true
}

// ProbeSourceURLForTest exposes resolveProbeURL for unit tests.
func ProbeSourceURLForTest(entry catalog.RegistryEntry) (string, error) {
	url, err := resolveProbeURL(entry)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(url) == "" {
		return "", fmt.Errorf("empty probe url")
	}
	return url, nil
}

// ProbeAll probes every catalog dataset with ingest-compatible URLs and headers.
func ProbeAll(ctx context.Context, reg *catalog.Registry, concurrency int) []sourceprobe.DatasetProbeOutcome {
	entries := reg.Entries()
	results := make([]sourceprobe.DatasetProbeOutcome, len(entries))
	if concurrency <= 0 {
		concurrency = 6
	}

	client := sourceprobe.NewClient()
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, entry := range entries {
		wg.Add(1)
		go func(idx int, entry catalog.RegistryEntry) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[idx] = probeDataset(ctx, client, entry)
		}(i, entry)
	}

	wg.Wait()
	return results
}

func probeDataset(ctx context.Context, client *sourceprobe.Client, entry catalog.RegistryEntry) sourceprobe.DatasetProbeOutcome {
	outcome := sourceprobe.DatasetProbeOutcome{
		DatasetID: entry.DatasetID.String(),
		ProbedAt:  time.Now().UTC(),
	}

	sourceProbe := ProbeCatalogEntry(ctx, client, entry)
	outcome.Endpoints = append(outcome.Endpoints, sourceProbe)

	if portalProbe, ok := ProbePortalURL(ctx, client, entry, sourceProbe.URL); ok {
		outcome.Endpoints = append(outcome.Endpoints, portalProbe)
	}

	source := sourceEndpointOutcome(&outcome)
	outcome.OverallOK = source != nil && source.Status == sourceprobe.ProbeOK
	return outcome
}

func sourceEndpointOutcome(outcome *sourceprobe.DatasetProbeOutcome) *sourceprobe.EndpointProbe {
	for i := range outcome.Endpoints {
		if outcome.Endpoints[i].Role == "source" {
			return &outcome.Endpoints[i]
		}
	}
	if len(outcome.Endpoints) > 0 {
		return &outcome.Endpoints[0]
	}
	return nil
}
