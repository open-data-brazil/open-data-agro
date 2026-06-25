package catalog

import "time"

// DatasetMetadata captures metadata for an embedded or ingested dataset.
type DatasetMetadata struct {
	DatasetID      DatasetID `json:"datasetId" yaml:"dataset_id"`
	CapturadoEm    time.Time `json:"capturadoEm" yaml:"capturado_em"`
	FonteOficial   string    `json:"fonteOficial" yaml:"fonte_oficial"`
	VersaoFonte    string    `json:"versaoFonte,omitempty" yaml:"versao_fonte,omitempty"`
	TotalRegistros int       `json:"totalRegistros" yaml:"total_registros"`
}
