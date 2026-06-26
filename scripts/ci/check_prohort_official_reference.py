#!/usr/bin/env python3
"""Verify Phase 11 OFFICIAL-REFERENCE documents Prohort column mappings."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
REF = ROOT / ".local" / "phases" / "11-conab-mercado" / "OFFICIAL-REFERENCE.md"

REQUIRED = (
    "conab.prohort-diario",
    "conab.prohort-mensal",
    "Profiler (live portal 2026-06-26)",
    "municipio_ceasa",
    "cod_ibge_municipio_ceasa",
    "make conab-mercado-prohort-mvp",
)

STALE = (
    "Catalog only",
    "Pending mappings",
)


def main() -> int:
    if not REF.is_file():
        print(f"check_prohort_official_reference: missing {REF}", file=sys.stderr)
        return 1

    text = REF.read_text(encoding="utf-8")
    errors: list[str] = []

    for phrase in STALE:
        if phrase in text:
            errors.append(f"stale phrase still present: {phrase!r}")

    for phrase in REQUIRED:
        if phrase not in text:
            errors.append(f"missing required phrase: {phrase!r}")

    if errors:
        print("check_prohort_official_reference: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print("check_prohort_official_reference: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
