package inmet

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	stationAPIAutomatic    = "https://apitempo.inmet.gov.br/estacoes/T"
	stationAPIConventional = "https://apitempo.inmet.gov.br/estacoes/M"
)

type stationAPIRecord struct {
	StationCode   string  `json:"CD_ESTACAO"`
	Name          string  `json:"DC_NOME"`
	State         string  `json:"SG_ESTADO"`
	Latitude      string  `json:"VL_LATITUDE"`
	Longitude     string  `json:"VL_LONGITUDE"`
	Altitude      string  `json:"VL_ALTITUDE"`
	Situation     string  `json:"CD_SITUACAO"`
	OperationFrom *string `json:"DT_INICIO_OPERACAO"`
}

// FetchStationCatalog downloads the live station catalog from apitempo.inmet.gov.br
// and returns legacy semicolon CSV bytes compatible with FlattenEstacoes.
func (c *Client) FetchStationCatalog(ctx context.Context, datasetID string) ([]byte, string, error) {
	apiURL, err := stationCatalogAPIURL(datasetID)
	if err != nil {
		return nil, "", err
	}

	result, err := c.Download(ctx, apiURL)
	if err != nil {
		return nil, "", err
	}

	payload, err := stationsJSONToCatalogCSV(datasetID, result.Body)
	if err != nil {
		return nil, "", err
	}
	return payload, apiURL, nil
}

func stationCatalogAPIURL(datasetID string) (string, error) {
	switch datasetID {
	case "inmet.estacoes-automaticas":
		return stationAPIAutomatic, nil
	case "inmet.estacoes-convencionais":
		return stationAPIConventional, nil
	default:
		return "", fmt.Errorf("unsupported station catalog dataset %s", datasetID)
	}
}

func stationsJSONToCatalogCSV(datasetID string, raw []byte) ([]byte, error) {
	var records []stationAPIRecord
	if err := json.Unmarshal(raw, &records); err != nil {
		return nil, fmt.Errorf("decode station catalog json: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("empty station catalog for %s", datasetID)
	}

	headers := []string{
		"REGIAO", "UF", "ESTACAO", "CODIGO WMO", "LATITUDE", "LONGITUDE",
		"ALTITUDE", "DATA DE FUNDACAO", "DATA DE INSTALACAO", "SITUACAO",
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ';'
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	written := 0
	for _, record := range records {
		code := strings.TrimSpace(record.StationCode)
		if code == "" {
			continue
		}
		if datasetID == "inmet.estacoes-automaticas" && strings.HasPrefix(strings.ToUpper(code), "S") {
			continue
		}

		installDate := formatStationInstallDate(record.OperationFrom)
		row := []string{
			ufToRegion(record.State),
			strings.TrimSpace(record.State),
			strings.TrimSpace(record.Name),
			code,
			formatStationDecimal(record.Latitude),
			formatStationDecimal(record.Longitude),
			formatStationDecimal(record.Altitude),
			installDate,
			installDate,
			strings.ToUpper(strings.TrimSpace(record.Situation)),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
		written++
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	if written == 0 {
		return nil, fmt.Errorf("no station rows mapped for %s", datasetID)
	}
	return buf.Bytes(), nil
}

func formatStationDecimal(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return strings.ReplaceAll(value, ".", ",")
}

func formatStationInstallDate(raw *string) string {
	if raw == nil {
		return ""
	}
	value := strings.TrimSpace(*raw)
	if value == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		if len(value) >= 10 {
			parts := strings.Split(value[:10], "-")
			if len(parts) == 3 {
				return parts[2] + "/" + parts[1] + "/" + parts[0]
			}
		}
		return value
	}
	return parsed.Format("02/01/2006")
}

func ufToRegion(uf string) string {
	switch strings.ToUpper(strings.TrimSpace(uf)) {
	case "AC", "AP", "AM", "PA", "RO", "RR", "TO":
		return "Norte"
	case "AL", "BA", "CE", "MA", "PB", "PE", "PI", "RN", "SE":
		return "Nordeste"
	case "DF", "GO", "MS", "MT":
		return "Centro-Oeste"
	case "ES", "MG", "RJ", "SP":
		return "Sudeste"
	case "PR", "RS", "SC":
		return "Sul"
	default:
		return ""
	}
}
