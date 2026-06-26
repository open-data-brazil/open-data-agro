#!/usr/bin/env python3
"""Verify docs/NEW-PROJECT-CHECKLIST.md reflects Go local-first post-collection state."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
CHECKLIST = ROOT / "docs" / "NEW-PROJECT-CHECKLIST.md"

STALE_PHRASES = (
    "TypeScript 5+",
    "pnpm",
    "before writing the first line of code",
    "Per-dataset rows` — fill as implemented",
)

REQUIRED_PHRASES = (
    "Go 1.22+",
    "make ci-collection-full-mvp",
    "check_official_sources_status.py",
    "47 catalog datasets",
)


def main() -> int:
    if not CHECKLIST.is_file():
        print(f"check_new_project_checklist: missing {CHECKLIST}", file=sys.stderr)
        return 1

    text = CHECKLIST.read_text(encoding="utf-8")
    errors: list[str] = []

    for phrase in STALE_PHRASES:
        if phrase in text:
            errors.append(f"stale phrase still present: {phrase!r}")

    for phrase in REQUIRED_PHRASES:
        if phrase not in text:
            errors.append(f"missing required phrase: {phrase!r}")

    if errors:
        print("check_new_project_checklist: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print("check_new_project_checklist: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
