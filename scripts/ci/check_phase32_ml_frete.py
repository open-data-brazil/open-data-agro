#!/usr/bin/env python3
"""Verify Phase 32 frete deliverables."""

from __future__ import annotations

import json
import os
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]


def main() -> int:
    errors: list[str] = []
    export_dir = Path(
        os.environ.get("ML_FRETE_EXPORT_DIR", ROOT / ".local" / "ml" / "export_frete")
    )
    train_dir = Path(
        os.environ.get("ML_FRETE_TRAIN_DIR", ROOT / ".local" / "ml" / "train_frete")
    )
    reports = Path(os.environ.get("ML_REPORTS_DIR", ROOT / ".local" / "ml" / "reports"))
    registry = Path(os.environ.get("ML_FRETE_REGISTRY_DIR", ROOT / ".local" / "ml" / "registry"))

    for p in (
        export_dir / "frete_training.parquet",
        export_dir / "manifest.json",
        train_dir / "frete_lgbm_metrics.json",
        reports / "frete_shap_ablation.md",
        reports / "frete_early_warning_eval.md",
        registry / "frete_lgbm_v1.json",
    ):
        if not p.is_file():
            errors.append(f"missing {p}")

    metrics_path = train_dir / "frete_lgbm_metrics.json"
    if metrics_path.is_file():
        m = json.loads(metrics_path.read_text(encoding="utf-8"))
        if m.get("success_gate", {}).get("passed") is not True:
            errors.append("frete B3-style gate not passed")
        if m.get("n_test", 0) < 50:
            errors.append("n_test too small")

    phase = ROOT / ".local" / "phases" / "32-real-ml-hardware-tests" / "TASKS.md"
    if not phase.is_file():
        # .local may be absent in CI checkout — only require when present
        pass

    if errors:
        print("check_phase32_ml_frete: FAIL", file=sys.stderr)
        for e in errors:
            print(f"  - {e}", file=sys.stderr)
        return 1
    print("check_phase32_ml_frete: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
