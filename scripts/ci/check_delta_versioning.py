#!/usr/bin/env python3
"""Verify native Delta Lake versioning is wired for CI."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
RUNNER = ROOT / "scripts/ci/run_delta_versioning.sh"
MAKEFILE = ROOT / "Makefile"
CI_YML = ROOT / ".github/workflows/ci.yml"
PROMOTE = ROOT / "scripts/delta/promote.py"
SILVER_README = ROOT / "lake/silver/README.md"


def main() -> int:
    errors: list[str] = []

    if not RUNNER.is_file():
        errors.append(f"missing runner: {RUNNER}")

    if not PROMOTE.is_file():
        errors.append(f"missing promote script: {PROMOTE}")
    else:
        promote = PROMOTE.read_text(encoding="utf-8")
        if "retention_configuration" not in promote:
            errors.append("promote.py must define retention_configuration(min_versions)")

    if not SILVER_README.is_file():
        errors.append(f"missing {SILVER_README}")
    else:
        readme = SILVER_README.read_text(encoding="utf-8")
        if "version => 0" not in readme:
            errors.append("lake/silver/README.md must document time travel")

    if not MAKEFILE.is_file():
        errors.append(f"missing {MAKEFILE}")
    else:
        makefile = MAKEFILE.read_text(encoding="utf-8")
        if "ci-delta-versioning:" not in makefile:
            errors.append("Makefile missing ci-delta-versioning target")

    if not CI_YML.is_file():
        errors.append(f"missing {CI_YML}")
    else:
        ci = CI_YML.read_text(encoding="utf-8")
        if "ci-delta-versioning" not in ci and "run_delta_versioning.sh" not in ci:
            errors.append("ci.yml must run ci-delta-versioning")

    if errors:
        for err in errors:
            print(f"check_delta_versioning: {err}", file=sys.stderr)
        return 1

    print("check_delta_versioning: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
