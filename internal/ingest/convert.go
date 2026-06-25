package ingest

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/parquet-go/parquet-go"
	"github.com/xuri/excelize/v2"
)

// ConvertToParquet transforms a source file into bronze parquet bytes preserving column names.
func ConvertToParquet(entry catalog.RegistryEntry, raw []byte) ([]byte, int, error) {
	switch entry.Format {
	case catalog.FormatCSV, catalog.FormatTXT:
		return convertDelimitedToParquet(raw, delimiterFor(entry))
	case catalog.FormatXLSX:
		return convertXLSXToParquet(raw)
	default:
		return nil, 0, fmt.Errorf("unsupported format %q for %s", entry.Format, entry.DatasetID)
	}
}

func delimiterFor(entry catalog.RegistryEntry) rune {
	if entry.Delimiter == "" {
		return ';'
	}
	runes := []rune(entry.Delimiter)
	if len(runes) == 0 {
		return ';'
	}
	return runes[0]
}

func convertDelimitedToParquet(raw []byte, delimiter rune) ([]byte, int, error) {
	reader := csv.NewReader(bytes.NewReader(raw))
	reader.Comma = delimiter
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		return nil, 0, fmt.Errorf("read header: %w", err)
	}
	headers = normalizeHeaders(headers)
	if len(headers) == 0 {
		return nil, 0, fmt.Errorf("empty header row")
	}

	var rows [][]string
	for {
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, 0, fmt.Errorf("read row: %w", readErr)
		}
		rows = append(rows, padRecord(record, len(headers)))
	}

	return writeStringTable(headers, rows)
}

func convertXLSXToParquet(raw []byte) ([]byte, int, error) {
	book, err := excelize.OpenReader(bytes.NewReader(raw))
	if err != nil {
		return nil, 0, fmt.Errorf("open workbook: %w", err)
	}
	defer func() { _ = book.Close() }()

	sheet := book.GetSheetName(0)
	if sheet == "" {
		return nil, 0, fmt.Errorf("workbook has no sheets")
	}

	table, err := book.GetRows(sheet)
	if err != nil {
		return nil, 0, fmt.Errorf("read sheet %q: %w", sheet, err)
	}
	if len(table) == 0 {
		return nil, 0, fmt.Errorf("sheet %q is empty", sheet)
	}

	headers := normalizeHeaders(table[0])
	var rows [][]string
	for _, record := range table[1:] {
		rows = append(rows, padRecord(record, len(headers)))
	}

	return writeStringTable(headers, rows)
}

func normalizeHeaders(headers []string) []string {
	out := make([]string, len(headers))
	seen := make(map[string]int)
	for i, header := range headers {
		name := strings.TrimSpace(header)
		if name == "" {
			name = fmt.Sprintf("column_%d", i+1)
		}
		if count, ok := seen[name]; ok {
			seen[name] = count + 1
			name = fmt.Sprintf("%s_%d", name, count+1)
		} else {
			seen[name] = 1
		}
		out[i] = name
	}
	return out
}

func padRecord(record []string, width int) []string {
	if len(record) >= width {
		return record[:width]
	}
	out := make([]string, width)
	copy(out, record)
	return out
}

func writeStringTable(headers []string, rows [][]string) ([]byte, int, error) {
	group := parquet.Group{}
	for _, header := range headers {
		group[header] = parquet.String()
	}
	schema := parquet.NewSchema("bronze", group)

	buf := new(bytes.Buffer)
	writer := parquet.NewGenericWriter[map[string]any](buf, schema)

	batch := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		item := make(map[string]any, len(headers))
		for i, header := range headers {
			item[header] = row[i]
		}
		batch = append(batch, item)
	}

	if _, err := writer.Write(batch); err != nil {
		return nil, 0, fmt.Errorf("write parquet rows: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, 0, fmt.Errorf("close parquet writer: %w", err)
	}

	return buf.Bytes(), len(rows), nil
}
