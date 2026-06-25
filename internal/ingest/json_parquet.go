package ingest

import (
	"fmt"
	"os"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/ibge"
)

func convertJSONFileToParquet(entry catalog.RegistryEntry, path string) ([]byte, int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("read staged json: %w", err)
	}
	return convertJSONToParquet(entry, raw)
}

func convertJSONToParquet(entry catalog.RegistryEntry, raw []byte) ([]byte, int, error) {
	agency, _, err := catalog.SplitDatasetID(entry.DatasetID.String())
	if err != nil {
		return nil, 0, err
	}
	if agency != "ibge" {
		return nil, 0, fmt.Errorf("json ingest not implemented for agency %q", agency)
	}

	headers, rows, err := ibge.FlattenIBGEJSON(entry.DatasetID.String(), raw)
	if err != nil {
		return nil, 0, err
	}
	return writeStringTable(headers, rows)
}
