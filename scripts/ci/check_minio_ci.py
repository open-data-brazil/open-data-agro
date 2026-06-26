#!/usr/bin/env python3
"""Verify MinIO integration is wired for CI (Makefile + GitHub Actions)."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
MAKEFILE = ROOT / "Makefile"
CI_YML = ROOT / ".github/workflows/ci.yml"
RUNNER = ROOT / "scripts/ci/run_minio_integration.sh"


def main() -> int:
    errors: list[str] = []

    if not RUNNER.is_file():
        errors.append(f"missing runner script: {RUNNER}")
    elif not RUNNER.read_text(encoding="utf-8").startswith("#!/usr/bin/env bash"):
        errors.append(f"runner must be bash: {RUNNER}")

    if not MAKEFILE.is_file():
        errors.append(f"missing {MAKEFILE}")
    else:
        makefile = MAKEFILE.read_text(encoding="utf-8")
        if "ci-minio:" not in makefile:
            errors.append("Makefile missing ci-minio target")
        if "run_minio_integration.sh" not in makefile:
            errors.append("Makefile ci-minio must invoke run_minio_integration.sh")

    if not CI_YML.is_file():
        errors.append(f"missing {CI_YML}")
    else:
        ci = CI_YML.read_text(encoding="utf-8")
        if "make ci-minio" not in ci and "run_minio_integration.sh" not in ci:
            errors.append("ci.yml must run make ci-minio or run_minio_integration.sh")
        if "MINIO_INTEGRATION" not in ci and "make ci-minio" not in ci:
            errors.append("ci.yml must exercise MinIO integration tests")

    if errors:
        for err in errors:
            print(f"check_minio_ci: {err}", file=sys.stderr)
        return 1

    print("check_minio_ci: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
