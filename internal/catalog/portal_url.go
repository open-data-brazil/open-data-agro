package catalog

import "strings"

// CONABSourcePortalURL is the mandatory portal URL for all CONAB catalog entries.
const CONABSourcePortalURL = "https://portaldeinformacoes.conab.gov.br/download-arquivos.html"

// ANPCombustiveisPortalURL is the ANP weekly fuel price survey listing page.
const ANPCombustiveisPortalURL = "https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas"

// SourcePortalURL returns the portal URL for a registry entry.
func SourcePortalURL(datasetID DatasetID) string {
	if strings.HasPrefix(datasetID.String(), "conab.") {
		return CONABSourcePortalURL
	}
	if strings.HasPrefix(datasetID.String(), "anp.") {
		return ANPCombustiveisPortalURL
	}
	if strings.HasPrefix(datasetID.String(), "ibge.pam-") {
		return IBGESIDRAPAMURL
	}
	if strings.HasPrefix(datasetID.String(), "ibge.lspa-") {
		return IBGESIDRALSPAURL
	}
	if strings.HasPrefix(datasetID.String(), "ibge.") {
		return IBGELocalidadesDocsURL
	}
	if strings.HasPrefix(datasetID.String(), "inmet.") {
		return INMETBDMEPPortalURL
	}
	if strings.HasPrefix(datasetID.String(), "bcb.") {
		return BCBDadosAbertosURL
	}
	if strings.HasPrefix(datasetID.String(), "cepea.") {
		return CEPEAPortalURL
	}
	if strings.HasPrefix(datasetID.String(), "mdic.") {
		return MDICComexStatPortalURL
	}
	if strings.HasPrefix(datasetID.String(), "antt.") {
		return ANTTDadosAbertosPortalURL
	}
	if strings.HasPrefix(datasetID.String(), "mapa.") {
		return MAPADadosAbertosPortalURL
	}
	if strings.HasPrefix(datasetID.String(), "b3.") {
		return B3MarketDataPortalURL
	}
	if strings.HasPrefix(datasetID.String(), "usda.") {
		return USDAFASPSDPortalURL
	}
	if strings.HasPrefix(datasetID.String(), "fao.") {
		return FAOFaostatPortalURL
	}
	if strings.HasPrefix(datasetID.String(), "worldbank.") {
		return WorldBankCommodityMarketsURL
	}
	if strings.HasPrefix(datasetID.String(), "noaa.") {
		return NOAANCEIPortalURL
	}
	if strings.HasPrefix(datasetID.String(), "eia.") {
		return EIAOpenDataPortalURL
	}
	return ""
}

// NOAANCEIPortalURL is the NOAA NCEI climate monitoring portal URL.
const NOAANCEIPortalURL = "https://www.ncei.noaa.gov/access/monitoring/"

// EIAOpenDataPortalURL is the U.S. EIA Open Data portal URL.
const EIAOpenDataPortalURL = "https://www.eia.gov/opendata/"

// WorldBankCommodityMarketsURL is the World Bank commodity markets portal URL.
const WorldBankCommodityMarketsURL = "https://www.worldbank.org/en/research/commodity-markets"

// FAOFaostatPortalURL is the FAO FAOSTAT data portal URL.
const FAOFaostatPortalURL = "https://www.fao.org/faostat/en/#data"

// USDAFASPSDPortalURL is the USDA FAS PSD Online portal URL.
const USDAFASPSDPortalURL = "https://apps.fas.usda.gov/psdonline/"

// B3MarketDataPortalURL is the B3 market data portal URL.
const B3MarketDataPortalURL = "https://www.b3.com.br/pt_br/market-data-e-indices/servicos-de-dados/market-data/historico/boletins-diarios/pesquisa-por-pregao/"

// MAPADadosAbertosPortalURL is the MAPA open data portal URL.
const MAPADadosAbertosPortalURL = "https://dados.agricultura.gov.br/"

// ANTTDadosAbertosPortalURL is the ANTT open data portal URL.
const ANTTDadosAbertosPortalURL = "https://dados.antt.gov.br/"

// MDICComexStatPortalURL is the MDIC Comex Stat portal and API documentation URL.
const MDICComexStatPortalURL = "https://comexstat.mdic.gov.br/"

// CEPEAPortalURL is the CEPEA/ESALQ-USP indicators portal URL.
const CEPEAPortalURL = "https://www.cepea.org.br/"

// BCBDadosAbertosURL is the BCB open data portal URL.
const BCBDadosAbertosURL = "https://dadosabertos.bcb.gov.br/"

// IBGELocalidadesDocsURL is the IBGE Localidades API documentation URL.
const IBGELocalidadesDocsURL = "https://servicodados.ibge.gov.br/api/docs/localidades"

// IBGESIDRAPAMURL is the IBGE SIDRA PAM survey portal URL.
const IBGESIDRAPAMURL = "https://sidra.ibge.gov.br/pesquisa/pam"

// IBGESIDRALSPAURL is the IBGE SIDRA LSPA survey portal URL.
const IBGESIDRALSPAURL = "https://sidra.ibge.gov.br/pesquisa/lspa"

// INMETBDMEPPortalURL is the INMET BDMEP historical data portal URL.
const INMETBDMEPPortalURL = "https://bdmep.inmet.gov.br/"
