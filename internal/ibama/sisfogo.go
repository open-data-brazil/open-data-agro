package ibama

import (
	"context"
	"fmt"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// FetchSISFOGOSnapshot downloads the SISFOGO fire occurrence ROI CSV.
func (c *Client) FetchSISFOGOSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	sourceURL, err := ResolveURL(entry)
	if err != nil {
		return nil, "", err
	}
	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}
	return NormalizeCSV(result.Body), sourceURL, nil
}

// FetchLicencasSnapshot downloads federal environmental licenses (SISLIC).
func (c *Client) FetchLicencasSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	sourceURL, err := ResolveURL(entry)
	if err != nil {
		return nil, "", err
	}
	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}
	return NormalizeCSV(result.Body), sourceURL, nil
}

// FetchAutosSnapshot downloads one yearly autos CSV extracted from the official ZIP bundle.
func (c *Client) FetchAutosSnapshot(ctx context.Context, entry catalog.RegistryEntry, year int) ([]byte, string, error) {
	if path := autosBulkPath(); path != "" {
		body, member, err := extractAutosYearFromFile(path, year)
		if err != nil {
			return nil, "", err
		}
		return NormalizeCSV(body), path + " (" + member + ")", nil
	}

	sourceURL, err := ResolveURL(entry)
	if err != nil {
		return nil, "", err
	}
	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}
	body, member, err := extractAutosYearFromZIP(result.Body, year)
	if err != nil {
		return nil, "", err
	}
	return NormalizeCSV(body), fmt.Sprintf("%s (%s)", sourceURL, member), nil
}
