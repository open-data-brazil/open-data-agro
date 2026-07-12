#!/usr/bin/env python3
"""Verify Phase 30 ML training dataset deliverables."""

from __future__ import annotations

import json
import os
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
UC = ROOT / "docs" / "use-cases" / "UC-ML-001-price-early-warning.md"
REFRESH = ROOT / "docs" / "REFRESH-POLICY.md"
ROADMAP = ROOT / "docs" / "ROADMAP.md"


def main() -> int:
    errors: list[str] = []
    export_dir = Path(os.environ.get("ML_EXPORT_DIR", ROOT / ".local" / "ml" / "export"))
    baselines_dir = Path(os.environ.get("ML_BASELINES_DIR", ROOT / ".local" / "ml" / "baselines"))
    corr_dir = Path(os.environ.get("ML_CORR_DIR", ROOT / ".local" / "ml" / "corr"))

    for path in (
        export_dir / "soy_daily_training.parquet",
        export_dir / "manifest.json",
        baselines_dir / "soy_daily_baselines.json",
        corr_dir / "soy_daily_corr.json",
        UC,
        REFRESH,
    ):
        if not path.is_file():
            errors.append(f"missing {path}")

    if UC.is_file():
        text = UC.read_text(encoding="utf-8")
        for phrase in (
            "Walk-forward",
            "y_ret_30d",
            "Success gate",
            "seasonal_naive",
            "Feature list",
        ):
            if phrase not in text:
                errors.append(f"UC-ML-001 missing: {phrase!r}")

    if REFRESH.is_file():
        refresh = REFRESH.read_text(encoding="utf-8")
        if "retrain" not in refresh.lower() and "ML retrain" not in refresh:
            errors.append("REFRESH-POLICY missing ML retrain note")

    if ROADMAP.is_file():
        roadmap = ROADMAP.read_text(encoding="utf-8")
        if "Phase 30" not in roadmap and "ml-dataset-export" not in roadmap:
            errors.append("ROADMAP missing Phase 30 / ml-dataset-export")

    manifest_path = export_dir / "manifest.json"
    if manifest_path.is_file():
        manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
        if manifest.get("row_count", 0) < 100:
            errors.append(f"export row_count too small for CI: {manifest.get('row_count')}")
        if "y_ret_30d" not in (manifest.get("columns") or []):
            errors.append("manifest missing y_ret_30d column")

    if errors:
        print("check_phase30_ml_dataset: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print("check_phase30_ml_dataset: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
