#!/usr/bin/env python3
"""Verify IMPLEMENTATION-PLAN.md reflects full E2E for Phases 15-19."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
PLAN = ROOT / ".local" / "IMPLEMENTATION-PLAN.md"

STALE_PHRASES = (
    "ingest only",
    "apenas bronze + qualidade + staging",
    "Outras fontes — coleta (Phases 15–19)",
)

REQUIRED_PHRASES = (
    "full E2E local",
    "make ci-collection-full-mvp",
    "make ibge-localidades-mvp",
    "make collection-macro-mvp",
)


def main() -> int:
    if not PLAN.is_file():
        print(f"check_implementation_plan: missing {PLAN}", file=sys.stderr)
        return 1

    text = PLAN.read_text(encoding="utf-8")
    errors: list[str] = []

    for phrase in STALE_PHRASES:
        if phrase in text:
            errors.append(f"stale phrase still present: {phrase!r}")

    for phrase in REQUIRED_PHRASES:
        if phrase not in text:
            errors.append(f"missing required phrase: {phrase!r}")

    if errors:
        print("check_implementation_plan: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print("check_implementation_plan: PASS (Phases 15-19 documented as full E2E)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
