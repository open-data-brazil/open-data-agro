package ingest

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/sourceprobe"
)

// ProbeCatalogEntry resolves the operational URL and probes a sample request (ingest-compatible).
func ProbeCatalogEntry(ctx context.Context, client *sourceprobe.Client, entry catalog.RegistryEntry) sourceprobe.EndpointProbe {
	spec, err := BuildProbeSpec(entry)
	if err != nil {
		return sourceprobe.EndpointProbe{
			Role:   "source",
			URL:    strings.TrimSpace(entry.SourceURL),
			Status: sourceprobe.ProbeResolveError,
			Error:  err.Error(),
		}
	}

	var result *sourceprobe.SampleResult
	var probeErr error
	if spec.Method == "POST" || spec.Body != "" {
		result, probeErr = client.ProbeHTTP(ctx, sourceprobe.HTTPRequest{
			Method:  spec.Method,
			URL:     spec.URL,
			Headers: spec.Headers,
			Body:    spec.Body,
		})
	} else {
		result, probeErr = client.ProbeURLWithHeaders(ctx, spec.URL, spec.Headers)
	}
	if probeErr != nil {
		return sourceprobe.EndpointProbe{
			Role:   "source",
			URL:    spec.URL,
			Status: sourceprobe.ProbeUnavailable,
			Error:  probeErr.Error(),
		}
	}

	return sourceprobe.EndpointProbe{
		Role:         "source",
		URL:          spec.URL,
		Status:       sourceprobe.ProbeOK,
		HTTPStatus:   result.HTTPStatus,
		SampleBytes:  result.SampleBytes,
		SampleSHA256: result.SampleSHA256,
		LastModified: result.LastModified,
		ContentType:  result.ContentType,
	}
}

func probeHeaders(entry catalog.RegistryEntry) map[string]string {
	agency, _, err := catalog.SplitDatasetID(entry.DatasetID.String())
	if err != nil {
		return map[string]string{"Accept": "*/*"}
	}

	switch agency {
	case "bcb", "eia", "fred", "nasa", "worldbank", "oecd-fao", "eurostat", "argentina", "jrc", "fao", "ons", "inpe", "ipea", "bndes", "aneel", "mdic", "ibge", "cftc", "copernicus", "sagis", "mexico":
		return map[string]string{"Accept": "application/json"}
	case "cepea", "japan":
		return map[string]string{"Accept": "text/html,application/xhtml+xml,*/*"}
	case "conab", "anp", "antt", "antaq", "mapa", "b3", "dnit", "suframa", "transportes", "inmet", "ana", "noaa", "igc", "usda", "un", "wto":
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

// ProbeSourceURLForTest exposes BuildProbeSpec for unit tests.
func ProbeSourceURLForTest(entry catalog.RegistryEntry) (string, error) {
	spec, err := BuildProbeSpec(entry)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(spec.URL) == "" {
		return "", fmt.Errorf("empty probe url")
	}
	return spec.URL, nil
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
