#!/usr/bin/env python3
"""Verify benchmark profile paths and dataset counts match documentation."""

from __future__ import annotations

import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
PROFILES = ROOT / "scripts" / "benchmark" / "profiles"
README = ROOT / ".local" / "benchmark" / "README.md"

EXPECTED = {
    "fast10.json": 16,
    "fast10-stress.json": 18,
}


def count_datasets(profile: dict) -> int:
    return len(profile["datasets"])


def main() -> int:
    errors: list[str] = []

    for name, expected in EXPECTED.items():
        path = PROFILES / name
        if not path.is_file():
            errors.append(f"missing profile: {path}")
            continue
        profile = json.loads(path.read_text(encoding="utf-8"))
        got = count_datasets(profile)
        if got != expected:
            errors.append(f"{name}: expected {expected} datasets, got {got}")

    readme = README.read_text(encoding="utf-8") if README.is_file() else ""
    if ".local/benchmark/profiles/" in readme:
        errors.append("README still references .local/benchmark/profiles/ (use scripts/benchmark/profiles/)")
    if "**16**" not in readme or "**18**" not in readme:
        errors.append("README missing fast10=16 / fast10-stress=18 dataset counts")
    if "benchmark-ingestor-fast10-stress" not in readme:
        errors.append("README missing make benchmark-ingestor-fast10-stress")

    if errors:
        print("check_benchmark_profiles: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print(
        f"check_benchmark_profiles: PASS (fast10={EXPECTED['fast10.json']}, "
        f"fast10-stress={EXPECTED['fast10-stress.json']})"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
