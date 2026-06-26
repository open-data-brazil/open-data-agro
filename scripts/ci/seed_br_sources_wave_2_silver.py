#!/usr/bin/env python3
"""Seed minimal silver Delta for BR sources wave 2 CI (Phase 40)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
from deltalake import write_deltalake


def write_table(root: Path, agency: str, table: str, data: pa.Table) -> None:
    path = root / "silver" / agency / table
    path.parent.mkdir(parents=True, exist_ok=True)
    write_deltalake(str(path), data, mode="overwrite")


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_mapa__agrofit_produtos_formulados").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    agrofit = pa.table(
        {
            "NR_REGISTRO": ["35523", "35523"],
            "MARCA_COMERCIAL": ["KBR-829M1-02", "KBR-829M1-02"],
            "FORMULACAO": ["Nematóides vivos", "Nematóides vivos"],
            "INGREDIENTE_ATIVO": ["Heterorhabditis bacteriophora", "Heterorhabditis bacteriophora"],
            "TITULAR_DE_REGISTRO": ["Koppert do Brasil Holding S.A.", "Koppert do Brasil Holding S.A."],
            "CLASSE": ["Agente Biológico de Controle", "Agente Biológico de Controle"],
            "MODO_DE_ACAO": ["", ""],
            "CULTURA": ["Todas as culturas", "Todas as culturas"],
            "PRAGA_NOME_CIENTIFICO": ["Scaptocoris castanea", "Spodoptera frugiperda"],
            "PRAGA_NOME_COMUM": ["Percevejo-castanho", "Lagarta-militar"],
            "EMPRESA_PAIS_TIPO": ["BRASIL", "BRASIL"],
            "CLASSE_TOXICOLOGICA": ["Não Classificado", "Não Classificado"],
            "CLASSE_AMBIENTAL": ["Produto Pouco Perigoso", "Produto Pouco Perigoso"],
            "ORGANICOS": ["NAO", "NAO"],
            "SITUACAO": ["TRUE", "TRUE"],
            "_dataset_id": ["mapa.agrofit-produtos-formulados", "mapa.agrofit-produtos-formulados"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "mapa", "agrofit_produtos_formulados", agrofit)

    print(f"seeded BR sources wave 2 silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
