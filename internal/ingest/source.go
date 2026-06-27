package ingest

import (
	"context"
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/anp"
	"github.com/open-data-brazil/open-data-agro/internal/aneel"
	"github.com/open-data-brazil/open-data-agro/internal/antt"
	"github.com/open-data-brazil/open-data-agro/internal/antaq"
	"github.com/open-data-brazil/open-data-agro/internal/ana"
	"github.com/open-data-brazil/open-data-agro/internal/b3"
	"github.com/open-data-brazil/open-data-agro/internal/bcb"
	"github.com/open-data-brazil/open-data-agro/internal/bndes"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/cepea"
	"github.com/open-data-brazil/open-data-agro/internal/conab"
	"github.com/open-data-brazil/open-data-agro/internal/dnit"
	"github.com/open-data-brazil/open-data-agro/internal/ibge"
	"github.com/open-data-brazil/open-data-agro/internal/inmet"
	"github.com/open-data-brazil/open-data-agro/internal/ipea"
	"github.com/open-data-brazil/open-data-agro/internal/mapa"
	"github.com/open-data-brazil/open-data-agro/internal/mdic"
	"github.com/open-data-brazil/open-data-agro/internal/usda"
	"github.com/open-data-brazil/open-data-agro/internal/fao"
	"github.com/open-data-brazil/open-data-agro/internal/worldbank"
	"github.com/open-data-brazil/open-data-agro/internal/noaa"
	"github.com/open-data-brazil/open-data-agro/internal/eia"
	"github.com/open-data-brazil/open-data-agro/internal/igc"
	"github.com/open-data-brazil/open-data-agro/internal/eurostat"
	"github.com/open-data-brazil/open-data-agro/internal/argentina"
	"github.com/open-data-brazil/open-data-agro/internal/embrapa"
	"github.com/open-data-brazil/open-data-agro/internal/ibama"
	"github.com/open-data-brazil/open-data-agro/internal/inpe"
	"github.com/open-data-brazil/open-data-agro/internal/ons"
	"github.com/open-data-brazil/open-data-agro/internal/suframa"
	"github.com/open-data-brazil/open-data-agro/internal/transportes"
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
	"github.com/open-data-brazil/open-data-agro/internal/oecd"
)

// SourceOptions carries optional ingest parameters (PAM crop/year/UF filters, INMET year).
type SourceOptions struct {
	Crop     string
	FromYear int
	ToYear   int
	Year     int
	FromDate string
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
	case "antt":
		return antt.ResolveURL(entry)
	case "aneel":
		return aneel.ResolveURL(entry)
	case "bndes":
		return bndes.ResolveURL(entry)
	case "dnit":
		return dnit.ResolveURL(entry)
	case "suframa":
		return suframa.ResolveURL(entry)
	case "transportes":
		return transportes.ResolveURL(entry)
	case "ons":
		return ons.ResolveURL(entry)
	case "inpe":
		return inpe.ResolveURL(entry)
	case "ipea":
		return ipea.ResolveURL(entry)
	case "ibge":
		if strings.HasPrefix(entry.DatasetID.String(), "ibge.pam-") {
			return ibge.ResolvePAMURL(entry)
		}
		if strings.HasPrefix(entry.DatasetID.String(), "ibge.lspa-") {
			return ibge.ResolveLSPAURL(entry)
		}
		if strings.HasPrefix(entry.DatasetID.String(), "ibge.pevs-") {
			return ibge.ResolvePEVSURL(entry)
		}
		if strings.HasPrefix(entry.DatasetID.String(), "ibge.ppm-") {
			return ibge.ResolvePPMURL(entry)
		}
		if strings.HasPrefix(entry.DatasetID.String(), "ibge.censo-agro-") {
			return ibge.ResolveCensoAgroURL(entry)
		}
		if strings.HasPrefix(entry.DatasetID.String(), "ibge.pnad-continua-") ||
			strings.HasPrefix(entry.DatasetID.String(), "ibge.pnad-rural-") {
			return ibge.ResolvePNADRuralURL(entry)
		}
		return ibge.ResolveURL(entry)
	case "inmet":
		return inmet.ResolveURL(entry)
	case "bcb":
		return bcb.ResolveURL(entry)
	case "cepea":
		return cepea.ResolveURL(entry)
	case "mdic":
		return mdic.ResolveURL(entry)
	case "mapa":
		return mapa.ResolveURL(entry)
	case "b3":
		return b3.ResolveURL(entry)
	case "usda":
		return usda.ResolveURL(entry)
	case "fao":
		return fao.ResolveURL(entry)
	case "worldbank":
		return worldbank.ResolveURL(entry)
	case "noaa":
		return noaa.ResolveURL(entry)
	case "eia":
		return eia.ResolveURL(entry)
	case "igc":
		return igc.ResolveURL(entry)
	case "ana":
		return ana.ResolveURL(entry)
	case "antaq":
		return antaq.ResolveURL(entry)
	case "embrapa":
		return embrapa.ResolveURL(entry)
	case "ibama":
		return ibama.ResolveURL(entry)
	case "un":
		return un.ResolveURL(entry)
	case "eurostat":
		return eurostat.ResolveURL(entry)
	case "argentina":
		return argentina.ResolveURL(entry)
	case "oecd-fao":
		return oecd.ResolveURL(entry)
	case "cftc":
		return cftc.ResolveURL(entry)
	case "jrc":
		return jrc.ResolveURL(entry)
	case "wto":
		return wto.ResolveURL(entry)
	case "fred":
		return fred.ResolveURL(entry)
	case "nasa":
		return nasa.ResolveURL(entry)
	case "sagis":
		return sagis.ResolveURL(entry)
	case "japan":
		return japan.ResolveURL(entry)
	case "mexico":
		return mexico.ResolveURL(entry)
	case "copernicus":
		return copernicus.ResolveURL(entry)
	default:
		return "", fmt.Errorf("unsupported agency %q for dataset %s", agency, entry.DatasetID)
	}
}

