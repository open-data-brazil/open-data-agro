package b3

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

const sprdXMLPrefix = "BVBG.187"

func parseSPRDZip(raw []byte, commodityPrefix string) ([]futuroRow, error) {
	xmlBytes, err := extractSPRDXML(raw)
	if err != nil {
		return nil, err
	}
	return parseSPRDXML(xmlBytes, commodityPrefix)
}

func extractSPRDXML(raw []byte) ([]byte, error) {
	xmlBytes, err := findXMLInZip(raw)
	if err == nil {
		return xmlBytes, nil
	}

	outer, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, fmt.Errorf("open sprd zip: %w", err)
	}

	for _, file := range outer.File {
		if !strings.HasSuffix(strings.ToLower(file.Name), ".zip") {
			continue
		}
		innerRaw, err := readZipEntry(file)
		if err != nil {
			return nil, err
		}
		xmlBytes, err := findXMLInZip(innerRaw)
		if err == nil {
			return xmlBytes, nil
		}
	}

	return nil, fmt.Errorf("no BVBG.187 xml in sprd archive")
}

func findXMLInZip(raw []byte) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, err
	}
	for _, file := range reader.File {
		name := strings.ToUpper(file.Name)
		if !strings.Contains(name, sprdXMLPrefix) || !strings.HasSuffix(name, ".XML") {
			continue
		}
		return readZipEntry(file)
	}
	return nil, fmt.Errorf("xml not found")
}

func readZipEntry(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = rc.Close() }()
	return io.ReadAll(io.LimitReader(rc, 128<<20))
}

func parseSPRDXML(raw []byte, commodityPrefix string) ([]futuroRow, error) {
	prefix := strings.ToUpper(strings.TrimSpace(commodityPrefix))
	if prefix == "" {
		return nil, fmt.Errorf("empty commodity prefix")
	}

	decoder := xml.NewDecoder(bytes.NewReader(raw))
	var (
		inReport bool
		current  futuroRow
		rows     []futuroRow
	)

	flush := func() {
		if !inReport || current.Symbol == "" {
			return
		}
		if !strings.HasPrefix(strings.ToUpper(current.Symbol), prefix) {
			return
		}
		current.Commodity = prefix
		current.MaturityCode = strings.TrimPrefix(strings.ToUpper(current.Symbol), prefix)
		rows = append(rows, current)
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse sprd xml: %w", err)
		}

		switch elem := token.(type) {
		case xml.StartElement:
			name := localXMLName(elem.Name.Local)
			switch name {
			case "PricRpt":
				flush()
				inReport = true
				current = futuroRow{}
			case "Dt":
				if inReport {
					current.pendingDate = true
				}
			case "TckrSymb":
				if inReport {
					current.pendingSymbol = true
				}
			case "AdjstdQt":
				if inReport {
					current.Currency = strings.TrimSpace(attrValue(elem.Attr, "Ccy"))
					current.pendingAdj = true
				}
			case "PrvsAdjstdQt":
				if inReport {
					current.pendingPrev = true
				}
			}
		case xml.EndElement:
			if localXMLName(elem.Name.Local) == "PricRpt" {
				flush()
				inReport = false
				current = futuroRow{}
			}
		case xml.CharData:
			text := strings.TrimSpace(string(elem))
			if text == "" || !inReport {
				continue
			}
			switch {
			case current.pendingDate:
				current.RefDate = text
				current.pendingDate = false
			case current.pendingSymbol:
				current.Symbol = strings.ToUpper(text)
				current.pendingSymbol = false
			case current.pendingAdj:
				current.Price = text
				current.pendingAdj = false
			case current.pendingPrev:
				current.PreviousPrice = text
				current.pendingPrev = false
			}
		}
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no %s rows in sprd xml", prefix)
	}
	return rows, nil
}

func localXMLName(name string) string {
	if i := strings.LastIndex(name, "}"); i >= 0 {
		return name[i+1:]
	}
	return name
}

func attrValue(attrs []xml.Attr, key string) string {
	for _, attr := range attrs {
		if localXMLName(attr.Name.Local) == key || attr.Name.Local == key {
			return attr.Value
		}
	}
	return ""
}
