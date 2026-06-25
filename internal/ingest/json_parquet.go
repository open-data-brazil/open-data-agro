package ingest

import (
	"fmt"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/ibge"
)

func convertJSONToParquet(entry catalog.RegistryEntry, raw []byte) ([]byte, int, error) {
	agency, _, err := catalog.SplitDatasetID(entry.DatasetID.String())
	if err != nil {
		return nil, 0, err
	}
	if agency != "ibge" {
		return nil, 0, fmt.Errorf("json ingest not implemented for agency %q", agency)
	}

	headers, rows, err := ibge.FlattenLocalidades(entry.DatasetID.String(), raw)
	if err != nil {
		return nil, 0, err
	}
	return writeStringTable(headers, rows)
}
