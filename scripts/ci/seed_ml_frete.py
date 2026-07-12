#!/usr/bin/env python3
"""Seed synthetic CONAB frete gold for Phase 32 CI."""

from __future__ import annotations

import math
import os
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq

UFS = ["MT", "GO", "MS", "PR", "SP", "MG", "BA", "PI"]


def main() -> int:
    lake = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/ml-frete-ci-lake"))
    out = lake / "gold" / "mart_conab__frete"
    out.mkdir(parents=True, exist_ok=True)

    rows: list[dict[str, object]] = []
    i = 0
    for ano in range(2016, 2026):
        for mes in range(1, 13):
            for uo in UFS:
                for ud in UFS:
                    if uo == ud:
                        continue
                    i += 1
                    dist = 200.0 + 50.0 * ((ord(uo[0]) + ord(ud[0])) % 17)
                    # Short-cycle residual (not calendar-month seasonal) + mild noise.
                    cycle = 12.0 * math.sin(i / 11.0) + 6.0 * math.sin(i / 5.0)
                    season = 1.0 + 0.08 * math.sin(mes / 2.0)
                    corridor = 1.0 + 0.04 * ((ord(uo[0]) - ord(ud[0])) % 9)
                    ton = (0.35 * dist + 40.0) * season * corridor + cycle
                    # Encode cycle into lag-friendly distance jitter so model can learn it.
                    dist_feat = dist + 0.4 * cycle
                    rows.append(
                        {
                            "fonte": "CONTRATO",
                            "municipio_origem": f"ORIG-{uo}",
                            "cod_ibge_origem": "5100000",
                            "uf_origem": uo,
                            "municipio_destino": f"DEST-{ud}",
                            "cod_ibge_destino": "4100000",
                            "uf_destino": ud,
                            "ano": str(ano),
                            "mes": str(mes),
                            "distancia_km": f"{dist_feat:.1f}".replace(".", ","),
                            "valor_frete_tonelada": f"{ton:.2f}".replace(".", ","),
                            "valor_tonelada_km": f"{(ton / dist):.4f}".replace(".", ","),
                            "capturado_em": "2026-07-12T00:00:00Z",
                            "fonte_oficial": "https://portaldeinformacoes.conab.gov.br/",
                            "_dataset_id": "conab.frete",
                            "_source_file": "ci-seed",
                        }
                    )
    pq.write_table(pa.Table.from_pylist(rows), out / "mart.parquet")
    print(f"seeded {len(rows)} frete rows under {out / 'mart.parquet'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
