package ingest

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/parquet-go/parquet-go"
	"github.com/parquet-go/parquet-go/compress/snappy"
)

// Target ~128MB row groups for DuckDB scans (approximate via row count).
const bronzeMaxRowsPerRowGroup = 262144

// ConvertToParquet transforms in-memory source bytes into bronze parquet.
func ConvertToParquet(entry catalog.RegistryEntry, raw []byte) ([]byte, int, error) {
	return convertToParquet(entry, raw, "")
}

// ConvertToParquetFromFile transforms a staged download file into bronze parquet.
func ConvertToParquetFromFile(entry catalog.RegistryEntry, path string) ([]byte, int, error) {
	return convertToParquet(entry, nil, path)
}

func convertToParquet(entry catalog.RegistryEntry, raw []byte, path string) ([]byte, int, error) {
	switch entry.Format {
	case catalog.FormatCSV, catalog.FormatTXT:
		if path != "" {
			return convertDelimitedFileToParquet(path, entry)
		}
		return convertDelimitedToParquet(raw, delimiterFor(entry))
	case catalog.FormatJSON:
		if path != "" {
			return convertJSONFileToParquet(entry, path)
		}
		return convertJSONToParquet(entry, raw)
	case catalog.FormatXLS:
		if path != "" {
			return convertXLSFileToParquet(entry, path)
		}
		return convertXLSToParquet(entry, raw)
	case catalog.FormatXLSX:
		if path != "" {
			return convertXLSXFileToParquet(entry, path)
		}
		return convertXLSXToParquet(entry, raw)
	default:
		return nil, 0, fmt.Errorf("unsupported format %q for %s", entry.Format, entry.DatasetID)
	}
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
	writer := parquet.NewGenericWriter[map[string]any](buf, schema,
		parquet.Compression(&snappy.Codec{}),
		parquet.MaxRowsPerRowGroup(bronzeMaxRowsPerRowGroup),
	)

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
