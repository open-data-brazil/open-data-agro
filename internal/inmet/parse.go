package inmet

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	datePatterns = []string{
		"2006-01-02",
		"2006/01/02",
		"02/01/2006",
	}
	variablePatterns = []struct {
		pattern *regexp.Regexp
		name    string
		agg     string // sum or mean
	}{
		{regexp.MustCompile(`(?i)precipita`), "precipitacao", "sum"},
		{regexp.MustCompile(`(?i)temperatura do ar`), "temperatura_ar", "mean"},
		{regexp.MustCompile(`(?i)temperatura m[aá]x`), "temperatura_maxima", "mean"},
		{regexp.MustCompile(`(?i)temperatura m[ií]n`), "temperatura_minima", "mean"},
		{regexp.MustCompile(`(?i)umidade relativa do ar`), "umidade_relativa", "mean"},
		{regexp.MustCompile(`(?i)vento, velocidade`), "vento_velocidade", "mean"},
		{regexp.MustCompile(`(?i)radia`), "radiacao", "mean"},
		{regexp.MustCompile(`(?i)press[aã]o atmosf[eé]rica ao n[ií]vel`), "pressao_atmosferica", "mean"},
	}
)

type stationMeta struct {
	Code      string
	Name      string
	Region    string
	State     string
	Latitude  string
	Longitude string
	Situation string
}

func parseStationHeader(lines []string) stationMeta {
	meta := stationMeta{}
	for _, line := range lines {
		parts := strings.Split(line, ";")
		if len(parts) < 2 {
			continue
		}
		label := strings.ToUpper(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		switch {
		case strings.Contains(label, "REGI"):
			meta.Region = value
		case strings.Contains(label, "UF"):
			meta.State = value
		case strings.Contains(label, "ESTA"):
			meta.Name = value
		case strings.Contains(label, "CODIGO") || strings.Contains(label, "WMO"):
			meta.Code = value
		case strings.Contains(label, "LAT"):
			meta.Latitude = normalizeDecimal(value)
		case strings.Contains(label, "LONG"):
			meta.Longitude = normalizeDecimal(value)
		case strings.Contains(label, "SITU"):
			meta.Situation = value
		}
	}
	return meta
}

func normalizeDecimal(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return strings.ReplaceAll(value, ",", ".")
}

func parseNumeric(raw string) (float64, bool) {
	if IsMissingValue(raw) {
		return 0, false
	}
	value := strings.ReplaceAll(strings.TrimSpace(raw), ",", ".")
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return num, true
}

func normalizeDate(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty date")
	}
	for _, layout := range datePatterns {
		t, err := time.Parse(layout, raw)
		if err == nil {
			return t.Format("2006-01-02"), nil
		}
	}
	return "", fmt.Errorf("unparseable date %q", raw)
}

func mapVariableColumn(header string) (string, string, bool) {
	header = strings.TrimSpace(header)
	for _, item := range variablePatterns {
		if item.pattern.MatchString(header) {
			return item.name, item.agg, true
		}
	}
	return "", "", false
}

func readLatinCSV(raw []byte) ([][]string, error) {
	reader := csv.NewReader(bufio.NewReader(bytes.NewReader(raw)))
	reader.Comma = ';'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}

func findDataHeaderIndex(lines []string) int {
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		upper := strings.ToUpper(trimmed)
		if strings.HasPrefix(upper, "DATA;") || strings.HasPrefix(upper, "DATA (") {
			return i
		}
	}
	return -1
}

func splitStationFile(raw []byte) (stationMeta, []byte, error) {
	text := string(raw)
	// Normalize line endings for header scan.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) < 10 {
		return stationMeta{}, nil, fmt.Errorf("station csv too short")
	}

	headerIdx := findDataHeaderIndex(lines)
	if headerIdx < 0 {
		return stationMeta{}, nil, fmt.Errorf("data header not found")
	}

	meta := parseStationHeader(lines[:headerIdx])
	body := []byte(strings.Join(lines[headerIdx:], "\n"))
	return meta, body, nil
}

type dailyAccumulator struct {
	sum   float64
	count int
	agg   string
}

func (a *dailyAccumulator) add(value float64, agg string) {
	a.agg = agg
	if agg == "sum" {
		a.sum += value
		a.count++
		return
	}
	a.sum += value
	a.count++
}

func (a *dailyAccumulator) value() (string, bool) {
	if a.count == 0 {
		return "", false
	}
	if a.agg == "sum" {
		return strconv.FormatFloat(a.sum, 'f', -1, 64), true
	}
	return strconv.FormatFloat(a.sum/float64(a.count), 'f', -1, 64), true
}

func parseStationDailyLong(meta stationMeta, dataRaw []byte, year int, allowedVars map[string]struct{}) ([][]string, error) {
	records, err := readLatinCSV(dataRaw)
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, nil
	}

	headers := records[0]
	colMap := make(map[int]struct {
		name string
		agg  string
	})
	for i, header := range headers {
		name, agg, ok := mapVariableColumn(header)
		if !ok {
			continue
		}
		if len(allowedVars) > 0 {
			if _, ok := allowedVars[name]; !ok {
				continue
			}
		}
		colMap[i] = struct {
			name string
			agg  string
		}{name: name, agg: agg}
	}
	if len(colMap) == 0 {
		return nil, fmt.Errorf("no supported variables in station %s", meta.Code)
	}

	dateCol := -1
	for i, header := range headers {
		if strings.EqualFold(strings.TrimSpace(header), "Data") || strings.Contains(strings.ToUpper(header), "DATA (YYYY") {
			dateCol = i
			break
		}
	}
	if dateCol < 0 {
		return nil, fmt.Errorf("date column not found for station %s", meta.Code)
	}

	type dayKey struct {
		date string
		vars string
	}
	accumulators := make(map[dayKey]*dailyAccumulator)

	for _, record := range records[1:] {
		if dateCol >= len(record) {
			continue
		}
		day, err := normalizeDate(record[dateCol])
		if err != nil {
			continue
		}
		if year > 0 && !strings.HasPrefix(day, strconv.Itoa(year)) {
			continue
		}

		for idx, mapped := range colMap {
			if idx >= len(record) {
				continue
			}
			value, ok := parseNumeric(record[idx])
			if !ok {
				continue
			}
			key := dayKey{date: day, vars: mapped.name}
			acc, exists := accumulators[key]
			if !exists {
				acc = &dailyAccumulator{agg: mapped.agg}
				accumulators[key] = acc
			}
			acc.add(value, mapped.agg)
		}
	}

	rows := make([][]string, 0, len(accumulators))
	yearStr := ""
	if year > 0 {
		yearStr = strconv.Itoa(year)
	}
	for key, acc := range accumulators {
		value, ok := acc.value()
		if !ok {
			continue
		}
		rows = append(rows, []string{
			meta.Code,
			key.date,
			key.vars,
			value,
			meta.State,
			yearStr,
		})
	}
	return rows, nil
}

func writeLongCSV(rows [][]string) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"cd_estacao", "data", "variavel", "valor", "uf", "ano"}); err != nil {
		return nil, err
	}
	if err := writer.WriteAll(rows); err != nil {
		return nil, err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