// DownloadSource fetches bytes for a catalog entry from its official portal.
func DownloadSource(ctx context.Context, entry catalog.RegistryEntry, conabClient *conab.Client, anpClient *anp.Client, anttClient *antt.Client, aneelClient *aneel.Client, bndesClient *bndes.Client, ibgeClient *ibge.Client, inmetClient *inmet.Client, bcbClient *bcb.Client, cepeaClient *cepea.Client, mdicClient *mdic.Client, mapaClient *mapa.Client, b3Client *b3.Client, usdaClient *usda.Client, faoClient *fao.Client, worldbankClient *worldbank.Client, noaaClient *noaa.Client, eiaClient *eia.Client, igcClient *igc.Client, anaClient *ana.Client, antaqClient *antaq.Client, dnitClient *dnit.Client, ipeaClient *ipea.Client, eurostatClient *eurostat.Client, argentinaClient *argentina.Client, oecdClient *oecd.Client, unClient *un.Client, cftcClient *cftc.Client, jrcClient *jrc.Client, wtoClient *wto.Client, fredClient *fred.Client, nasaClient *nasa.Client, sagisClient *sagis.Client, japanClient *japan.Client, mexicoClient *mexico.Client, copernicusClient *copernicus.Client, suframaClient *suframa.Client, transportesClient *transportes.Client, onsClient *ons.Client, inpeClient *inpe.Client, ibamaClient *ibama.Client, embrapaClient *embrapa.Client, opts SourceOptions) (*SourceDownload, error) {
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

	if agency == "ibge" && strings.HasPrefix(entry.DatasetID.String(), "ibge.lspa-") {
		body, sourceURL, err := ibgeClient.FetchLSPASnapshot(ctx, entry, ibge.LSPAFetchOptions{
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

	if agency == "ibge" && strings.HasPrefix(entry.DatasetID.String(), "ibge.pevs-") {
		body, sourceURL, err := ibgeClient.FetchPEVSSnapshot(ctx, entry, ibge.PEVSFetchOptions{
			FromYear: opts.FromYear,
			ToYear:   opts.ToYear,
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

	if agency == "ibge" && strings.HasPrefix(entry.DatasetID.String(), "ibge.ppm-") {
		body, sourceURL, err := ibgeClient.FetchPPMSnapshot(ctx, entry, ibge.PPMFetchOptions{
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

	if agency == "ibge" && strings.HasPrefix(entry.DatasetID.String(), "ibge.censo-agro-") {
		body, sourceURL, err := ibgeClient.FetchCensoAgroSnapshot(ctx, entry)
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

	if agency == "ibge" && (strings.HasPrefix(entry.DatasetID.String(), "ibge.pnad-continua-") ||
		strings.HasPrefix(entry.DatasetID.String(), "ibge.pnad-rural-")) {
		body, sourceURL, err := ibgeClient.FetchPNADRuralSnapshot(ctx, entry, ibge.PNADFetchOptions{
			UFs: opts.UFs,
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

	if agency == "inpe" {
		body, sourceURL, err := inpeClient.FetchDETERSnapshot(ctx, entry)
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

	if agency == "ipea" {
		body, sourceURL, err := ipeaClient.FetchSeriesSnapshot(ctx, entry)
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

	if agency == "cepea" {
		body, sourceURL, err := cepeaClient.FetchIndicadorSnapshot(ctx, entry, opts.FromDate)
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

	if agency == "mdic" {
		body, sourceURL, err := mdicClient.FetchComexSnapshot(ctx, entry, opts.FromDate)
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

	if agency == "b3" {
		body, sourceURL, err := b3Client.FetchFuturoSnapshot(ctx, entry, opts.FromDate)
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

	if agency == "usda" {
		if entry.DatasetID.String() == "usda.wasde" {
			body, sourceURL, err := usdaClient.FetchWASDESnapshot(ctx, entry, opts.FromDate)
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
		if entry.DatasetID.String() == "usda.gats-trade" {
			body, sourceURL, err := usdaClient.FetchGATSSnapshot(ctx, entry, opts.FromDate)
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
		body, sourceURL, err := usdaClient.FetchPSDSnapshot(ctx, entry, opts.FromDate)
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

	if agency == "fao" {
		body, sourceURL, err := faoClient.FetchSnapshot(ctx, entry, opts.FromDate)
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

	if agency == "worldbank" {
		body, sourceURL, err := worldbankClient.FetchSnapshot(ctx, entry, opts.FromDate)
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

	if agency == "noaa" {
		body, sourceURL, err := noaaClient.FetchClimateSnapshot(ctx, entry, opts.FromDate)
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

	if agency == "eia" {
		body, sourceURL, err := eiaClient.FetchPetroleumSnapshot(ctx, entry, opts.FromDate)
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

	if agency == "igc" {
		body, sourceURL, err := igcClient.FetchGOISnapshot(ctx, entry)
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

	if agency == "ana" {
		if entry.DatasetID.String() == "ana.pluviometria-redes" {
			body, sourceURL, err := anaClient.FetchPluviometriaSnapshot(ctx, entry, ana.HidrologiaFetchOptions{
				DataInicio: opts.FromDate,
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
		body, sourceURL, err := anaClient.FetchHidrologiaSnapshot(ctx, entry, ana.HidrologiaFetchOptions{
			DataInicio: opts.FromDate,
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

	if agency == "eurostat" {
		body, sourceURL, err := eurostatClient.FetchAgPricesSnapshot(ctx, entry, opts.FromDate)
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

	if agency == "argentina" {
		if entry.DatasetID.String() == "argentina.magyp-producion-granos" {
			body, sourceURL, err := argentinaClient.FetchGranosSnapshot(ctx, entry, opts.FromDate)
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
		body, sourceURL, err := argentinaClient.FetchCambioSnapshot(ctx, entry, opts.FromDate)
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

	if agency == "oecd-fao" {
		body, sourceURL, err := oecdClient.FetchAgOutlookSnapshot(ctx, entry)
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

	if agency == "cftc" {
		body, sourceURL, err := cftcClient.FetchCOTAgriculturalSnapshot(ctx, entry)
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

	if agency == "jrc" {
		body, sourceURL, err := jrcClient.FetchMARSCropYieldSnapshot(ctx, entry)
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

	if agency == "wto" {
		body, sourceURL, err := wtoClient.FetchITSTradeSnapshot(ctx, entry)
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

	if agency == "fred" {
		body, sourceURL, err := fredClient.FetchCommodityIndexesSnapshot(ctx, entry)
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

	if agency == "nasa" {
		body, sourceURL, err := nasaClient.FetchPOWERAgroSnapshot(ctx, entry)
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

	if agency == "sagis" {
		body, sourceURL, err := sagisClient.FetchGrainSupplySnapshot(ctx, entry)
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

	if agency == "embrapa" {
		body, sourceURL, err := embrapaClient.FetchAgroAPIAgrofitSnapshot(ctx, entry)
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

	if agency == "ibama" {
		switch entry.DatasetID.String() {
		case "ibama.sisfogo-incendios":
			body, sourceURL, err := ibamaClient.FetchSISFOGOSnapshot(ctx, entry)
			if err != nil {
				return nil, err
			}
			return &SourceDownload{
				Body:          body,
				ContentType:   "text/csv",
				ContentLength: int64(len(body)),
				SourceURL:     sourceURL,
			}, nil
		case "ibama.licencas-ambientais":
			body, sourceURL, err := ibamaClient.FetchLicencasSnapshot(ctx, entry)
			if err != nil {
				return nil, err
			}
			return &SourceDownload{
				Body:          body,
				ContentType:   "text/csv",
				ContentLength: int64(len(body)),
				SourceURL:     sourceURL,
			}, nil
		case "ibama.autos-infracao":
			body, sourceURL, err := ibamaClient.FetchAutosSnapshot(ctx, entry, opts.FromYear)
			if err != nil {
				return nil, err
			}
			return &SourceDownload{
				Body:          body,
				ContentType:   "text/csv",
				ContentLength: int64(len(body)),
				SourceURL:     sourceURL,
			}, nil
		}
	}

	if agency == "japan" {
		body, sourceURL, err := japanClient.FetchMAFFAgTradeSnapshot(ctx, entry)
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

	if agency == "mexico" {
		body, sourceURL, err := mexicoClient.FetchSIAPProduccionSnapshot(ctx, entry)
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

	if agency == "copernicus" {
		body, sourceURL, err := copernicusClient.FetchERA5AgroclimateSnapshot(ctx, entry)
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

	if agency == "un" {
		body, sourceURL, err := unClient.FetchComtradeSnapshot(ctx, entry, opts.FromDate)
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
	case "antt":
		result, err := anttClient.Download(ctx, sourceURL)
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
	case "aneel":
		result, err := aneelClient.Download(ctx, sourceURL)
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
	case "bndes":
		result, err := bndesClient.Download(ctx, sourceURL)
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
	case "antaq":
		result, err := antaqClient.Download(ctx, sourceURL)
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
	case "mapa":
		if strings.HasPrefix(entry.DatasetID.String(), "mapa.sipeagro-") {
			body, sourceURL, err := mapaClient.FetchSIPEAGROSnapshot(ctx, entry)
			if err != nil {
				return nil, err
			}
			return &SourceDownload{
				Body:          body,
				ContentType:   "text/csv",
				ContentLength: int64(len(body)),
				SourceURL:     sourceURL,
			}, nil
		}
		if entry.DatasetID.String() == "mapa.sisser-seguro-rural" {
			body, sourceURL, err := mapaClient.FetchSISSERSnapshot(ctx, entry)
			if err != nil {
				return nil, err
			}
			return &SourceDownload{
				Body:          body,
				ContentType:   "text/csv",
				ContentLength: int64(len(body)),
				SourceURL:     sourceURL,
			}, nil
		}
		result, err := mapaClient.Download(ctx, sourceURL)
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
	case "dnit":
		result, err := dnitClient.Download(ctx, sourceURL)
		if err != nil {
			return nil, err
		}
		body, err := dnit.PrepareCSV(result.Body)
		if err != nil {
			return nil, err
		}
		return &SourceDownload{
			Body:          body,
			ContentType:   result.ContentType,
			LastModified:  result.LastModified,
			ContentLength: int64(len(body)),
			SourceURL:     sourceURL,
		}, nil
	case "transportes":
		if entry.DatasetID.String() == "transportes.mtr-bit-malha-shapefile" {
			body, sourceURL, err := transportesClient.FetchMTRBITShapefileSnapshot(ctx, entry)
			if err != nil {
				return nil, err
			}
			return &SourceDownload{
				Body:          body,
				ContentType:   "text/csv",
				ContentLength: int64(len(body)),
				SourceURL:     sourceURL,
			}, nil
		}
		result, err := transportesClient.Download(ctx, sourceURL)
		if err != nil {
			return nil, err
		}
		body, err := transportes.PrepareCSV(result.Body)
		if err != nil {
			return nil, err
		}
		return &SourceDownload{
			Body:          body,
			ContentType:   result.ContentType,
			LastModified:  result.LastModified,
			ContentLength: int64(len(body)),
			SourceURL:     sourceURL,
		}, nil
	case "suframa":
		result, err := suframaClient.Download(ctx, sourceURL)
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
	case "ons":
		result, err := onsClient.Download(ctx, sourceURL)
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
		body, sourceURL, err := client.FetchStationCatalog(ctx, datasetID)
		if err != nil {
			return nil, err
		}
		return &SourceDownload{
			Body:          body,
			ContentType:   "text/csv",
			ContentLength: int64(len(body)),
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
	case "inmet.sequia-monitor":
		body, sourceURL, err := client.FetchSecaMonitorSnapshot(ctx, entry.SourceURL)
		if err != nil {
			return nil, err
		}
		return &SourceDownload{
			Body:          body,
			ContentType:   "application/json",
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
