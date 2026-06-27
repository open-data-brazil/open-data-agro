package mapa

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const sipeagroPackageID = "sipeagro"

var sipeagroProdutoLines = []string{
	"Fertilizantes",
	"Qualidade Vegetal",
	"Produto Veterinário",
	"Vinhos e Bebidas",
	"Alimentação Animal",
}

var sipeagroEstabelecimentoLines = append([]string{
	"Material de Multiplicação Animal",
	"Aves de Reprodução",
	"Aviação Agrícola - Registro",
}, sipeagroProdutoLines...)

var sipeagroCanonicalHeaders = []string{
	"linha_produto",
	"uf",
	"municipio",
	"numero_registro_estabelecimento",
	"status_registro",
	"cpf_cnpj",
	"razao_social",
	"nome_fantasia",
	"area_atuacao",
	"atividade",
	"classificacao",
	"caracteristica_adicional",
	"especie",
}

// FetchSIPEAGROSnapshot downloads and merges SIPEAGRO CKAN CSV resources.
func (c *Client) FetchSIPEAGROSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	lines, err := sipeagroResourceNames(entry.DatasetID.String())
	if err != nil {
		return nil, "", err
	}

	resources, err := ListCKANResources(ctx, sipeagroPackageID)
	if err != nil {
		return nil, "", err
	}

	byName := make(map[string]string, len(resources))
	for _, res := range resources {
		if strings.ToUpper(strings.TrimSpace(res.Format)) != "CSV" {
			continue
		}
		if strings.TrimSpace(res.URL) == "" {
			continue
		}
		byName[strings.TrimSpace(res.Name)] = res.URL
	}

	var merged [][]string
	var requestURLs []string
	for _, line := range lines {
		url, ok := byName[line]
		if !ok {
			return nil, "", fmt.Errorf("sipeagro resource %q not found in CKAN package", line)
		}
		rows, err := c.downloadNormalizedSIPEAGRORows(ctx, url, line)
		if err != nil {
			return nil, "", err
		}
		merged = append(merged, rows...)
		requestURLs = append(requestURLs, url)
	}
	if len(merged) == 0 {
		return nil, "", fmt.Errorf("sipeagro returned no rows for %s", entry.DatasetID)
	}

	payload, err := encodeCSV(sipeagroCanonicalHeaders, merged)
	if err != nil {
		return nil, "", err
	}
	sourceURL := fmt.Sprintf("%s?id=%s (resources: %d)", defaultCKANPackageShowURL, sipeagroPackageID, len(requestURLs))
	return payload, sourceURL, nil
}

func sipeagroResourceNames(datasetID string) ([]string, error) {
	switch datasetID {
	case "mapa.sipeagro-estabelecimentos":
		return sipeagroEstabelecimentoLines, nil
	case "mapa.sipeagro-produtos":
		return sipeagroProdutoLines, nil
	default:
		return nil, fmt.Errorf("unsupported sipeagro dataset %s", datasetID)
	}
}

func resolveSIPEAGROURL(entry catalog.RegistryEntry) (string, error) {
	return fmt.Sprintf("%s?id=%s", defaultCKANPackageShowURL, sipeagroPackageID), nil
}

func (c *Client) downloadNormalizedSIPEAGRORows(ctx context.Context, sourceURL, lineName string) ([][]string, error) {
	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, err
	}
	return parseSIPEAGROCSV(result.Body, lineName)
}

func parseSIPEAGROCSV(raw []byte, lineName string) ([][]string, error) {
	raw = decodeLatin1(raw)
	reader := csv.NewReader(bytes.NewReader(raw))
	reader.Comma = ';'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read sipeagro header: %w", err)
	}
	index := buildSIPEAGROColumnIndex(headers)

	var rows [][]string
	for {
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, fmt.Errorf("read sipeagro row: %w", readErr)
		}
		rows = append(rows, []string{
			lineName,
			fieldAt(record, index.uf),
			fieldAt(record, index.municipio),
			fieldAt(record, index.numeroRegistro),
			fieldAt(record, index.status),
			fieldAt(record, index.cpfCnpj),
			fieldAt(record, index.razaoSocial),
			fieldAt(record, index.nomeFantasia),
			fieldAt(record, index.areaAtuacao),
			fieldAt(record, index.atividade),
			fieldAt(record, index.classificacao),
			fieldAt(record, index.caracteristica),
			fieldAt(record, index.especie),
		})
	}
	return rows, nil
}

type sipeagroColumnIndex struct {
	uf            int
	municipio     int
	numeroRegistro int
	status        int
	cpfCnpj       int
	razaoSocial   int
	nomeFantasia  int
	areaAtuacao   int
	atividade     int
	classificacao int
	caracteristica int
	especie       int
}

func buildSIPEAGROColumnIndex(headers []string) sipeagroColumnIndex {
	normalized := make([]string, len(headers))
	for i, header := range headers {
		normalized[i] = strings.ToUpper(strings.TrimSpace(header))
	}
	return sipeagroColumnIndex{
		uf:             findColumn(normalized, "UNIDADE_DA_FEDERACAO", "UF"),
		municipio:      findColumn(normalized, "MUNICIPIO"),
		numeroRegistro: findColumn(normalized, "NUMERO_REGISTRO_ESTABELECIMENTO"),
		status:         findColumn(normalized, "STATUS_DO_REGISTRO", "STATUS_REGISTRO", "STATUS_REGISTRO_ESTABELECIMENTO"),
		cpfCnpj:        findColumn(normalized, "CNPJ", "CPF_CNPJ"),
		razaoSocial:    findColumn(normalized, "RAZAO_SOCIAL"),
		nomeFantasia:   findColumn(normalized, "NOME_FANTASIA"),
		areaAtuacao:    findColumn(normalized, "AREA_ATUACAO"),
		atividade:      findColumn(normalized, "ATIVIDADE"),
		classificacao:  findColumn(normalized, "CLASSIFICACAO"),
		caracteristica: findColumn(normalized, "CARACTERISTICA_ADICIONAL"),
		especie:        findColumn(normalized, "ESPECIE"),
	}
}

func findColumn(headers []string, names ...string) int {
	for _, name := range names {
		for i, header := range headers {
			if header == name {
				return i
			}
		}
	}
	return -1
}

func fieldAt(record []string, index int) string {
	if index < 0 || index >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[index])
}

func decodeLatin1(raw []byte) []byte {
	reader := transform.NewReader(bytes.NewReader(raw), charmap.ISO8859_1.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return raw
	}
	return decoded
}

func encodeCSV(headers []string, rows [][]string) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ';'
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MergeSIPEAGROSampleCSV merges sample CSV files for golden tests.
func MergeSIPEAGROSampleCSV(datasetID string, samples map[string][]byte) ([]byte, error) {
	lines, err := sipeagroResourceNames(datasetID)
	if err != nil {
		return nil, err
	}
	var merged [][]string
	for _, line := range lines {
		raw, ok := samples[line]
		if !ok {
			continue
		}
		rows, err := parseSIPEAGROCSV(raw, line)
		if err != nil {
			return nil, err
		}
		merged = append(merged, rows...)
	}
	if len(merged) == 0 {
		return nil, fmt.Errorf("no sipeagro sample rows for %s", datasetID)
	}
	return encodeCSV(sipeagroCanonicalHeaders, merged)
}

// SortSIPEAGROLines returns resource line names in stable order for tests.
func SortSIPEAGROLines(lines []string) []string {
	out := append([]string(nil), lines...)
	sort.Strings(out)
	return out
}
