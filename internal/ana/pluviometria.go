package ana

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// FetchPluviometriaSnapshot downloads daily rainfall series (HidroWeb tipoDados=2) or fixture fallback.
func (c *Client) FetchPluviometriaSnapshot(ctx context.Context, entry catalog.RegistryEntry, opts HidrologiaFetchOptions) ([]byte, string, error) {
	tipoDados := strings.TrimSpace(entry.ANATipoDados)
	if tipoDados == "" {
		tipoDados = "2"
	}
	rainEntry := entry
	rainEntry.ANATipoDados = tipoDados

	body, sourceURL, err := c.FetchHidrologiaSnapshot(ctx, rainEntry, opts)
	if err == nil && len(body) > 2 {
		return body, sourceURL, nil
	}

	if path := strings.TrimSpace(os.Getenv("ANA_PLUVIOMETRIA_FIXTURE_PATH")); path != "" {
		raw, readErr := os.ReadFile(path)
		if readErr == nil && len(raw) > 2 {
			return raw, path + " (ANA_PLUVIOMETRIA_FIXTURE_PATH)", nil
		}
	}

	rows := embeddedPluviometriaSample()
	payload, marshalErr := json.Marshal(rows)
	if marshalErr != nil {
		return nil, "", marshalErr
	}
	note := defaultServiceBaseURL + " (fixture — live rain series often empty on legacy HidroWeb; complement INMET BDMEP)"
	return payload, note, nil
}

func embeddedPluviometriaSample() []hidrologiaRow {
	return []hidrologiaRow{
		{
			StationCode:      "87017001",
			ConsistencyLevel: "1",
			DataType:         "2",
			ObservedAt:       "2024-06-01 00:00:00",
			DailyMean:        "12.4",
		},
		{
			StationCode:      "87017001",
			ConsistencyLevel: "1",
			DataType:         "2",
			ObservedAt:       "2024-06-02 00:00:00",
			DailyMean:        "0.0",
		},
	}
}
