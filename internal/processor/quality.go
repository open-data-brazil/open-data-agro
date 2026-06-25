package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/lake"
)

// QualityOptions controls a bronze checkpoint run.
type QualityOptions struct {
	DatasetID  string
	Checkpoint string
}

// QualityResult summarizes a checkpoint validation.
type QualityResult struct {
	DatasetID               string
	Checkpoint              string
	BronzeDir               string
	Success                 bool
	EvaluatedExpectations   int
	SuccessfulExpectations  int
	UnsuccessfulExpectations int
}

// QualityGate runs Great Expectations bronze checkpoints before promotion.
type QualityGate struct {
	cfg      config.LakeConfig
	registry *catalog.Registry
	python   string
	script   string
}

// NewQualityGate wires quality validation dependencies.
func NewQualityGate(cfg config.LakeConfig, registry *catalog.Registry) *QualityGate {
	return &QualityGate{
		cfg:      cfg,
		registry: registry,
		python:   envOr("PYTHON", "python3"),
		script:   defaultQualityScript(),
	}
}

// RunBronzeCheckpoint validates local bronze Parquet against the dataset suite.
func (q *QualityGate) RunBronzeCheckpoint(ctx context.Context, opts QualityOptions) (*QualityResult, error) {
	if _, err := q.registry.Require(opts.DatasetID); err != nil {
		return nil, err
	}

	checkpoint := opts.Checkpoint
	if checkpoint == "" {
		var ok bool
		checkpoint, ok = bronzeCheckpointForDataset(opts.DatasetID)
		if !ok {
			return nil, fmt.Errorf("no bronze checkpoint configured for dataset %s", opts.DatasetID)
		}
	}

	lakeRoot := lake.NormalizeRoot(q.cfg.LakeLocalRoot)
	bronzeDir, err := lake.BronzeDir(lakeRoot, opts.DatasetID)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(bronzeDir); err != nil {
		return nil, fmt.Errorf("bronze directory missing (%s): run ingestor first", bronzeDir)
	}

	if _, err := os.Stat(q.script); err != nil {
		return nil, fmt.Errorf("quality script not found at %s", q.script)
	}

	cmd := exec.CommandContext(ctx, q.python, q.script,
		"--checkpoint", checkpoint,
		"--bronze-dir", bronzeDir,
	)
	cmd.Env = append(os.Environ(), "LAKE_LOCAL_ROOT="+lakeRoot)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("quality checkpoint: %w\n%s", err, strings.TrimSpace(string(out)))
	}

	var payload struct {
		Success                  bool   `json:"success"`
		EvaluatedExpectations    int    `json:"evaluated_expectations"`
		SuccessfulExpectations   int    `json:"successful_expectations"`
		UnsuccessfulExpectations int    `json:"unsuccessful_expectations"`
		BronzeDir                string `json:"bronze_dir"`
	}
	if err := json.Unmarshal(lastJSONLine(out), &payload); err != nil {
		return nil, fmt.Errorf("parse quality output: %w (raw=%q)", err, string(out))
	}

	result := &QualityResult{
		DatasetID:                opts.DatasetID,
		Checkpoint:               checkpoint,
		BronzeDir:                payload.BronzeDir,
		Success:                  payload.Success,
		EvaluatedExpectations:    payload.EvaluatedExpectations,
		SuccessfulExpectations:   payload.SuccessfulExpectations,
		UnsuccessfulExpectations: payload.UnsuccessfulExpectations,
	}
	if !payload.Success {
		return result, fmt.Errorf(
			"bronze quality failed: %d/%d expectations passed",
			payload.SuccessfulExpectations,
			payload.EvaluatedExpectations,
		)
	}
	return result, nil
}

func bronzeCheckpointForDataset(datasetID string) (string, bool) {
	switch datasetID {
	case "conab.estimativa-graos":
		return "bronze_conab_estimativa_graos", true
	case "conab.serie-historica-graos":
		return "bronze_conab_serie_historica_graos", true
	case "conab.oferta-demanda":
		return "bronze_conab_oferta_demanda", true
	case "conab.estoques-publicos":
		return "bronze_conab_estoques_publicos", true
	case "conab.operacoes-comercializacao":
		return "bronze_conab_operacoes_comercializacao", true
	case "conab.vendas-balcao":
		return "bronze_conab_vendas_balcao", true
	case "anp.combustiveis-precos-medios-municipios":
		return "bronze_anp_combustiveis_precos_medios_municipios", true
	case "anp.combustiveis-precos-postos":
		return "bronze_anp_combustiveis_precos_postos", true
	case "conab.armazenagem":
		return "bronze_conab_armazenagem", true
	case "conab.frete":
		return "bronze_conab_frete", true
	case "conab.serie-historica-capacidade-estatica":
		return "bronze_conab_serie_historica_capacidade_estatica", true
	case "conab.alimenta-brasil-entregas":
		return "bronze_conab_alimenta_brasil_entregas", true
	case "conab.alimenta-brasil-propostas":
		return "bronze_conab_alimenta_brasil_propostas", true
	case "ibge.localidades-municipios":
		return "bronze_ibge_localidades_municipios", true
	case "ibge.localidades-ufs":
		return "bronze_ibge_localidades_ufs", true
	case "ibge.localidades-regioes":
		return "bronze_ibge_localidades_regioes", true
	case "ibge.localidades-mesorregioes":
		return "bronze_ibge_localidades_mesorregioes", true
	case "ibge.localidades-microrregioes":
		return "bronze_ibge_localidades_microrregioes", true
	case "ibge.pam-area-quantidade":
		return "bronze_ibge_pam_area_quantidade", true
	case "ibge.pam-rendimento-valor":
		return "bronze_ibge_pam_rendimento_valor", true
	case "ibge.pam-estabelecimentos":
		return "bronze_ibge_pam_estabelecimentos", true
	default:
		return "", false
	}
}

func defaultQualityScript() string {
	if v := strings.TrimSpace(os.Getenv("QUALITY_CHECKPOINT_SCRIPT")); v != "" {
		return v
	}
	if root, err := findModuleRoot(); err == nil {
		return filepath.Join(root, "scripts", "quality", "run_checkpoint.py")
	}
	return filepath.Join("scripts", "quality", "run_checkpoint.py")
}
