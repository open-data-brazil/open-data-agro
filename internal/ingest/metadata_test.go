package ingest

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestPartitionMetadataKey(t *testing.T) {
	t.Parallel()

	key, err := PartitionMetadataKey("conab.estimativa-graos", parseDate(t, "2026-06-25"))
	if err != nil {
		t.Fatalf("PartitionMetadataKey: %v", err)
	}
	want := "bronze/conab/estimativa-graos/ingest_date=2026-06-25/_metadata.json"
	if key != want {
		t.Fatalf("got %q want %q", key, want)
	}
}

func TestMarshalBronzeMetadata(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
	}
	fp := Fingerprint{
		SHA256:     "abc123",
		RowCount:   42,
		IngestDate: time.Date(2026, 6, 25, 12, 0, 0, 0, time.UTC),
	}
	meta := NewBronzeMetadata(entry, fp, "https://portaldeinformacoes.conab.gov.br/example", "local")
	data, err := MarshalBronzeMetadata(meta)
	if err != nil {
		t.Fatalf("MarshalBronzeMetadata: %v", err)
	}

	var decoded BronzeMetadata
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if decoded.DatasetID != "conab.estimativa-graos" {
		t.Fatalf("dataset_id: %q", decoded.DatasetID)
	}
	if decoded.SHA256 != "abc123" || decoded.RowCount != 42 {
		t.Fatalf("unexpected payload: %+v", decoded)
	}
	if !strings.Contains(string(data), "portaldeinformacoes.conab.gov.br") {
		t.Fatalf("expected CONAB portal in source_url: %s", data)
	}
}
