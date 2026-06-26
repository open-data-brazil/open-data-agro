#!/usr/bin/env python3
"""Verify docs/NEW-DATASET-CHECKLIST.md exists and reference exemplar paths are present."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
CHECKLIST = ROOT / "docs" / "NEW-DATASET-CHECKLIST.md"

REQUIRED_SECTIONS = (
    "## 1. Sample fixtures",
    "## 2. Golden ingest test",
    "## 3. Great Expectations bronze suite",
    "## 4. dbt staging and marts",
    "## 5. CI silver seed",
    "## 6. Makefile MVP gate",
    "## 7. Official sources catalog",
    "## 8. Phase reference doc",
    "## 9. Changelog",
    "## 10. Validation commands",
    "check_official_sources_status.py",
    "conab.estimativa-graos",
)

# Paths cited as exemplar for conab.estimativa-graos — proves checklist matches repo reality.
EXEMPLAR_PATHS = (
    "internal/conab/testdata",
    "expectations/suites/bronze/conab/estimativa_graos.json",
    "internal/processor/quality.go",
    "dbt/models/staging/conab/_conab__sources.yml",
    "scripts/ci/seed_dbt_silver.py",
    "docs/OFFICIAL-SOURCES.md",
)


def main() -> int:
    errors: list[str] = []

    if not CHECKLIST.is_file():
        errors.append(f"missing {CHECKLIST}")
    else:
        text = CHECKLIST.read_text(encoding="utf-8")
        for section in REQUIRED_SECTIONS:
            if section not in text:
                errors.append(f"checklist missing: {section!r}")

    for rel in EXEMPLAR_PATHS:
        path = ROOT / rel
        if not path.exists():
            errors.append(f"exemplar path missing: {rel}")

    if errors:
        for err in errors:
            print(f"check_new_dataset_checklist: {err}", file=sys.stderr)
        return 1

    print("check_new_dataset_checklist: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
