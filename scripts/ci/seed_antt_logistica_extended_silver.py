#!/usr/bin/env python3
"""Seed minimal silver Delta for ANTT extended logistics CI (local-first)."""

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
    (lake_root / "gold" / "mart_antt__volume_trafego_pedagio").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_antt__receita_por_praca").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    volume = pa.table(
        {
            "concessionaria": ["AUTOPISTA FERNÃO DIAS", "NOVADUTRA"],
            "mes_ano": ["01/2026", "01/2026"],
            "sentido": ["Crescente", "CrescenteDecrescente"],
            "praca": ["3 (Cambuí)", "Itatiaia"],
            "tipo_cobranca": ["Automática", "Automática"],
            "categoria_eixo": ["2", "2"],
            "tipo_de_veiculo": ["Passeio", "Passeio"],
            "volume_total": ["12500,00", "45000,00"],
            "_dataset_id": ["antt.volume-trafego-pedagio", "antt.volume-trafego-pedagio"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "antt", "volume_trafego_pedagio", volume)

    receita = pa.table(
        {
            "Concessionaria": ["AUTOPISTA FERNÃO DIAS", "NOVADUTRA"],
            "Praca_de_pedagio": ["3 (Cambuí)", "Itatiaia"],
            "Ano_PNV_SNV": ["2009", "1998"],
            "UF": ["Minas Gerais", "Rio de Janeiro"],
            "Rodovia": ["BR-381/MG", "BR-116/RJ"],
            "Km_m": ["900,9", "318,900"],
            "Tipo_de_Pista": ["Principal", "Principal"],
            "Sentido": ["CrescenteDecrescente", "CrescenteDecrescente"],
            "Municipio": ["Cambuí", "Itatiaia"],
            "Direcao": ["SulNorte", "Sul/Norte"],
            "Latitude": ["-22,628487", "-22,494913"],
            "Longitude": ["-46,07789", "-44,569597"],
            "Data_da_Ativacao": ["23/03/2009", "01/08/1996"],
            "Mes_ano": ["01/2025", "01/2025"],
            "Receita_Praca_de_Pedagio": ["10956547,5", "15552503,29"],
            "_dataset_id": ["antt.receita-por-praca", "antt.receita-por-praca"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "antt", "receita_por_praca", receita)

    print(f"seeded ANTT extended logistics silver under {lake_root / 'silver' / 'antt'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
