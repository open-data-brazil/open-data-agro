package processor

import (
	"context"
	"fmt"

	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/lake"
)

// PreviewPromote runs DuckDB promotion preview SQL and returns the staged row count.
func PreviewPromote(ctx context.Context, cfg config.LakeConfig, datasetID, ingestDate string) (int, error) {
	duck, err := NewDuckDB(cfg.DuckDBPath)
	if err != nil {
		return 0, err
	}

	smoker, err := NewSmoker(cfg, nil)
	if err != nil {
		return 0, err
	}
	bronzeURI, err := smoker.bronzeURI(SmokeOptions{DatasetID: datasetID, IngestDate: ingestDate})
	if err != nil {
		return 0, err
	}

	silverRoot := lake.NormalizeRoot(cfg.DeltaStoragePath)
	silverDir, err := lake.SilverTableDir(silverRoot, datasetID)
	if err != nil {
		return 0, err
	}

	scriptPath, err := ScriptPath("promote_bronze_to_silver.sql")
	if err != nil {
		return 0, err
	}

	vars := map[string]string{
		"bronze_uri": bronzeURI,
		"dataset_id": datasetID,
		"silver_dir": silverDir,
	}
	sql, err := loadScript(scriptPath, vars)
	if err != nil {
		return 0, err
	}

	out, err := duck.RunSQL(ctx, smoker.extensionSetup()+sql)
	if err != nil {
		return 0, fmt.Errorf("promotion preview: %w", err)
	}
	return parseCountCSV(out)
}
