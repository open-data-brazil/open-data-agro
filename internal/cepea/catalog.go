package cepea

import "fmt"

const (
	PortalURL        = "https://www.cepea.org.br/"
	LicenseURL       = "https://www.cepea.org.br/br/licenca-de-uso-de-dados.aspx"
	DefaultLicense   = "CC BY-NC 4.0"
	DefaultFonteTipo = "referencia_mercado"
	MirrorBaseURL    = "https://www.noticiasagricolas.com.br/cotacoes"
)

// IndicatorSpec maps a catalog dataset to CEPEA product and trading location.
type IndicatorSpec struct {
	ProductSlug string
	Praca       string
	MirrorPath  string
	PortalPath  string
	Unidade     string
}

var indicatorSpecs = map[string]IndicatorSpec{
	"cepea.soja-paranagua": {
		ProductSlug: "soja",
		Praca:       "Paranaguá",
		MirrorPath:  "soja/soja-indicador-cepea-esalq-porto-paranagua",
		PortalPath:  "soja.aspx",
		Unidade:     "R$/saca 60 kg",
	},
	"cepea.soja-parana": {
		ProductSlug: "soja",
		Praca:       "Paraná",
		MirrorPath:  "soja/indicador-cepea-esalq-soja-parana",
		PortalPath:  "soja.aspx",
		Unidade:     "R$/saca 60 kg",
	},
	"cepea.milho": {
		ProductSlug: "milho",
		Praca:       "Campinas",
		MirrorPath:  "milho/indicador-cepea-esalq-milho",
		PortalPath:  "milho.aspx",
		Unidade:     "R$/saca 60 kg",
	},
	"cepea.boi-gordo": {
		ProductSlug: "boi-gordo",
		Praca:       "São Paulo",
		MirrorPath:  "boi-gordo/boi-gordo-indicador-esalq-bmf",
		PortalPath:  "boi-gordo.aspx",
		Unidade:     "R$/@",
	},
}

// IndicatorSpecFor returns the CEPEA indicator mapping for a dataset ID.
func IndicatorSpecFor(datasetID string) (IndicatorSpec, error) {
	spec, ok := indicatorSpecs[datasetID]
	if !ok {
		return IndicatorSpec{}, fmt.Errorf("unknown cepea dataset %q", datasetID)
	}
	return spec, nil
}

// MirrorURL builds the Notícias Agrícolas mirror URL for a dataset.
func MirrorURL(datasetID string) (string, error) {
	spec, err := IndicatorSpecFor(datasetID)
	if err != nil {
		return "", err
	}
	return MirrorBaseURL + "/" + spec.MirrorPath, nil
}

// PortalIndicatorURL builds the official CEPEA indicator page URL.
func PortalIndicatorURL(datasetID string) (string, error) {
	spec, err := IndicatorSpecFor(datasetID)
	if err != nil {
		return "", err
	}
	return PortalURL + "br/indicador/" + spec.PortalPath, nil
}
