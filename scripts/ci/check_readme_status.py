#!/usr/bin/env python3
"""Verify README.md reflects post-collection sprint status."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
README = ROOT / "README.md"

STALE_PHRASES = (
    "Phase 0 Go platform scaffold",
)

REQUIRED_PHRASES = (
    "Collection sprint complete",
    "make ci-collection-full-mvp",
    "make ci-go",
    "docs/ROADMAP.md",
)


def main() -> int:
    if not README.is_file():
        print(f"check_readme_status: missing {README}", file=sys.stderr)
        return 1

    text = README.read_text(encoding="utf-8")
    errors: list[str] = []

    for phrase in STALE_PHRASES:
        if phrase in text:
            errors.append(f"stale phrase still present: {phrase!r}")

    for phrase in REQUIRED_PHRASES:
        if phrase not in text:
            errors.append(f"missing required phrase: {phrase!r}")

    if errors:
        print("check_readme_status: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print("check_readme_status: PASS (README reflects collection sprint exit)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
