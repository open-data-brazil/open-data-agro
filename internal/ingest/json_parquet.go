package ingest

import (
	"fmt"
	"os"

	"github.com/open-data-brazil/open-data-agro/internal/b3"
	"github.com/open-data-brazil/open-data-agro/internal/bcb"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/cepea"
	"github.com/open-data-brazil/open-data-agro/internal/ibge"
	"github.com/open-data-brazil/open-data-agro/internal/mdic"
	"github.com/open-data-brazil/open-data-agro/internal/usda"
	"github.com/open-data-brazil/open-data-agro/internal/fao"
	"github.com/open-data-brazil/open-data-agro/internal/worldbank"
	"github.com/open-data-brazil/open-data-agro/internal/noaa"
	"github.com/open-data-brazil/open-data-agro/internal/eia"
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

	switch agency {
	case "ibge":
		headers, rows, err := ibge.FlattenIBGEJSON(entry.DatasetID.String(), raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "bcb":
		headers, rows, err := bcb.FlattenSGS(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "cepea":
		headers, rows, err := cepea.FlattenIndicador(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "mdic":
		headers, rows, err := mdic.FlattenComex(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "b3":
		headers, rows, err := b3.FlattenFuturo(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "usda":
		headers, rows, err := usda.FlattenPSD(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "fao":
		headers, rows, err := fao.Flatten(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "worldbank":
		headers, rows, err := worldbank.Flatten(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "noaa":
		headers, rows, err := noaa.FlattenClimate(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "eia":
		headers, rows, err := eia.FlattenPetroleum(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	default:
		return nil, 0, fmt.Errorf("json ingest not implemented for agency %q", agency)
	}
}
