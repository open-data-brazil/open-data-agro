#!/usr/bin/env python3
"""Verify DATA-HISTORY-RANGES.md exists and references core agencies (Phase 33)."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
DOC = ROOT / "docs" / "DATA-HISTORY-RANGES.md"
REFRESH = ROOT / "docs" / "REFRESH-POLICY.md"
ROADMAP = ROOT / "docs" / "ROADMAP.md"

REQUIRED_SECTIONS = [
    "## CONAB",
    "## BCB SGS",
    "## INMET",
    "## IBGE",
    "## CEPEA",
    "## MDIC",
]

REQUIRED_DATASETS = [
    "conab.estimativa-graos",
    "bcb.sgs-ptax-usd-venda",
    "ibge.pam-area-quantidade",
    "inmet.bdmep-diario",
    "cepea.soja-paranagua",
    "mdic.comex-exportacao-ncm-mes",
]


def main() -> int:
    errors: list[str] = []

    if not DOC.is_file():
        errors.append(f"missing {DOC}")
    else:
        text = DOC.read_text(encoding="utf-8")
        for section in REQUIRED_SECTIONS:
            if section not in text:
                errors.append(f"DATA-HISTORY-RANGES missing section: {section}")
        for ds in REQUIRED_DATASETS:
            if ds not in text:
                errors.append(f"DATA-HISTORY-RANGES missing dataset: {ds}")

    if REFRESH.is_file():
        refresh = REFRESH.read_text(encoding="utf-8")
        if "DATA-HISTORY-RANGES" not in refresh:
            errors.append("REFRESH-POLICY.md must link to DATA-HISTORY-RANGES.md")
    else:
        errors.append(f"missing {REFRESH}")

    if ROADMAP.is_file():
        roadmap = ROADMAP.read_text(encoding="utf-8")
        if "DATA-HISTORY-RANGES" not in roadmap and "Phase 33" not in roadmap:
            errors.append("ROADMAP.md should mention Phase 33 or DATA-HISTORY-RANGES")
    else:
        errors.append(f"missing {ROADMAP}")

    if errors:
        for err in errors:
            print(f"ERROR: {err}", file=sys.stderr)
        return 1

    print("check_data_history_ranges: ok")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
