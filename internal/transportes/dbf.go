package transportes

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// ParseDBFAttributes reads DBF field values into CSV bytes (no geometry).
func ParseDBFAttributes(dbfBytes []byte) ([]byte, error) {
	fields, records, err := parseDBF(dbfBytes)
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("dbf has no fields")
	}

	headers := make([]string, len(fields))
	for i, field := range fields {
		headers[i] = field.name
	}

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type dbfField struct {
	name string
	size uint8
}

func parseDBF(raw []byte) ([]dbfField, [][]string, error) {
	if len(raw) < 32 {
		return nil, nil, fmt.Errorf("dbf header too short")
	}
	numRecords := int(uint32(raw[4]) | uint32(raw[5])<<8 | uint32(raw[6])<<16 | uint32(raw[7])<<24)
	headerSize := int(uint16(raw[8]) | uint16(raw[9])<<8)
	recordSize := int(uint16(raw[10]) | uint16(raw[11])<<8)
	if headerSize <= 32 || recordSize <= 1 || headerSize > len(raw) {
		return nil, nil, fmt.Errorf("invalid dbf header sizes")
	}

	pos := 32
	var fields []dbfField
	for pos < headerSize-1 {
		if raw[pos] == 0x0D {
			break
		}
		name := strings.TrimSpace(string(raw[pos : pos+11]))
		size := raw[pos+16]
		fields = append(fields, dbfField{name: name, size: size})
		pos += 32
	}
	if len(fields) == 0 {
		return nil, nil, fmt.Errorf("dbf has no field descriptors")
	}

	dataStart := headerSize
	rows := make([][]string, 0, numRecords)
	for i := 0; i < numRecords; i++ {
		start := dataStart + i*recordSize
		end := start + recordSize
		if end > len(raw) {
			break
		}
		rec := raw[start:end]
		if len(rec) == 0 || rec[0] == 0x2A {
			continue
		}
		off := 1
		row := make([]string, len(fields))
		for j, field := range fields {
			chunk := rec[off : off+int(field.size)]
			row[j] = strings.TrimSpace(decodeDBFText(chunk))
			off += int(field.size)
		}
		rows = append(rows, row)
	}
	return fields, rows, nil
}

func decodeDBFText(chunk []byte) string {
	if utf8Valid(chunk) {
		return string(chunk)
	}
	return string(chunk)
}

func utf8Valid(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] >= 0x80 {
			return false
		}
	}
	return true
}

func ioReadAllLimit(r io.Reader, limit int64) ([]byte, error) {
	return io.ReadAll(io.LimitReader(r, limit))
}
