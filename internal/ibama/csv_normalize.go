package ibama

import (
	"bytes"
	"io"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// NormalizeCSV converts IBAMA portal CSV bytes to UTF-8 and strips a UTF-8 BOM.
func NormalizeCSV(raw []byte) []byte {
	raw = stripUTF8BOM(raw)
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

func stripUTF8BOM(raw []byte) []byte {
	if len(raw) >= 3 && raw[0] == 0xEF && raw[1] == 0xBB && raw[2] == 0xBF {
		return raw[3:]
	}
	return raw
}
