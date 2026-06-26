#!/usr/bin/env python3
"""Verify R2 production runbook and env validation wiring."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
RUNBOOK = ROOT / "infra/r2/README.md"
VALIDATOR = ROOT / "scripts/deploy/validate_r2_env.sh"
MAKEFILE = ROOT / "Makefile"
ENV_EXAMPLE = ROOT / ".env.example"

RUNBOOK_SECTIONS = (
    "## Production deploy runbook",
    "### 3. Environment variables",
    "### 5. Verification",
    "### 6. Rollback",
    "make validate-r2-env",
    "make ci-validate-r2-env",
)

ENV_VARS = (
    "STORAGE_MODE=r2",
    "R2_ACCOUNT_ID",
    "R2_ACCESS_KEY_ID",
    "R2_SECRET_ACCESS_KEY",
    "R2_BUCKET",
    "R2_ENDPOINT",
)


def main() -> int:
    errors: list[str] = []

    if not RUNBOOK.is_file():
        errors.append(f"missing runbook: {RUNBOOK}")
    else:
        text = RUNBOOK.read_text(encoding="utf-8")
        for section in RUNBOOK_SECTIONS:
            if section not in text:
                errors.append(f"runbook missing section or target: {section!r}")

    if not VALIDATOR.is_file():
        errors.append(f"missing validator script: {VALIDATOR}")
    elif not VALIDATOR.read_text(encoding="utf-8").startswith("#!/usr/bin/env bash"):
        errors.append("validate_r2_env.sh must be bash")

    if not ENV_EXAMPLE.is_file():
        errors.append(f"missing {ENV_EXAMPLE}")
    else:
        env_text = ENV_EXAMPLE.read_text(encoding="utf-8")
        for var in ENV_VARS:
            if var not in env_text:
                errors.append(f".env.example missing {var!r}")

    if not MAKEFILE.is_file():
        errors.append(f"missing {MAKEFILE}")
    else:
        makefile = MAKEFILE.read_text(encoding="utf-8")
        for target in ("validate-r2-env:", "ci-validate-r2-env:", "validate-r2-env-live:"):
            if target not in makefile:
                errors.append(f"Makefile missing {target}")

    if errors:
        for err in errors:
            print(f"check_r2_runbook: {err}", file=sys.stderr)
        return 1

    print("check_r2_runbook: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
