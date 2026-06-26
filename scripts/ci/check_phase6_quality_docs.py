#!/usr/bin/env python3
"""Verify Phase 6 docs describe GE vs validate_codigo_ibge split."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
PHASE_README = ROOT / ".local" / "phases" / "06-quality-great-expectations" / "README.md"
QUALITY_README = ROOT / "scripts" / "quality" / "README.md"

REQUIRED_IN_PHASE = (
    "GE vs `validate_codigo_ibge.py`",
    "Referential integrity",
    "make ci-validate-codigo-ibge",
    "deferred",
)

REQUIRED_IN_QUALITY = (
    "validate_codigo_ibge.py",
    "make ci-validate-codigo-ibge",
    "Great Expectations",
)


def main() -> int:
    errors: list[str] = []

    if not PHASE_README.is_file():
        errors.append(f"missing {PHASE_README}")
    else:
        phase_text = PHASE_README.read_text(encoding="utf-8")
        for phrase in REQUIRED_IN_PHASE:
            if phrase not in phase_text:
                errors.append(f"phase README missing: {phrase!r}")

    if not QUALITY_README.is_file():
        errors.append(f"missing {QUALITY_README}")
    else:
        quality_text = QUALITY_README.read_text(encoding="utf-8")
        for phrase in REQUIRED_IN_QUALITY:
            if phrase not in quality_text:
                errors.append(f"quality README missing: {phrase!r}")

    if errors:
        print("check_phase6_quality_docs: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print("check_phase6_quality_docs: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
