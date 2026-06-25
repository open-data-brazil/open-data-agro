package ingest

import (
	"fmt"
	"os"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/inmet"
)

func convertINMETCSVToParquet(entry catalog.RegistryEntry, raw []byte) ([]byte, int, error) {
	headers, rows, err := inmet.FlattenINMETCSV(entry.DatasetID.String(), raw)
	if err != nil {
		return nil, 0, err
	}
	return writeStringTable(headers, rows)
}

func convertINMETCSVFileToParquet(entry catalog.RegistryEntry, path string) ([]byte, int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("read staged inmet csv: %w", err)
	}
	return convertINMETCSVToParquet(entry, raw)
}

func isINMETDataset(datasetID string) bool {
	agency, _, err := catalog.SplitDatasetID(datasetID)
	if err != nil {
		return false
	}
	return agency == "inmet"
}
