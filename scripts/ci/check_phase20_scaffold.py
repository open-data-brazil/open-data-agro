#!/usr/bin/env python3
"""Verify Phase 20 analytics-crossing scaffold exists."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
PHASE = ROOT / ".local" / "phases" / "20-analytics-crossing"
README = PHASE / "README.md"
TASKS = PHASE / "TASKS.md"
ROADMAP = ROOT / "docs" / "ROADMAP.md"

REQUIRED_README = (
    "DATA-CROSSING-VISION.md",
    "make ci-collection-full-mvp",
    "Explicit non-goals",
)


def main() -> int:
    errors: list[str] = []

    for path in (README, TASKS):
        if not path.is_file():
            errors.append(f"missing {path}")

    if README.is_file():
        text = README.read_text(encoding="utf-8")
        for phrase in REQUIRED_README:
            if phrase not in text:
                errors.append(f"README missing: {phrase!r}")

    if ROADMAP.is_file():
        roadmap = ROADMAP.read_text(encoding="utf-8")
        if "20-analytics-crossing" not in roadmap and "Phase 20" not in roadmap:
            errors.append("docs/ROADMAP.md missing Phase 20 reference")

    if errors:
        print("check_phase20_scaffold: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print("check_phase20_scaffold: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())
