#!/usr/bin/env python3
"""Unit tests for frete export helpers."""

from __future__ import annotations

import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.export_frete_dataset import build_training_rows, write_export  # noqa: E402
from ml.metrics import gate_passes  # noqa: E402


def test_parse_and_lags() -> None:
    raw = []
    for ano in (2020, 2021):
        for mes in range(1, 13):
            raw.append(
                {
                    "fonte": "CONTRATO",
                    "municipio_origem": "A-MT",
                    "cod_ibge_origem": "5103403",
                    "uf_origem": "MT",
                    "municipio_destino": "B-PR",
                    "cod_ibge_destino": "4115200",
                    "uf_destino": "PR",
                    "ano": str(ano),
                    "mes": str(mes),
                    "distancia_km": "1.200,5",
                    "valor_frete_tonelada": f"{100 + mes},5",
                    "valor_tonelada_km": "0,15",
                }
            )
    rows = build_training_rows(raw)
    assert len(rows) == 24
    # Second year January should have lag12 from prior January
    jan21 = next(r for r in rows if r["ano"] == 2021 and r["mes"] == 1)
    assert jan21["corridor_lag12_ton_avg"] is not None
    assert jan21["split"] == "train"


def test_write_export_min_rows() -> None:
    raw = []
    for i in range(600):
        ano = 2018 + (i % 8)
        mes = (i % 12) + 1
        raw.append(
            {
                "fonte": "CONTRATO",
                "municipio_origem": "A",
                "cod_ibge_origem": "5100000",
                "uf_origem": "MT",
                "municipio_destino": "B",
                "cod_ibge_destino": "4100000",
                "uf_destino": "PR",
                "ano": str(ano),
                "mes": str(mes),
                "distancia_km": "500,0",
                "valor_frete_tonelada": "120,0",
                "valor_tonelada_km": "0,24",
            }
        )
    rows = build_training_rows(raw)
    with tempfile.TemporaryDirectory() as tmp:
        path = write_export(rows, Path(tmp), Path("fake.parquet"))
        assert path.is_file()
        assert (Path(tmp) / "manifest.json").is_file()


def main() -> int:
    assert gate_passes(0.85, 1.0, 0.15)
    test_parse_and_lags()
    test_write_export_min_rows()
    print("test_ml_frete_unit: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
