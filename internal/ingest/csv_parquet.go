package ingest

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

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

func convertDelimitedFileToParquet(path string, entry catalog.RegistryEntry) ([]byte, int, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("read staged file: %w", err)
	}
	return convertDelimitedToParquet(raw, delimiterFor(entry))
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
