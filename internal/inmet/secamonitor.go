package inmet

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const defaultSecaMonitorURL = "https://apimsbr.ana.gov.br/rpc/v1/dados-tabulares-monitor"

var secaMonitorURL = defaultSecaMonitorURL

type secaMonitorResponse struct {
	Data struct {
		List []secaMonitorEntry `json:"list"`
	} `json:"data"`
}

type secaMonitorEntry struct {
	Areas []secaMonitorArea `json:"areas"`
	Mapa  secaMonitorMapa   `json:"mapa"`
}

type secaMonitorArea struct {
	ID       int    `json:"id"`
	Categoria string `json:"categoria"`
	Area     int    `json:"area"`
	TipoArea int    `json:"tipo_area"`
	AreaID   int    `json:"area_id"`
	MapaID   int    `json:"mapa_id"`
}

type secaMonitorMapa struct {
	ID  int `json:"id"`
	Ano int `json:"ano"`
	Mes int `json:"mes"`
}

// FetchSecaMonitorSnapshot downloads the ANA drought monitor tabular API payload.
func (c *Client) FetchSecaMonitorSnapshot(ctx context.Context, sourceURL string) ([]byte, string, error) {
	url := strings.TrimSpace(sourceURL)
	if url == "" {
		url = secaMonitorURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("seca monitor fetch: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("seca monitor status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return nil, "", err
	}
	if len(body) == 0 {
		return nil, "", fmt.Errorf("empty seca monitor response")
	}

	return body, url, nil
}

// FlattenSecaMonitor converts ANA drought monitor JSON into bronze rows.
func FlattenSecaMonitor(raw []byte) ([]string, [][]string, error) {
	var payload secaMonitorResponse
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, nil, fmt.Errorf("parse seca monitor json: %w", err)
	}

	headers := []string{
		"mapa_id",
		"ano",
		"mes",
		"categoria_seca",
		"area_km2",
		"area_id",
		"tipo_area",
	}

	out := make([][]string, 0)
	for _, entry := range payload.Data.List {
		ano := strconv.Itoa(entry.Mapa.Ano)
		mes := strconv.Itoa(entry.Mapa.Mes)
		mapaID := strconv.Itoa(entry.Mapa.ID)
		for _, area := range entry.Areas {
			out = append(out, []string{
				mapaID,
				ano,
				mes,
				strings.TrimSpace(area.Categoria),
				strconv.Itoa(area.Area),
				strconv.Itoa(area.AreaID),
				strconv.Itoa(area.TipoArea),
			})
		}
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no seca monitor rows to flatten")
	}
	return headers, out, nil
}
