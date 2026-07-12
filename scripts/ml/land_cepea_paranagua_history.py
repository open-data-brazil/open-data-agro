#!/usr/bin/env python3
"""Land CEPEA soja Paranaguá history bulk into bronze/silver/gold (Phase 32 R1)."""

from __future__ import annotations

import argparse
import json
import os
import sys
from datetime import datetime, timezone
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]


def _load_bulk(path: Path) -> list[dict[str, object]]:
    raw = json.loads(path.read_text(encoding="utf-8"))
    if isinstance(raw, list):
        return raw
    return list(raw.get("observations") or [])


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--bulk",
        default=os.environ.get(
            "CEPEA_BULK_PATH",
            str(ROOT / ".local" / "ml" / "bulk" / "cepea_soja_paranagua_history.json"),
        ),
    )
    parser.add_argument(
        "--lake-root",
        default=os.environ.get("LAKE_LOCAL_ROOT", str(ROOT / "lake")),
    )
    parser.add_argument("--min-rows", type=int, default=1000)
    args = parser.parse_args()

    bulk_path = Path(args.bulk)
    if not bulk_path.is_file():
        print(f"missing bulk: {bulk_path}", file=sys.stderr)
        return 1
    observations = _load_bulk(bulk_path)
    if len(observations) < args.min_rows:
        print(f"bulk too small: {len(observations)} < {args.min_rows}", file=sys.stderr)
        return 1

    lake = Path(args.lake_root)
    now = datetime.now(timezone.utc).isoformat()
    source = str(bulk_path)

    rows = []
    for obs in observations:
        data = str(obs["data"])[:10]
        rows.append(
            {
                "produto": "soja",
                "praca": "Paranaguá",
                "data": data,
                "preco_rs_sc": str(obs.get("preco_rs_sc") or ""),
                "variacao_dia_pct": str(obs.get("variacao_dia_pct") or ""),
                "preco_usd_sc": str(obs.get("preco_usd_sc") or ""),
                "ano": data[:4],
                "capturado_em": now,
                "fonte_oficial": "https://www.cepea.org.br/br/indicador/soja.aspx",
                "_dataset_id": "cepea.soja-paranagua",
                "_source_file": source,
                "_fill_source": obs.get("_fill_source", "bulk"),
            }
        )

    table = pa.Table.from_pylist(rows)

    bronze_dir = lake / "bronze" / "cepea" / "soja-paranagua" / f"ingest_date={now[:10]}"
    bronze_dir.mkdir(parents=True, exist_ok=True)
    pq.write_table(table, bronze_dir / "part-history.parquet")

    silver_dir = lake / "silver" / "cepea" / "soja_paranagua"
    silver_dir.mkdir(parents=True, exist_ok=True)
    # Prefer Delta when available; always write mart parquet for dbt gold rebuild.
    try:
        from deltalake import write_deltalake

        write_deltalake(str(silver_dir), table, mode="overwrite")
    except Exception as exc:  # noqa: BLE001
        print(f"delta write skipped: {exc}")
        pq.write_table(table, silver_dir / "part-00000.parquet")

    gold_dir = lake / "gold" / "mart_cepea__soja_paranagua"
    gold_dir.mkdir(parents=True, exist_ok=True)
    gold_cols = [
        "produto",
        "praca",
        "data",
        "preco_rs_sc",
        "variacao_dia_pct",
        "preco_usd_sc",
        "ano",
        "capturado_em",
        "fonte_oficial",
        "_dataset_id",
        "_source_file",
    ]
    gold = table.select([c for c in gold_cols if c in table.column_names])
    pq.write_table(gold, gold_dir / "mart.parquet")

    print(f"landed {len(rows)} CEPEA soja Paranaguá rows → {gold_dir / 'mart.parquet'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
