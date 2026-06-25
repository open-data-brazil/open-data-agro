package ingest

import (
	"context"
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/anp"
	"github.com/open-data-brazil/open-data-agro/internal/bcb"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/conab"
	"github.com/open-data-brazil/open-data-agro/internal/ibge"
	"github.com/open-data-brazil/open-data-agro/internal/inmet"
)

// SourceOptions carries optional ingest parameters (PAM crop/year/UF filters, INMET year).
type SourceOptions struct {
	Crop     string
	FromYear int
	ToYear   int
	Year     int
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
	case "inmet":
		return inmet.ResolveURL(entry)
	case "bcb":
		return bcb.ResolveURL(entry)
	default:
		return "", fmt.Errorf("unsupported agency %q for dataset %s", agency, entry.DatasetID)
	}
}

// DownloadSource fetches bytes for a catalog entry from its official portal.
func DownloadSource(ctx context.Context, entry catalog.RegistryEntry, conabClient *conab.Client, anpClient *anp.Client, ibgeClient *ibge.Client, inmetClient *inmet.Client, bcbClient *bcb.Client, opts SourceOptions) (*SourceDownload, error) {
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

	if agency == "inmet" {
		return downloadINMETSource(ctx, entry, inmetClient, opts)
	}

	if agency == "bcb" {
		body, sourceURL, err := bcbClient.FetchSGSSnapshot(ctx, entry)
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

func downloadINMETSource(ctx context.Context, entry catalog.RegistryEntry, client *inmet.Client, opts SourceOptions) (*SourceDownload, error) {
	datasetID := entry.DatasetID.String()

	switch datasetID {
	case "inmet.estacoes-automaticas", "inmet.estacoes-convencionais":
		sourceURL, err := inmet.ResolveURL(entry)
		if err != nil {
			return nil, err
		}
		result, err := client.Download(ctx, sourceURL)
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
	case "inmet.bdmep-diario", "inmet.pacote-anual-automaticas":
		year := resolveINMETYear(entry, opts)
		body, sourceURL, err := client.FetchBDMEPDailySnapshot(ctx, entry, inmet.BDMEPFetchOptions{
			Year: year,
			UFs:  opts.UFs,
		})
		if err != nil {
			return nil, err
		}
		return &SourceDownload{
			Body:          body,
			ContentType:   "text/csv",
			ContentLength: int64(len(body)),
			SourceURL:     sourceURL,
		}, nil
	case "inmet.bdmep-mensal":
		year := resolveINMETYear(entry, opts)
		body, sourceURL, err := client.FetchBDMEPMonthlySnapshot(ctx, entry, inmet.BDMEPFetchOptions{
			Year: year,
			UFs:  opts.UFs,
		})
		if err != nil {
			return nil, err
		}
		return &SourceDownload{
			Body:          body,
			ContentType:   "text/csv",
			ContentLength: int64(len(body)),
			SourceURL:     sourceURL,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported inmet dataset %s", datasetID)
	}
}

func resolveINMETYear(entry catalog.RegistryEntry, opts SourceOptions) int {
	if opts.Year > 0 {
		return opts.Year
	}
	if opts.FromYear > 0 {
		return opts.FromYear
	}
	if entry.PeriodEnd > 0 {
		return entry.PeriodEnd
	}
	return 0
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
