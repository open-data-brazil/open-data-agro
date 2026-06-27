package mexico

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

type produccionRow struct {
	StateCode   string `json:"state_code"`
	CropSlug    string `json:"crop_slug"`
	RefYear     string `json:"refyear"`
	AreaHa      string `json:"area_ha"`
	ProductionT string `json:"production_t"`
}

// Client downloads SIAP agricultural production statistics.
type Client struct{}

// NewClient creates a Mexico SIAP client.
func NewClient() *Client { return &Client{} }

// FetchSIAPProduccionSnapshot returns SIAP production rows (fixture until CKAN bulk stable).
func (c *Client) FetchSIAPProduccionSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	_ = ctx
	rows, err := embeddedProduccionSample()
	if err != nil {
		return nil, "", err
	}
	if path := strings.TrimSpace(os.Getenv("SIAP_BULK_PATH")); path != "" {
		if parsed, parseErr := parseProduccionJSONFile(path); parseErr == nil && len(parsed) > 0 {
			rows = parsed
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].StateCode+rows[i].CropSlug < rows[j].StateCode+rows[j].CropSlug
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = "https://www.gob.mx/siap"
	}
	return payload, sourceURL + " (fixture — open bulk API URL unstable 2026-06-26)", nil
}

func embeddedProduccionSample() ([]produccionRow, error) {
	return []produccionRow{
		{StateCode: "JAL", CropSlug: "milho", RefYear: "2023", AreaHa: "125000", ProductionT: "980000"},
		{StateCode: "SIN", CropSlug: "soja", RefYear: "2023", AreaHa: "45000", ProductionT: "120000"},
	}, nil
}

func parseProduccionJSONFile(path string) ([]produccionRow, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rows []produccionRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// FlattenSIAPProduccion converts merged SIAP JSON into canonical bronze columns.
func FlattenSIAPProduccion(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []produccionRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse mexico siap json: %w", err)
	}

	headers := []string{"state_code", "crop_slug", "refyear", "area_ha", "production_t"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.StateCode) == "" {
			continue
		}
		out = append(out, []string{
			row.StateCode, row.CropSlug, row.RefYear, row.AreaHa, row.ProductionT,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no mexico siap rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
