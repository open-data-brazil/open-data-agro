package inmet

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var stationCatalogColumns = map[string]string{
	"REGIAO":             "regiao",
	"UF":                 "uf",
	"ESTACAO":            "nome",
	"CODIGO WMO":         "cd_estacao",
	"CODIGO(WMO)":        "cd_estacao",
	"LATITUDE":           "latitude",
	"LONGITUDE":          "longitude",
	"SITUACAO":           "situacao",
	"ALTITUDE":           "altitude",
	"DATA DE FUNDACAO":   "data_fundacao",
	"DATA DE INSTALACAO": "data_instalacao",
}

// FlattenEstacoes converts an INMET station catalog CSV into canonical bronze columns.
func FlattenEstacoes(datasetID string, raw []byte) ([]string, [][]string, error) {
	decoded, err := decodeLatin1(raw)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(bytes.NewReader(decoded))
	reader.Comma = ';'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("read station catalog header: %w", err)
	}

	index := mapCanonicalHeaders(headers)
	outHeaders := []string{"cd_estacao", "nome", "latitude", "longitude", "uf", "situacao"}
	if datasetID == "inmet.estacoes-convencionais" {
		outHeaders = append(outHeaders, "regiao", "altitude")
	}

	var rows [][]string
	for {
		record, readErr := reader.Read()
		if readErr != nil {
			break
		}

		row := make([]string, len(outHeaders))
		for canonical, idx := range index {
			pos := indexOf(outHeaders, canonical)
			if pos < 0 || idx >= len(record) {
				continue
			}
			value := strings.TrimSpace(record[idx])
			if canonical == "latitude" || canonical == "longitude" {
				value = normalizeDecimal(value)
			}
			row[pos] = value
		}

		if strings.TrimSpace(row[0]) == "" {
			continue
		}
		if strings.TrimSpace(row[5]) == "" {
			row[5] = "DESCONHECIDA"
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("no station rows for %s", datasetID)
	}
	return outHeaders, rows, nil
}

func mapCanonicalHeaders(headers []string) map[string]int {
	out := make(map[string]int)
	for i, header := range headers {
		key := strings.ToUpper(strings.TrimSpace(header))
		if canonical, ok := stationCatalogColumns[key]; ok {
			out[canonical] = i
		}
	}
	return out
}

func indexOf(values []string, target string) int {
	for i, value := range values {
		if value == target {
			return i
		}
	}
	return -1
}

func decodeLatin1(raw []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(raw), charmap.ISO8859_1.NewDecoder())
	return io.ReadAll(reader)
}
