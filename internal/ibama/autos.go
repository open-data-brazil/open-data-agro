package ibama

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var autosYearPattern = regexp.MustCompile(`auto_infracao_(\d{4})\.csv$`)

func autosBulkPath() string {
	return strings.TrimSpace(os.Getenv("IBAMA_AUTOS_BULK_PATH"))
}

func extractAutosYearFromFile(path string, year int) ([]byte, string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	return extractAutosYearFromZIP(raw, year)
}

func extractAutosYearFromZIP(zipBytes []byte, year int) ([]byte, string, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, "", fmt.Errorf("parse ibama autos zip: %w", err)
	}

	type candidate struct {
		name string
		year int
	}
	var matches []candidate
	for _, file := range reader.File {
		name := strings.TrimSpace(file.Name)
		match := autosYearPattern.FindStringSubmatch(strings.ToLower(name))
		if len(match) < 2 {
			continue
		}
		y, convErr := strconv.Atoi(match[1])
		if convErr != nil {
			continue
		}
		if year > 0 && y != year {
			continue
		}
		matches = append(matches, candidate{name: name, year: y})
	}
	if len(matches) == 0 {
		if year > 0 {
			return nil, "", fmt.Errorf("no auto_infracao_%d.csv in ibama autos zip", year)
		}
		return nil, "", fmt.Errorf("no auto_infracao_YYYY.csv members in ibama autos zip")
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].year > matches[j].year
	})
	target := matches[0].name

	for _, file := range reader.File {
		if file.Name != target {
			continue
		}
		rc, openErr := file.Open()
		if openErr != nil {
			return nil, "", openErr
		}
		defer func() { _ = rc.Close() }()
		body, readErr := ioReadBytes(rc)
		if readErr != nil {
			return nil, "", readErr
		}
		return body, target, nil
	}
	return nil, "", fmt.Errorf("zip member %s not found", target)
}

func ioReadBytes(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	return buf.Bytes(), err
}
