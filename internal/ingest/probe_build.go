package ingest

import (
	"os"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/bcb"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/cepea"
	"github.com/open-data-brazil/open-data-agro/internal/ibge"
	"github.com/open-data-brazil/open-data-agro/internal/inmet"
	"github.com/open-data-brazil/open-data-agro/internal/inpe"
	"github.com/open-data-brazil/open-data-agro/internal/mdic"
	"github.com/open-data-brazil/open-data-agro/internal/oecd"
	"github.com/open-data-brazil/open-data-agro/internal/un"
	"github.com/open-data-brazil/open-data-agro/internal/usda"
	"github.com/open-data-brazil/open-data-agro/internal/wto"
)

// ProbeSpec describes an HTTP probe aligned with ingest fetch semantics.
type ProbeSpec struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

func BuildProbeSpec(entry catalog.RegistryEntry) (ProbeSpec, error) {
	datasetID := entry.DatasetID.String()
	agency, _, err := catalog.SplitDatasetID(datasetID)
	if err != nil {
		return ProbeSpec{}, err
	}

	if direct := strings.TrimSpace(entry.SourceURL); isDirectProbeURL(direct) {
		return ProbeSpec{URL: direct, Headers: probeHeaders(entry)}, nil
	}

	switch agency {
	case "bcb":
		url, err := bcb.BuildProbeURL(entry)
		if err != nil {
			return ProbeSpec{}, err
		}
		return ProbeSpec{URL: url, Headers: probeHeaders(entry)}, nil
	case "ibge":
		if url, err := ibge.BuildProbeURL(entry); err == nil {
			return ProbeSpec{URL: url, Headers: map[string]string{"Accept": "application/json"}}, nil
		}
	case "inmet":
		if strings.Contains(entry.SourceURL, "{year}") {
			url, err := inmet.ResolveAnnualZIPURL(entry, time.Now().UTC().Year()-1)
			if err != nil {
				return ProbeSpec{}, err
			}
			return ProbeSpec{URL: url, Headers: probeHeaders(entry)}, nil
		}
		if datasetID == "inmet.sequia-monitor" {
			return ProbeSpec{
				URL:     strings.TrimSpace(entry.SourceURL),
				Headers: map[string]string{"Accept": "application/json"},
			}, nil
		}
	case "mdic":
		body, url, err := mdic.BuildProbeRequest(entry)
		if err != nil {
			return ProbeSpec{}, err
		}
		return ProbeSpec{
			Method:  "POST",
			URL:     url,
			Body:    string(body),
			Headers: map[string]string{"Accept": "application/json", "Content-Type": "application/json"},
		}, nil
	case "inpe":
		url, err := inpe.BuildProbeURL(entry)
		if err != nil {
			return ProbeSpec{}, err
		}
		return ProbeSpec{URL: url, Headers: probeHeaders(entry)}, nil
	case "oecd", "oecd-fao":
		url, err := oecd.BuildMetadataProbeURL(entry)
		if err != nil {
			return ProbeSpec{}, err
		}
		return ProbeSpec{URL: url, Headers: map[string]string{"Accept": "application/json"}}, nil
	case "fao":
		if datasetID == "fao.giews-crop-prospects" {
			return ProbeSpec{
				URL:     "https://www.fao.org/giews/en/",
				Headers: map[string]string{"Accept": "text/html,application/xhtml+xml,*/*"},
			}, nil
		}
	case "japan":
		return ProbeSpec{
			URL: strings.TrimSpace(entry.PortalURL()),
			Headers: map[string]string{
				"Accept":          "text/html,application/xhtml+xml,*/*",
				"Accept-Language": "en-US,en;q=0.9",
			},
		}, nil
	case "usda":
		if datasetID == "usda.gats-trade" {
			return ProbeSpec{
				URL:     strings.TrimSpace(entry.SourceURL),
				Headers: map[string]string{"Accept": "application/json"},
			}, nil
		}
		if strings.HasPrefix(datasetID, "usda.psd-") {
			req, err := usda.BuildPSDProbeRequest(entry)
			if err != nil {
				return ProbeSpec{}, err
			}
			return ProbeSpec{Method: "POST", URL: req.URL, Headers: req.Headers, Body: req.Body}, nil
		}
	case "un":
		url, err := un.BuildProbeURL(entry)
		if err != nil {
			return ProbeSpec{}, err
		}
		return ProbeSpec{URL: url, Headers: map[string]string{"Accept": "application/json"}}, nil
	case "wto":
		if url, headers, err := wto.BuildProbeRequest(entry); err == nil {
			return ProbeSpec{URL: url, Headers: headers}, nil
		}
	}

	url, err := resolveLegacyProbeURL(entry, agency)
	if err != nil {
		return ProbeSpec{}, err
	}

	headers := probeHeaders(entry)
	if agency == "fao" && strings.Contains(strings.ToLower(url), "bulks-faostat") {
		headers["Range"] = "bytes=0-4095"
	}
	if agency == "japan" {
		headers["Accept"] = "text/html,application/xhtml+xml,*/*"
	}

	return ProbeSpec{URL: url, Headers: headers}, nil
}

func resolveLegacyProbeURL(entry catalog.RegistryEntry, agency string) (string, error) {
	if agency == "cepea" {
		if mirror, err := cepea.MirrorURL(entry.DatasetID.String()); err == nil {
			return mirror, nil
		}
	}

	url, err := ResolveSourceURL(entry)
	if err != nil {
		return "", err
	}

	if agency == "eia" {
		key := strings.TrimSpace(os.Getenv("EIA_API_KEY"))
		if key != "" && !strings.Contains(url, "api_key=") {
			sep := "?"
			if strings.Contains(url, "?") {
				sep = "&"
			}
			url = url + sep + "api_key=" + key
		}
	}

	if agency == "ana" && strings.Contains(strings.ToLower(url), ".asmx") {
		lower := strings.ToLower(url)
		if idx := strings.Index(lower, ".asmx"); idx >= 0 {
			url = url[:idx+5] + "?WSDL"
		}
	}

	return url, nil
}

func isDirectProbeURL(raw string) bool {
	if raw == "" {
		return false
	}
	lower := strings.ToLower(raw)
	return strings.Contains(lower, "/download/") ||
		strings.Contains(lower, "/api/3/action/package_show")
}
