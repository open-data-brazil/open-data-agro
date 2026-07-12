#!/usr/bin/env python3
"""Rebuild mart_ml__soy_daily_features locally from gold (Phase 32 R2).

Implements the UC-ML-001 daily spine joins without requiring a full dbt rebuild
when optional sources (INMET/B3/Comex) are sparse.
"""

from __future__ import annotations

import argparse
import math
import os
import sys
from datetime import date, datetime, timedelta
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]


def _as_date(v: object) -> date | None:
    if v is None:
        return None
    if isinstance(v, date) and not isinstance(v, datetime):
        return v
    s = str(v)[:10]
    try:
        return date.fromisoformat(s)
    except ValueError:
        return None


def _f(v: object) -> float | None:
    if v is None or v == "":
        return None
    if isinstance(v, (int, float)):
        return float(v)
    s = str(v).strip().replace(".", "").replace(",", ".") if "," in str(v) else str(v).strip()
    try:
        return float(s)
    except ValueError:
        return None


def _read_gold(lake: Path, mart: str) -> list[dict[str, object]]:
    path = lake / "gold" / mart / "mart.parquet"
    if not path.is_file():
        files = list((lake / "gold" / mart).rglob("*.parquet"))
        if not files:
            return []
        path = files[0]
    return pq.read_table(path).to_pylist()


def _month_end(d: date) -> date:
    if d.month == 12:
        return date(d.year, 12, 31)
    return date(d.year, d.month + 1, 1) - timedelta(days=1)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--lake-root", default=os.environ.get("LAKE_LOCAL_ROOT", str(ROOT / "lake")))
    parser.add_argument("--min-rows", type=int, default=1000)
    args = parser.parse_args()
    lake = Path(args.lake_root)

    cepea_raw = _read_gold(lake, "mart_cepea__soja_paranagua")
    cepea: list[dict[str, object]] = []
    for r in cepea_raw:
        d = _as_date(r.get("data"))
        p = _f(r.get("preco_rs_sc"))
        if d is None or p is None:
            continue
        cepea.append(
            {
                "data": d,
                "preco_rs_sc": p,
                "preco_usd_sc": _f(r.get("preco_usd_sc")),
                "variacao_dia_pct": _f(r.get("variacao_dia_pct")),
            }
        )
    cepea.sort(key=lambda r: r["data"])  # type: ignore[arg-type]
    if len(cepea) < args.min_rows:
        print(f"CEPEA gold too short: {len(cepea)} < {args.min_rows}", file=sys.stderr)
        return 1

    ptax_map: dict[date, float] = {}
    for r in _read_gold(lake, "mart_bcb__sgs_ptax_usd_venda"):
        d = _as_date(r.get("data"))
        v = _f(r.get("valor"))
        if d is not None and v is not None:
            ptax_map[d] = v

    frete_rows = []
    for r in _read_gold(lake, "mart_conab__frete"):
        if str(r.get("uf_origem", "")).upper() != "MT" or str(r.get("uf_destino", "")).upper() != "PR":
            continue
        try:
            ano = int(str(r.get("ano")).strip())
            mes = int(str(r.get("mes")).strip())
        except (TypeError, ValueError):
            continue
        ton = _f(r.get("valor_frete_tonelada"))
        tkm = _f(r.get("valor_tonelada_km"))
        if ton is None:
            continue
        ref = date(ano, mes, 1)
        as_of = _month_end(ref) + timedelta(days=14)
        frete_rows.append((as_of, ton, tkm if tkm is not None else ton / max(_f(r.get("distancia_km")) or 1.0, 1.0)))
    frete_rows.sort()

    out_rows: list[dict[str, object]] = []
    prices = [float(r["preco_rs_sc"]) for r in cepea]  # type: ignore[arg-type]
    for i, row in enumerate(cepea):
        d: date = row["data"]  # type: ignore[assignment]
        price = float(row["preco_rs_sc"])  # type: ignore[arg-type]
        lag1 = prices[i - 1] if i >= 1 else None
        lag7 = prices[i - 7] if i >= 7 else None
        window = prices[max(0, i - 6) : i + 1]
        vol = None
        if len(window) >= 2:
            m = sum(window) / len(window)
            vol = math.sqrt(sum((x - m) ** 2 for x in window) / (len(window) - 1))

        # ASOF ptax
        ptax = None
        as_of_ptax = None
        for pd_ in sorted(ptax_map):
            if pd_ <= d:
                ptax = ptax_map[pd_]
                as_of_ptax = pd_
            else:
                break

        usd = row.get("preco_usd_sc")
        if usd is None and ptax:
            usd = price / ptax

        # ASOF frete
        frete_ton = frete_tkm = None
        as_of_frete = None
        for as_of, ton, tkm in frete_rows:
            if as_of <= d:
                frete_ton, frete_tkm, as_of_frete = ton, tkm, as_of
            else:
                break

        ret1 = (price / lag1 - 1.0) if lag1 else None
        ret7 = (price / lag7 - 1.0) if lag7 else None
        var = row.get("variacao_dia_pct")
        if var is None and ret1 is not None:
            var = ret1 * 100.0

        out_rows.append(
            {
                "produto_slug": "soja",
                "data": d,
                "label_date": d,
                "preco_rs_sc": price,
                "preco_usd_sc": usd,
                "variacao_dia_pct": var,
                "cepea_ret_1d": ret1,
                "cepea_ret_7d": ret7,
                "cepea_vol_7d": vol,
                "ptax_usd_venda": ptax,
                "b3_front_price": None,
                "b3_front_symbol": None,
                "basis_cepea_usd_minus_b3": None,
                "comex_soja_fob_usd": None,
                "comex_soja_kg": None,
                "frete_mt_pr_ton_avg": frete_ton,
                "frete_mt_pr_tkm_avg": frete_tkm,
                "inmet_national_avg": None,
                "as_of_cepea": d,
                "as_of_ptax": as_of_ptax,
                "as_of_b3": None,
                "as_of_comex": None,
                "as_of_frete": as_of_frete,
                "as_of_inmet": None,
                "feature_as_of_max": d,
            }
        )

    out_dir = lake / "gold" / "mart_ml__soy_daily_features"
    out_dir.mkdir(parents=True, exist_ok=True)
    pq.write_table(pa.Table.from_pylist(out_rows), out_dir / "mart.parquet")
    print(f"wrote {len(out_rows)} rows → {out_dir / 'mart.parquet'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
