#!/usr/bin/env python3
"""Verify Phase 20 analytics-crossing deliverables (docs + gold mart)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
PHASE = ROOT / ".local" / "phases" / "20-analytics-crossing"
README = PHASE / "README.md"
TASKS = PHASE / "TASKS.md"
UC = ROOT / "docs" / "use-cases" / "UC-ML-001-price-early-warning.md"
ADR = ROOT / "docs" / "adr" / "005-cross-source-dbt-analytics.md"
ROADMAP = ROOT / "docs" / "ROADMAP.md"
GLOSSARY = ROOT / "docs" / "GLOSSARY.md"


def main() -> int:
    errors: list[str] = []
    lake = Path(os.environ.get("LAKE_LOCAL_ROOT", ROOT / "lake"))
    mart = lake / "gold" / "mart_ml__soy_daily_features" / "mart.parquet"
    dim = lake / "gold" / "dim_municipio" / "mart.parquet"

    for path in (README, TASKS, UC, ADR, ROADMAP, GLOSSARY):
        if not path.is_file():
            errors.append(f"missing {path.relative_to(ROOT)}")

    if UC.is_file():
        text = UC.read_text(encoding="utf-8")
        for phrase in (
            "preco_rs_sc",
            "7 / 30 / 90",
            "CC BY-NC",
            "as_of",
            "Publication lags",
        ):
            if phrase not in text:
                errors.append(f"UC-ML-001 missing: {phrase!r}")

    if ADR.is_file():
        adr = ADR.read_text(encoding="utf-8")
        for phrase in ("int_cross__", "mart_ml__", "No bronze changes"):
            if phrase not in adr:
                errors.append(f"ADR 005 missing: {phrase!r}")

    if GLOSSARY.is_file():
        gloss = GLOSSARY.read_text(encoding="utf-8")
        if "Analytics resampling" not in gloss and "analytics resampling" not in gloss.lower():
            errors.append("GLOSSARY missing analytics resampling policy")
        if "dim_municipio" not in gloss:
            errors.append("GLOSSARY missing dim_municipio")

    if ROADMAP.is_file():
        roadmap = ROADMAP.read_text(encoding="utf-8")
        if "Phase 20" not in roadmap:
            errors.append("docs/ROADMAP.md missing Phase 20 reference")
        if "mart_ml__soy_daily_features" not in roadmap and "analytics-crossing" not in roadmap:
            errors.append("docs/ROADMAP.md missing Phase 20 analytics deliverable")

    if not mart.is_file():
        errors.append(f"missing gold mart {mart}")
    if not dim.is_file():
        errors.append(f"missing gold dim {dim}")

    if errors:
        print("check_phase20_scaffold: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print("check_phase20_scaffold: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
