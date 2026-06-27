package suframa

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultBasePageURL = "https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados"

var xlsxLinkPattern = regexp.MustCompile(`href="([^"]+\.xlsx)"`)

// ResolveURL returns the latest SUFRAMA ZFM trade spreadsheet URL.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw != "" && strings.HasSuffix(strings.ToLower(raw), ".xlsx") {
		return validateGovURL(entry, raw)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return resolveLatestXLSXURL(ctx)
}

func validateGovURL(entry catalog.RegistryEntry, raw string) (string, error) {
	if !strings.Contains(strings.ToLower(raw), "gov.br") {
		return "", fmt.Errorf("source_url for %s must be on gov.br", entry.DatasetID)
	}
	return raw, nil
}

func resolveLatestXLSXURL(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultBasePageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("suframa portal fetch: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("suframa portal status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return "", err
	}

	html := string(body)
	matches := xlsxLinkPattern.FindAllStringSubmatch(html, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("no xlsx resources on suframa base-de-dados page")
	}

	type candidate struct {
		url  string
		year int
	}
	var candidates []candidate
	yearPattern := regexp.MustCompile(`(20\d{2})`)
	for _, match := range matches {
		link := match[1]
		if !strings.Contains(strings.ToLower(link), "notasfiscais") &&
			!strings.Contains(strings.ToLower(link), "nota") {
			continue
		}
		year := 0
		if y := yearPattern.FindString(link); y != "" {
			if parsed, err := strconv.Atoi(y); err == nil {
				year = parsed
			}
		}
		if !strings.HasPrefix(link, "http") {
			link = "https://www.gov.br" + link
		}
		candidates = append(candidates, candidate{url: link, year: year})
	}

	if len(candidates) == 0 {
		link := matches[len(matches)-1][1]
		if !strings.HasPrefix(link, "http") {
			link = "https://www.gov.br" + link
		}
		return link, nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].year > candidates[j].year
	})

	return candidates[0].url, nil
}
