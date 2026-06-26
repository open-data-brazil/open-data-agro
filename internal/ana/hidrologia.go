package ana

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

type hidrologiaRow struct {
	StationCode         string `json:"station_code"`
	ConsistencyLevel    string `json:"consistency_level"`
	DataType            string `json:"data_type"`
	ObservedAt          string `json:"observed_at"`
	DailyMean           string `json:"daily_mean"`
	AcquisitionMethod   string `json:"acquisition_method,omitempty"`
	MaxValue            string `json:"max_value,omitempty"`
	MinValue            string `json:"min_value,omitempty"`
	MeanValue           string `json:"mean_value,omitempty"`
}

type serieHistoricaXML struct {
	EstacaoCodigo         string `xml:"EstacaoCodigo"`
	NivelConsistencia     string `xml:"NivelConsistencia"`
	DataHora              string `xml:"DataHora"`
	MediaDiaria           string `xml:"MediaDiaria"`
	MetodoObtencaoVazoes  string `xml:"MetodoObtencaoVazoes"`
	Maxima                string `xml:"Maxima"`
	Minima                string `xml:"Minima"`
	Media                 string `xml:"Media"`
}

type dataTableXML struct {
	Records []serieHistoricaXML `xml:"diffgram>DocumentElement>SerieHistorica"`
}

// FetchHidrologiaSnapshot downloads ANA HidroWeb daily series for configured stations.
func (c *Client) FetchHidrologiaSnapshot(ctx context.Context, entry catalog.RegistryEntry, opts HidrologiaFetchOptions) ([]byte, string, error) {
	stations := entry.ANAStationCodes
	if len(stations) == 0 {
		return nil, "", fmt.Errorf("dataset %s has no ana_station_codes", entry.DatasetID)
	}

	tipoDados := strings.TrimSpace(entry.ANATipoDados)
	if tipoDados == "" {
		tipoDados = "3"
	}
	nivel := strings.TrimSpace(entry.ANANivelConsistencia)
	if nivel == "" {
		nivel = "2"
	}

	dataInicio := strings.TrimSpace(opts.DataInicio)
	if dataInicio == "" {
		dataInicio = strings.TrimSpace(entry.StartDate)
	}
	if dataInicio == "" {
		dataInicio = "01/01/2024"
	}
	dataFim := strings.TrimSpace(opts.DataFim)
	if dataFim == "" {
		dataFim = time.Now().Format("02/01/2006")
	}

	baseURL := strings.TrimSpace(entry.SourceURL)
	if baseURL == "" {
		baseURL = defaultServiceBaseURL
	}

	var merged []hidrologiaRow
	var firstURL string
	for _, station := range stations {
		station = strings.TrimSpace(station)
		if station == "" {
			continue
		}
		query := url.Values{
			"codEstacao":         {station},
			"dataInicio":         {dataInicio},
			"dataFim":            {dataFim},
			"tipoDados":          {tipoDados},
			"nivelConsistencia":  {nivel},
		}
		sourceURL := baseURL + "?" + query.Encode()
		if firstURL == "" {
			firstURL = sourceURL
		}

		raw, err := c.download(ctx, sourceURL)
		if err != nil {
			return nil, "", fmt.Errorf("fetch ana station %s: %w", station, err)
		}

		rows, err := parseHidrologiaXML(raw, tipoDados)
		if err != nil {
			return nil, "", err
		}
		merged = append(merged, rows...)
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("no ana hidrologia rows for %s", entry.DatasetID)
	}

	payload, err := json.Marshal(merged)
	if err != nil {
		return nil, "", err
	}
	return payload, firstURL, nil
}

// HidrologiaFetchOptions controls ANA series date windows.
type HidrologiaFetchOptions struct {
	DataInicio string
	DataFim    string
}

func parseHidrologiaXML(raw []byte, tipoDados string) ([]hidrologiaRow, error) {
	var table dataTableXML
	if err := xml.Unmarshal(raw, &table); err != nil {
		return nil, fmt.Errorf("parse ana hidrologia xml: %w", err)
	}

	rows := make([]hidrologiaRow, 0, len(table.Records))
	for _, rec := range table.Records {
		if strings.TrimSpace(rec.DataHora) == "" {
			continue
		}
		rows = append(rows, hidrologiaRow{
			StationCode:       strings.TrimSpace(rec.EstacaoCodigo),
			ConsistencyLevel:  strings.TrimSpace(rec.NivelConsistencia),
			DataType:          tipoDados,
			ObservedAt:        strings.TrimSpace(rec.DataHora),
			DailyMean:         strings.TrimSpace(rec.MediaDiaria),
			AcquisitionMethod: strings.TrimSpace(rec.MetodoObtencaoVazoes),
			MaxValue:          strings.TrimSpace(rec.Maxima),
			MinValue:          strings.TrimSpace(rec.Minima),
			MeanValue:         strings.TrimSpace(rec.Media),
		})
	}
	return rows, nil
}

// FlattenHidrologia converts merged hydrology JSON into canonical bronze columns.
func FlattenHidrologia(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []hidrologiaRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse ana hidrologia json: %w", err)
	}

	headers := []string{
		"station_code",
		"consistency_level",
		"data_type",
		"observed_at",
		"daily_mean",
		"acquisition_method",
		"max_value",
		"min_value",
		"mean_value",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.ObservedAt) == "" {
			continue
		}
		out = append(out, []string{
			row.StationCode,
			row.ConsistencyLevel,
			row.DataType,
			row.ObservedAt,
			row.DailyMean,
			row.AcquisitionMethod,
			row.MaxValue,
			row.MinValue,
			row.MeanValue,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no ana hidrologia rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

// ResolveURL returns the ANA HidroWeb SOAP endpoint for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return defaultServiceBaseURL, nil
	}
	if !strings.Contains(strings.ToLower(raw), "ana.gov.br") {
		return "", fmt.Errorf("source_url for %s must be on ana.gov.br", entry.DatasetID)
	}
	return raw, nil
}
