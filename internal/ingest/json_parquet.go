package ingest

import (
	"fmt"
	"os"

	"github.com/open-data-brazil/open-data-agro/internal/b3"
	"github.com/open-data-brazil/open-data-agro/internal/bcb"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/cepea"
	"github.com/open-data-brazil/open-data-agro/internal/ibge"
	"github.com/open-data-brazil/open-data-agro/internal/ipea"
	"github.com/open-data-brazil/open-data-agro/internal/inmet"
	"github.com/open-data-brazil/open-data-agro/internal/mdic"
	"github.com/open-data-brazil/open-data-agro/internal/usda"
	"github.com/open-data-brazil/open-data-agro/internal/fao"
	"github.com/open-data-brazil/open-data-agro/internal/worldbank"
	"github.com/open-data-brazil/open-data-agro/internal/noaa"
	"github.com/open-data-brazil/open-data-agro/internal/eia"
	"github.com/open-data-brazil/open-data-agro/internal/igc"
	"github.com/open-data-brazil/open-data-agro/internal/inpe"
	"github.com/open-data-brazil/open-data-agro/internal/ana"
	"github.com/open-data-brazil/open-data-agro/internal/eurostat"
	"github.com/open-data-brazil/open-data-agro/internal/argentina"
	"github.com/open-data-brazil/open-data-agro/internal/oecd"
	"github.com/open-data-brazil/open-data-agro/internal/un"
	"github.com/open-data-brazil/open-data-agro/internal/cftc"
	"github.com/open-data-brazil/open-data-agro/internal/jrc"
	"github.com/open-data-brazil/open-data-agro/internal/wto"
	"github.com/open-data-brazil/open-data-agro/internal/fred"
	"github.com/open-data-brazil/open-data-agro/internal/nasa"
	"github.com/open-data-brazil/open-data-agro/internal/sagis"
	"github.com/open-data-brazil/open-data-agro/internal/japan"
	"github.com/open-data-brazil/open-data-agro/internal/mexico"
	"github.com/open-data-brazil/open-data-agro/internal/copernicus"
	"github.com/open-data-brazil/open-data-agro/internal/embrapa"
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
		if entry.DatasetID.String() == "usda.wasde" {
			headers, rows, err := usda.FlattenWASDE(entry, raw)
			if err != nil {
				return nil, 0, err
			}
			return writeStringTable(headers, rows)
		}
		if entry.DatasetID.String() == "usda.gats-trade" {
			headers, rows, err := usda.FlattenGATS(entry, raw)
			if err != nil {
				return nil, 0, err
			}
			return writeStringTable(headers, rows)
		}
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
	case "igc":
		headers, rows, err := igc.FlattenGOI(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "ana":
		headers, rows, err := ana.FlattenHidrologia(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "eurostat":
		headers, rows, err := eurostat.FlattenAgPrices(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "argentina":
		if entry.DatasetID.String() == "argentina.magyp-producion-granos" {
			headers, rows, err := argentina.FlattenGranos(entry, raw)
			if err != nil {
				return nil, 0, err
			}
			return writeStringTable(headers, rows)
		}
		headers, rows, err := argentina.FlattenCambio(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "oecd-fao":
		headers, rows, err := oecd.FlattenAgOutlook(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "inpe":
		headers, rows, err := inpe.FlattenDETER(raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "ipea":
		headers, rows, err := ipea.FlattenSeries(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "inmet":
		if entry.DatasetID.String() == "inmet.sequia-monitor" {
			headers, rows, err := inmet.FlattenSecaMonitor(raw)
			if err != nil {
				return nil, 0, err
			}
			return writeStringTable(headers, rows)
		}
		return nil, 0, fmt.Errorf("json ingest not implemented for inmet dataset %s", entry.DatasetID)
	case "un":
		headers, rows, err := un.FlattenComtrade(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "cftc":
		headers, rows, err := cftc.FlattenCOTAgricultural(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "jrc":
		headers, rows, err := jrc.FlattenMARSCropYield(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "wto":
		headers, rows, err := wto.FlattenITSTrade(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "fred":
		headers, rows, err := fred.FlattenCommodityIndexes(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "nasa":
		headers, rows, err := nasa.FlattenPOWERAgro(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "sagis":
		headers, rows, err := sagis.FlattenGrainSupply(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "japan":
		headers, rows, err := japan.FlattenMAFFAgTrade(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "mexico":
		headers, rows, err := mexico.FlattenSIAPProduccion(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "copernicus":
		headers, rows, err := copernicus.FlattenERA5Agroclimate(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	case "embrapa":
		headers, rows, err := embrapa.FlattenAgroAPIAgrofit(entry, raw)
		if err != nil {
			return nil, 0, err
		}
		return writeStringTable(headers, rows)
	default:
		return nil, 0, fmt.Errorf("json ingest not implemented for agency %q", agency)
	}
}
