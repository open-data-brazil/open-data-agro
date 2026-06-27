package usda

import (
	"fmt"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// PSDProbeRequest describes a SOAP POST health check for PSD datasets.
type PSDProbeRequest struct {
	URL     string
	Headers map[string]string
	Body    string
}

// BuildPSDProbeRequest returns a SOAP POST probe matching ingest fetch.
func BuildPSDProbeRequest(entry catalog.RegistryEntry) (PSDProbeRequest, error) {
	code := strings.TrimSpace(entry.PSDCommodityCode)
	if code == "" {
		return PSDProbeRequest{}, fmt.Errorf("dataset %s missing psd_commodity_code", entry.DatasetID)
	}

	year := time.Now().UTC().Year() - 1
	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <getDatabyCommodityPerYear xmlns="http://www.fas.usda.gov/wsfaspsd/">
      <strCommodityCode>%s</strCommodityCode>
      <strYear>%d</strYear>
    </getDatabyCommodityPerYear>
  </soap:Body>
</soap:Envelope>`, code, year)

	return PSDProbeRequest{
		URL: psdSOAPEndpoint,
		Headers: map[string]string{
			"Content-Type": "text/xml; charset=utf-8",
			"SOAPAction":   soapActionPerYear,
			"Accept":       "*/*",
		},
		Body: envelope,
	}, nil
}
