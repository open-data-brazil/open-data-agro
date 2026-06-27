package transportes

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultBaseFerroZIPURL = "https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip"

// FetchMTRBITShapefileSnapshot extracts railway network attribute table (DBF) from MTR BIT ZIP.
func (c *Client) FetchMTRBITShapefileSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	if path := strings.TrimSpace(os.Getenv("TRANSPORTES_SHAPEFILE_BULK_PATH")); path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, "", err
		}
		csvBytes, member, extractErr := extractDBFFromZIP(raw)
		if extractErr != nil {
			return nil, "", extractErr
		}
		return csvBytes, path + " (" + member + ")", nil
	}

	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = defaultBaseFerroZIPURL
	}
	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}
	csvBytes, member, err := extractDBFFromZIP(result.Body)
	if err != nil {
		return nil, "", err
	}
	return csvBytes, fmt.Sprintf("%s (%s attributes)", sourceURL, member), nil
}

func extractDBFFromZIP(zipBytes []byte) ([]byte, string, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, "", fmt.Errorf("parse transportes shapefile zip: %w", err)
	}

	var dbfName string
	for _, file := range reader.File {
		name := strings.ToLower(strings.TrimSpace(file.Name))
		if strings.HasSuffix(name, ".dbf") && !strings.Contains(name, "xml") {
			dbfName = file.Name
			break
		}
	}
	if dbfName == "" {
		return nil, "", fmt.Errorf("no .dbf member in transportes shapefile zip")
	}

	for _, file := range reader.File {
		if file.Name != dbfName {
			continue
		}
		rc, openErr := file.Open()
		if openErr != nil {
			return nil, "", openErr
		}
		defer func() { _ = rc.Close() }()
		dbfBytes, readErr := ioReadAllLimit(rc, 64<<20)
		if readErr != nil {
			return nil, "", readErr
		}
		csvBytes, convErr := ParseDBFAttributes(dbfBytes)
		if convErr != nil {
			return nil, "", convErr
		}
		return csvBytes, dbfName, nil
	}
	return nil, "", fmt.Errorf("dbf member %s not found", dbfName)
}
