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
	return ""
}

// CEPEAPortalURL is the CEPEA/ESALQ-USP indicators portal URL.
const CEPEAPortalURL = "https://www.cepea.org.br/"

// BCBDadosAbertosURL is the BCB open data portal URL.
const BCBDadosAbertosURL = "https://dadosabertos.bcb.gov.br/"

// IBGELocalidadesDocsURL is the IBGE Localidades API documentation URL.
const IBGELocalidadesDocsURL = "https://servicodados.ibge.gov.br/api/docs/localidades"

// IBGESIDRAPAMURL is the IBGE SIDRA PAM survey portal URL.
const IBGESIDRAPAMURL = "https://sidra.ibge.gov.br/pesquisa/pam"

// INMETBDMEPPortalURL is the INMET BDMEP historical data portal URL.
const INMETBDMEPPortalURL = "https://bdmep.inmet.gov.br/"
