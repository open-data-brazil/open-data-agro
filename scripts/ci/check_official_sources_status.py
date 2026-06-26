#!/usr/bin/env python3
"""Verify docs/OFFICIAL-SOURCES.md status column matches catalog reality."""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
SOURCES = ROOT / "docs" / "OFFICIAL-SOURCES.md"
CATALOG = ROOT / "configs" / "catalog"
VISION = ROOT / ".local" / "DATA-CROSSING-VISION.md"

# João P1 collection rows from DATA-CROSSING-VISION (concrete dataset_id or prefix).
VISION_DATASET_IDS = {
    "conab.estimativa-graos",
    "conab.serie-historica-graos",
    "conab.custo-producao",
    "conab.oferta-demanda",
    "conab.frete",
    "conab.armazenagem",
    "conab.serie-historica-capacidade-estatica",
    "ibge.localidades-municipios",
    "bcb.sgs-ipca",
    "bcb.sgs-ptax-usd-venda",
    "cepea.soja-paranagua",
    "ibge.pam-area-quantidade",
}

VISION_PREFIXES = (
    "conab.precos-agropecuarios-",
    "inmet.",
)


def catalog_dataset_ids() -> set[str]:
    ids: set[str] = set()
    for path in CATALOG.rglob("*.yaml"):
        if path.name == "registry.yaml":
            continue
        text = path.read_text(encoding="utf-8")
        ids.update(re.findall(r"dataset_id:\s*([a-z0-9._-]+)", text))
    return ids


def parse_official_sources() -> dict[str, str]:
    text = SOURCES.read_text(encoding="utf-8")
    rows: dict[str, str] = {}
    for line in text.splitlines():
        m = re.match(r"\|\s*`([^`]+)`\s*\|[^|]+\|\s*\*\*([^*]+)\*\*\s*\|", line)
        if m:
            rows[m.group(1)] = m.group(2).strip()
    return rows


def vision_covered(ids: set[str]) -> set[str]:
    covered: set[str] = set()
    for dataset_id in ids:
        if dataset_id in VISION_DATASET_IDS:
            covered.add(dataset_id)
            continue
        for prefix in VISION_PREFIXES:
            if dataset_id.startswith(prefix):
                covered.add(dataset_id)
                break
    return covered


def main() -> int:
    catalog_ids = catalog_dataset_ids()
    statuses = parse_official_sources()
    errors: list[str] = []

    doc_text = SOURCES.read_text(encoding="utf-8")
    if "full pipeline" in doc_text.lower():
        errors.append('docs/OFFICIAL-SOURCES.md still contains "full pipeline"')

    missing = sorted(catalog_ids - statuses.keys())
    if missing:
        errors.append(f"catalog datasets missing from OFFICIAL-SOURCES: {missing}")

    extra = sorted(statuses.keys() - catalog_ids)
    if extra:
        errors.append(f"OFFICIAL-SOURCES rows not in catalog: {extra}")

    for dataset_id in sorted(catalog_ids):
        status = statuses.get(dataset_id, "")
        if "implemented" not in status.lower():
            errors.append(f"{dataset_id}: status must contain 'implemented', got {status!r}")

    vision_ids = vision_covered(catalog_ids)
    for dataset_id in sorted(vision_ids):
        status = statuses.get(dataset_id, "")
        if "implemented" not in status.lower():
            errors.append(f"DATA-CROSSING-VISION row {dataset_id}: not implemented in OFFICIAL-SOURCES")

    if not VISION.exists():
        errors.append(f"missing {VISION}")

    if errors:
        print("check_official_sources_status: FAIL", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print(f"check_official_sources_status: PASS ({len(catalog_ids)} catalog datasets, {len(vision_ids)} vision rows)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
