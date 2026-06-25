package catalog

import "strings"

// CONABSourcePortalURL is the mandatory portal URL for all CONAB catalog entries.
const CONABSourcePortalURL = "https://portaldeinformacoes.conab.gov.br/download-arquivos.html"

// SourcePortalURL returns the portal URL for a registry entry.
func SourcePortalURL(datasetID DatasetID) string {
	if strings.HasPrefix(datasetID.String(), "conab.") {
		return CONABSourcePortalURL
	}
	return ""
}
