package inmet

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestStationsJSONToCatalogCSVAutomatic(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "estacoes_automaticas.api.sample.json")
	payload, err := stationsJSONToCatalogCSV("inmet.estacoes-automaticas", raw)
	if err != nil {
		t.Fatalf("stationsJSONToCatalogCSV: %v", err)
	}

	headers, rows, err := FlattenEstacoes("inmet.estacoes-automaticas", payload)
	if err != nil {
		t.Fatalf("FlattenEstacoes: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows: got %d want 2", len(rows))
	}
	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["cd_estacao"]]; got != "A901" {
		t.Fatalf("cd_estacao: got %q want A901", got)
	}
}

func TestWriteLongCSVRoundTrip(t *testing.T) {
	t.Parallel()

	rows := [][]string{{"A901", "2023-01-01", "precipitacao", "10.5", "MT", "2023"}}
	payload, err := writeLongCSV(rows)
	if err != nil {
		t.Fatalf("writeLongCSV: %v", err)
	}
	records, err := readLatinCSV(payload)
	if err != nil {
		t.Fatalf("readLatinCSV: %v", err)
	}
	if len(records) != 2 || len(records[1]) != 6 {
		t.Fatalf("records: %+v", records)
	}
}

func TestAggregateDailyLongToMonthlyFromSampleZIP(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "bdmep_2023_mt.sample.zip")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("inmet.bdmep-mensal"),
		ClimateVariables: []string{
			"precipitacao", "temperatura_ar", "temperatura_maxima", "temperatura_minima",
			"umidade_relativa", "vento_velocidade",
		},
	}

	dailyRows, err := parseAnnualZIPToDailyLong(raw, 2023, map[string]struct{}{"MT": {}}, allowedClimateVariables(entry))
	if err != nil {
		t.Fatalf("parseAnnualZIPToDailyLong: %v", err)
	}
	dailyPayload, err := writeLongCSV(dailyRows)
	if err != nil {
		t.Fatalf("writeLongCSV: %v", err)
	}

	monthlyPayload, err := aggregateDailyLongCSVToMonthly(dailyPayload, entry.DatasetID.String())
	if err != nil {
		t.Fatalf("aggregateDailyLongCSVToMonthly: %v", err)
	}
	_, rows, err := flattenLongCSV(monthlyPayload, []string{"cd_estacao", "mes", "variavel", "valor", "uf", "ano"})
	if err != nil {
		t.Fatalf("flattenLongCSV: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected monthly rows")
	}
}

func TestUFToRegion(t *testing.T) {
	t.Parallel()

	if got := ufToRegion("MT"); got != "Centro-Oeste" {
		t.Fatalf("ufToRegion(MT): got %q", got)
	}
	if got := ufToRegion("RS"); got != "Sul" {
		t.Fatalf("ufToRegion(RS): got %q", got)
	}
}
