package ingest

import (
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/conab"
)

// Citation is required source attribution on exported bronze datasets.
type Citation struct {
	FonteOficial  string `json:"fonteOficial"`
	Agencia       string `json:"agencia"`
	DatasetPortal string `json:"datasetPortal"`
}

// CitationFor builds the official source citation for a registry entry.
func CitationFor(entry catalog.RegistryEntry) Citation {
	return Citation{
		FonteOficial:  conab.PortalDownloadPage,
		Agencia:       "CONAB",
		DatasetPortal: entry.PortalLabel,
	}
}
