package anp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var lpcLinkPattern = regexp.MustCompile(`href="([^"]*arquivos-lpc/[^"]+\.xlsx)"`)

// ResolveLatestLPCURL scrapes the ANP LPC listing page for the newest weekly XLSX file.
func ResolveLatestLPCURL(listingURL, filePrefix string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, listingURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch LPC listing: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("LPC listing status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return "", err
	}

	html := string(body)
	matches := lpcLinkPattern.FindAllStringSubmatch(html, -1)
	prefix := strings.TrimSpace(filePrefix)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		href := strings.ReplaceAll(match[1], "&amp;", "&")
		if !strings.Contains(href, prefix) {
			continue
		}
		if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
			return href, nil
		}
		if strings.HasPrefix(href, "/") {
			return "https://www.gov.br" + href, nil
		}
		return "https://www.gov.br/" + strings.TrimPrefix(href, "./"), nil
	}

	return "", fmt.Errorf("no %s .xlsx link found on %s", prefix, listingURL)
}
