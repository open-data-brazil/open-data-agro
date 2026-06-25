package ibge

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// FlattenLocalidades converts an IBGE Localidades JSON snapshot into a string table.
func FlattenLocalidades(datasetID string, raw []byte) ([]string, [][]string, error) {
	var payload any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, nil, fmt.Errorf("parse json: %w", err)
	}

	records, err := asObjectSlice(payload)
	if err != nil {
		return nil, nil, err
	}

	switch datasetID {
	case "ibge.localidades-municipios":
		return flattenMunicipios(records)
	case "ibge.localidades-ufs":
		return flattenUFs(records)
	case "ibge.localidades-regioes":
		return flattenRegioes(records)
	case "ibge.localidades-mesorregioes":
		return flattenMesorregioes(records)
	case "ibge.localidades-microrregioes":
		return flattenMicrorregioes(records)
	default:
		return nil, nil, fmt.Errorf("unsupported ibge localidades dataset %s", datasetID)
	}
}

func flattenMunicipios(records []map[string]any) ([]string, [][]string, error) {
	headers := []string{
		"codigo_ibge", "nome", "sigla_uf", "codigo_uf", "codigo_regiao", "nome_regiao",
	}
	var rows [][]string
	for _, record := range records {
		uf, ok := municipalityUF(record)
		if !ok {
			return nil, nil, fmt.Errorf("municipio %v missing UF metadata", record["id"])
		}
		regiao, _ := uf["regiao"].(map[string]any)
		rows = append(rows, []string{
			formatID(record["id"]),
			stringField(record["nome"]),
			stringField(uf["sigla"]),
			formatID(uf["id"]),
			formatID(regiao["id"]),
			stringField(regiao["nome"]),
		})
	}
	return headers, rows, nil
}

func flattenUFs(records []map[string]any) ([]string, [][]string, error) {
	headers := []string{
		"codigo_uf", "sigla_uf", "nome", "codigo_regiao", "sigla_regiao", "nome_regiao",
	}
	var rows [][]string
	for _, record := range records {
		regiao, _ := record["regiao"].(map[string]any)
		rows = append(rows, []string{
			formatID(record["id"]),
			stringField(record["sigla"]),
			stringField(record["nome"]),
			formatID(regiao["id"]),
			stringField(regiao["sigla"]),
			stringField(regiao["nome"]),
		})
	}
	return headers, rows, nil
}

func flattenRegioes(records []map[string]any) ([]string, [][]string, error) {
	headers := []string{"codigo_regiao", "sigla_regiao", "nome"}
	var rows [][]string
	for _, record := range records {
		rows = append(rows, []string{
			formatID(record["id"]),
			stringField(record["sigla"]),
			stringField(record["nome"]),
		})
	}
	return headers, rows, nil
}

func flattenMesorregioes(records []map[string]any) ([]string, [][]string, error) {
	headers := []string{
		"codigo_mesorregiao", "nome", "codigo_uf", "sigla_uf", "nome_uf",
		"codigo_regiao", "sigla_regiao", "nome_regiao",
	}
	var rows [][]string
	for _, record := range records {
		uf, _ := record["UF"].(map[string]any)
		regiao, _ := uf["regiao"].(map[string]any)
		rows = append(rows, []string{
			formatID(record["id"]),
			stringField(record["nome"]),
			formatID(uf["id"]),
			stringField(uf["sigla"]),
			stringField(uf["nome"]),
			formatID(regiao["id"]),
			stringField(regiao["sigla"]),
			stringField(regiao["nome"]),
		})
	}
	return headers, rows, nil
}

func flattenMicrorregioes(records []map[string]any) ([]string, [][]string, error) {
	headers := []string{
		"codigo_microrregiao", "nome", "codigo_mesorregiao", "nome_mesorregiao",
		"codigo_uf", "sigla_uf", "nome_uf",
	}
	var rows [][]string
	for _, record := range records {
		meso, _ := record["mesorregiao"].(map[string]any)
		uf, _ := meso["UF"].(map[string]any)
		rows = append(rows, []string{
			formatID(record["id"]),
			stringField(record["nome"]),
			formatID(meso["id"]),
			stringField(meso["nome"]),
			formatID(uf["id"]),
			stringField(uf["sigla"]),
			stringField(uf["nome"]),
		})
	}
	return headers, rows, nil
}

func municipalityUF(record map[string]any) (map[string]any, bool) {
	if micro, ok := record["microrregiao"].(map[string]any); ok && micro != nil {
		if meso, ok := micro["mesorregiao"].(map[string]any); ok && meso != nil {
			if uf, ok := meso["UF"].(map[string]any); ok && uf != nil {
				return uf, true
			}
		}
	}
	if imediata, ok := record["regiao-imediata"].(map[string]any); ok && imediata != nil {
		if inter, ok := imediata["regiao-intermediaria"].(map[string]any); ok && inter != nil {
			if uf, ok := inter["UF"].(map[string]any); ok && uf != nil {
				return uf, true
			}
		}
	}
	return nil, false
}

func asObjectSlice(payload any) ([]map[string]any, error) {
	switch typed := payload.(type) {
	case []any:
		out := make([]map[string]any, 0, len(typed))
		for i, item := range typed {
			obj, ok := item.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("record %d is not an object", i)
			}
			out = append(out, obj)
		}
		return out, nil
	case map[string]any:
		return []map[string]any{typed}, nil
	default:
		return nil, fmt.Errorf("expected json array or object")
	}
}

func formatID(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(typed)
	case float64:
		return strconv.FormatInt(int64(typed), 10)
	case json.Number:
		return typed.String()
	default:
		return fmt.Sprint(value)
	}
}

func stringField(value any) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}
