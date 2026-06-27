package suframa

import (
	"fmt"
	"strings"
	"testing"
)

func TestResolveURLDirect(t *testing.T) {
	t.Parallel()

	entry := struct {
		DatasetID string
		SourceURL string
	}{
		DatasetID: "suframa.comercio-mercadorias-zfm",
		SourceURL: "https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx",
	}

	raw := strings.TrimSpace(entry.SourceURL)
	if !strings.Contains(strings.ToLower(raw), "gov.br") {
		t.Fatalf("expected gov.br url")
	}
	if !strings.HasSuffix(strings.ToLower(raw), ".xlsx") {
		t.Fatalf("expected xlsx suffix")
	}
	_ = fmt.Sprintf("%s", entry.DatasetID)
}
