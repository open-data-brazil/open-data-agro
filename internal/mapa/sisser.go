package mapa

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const sisserPackageID = "sisser3"

// FetchSISSERSnapshot downloads and merges all PSR CSV resources from SISSER CKAN package.
func (c *Client) FetchSISSERSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	resources, err := ListCKANResources(ctx, sisserPackageID)
	if err != nil {
		return nil, "", err
	}

	var psrResources []ckanResource
	for _, res := range resources {
		if strings.ToUpper(strings.TrimSpace(res.Format)) != "CSV" {
			continue
		}
		if !strings.Contains(strings.ToUpper(res.Name), "PSR") {
			continue
		}
		if strings.TrimSpace(res.URL) == "" {
			continue
		}
		psrResources = append(psrResources, res)
	}
	if len(psrResources) == 0 {
		return nil, "", fmt.Errorf("no PSR CSV resources in sisser package")
	}

	sort.Slice(psrResources, func(i, j int) bool {
		return psrResources[i].Name < psrResources[j].Name
	})

	var headers []string
	var merged [][]string
	var requestURLs []string

	for _, res := range psrResources {
		download, err := c.Download(ctx, res.URL)
		if err != nil {
			return nil, "", fmt.Errorf("download %s: %w", res.Name, err)
		}
		fileHeaders, rows, err := parseSISSERCSV(download.Body, res.Name)
		if err != nil {
			return nil, "", err
		}
		if len(headers) == 0 {
			headers = append([]string{"periodo_arquivo"}, fileHeaders...)
		}
		merged = append(merged, rows...)
		requestURLs = append(requestURLs, res.URL)
	}
	if len(merged) == 0 {
		return nil, "", fmt.Errorf("sisser returned no rows for %s", entry.DatasetID)
	}

	payload, err := encodeCSV(headers, merged)
	if err != nil {
		return nil, "", err
	}
	sourceURL := fmt.Sprintf("%s?id=%s (resources: %d)", defaultCKANPackageShowURL, sisserPackageID, len(requestURLs))
	return payload, sourceURL, nil
}

func resolveSISSERURL(entry catalog.RegistryEntry) (string, error) {
	return fmt.Sprintf("%s?id=%s", defaultCKANPackageShowURL, sisserPackageID), nil
}

func parseSISSERCSV(raw []byte, periodLabel string) ([]string, [][]string, error) {
	raw = decodeLatin1(raw)
	reader := csv.NewReader(bytes.NewReader(raw))
	reader.Comma = ';'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("read sisser header: %w", err)
	}

	var rows [][]string
	for {
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, nil, fmt.Errorf("read sisser row: %w", readErr)
		}
		row := make([]string, 0, len(record)+1)
		row = append(row, periodLabel)
		row = append(row, record...)
		rows = append(rows, row)
	}
	return headers, rows, nil
}
