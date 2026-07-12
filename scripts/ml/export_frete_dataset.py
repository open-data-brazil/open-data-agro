#!/usr/bin/env python3
"""Export CONAB frete training dataset with corridor lags + walk-forward split."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import subprocess
import sys
from datetime import date, datetime, timezone
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.config_frete import (  # noqa: E402
    MIN_EXPORT_ROWS,
    SCHEMA_VERSION,
    SPLIT_TEST_YEAR_START,
    SPLIT_TRAIN_YEAR_END,
    SPLIT_VAL_YEAR_END,
    SPLIT_VAL_YEAR_START,
    TARGET_COL,
)


def _git_sha() -> str:
    try:
        return subprocess.check_output(
            ["git", "rev-parse", "HEAD"], cwd=ROOT, stderr=subprocess.DEVNULL, text=True
        ).strip()
    except (subprocess.CalledProcessError, FileNotFoundError):
        return "unknown"


def _br_float(raw: object) -> float | None:
    if raw is None:
        return None
    if isinstance(raw, (int, float)):
        return float(raw)
    s = str(raw).strip().replace(".", "").replace(",", ".")
    if not s:
        return None
    try:
        return float(s)
    except ValueError:
        return None


def _assign_split(ano: int) -> str:
    if ano <= SPLIT_TRAIN_YEAR_END:
        return "train"
    if SPLIT_VAL_YEAR_START <= ano <= SPLIT_VAL_YEAR_END:
        return "val"
    if ano >= SPLIT_TEST_YEAR_START:
        return "test"
    return "holdout"


def _fonte_code(fonte: str | None) -> int:
    key = (fonte or "").strip().upper()
    mapping = {"CONTRATO": 0, "OFERTA": 1, "PESQUISA": 2}
    return mapping.get(key, 3)


def load_frete_mart(path: Path) -> list[dict[str, object]]:
    if not path.is_file():
        raise FileNotFoundError(f"missing frete mart: {path}")
    return pq.read_table(path).to_pylist()


def build_training_rows(raw: list[dict[str, object]]) -> list[dict[str, object]]:
    parsed: list[dict[str, object]] = []
    for row in raw:
        ano = row.get("ano")
        mes = row.get("mes")
        try:
            ano_i = int(str(ano).strip())
            mes_i = int(str(mes).strip())
        except (TypeError, ValueError):
            continue
        if mes_i < 1 or mes_i > 12:
            continue
        ton = _br_float(row.get("valor_frete_tonelada"))
        tkm = _br_float(row.get("valor_tonelada_km"))
        dist = _br_float(row.get("distancia_km"))
        if ton is None or dist is None or dist <= 0:
            continue
        uf_o = str(row.get("uf_origem") or "").strip().upper()
        uf_d = str(row.get("uf_destino") or "").strip().upper()
        if len(uf_o) != 2 or len(uf_d) != 2:
            continue
        import calendar
        from datetime import timedelta

        refmonth = date(ano_i, mes_i, 1)
        # as_of ≈ month_end + 14d (matches UC-ML-001 frete lag)
        month_end = date(ano_i, mes_i, calendar.monthrange(ano_i, mes_i)[1])
        as_of = month_end + timedelta(days=14)
        parsed.append(
            {
                "fonte": row.get("fonte"),
                "fonte_code": _fonte_code(str(row.get("fonte") or "")),
                "municipio_origem": row.get("municipio_origem"),
                "municipio_destino": row.get("municipio_destino"),
                "cod_ibge_origem": str(row.get("cod_ibge_origem") or "").zfill(7),
                "cod_ibge_destino": str(row.get("cod_ibge_destino") or "").zfill(7),
                "uf_origem": uf_o,
                "uf_destino": uf_d,
                "ano": ano_i,
                "mes": mes_i,
                "refmonth": refmonth.isoformat(),
                "as_of": as_of.isoformat(),
                "distancia_km": dist,
                TARGET_COL: ton,
                "valor_tonelada_km": tkm if tkm is not None else (ton / dist),
                "split": _assign_split(ano_i),
                "corridor_key": f"{uf_o}->{uf_d}",
            }
        )

    # Corridor lag features (same UF pair, prior months) — no future leakage.
    by_corridor: dict[str, list[dict[str, object]]] = {}
    for row in parsed:
        by_corridor.setdefault(str(row["corridor_key"]), []).append(row)
    for rows in by_corridor.values():
        rows.sort(key=lambda r: (int(r["ano"]), int(r["mes"])))  # type: ignore[arg-type]
        # Monthly mean ton for lag lookup
        month_avg: dict[tuple[int, int], float] = {}
        month_groups: dict[tuple[int, int], list[float]] = {}
        for r in rows:
            key = (int(r["ano"]), int(r["mes"]))  # type: ignore[arg-type]
            month_groups.setdefault(key, []).append(float(r[TARGET_COL]))  # type: ignore[arg-type]
        for key, vals in month_groups.items():
            month_avg[key] = sum(vals) / len(vals)

        for r in rows:
            y, m = int(r["ano"]), int(r["mes"])  # type: ignore[arg-type]
            # lag 1 month
            py, pm = (y - 1, 12) if m == 1 else (y, m - 1)
            r["corridor_lag1_ton_avg"] = month_avg.get((py, pm))
            # lag 12 months
            r["corridor_lag12_ton_avg"] = month_avg.get((y - 1, m))

    return parsed


def write_export(rows: list[dict[str, object]], out_dir: Path, source_path: Path) -> Path:
    if len(rows) < MIN_EXPORT_ROWS:
        raise ValueError(f"frete export too small: {len(rows)} < {MIN_EXPORT_ROWS}")

    # Drop helper corridor_key from parquet (keep in memory only).
    clean = [{k: v for k, v in r.items() if k != "corridor_key"} for r in rows]
    table = pa.Table.from_pylist(clean)
    out_dir.mkdir(parents=True, exist_ok=True)
    parquet_path = out_dir / "frete_training.parquet"
    pq.write_table(table, parquet_path)

    cols = table.column_names
    schema_hash = hashlib.sha256("|".join(cols).encode()).hexdigest()[:16]
    anos = [int(r["ano"]) for r in clean]  # type: ignore[arg-type]
    manifest = {
        "dataset": "ml.frete-training",
        "schema_version": SCHEMA_VERSION,
        "schema_hash": schema_hash,
        "row_count": len(clean),
        "columns": cols,
        "year_min": min(anos),
        "year_max": max(anos),
        "source_mart": str(source_path),
        "git_sha": _git_sha(),
        "exported_at": datetime.now(timezone.utc).isoformat(),
        "split_policy": {
            "train": f"ano <= {SPLIT_TRAIN_YEAR_END}",
            "val": f"{SPLIT_VAL_YEAR_START}..{SPLIT_VAL_YEAR_END}",
            "test": f"ano >= {SPLIT_TEST_YEAR_START}",
        },
        "target": TARGET_COL,
    }
    (out_dir / "manifest.json").write_text(
        json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    return parquet_path


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--lake-root",
        default=os.environ.get("LAKE_LOCAL_ROOT", str(ROOT / "lake")),
    )
    parser.add_argument(
        "--out-dir",
        default=os.environ.get("ML_FRETE_EXPORT_DIR", str(ROOT / ".local" / "ml" / "export_frete")),
    )
    args = parser.parse_args()

    source = Path(args.lake_root) / "gold" / "mart_conab__frete" / "mart.parquet"
    try:
        raw = load_frete_mart(source)
        rows = build_training_rows(raw)
        path = write_export(rows, Path(args.out_dir), source)
    except (FileNotFoundError, ValueError) as exc:
        print(str(exc), file=sys.stderr)
        return 1

    print(f"wrote {path} rows={len(rows)}")
    print(f"wrote {Path(args.out_dir) / 'manifest.json'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
