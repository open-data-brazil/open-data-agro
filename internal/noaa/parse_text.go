package noaa

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func parseONIASCII(raw []byte, startYear, endYear int) ([]ensoRow, error) {
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	var rows []ensoRow

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "SEAS") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		year, err := strconv.Atoi(fields[1])
		if err != nil || year < startYear || year > endYear {
			continue
		}
		rows = append(rows, ensoRow{
			RefYear:    strconv.Itoa(year),
			SeasonCode: fields[0],
			SSTTotal:   fields[2],
			Anomaly:    fields[3],
			IndexName:  "oni",
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan oni ascii: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no oni rows after filter")
	}
	return rows, nil
}

func parseGlobalTempCSV(raw []byte, startMonth, endMonth string) ([]globalTempRow, error) {
	reader := csv.NewReader(bytes.NewReader(raw))
	reader.Comment = '#'
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	var rows []globalTempRow
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read global temp csv: %w", err)
		}
		if len(record) < 2 {
			continue
		}
		date := strings.TrimSpace(record[0])
		if len(date) != 6 {
			continue
		}
		refmonth, ok := parseYYYYMM(date)
		if !ok {
			continue
		}
		if refmonth < startMonth || refmonth > endMonth {
			continue
		}
		value := strings.TrimSpace(record[1])
		if value == "" {
			continue
		}
		rows = append(rows, globalTempRow{
			RefMonth:   refmonth,
			Anomaly:    value,
			Unit:       "Degrees Celsius",
			BasePeriod: "1901-2000",
			IndexName:  "global_land_ocean",
		})
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no global temp rows after filter")
	}
	return rows, nil
}

func parseYYYYMM(raw string) (string, bool) {
	if len(raw) != 6 {
		return "", false
	}
	year, err := strconv.Atoi(raw[:4])
	if err != nil {
		return "", false
	}
	month, err := strconv.Atoi(raw[4:])
	if err != nil || month < 1 || month > 12 {
		return "", false
	}
	return fmt.Sprintf("%04d-%02d", year, month), true
}
