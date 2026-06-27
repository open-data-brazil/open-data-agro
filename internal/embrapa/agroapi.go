package embrapa

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

type agrofitRow struct {
	NumeroRegistro  string `json:"numero_registro"`
	MarcaComercial  string `json:"marca_comercial"`
	Situacao        string `json:"situacao"`
	Classe          string `json:"classe"`
	Formulacao      string `json:"formulacao"`
	IngredienteAtivo string `json:"ingrediente_ativo"`
}

// FetchAgroAPIAgrofitSnapshot returns Agrofit formulated products (live API or fixture).
func (c *Client) FetchAgroAPIAgrofitSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	baseURL := strings.TrimSpace(entry.SourceURL)
	if baseURL == "" {
		baseURL = defaultAgrofitAPIURL
	}

	token := strings.TrimSpace(os.Getenv("EMBRAPA_AGROAPI_KEY"))
	if token != "" {
		liveURL := baseURL + "?limit=100"
		raw, err := c.download(ctx, liveURL, token)
		if err == nil && len(raw) > 2 {
			normalized, normErr := normalizeAgrofitPayload(raw)
			if normErr == nil && len(normalized) > 0 {
				payload, marshalErr := json.Marshal(normalized)
				if marshalErr != nil {
					return nil, "", marshalErr
				}
				return payload, liveURL + " (live AgroAPI)", nil
			}
		}
	}

	if path := strings.TrimSpace(os.Getenv("EMBRAPA_AGROFIT_BULK_PATH")); path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, "", err
		}
		rows, err := parseAgrofitFixture(raw)
		if err != nil {
			return nil, "", err
		}
		payload, err := json.Marshal(rows)
		if err != nil {
			return nil, "", err
		}
		return payload, path + " (EMBRAPA_AGROFIT_BULK_PATH)", nil
	}

	rows, err := embeddedAgrofitSample()
	if err != nil {
		return nil, "", err
	}
	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	note := "https://www.agroapi.cnptia.embrapa.br/store/ (fixture — register for EMBRAPA_AGROAPI_KEY)"
	return payload, note, nil
}

func embeddedAgrofitSample() ([]agrofitRow, error) {
	return []agrofitRow{
		{
			NumeroRegistro:   "35523",
			MarcaComercial:   "KBR-829M1-02",
			Situacao:         "TRUE",
			Classe:           "Agente Biológico de Controle",
			Formulacao:       "Nematóides vivos",
			IngredienteAtivo: "Heterorhabditis bacteriophora",
		},
		{
			NumeroRegistro:   "12345",
			MarcaComercial:   "EXAMPLE AGROFIT",
			Situacao:         "TRUE",
			Classe:           "Herbicida",
			Formulacao:       "Concentrado emulsionável",
			IngredienteAtivo: "Glyphosate",
		},
	}, nil
}

func parseAgrofitFixture(raw []byte) ([]agrofitRow, error) {
	normalized, err := normalizeAgrofitPayload(raw)
	if err != nil {
		return nil, err
	}
	if len(normalized) == 0 {
		return nil, fmt.Errorf("empty embrapa agrofit fixture")
	}
	return normalized, nil
}

func normalizeAgrofitPayload(raw []byte) ([]agrofitRow, error) {
	var rows []agrofitRow
	if err := json.Unmarshal(raw, &rows); err == nil && len(rows) > 0 {
		return rows, nil
	}

	var envelope struct {
		Data []agrofitRow `json:"data"`
		Items []agrofitRow `json:"items"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("parse embrapa agrofit json: %w", err)
	}
	if len(envelope.Data) > 0 {
		return envelope.Data, nil
	}
	if len(envelope.Items) > 0 {
		return envelope.Items, nil
	}
	return nil, fmt.Errorf("no embrapa agrofit rows in payload")
}

// FlattenAgroAPIAgrofit converts Agrofit JSON into bronze columns.
func FlattenAgroAPIAgrofit(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	rows, err := parseAgrofitFixture(raw)
	if err != nil {
		return nil, nil, err
	}

	headers := []string{
		"numero_registro",
		"marca_comercial",
		"situacao",
		"classe",
		"formulacao",
		"ingrediente_ativo",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			row.NumeroRegistro,
			row.MarcaComercial,
			row.Situacao,
			row.Classe,
			row.Formulacao,
			row.IngredienteAtivo,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no embrapa agrofit rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

// ResolveURL returns the AgroAPI endpoint for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return defaultAgrofitAPIURL, nil
	}
	if !strings.Contains(strings.ToLower(raw), "embrapa.br") && !strings.Contains(strings.ToLower(raw), "cnptia.embrapa.br") {
		return "", fmt.Errorf("source_url for %s must be on embrapa.br", entry.DatasetID)
	}
	return raw, nil
}
