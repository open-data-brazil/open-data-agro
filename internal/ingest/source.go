package ingest

import (
	"context"
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/anp"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/conab"
	"github.com/open-data-brazil/open-data-agro/internal/ibge"
)

// SourceOptions carries optional ingest parameters (PAM crop/year/UF filters).
type SourceOptions struct {
	Crop     string
	FromYear int
	ToYear   int
	UFs      []string
}

// SourceDownload holds a fetched source file and the resolved download URL.
type SourceDownload struct {
	Body          []byte
	ContentType   string
	LastModified  string
	ContentLength int64
	SourceURL     string
}

// ResolveSourceURL returns the HTTPS URL to download for a catalog entry.
func ResolveSourceURL(entry catalog.RegistryEntry) (string, error) {
	agency, _, err := catalog.SplitDatasetID(entry.DatasetID.String())
	if err != nil {
		return "", err
	}

	switch agency {
	case "conab":
		return conab.ResolveURL(entry)
	case "anp":
		return anp.ResolveURL(entry)
	case "ibge":
		if strings.HasPrefix(entry.DatasetID.String(), "ibge.pam-") {
			return ibge.ResolvePAMURL(entry)
		}
		return ibge.ResolveURL(entry)
	default:
		return "", fmt.Errorf("unsupported agency %q for dataset %s", agency, entry.DatasetID)
	}
}

// DownloadSource fetches bytes for a catalog entry from its official portal.
func DownloadSource(ctx context.Context, entry catalog.RegistryEntry, conabClient *conab.Client, anpClient *anp.Client, ibgeClient *ibge.Client, opts SourceOptions) (*SourceDownload, error) {
	agency, _, err := catalog.SplitDatasetID(entry.DatasetID.String())
	if err != nil {
		return nil, err
	}

	if agency == "ibge" && strings.HasPrefix(entry.DatasetID.String(), "ibge.pam-") {
		body, sourceURL, err := ibgeClient.FetchPAMSnapshot(ctx, entry, ibge.PAMFetchOptions{
			Crop:     opts.Crop,
			FromYear: opts.FromYear,
			ToYear:   opts.ToYear,
			UFs:      opts.UFs,
		})
		if err != nil {
			return nil, err
		}
		return &SourceDownload{
			Body:          body,
			ContentType:   "application/json",
			ContentLength: int64(len(body)),
			SourceURL:     sourceURL,
		}, nil
	}

	sourceURL, err := ResolveSourceURL(entry)
	if err != nil {
		return nil, err
	}

	switch agency {
	case "conab":
		result, err := conabClient.Download(ctx, sourceURL)
		if err != nil {
			return nil, err
		}
		return &SourceDownload{
			Body:          result.Body,
			ContentType:   result.ContentType,
			LastModified:  result.LastModified,
			ContentLength: result.ContentLength,
			SourceURL:     sourceURL,
		}, nil
	case "anp":
		result, err := anpClient.Download(ctx, sourceURL)
		if err != nil {
			return nil, err
		}
		return &SourceDownload{
			Body:          result.Body,
			ContentType:   result.ContentType,
			LastModified:  result.LastModified,
			ContentLength: result.ContentLength,
			SourceURL:     sourceURL,
		}, nil
	case "ibge":
		result, err := ibgeClient.Download(ctx, sourceURL)
		if err != nil {
			return nil, err
		}
		return &SourceDownload{
			Body:          result.Body,
			ContentType:   result.ContentType,
			LastModified:  result.LastModified,
			ContentLength: result.ContentLength,
			SourceURL:     sourceURL,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported agency %q", agency)
	}
}

// AgencyForDataset returns the catalog agency prefix (conab, anp, …).
func AgencyForDataset(datasetID string) (string, error) {
	agency, _, err := catalog.SplitDatasetID(datasetID)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(agency) == "" {
		return "", fmt.Errorf("empty agency in %s", datasetID)
	}
	return agency, nil
}
