#!/usr/bin/env python3
"""Verify Phase 31 ML train / SHAP / early-warning deliverables."""

from __future__ import annotations

import json
import os
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
UC = ROOT / "docs" / "use-cases" / "UC-ML-001-price-early-warning.md"
ROADMAP = ROOT / "docs" / "ROADMAP.md"
CHANGELOG = ROOT / "CHANGELOG.md"


def main() -> int:
    errors: list[str] = []
    train_dir = Path(os.environ.get("ML_TRAIN_DIR", ROOT / ".local" / "ml" / "train"))
    reports_dir = Path(os.environ.get("ML_REPORTS_DIR", ROOT / ".local" / "ml" / "reports"))
    registry_dir = Path(os.environ.get("ML_REGISTRY_DIR", ROOT / ".local" / "ml" / "registry"))

    required = (
        train_dir / "soy_lgbm_metrics.json",
        reports_dir / "soy_shap_ablation.md",
        reports_dir / "soy_shap_ablation.json",
        reports_dir / "soy_early_warning_eval.md",
        reports_dir / "soy_early_warning_eval.json",
        registry_dir / "soy_lgbm_v1.json",
        UC,
    )
    for path in required:
        if not path.is_file():
            errors.append(f"missing {path}")

    metrics_path = train_dir / "soy_lgbm_metrics.json"
    if metrics_path.is_file():
        metrics = json.loads(metrics_path.read_text(encoding="utf-8"))
        if metrics.get("b3_gate_passed") is not True:
            errors.append("B3 gate not passed in soy_lgbm_metrics.json")
        results = metrics.get("results") or []
        horizons = {r.get("horizon_trading_days") for r in results}
        if horizons != {7, 30, 90}:
            errors.append(f"expected horizons 7/30/90, got {horizons}")
        for r in results:
            qs = r.get("quantiles") or {}
            if not {"p10", "p50", "p90"} <= set(qs):
                errors.append(f"missing quantiles for h={r.get('horizon_trading_days')}")

    if UC.is_file():
        text = UC.read_text(encoding="utf-8")
        for phrase in ("Phase 31", "LightGBM", "ml-train-soy", "early-warning"):
            if phrase not in text:
                errors.append(f"UC-ML-001 missing: {phrase!r}")

    if ROADMAP.is_file():
        roadmap = ROADMAP.read_text(encoding="utf-8")
        if "Phase 31" not in roadmap and "ml-train-soy" not in roadmap:
            errors.append("ROADMAP missing Phase 31 / ml-train-soy")

    if CHANGELOG.is_file():
        cl = CHANGELOG.read_text(encoding="utf-8")
        if "Phase 31" not in cl and "ml-train-soy" not in cl:
            errors.append("CHANGELOG missing Phase 31 entry")

    if errors:
        print("check_phase31_ml_train: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print("check_phase31_ml_train: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
