package usda

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

func parsePSDSoapResponse(raw []byte) ([]psdRow, error) {
	decoder := xml.NewDecoder(bytes.NewReader(raw))
	var (
		inCommodity bool
		current     psdRow
		pending     string
		rows        []psdRow
	)

	flush := func() {
		if !inCommodity {
			return
		}
		current.CommodityCode = strings.TrimSpace(current.CommodityCode)
		current.CountryCode = strings.TrimSpace(current.CountryCode)
		current.MarketingYear = strings.TrimSpace(current.MarketingYear)
		current.AttributeID = strings.TrimSpace(current.AttributeID)
		if current.CountryCode == "" || current.MarketingYear == "" || current.AttributeID == "" {
			return
		}
		rows = append(rows, current)
	}

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse psd soap xml: %w", err)
		}

		switch elem := token.(type) {
		case xml.StartElement:
			name := localXMLName(elem.Name.Local)
			switch name {
			case "Commodity":
				flush()
				inCommodity = true
				current = psdRow{}
			case "Commodity_code":
				pending = "commodity_code"
			case "Commodity_Description":
				pending = "commodity_name"
			case "Country_Code":
				pending = "country_code"
			case "Country_Name":
				pending = "country_name"
			case "Market_Year":
				pending = "marketing_year"
			case "Calendar_Year":
				pending = "calendar_year"
			case "Month":
				pending = "month"
			case "Attribute_Id":
				pending = "attribute_id"
			case "Attribute_Description":
				pending = "attribute_name"
			case "Unit_Id":
				pending = "unit_id"
			case "Unit_Description":
				pending = "unit_description"
			case "Value":
				pending = "value"
			}
		case xml.EndElement:
			if localXMLName(elem.Name.Local) == "Commodity" {
				flush()
				inCommodity = false
				current = psdRow{}
				pending = ""
			}
		case xml.CharData:
			text := strings.TrimSpace(string(elem))
			if text == "" || pending == "" || !inCommodity {
				continue
			}
			switch pending {
			case "commodity_code":
				current.CommodityCode = text
			case "commodity_name":
				current.CommodityName = strings.TrimSpace(text)
			case "country_code":
				current.CountryCode = text
			case "country_name":
				current.CountryName = strings.TrimSpace(text)
			case "marketing_year":
				current.MarketingYear = text
			case "calendar_year":
				current.CalendarYear = text
			case "month":
				current.Month = text
			case "attribute_id":
				current.AttributeID = text
			case "attribute_name":
				current.AttributeName = strings.TrimSpace(text)
			case "unit_id":
				current.UnitID = text
			case "unit_description":
				current.UnitDescription = strings.TrimSpace(text)
			case "value":
				current.Value = strings.TrimSpace(text)
			}
			pending = ""
		}
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no commodity rows in psd soap response")
	}
	return rows, nil
}

func localXMLName(name string) string {
	if i := strings.LastIndex(name, "}"); i >= 0 {
		return name[i+1:]
	}
	return name
}
