package ingest

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/dnit"
	"github.com/open-data-brazil/open-data-agro/internal/ibama"
	"github.com/open-data-brazil/open-data-agro/internal/transportes"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func convertDelimitedToParquet(raw []byte, delimiter rune) ([]byte, int, error) {
	raw = normalizeTextEncoding(raw)
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
	return convertDelimitedToParquet(preprocessDelimitedRaw(entry, raw), delimiterFor(entry))
}

func preprocessDelimitedRaw(entry catalog.RegistryEntry, raw []byte) []byte {
	agency, _, err := catalog.SplitDatasetID(entry.DatasetID.String())
	if err != nil {
		return raw
	}
	switch agency {
	case "dnit":
		return dnit.StripMetadataRows(raw)
	case "transportes":
		stripped, stripErr := transportes.PrepareCSV(raw)
		if stripErr != nil {
			return raw
		}
		return stripped
	case "ibama":
		return ibama.NormalizeCSV(raw)
	default:
		return raw
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

// normalizeTextEncoding converts ISO-8859-1 CONAB portal downloads to UTF-8 when needed.
func normalizeTextEncoding(raw []byte) []byte {
	if utf8.Valid(raw) {
		return raw
	}
	reader := transform.NewReader(bytes.NewReader(raw), charmap.ISO8859_1.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return raw
	}
	return decoded
}
