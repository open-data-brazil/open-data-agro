#!/usr/bin/env python3
"""Seed minimal silver Delta for BR sources wave 5 env/logistics CI (Phase 54)."""

from __future__ import annotations

import os
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
    gold_tables = [
        "mart_ibama__sisfogo_incendios",
        "mart_ibama__licencas_ambientais",
        "mart_ibama__autos_infracao",
        "mart_ana__pluviometria_redes",
        "mart_embrapa__agroapi_agrofit",
        "mart_transportes__mtr_bit_malha_shapefile",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-27T12:00:00Z"

    sisfogo = pa.table(
        {
            "SEQ OCORRENCIA INCENDIO": ["21889", "17654"],
            "UF": ["BA", "BA"],
            "MUNICIPIO": ["IRAMAIA", "PORTO SEGURO"],
            "TIPO LOCALIDADE": ["PROPRIEDADE RURAL", "UNIDADE DE CONSERVAÇÃO"],
            "_dataset_id": ["ibama.sisfogo-incendios", "ibama.sisfogo-incendios"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibama", "sisfogo_incendios", sisfogo)

    licencas = pa.table(
        {
            "NUM_LICENCA": ["27188256/2026", "26952267/2026"],
            "DES_TIPOLICENCA": ["Anuência", "Anuência"],
            "NOM_EMPREENDIMENTO": ["Gasoduto Cacimbas - Catu", "Duplicação BR101"],
            "NOM_PESSOA": ["TAG", "ECO-101"],
            "_dataset_id": ["ibama.licencas-ambientais", "ibama.licencas-ambientais"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibama", "licencas_ambientais", licencas)

    autos = pa.table(
        {
            "SEQ_AUTO_INFRACAO": ["645774", "759674"],
            "NUM_AUTO_INFRACAO": ["961297", "145566"],
            "UF": ["SP", "MG"],
            "MUNICIPIO": ["MOGI DAS CRUZES", "TRES MARIAS"],
            "_dataset_id": ["ibama.autos-infracao", "ibama.autos-infracao"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibama", "autos_infracao", autos)

    pluvio = pa.table(
        {
            "station_code": ["87017001", "87017001"],
            "consistency_level": ["1", "1"],
            "data_type": ["2", "2"],
            "observed_at": ["2024-06-01 00:00:00", "2024-06-02 00:00:00"],
            "daily_mean": ["12.4", "0.0"],
            "acquisition_method": ["", ""],
            "max_value": ["", ""],
            "min_value": ["", ""],
            "mean_value": ["", ""],
            "_dataset_id": ["ana.pluviometria-redes", "ana.pluviometria-redes"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ana", "pluviometria_redes", pluvio)

    agrofit = pa.table(
        {
            "numero_registro": ["35523", "12345"],
            "marca_comercial": ["KBR-829M1-02", "EXAMPLE AGROFIT"],
            "situacao": ["TRUE", "TRUE"],
            "classe": ["Agente Biológico de Controle", "Herbicida"],
            "formulacao": ["Nematóides vivos", "Concentrado emulsionável"],
            "ingrediente_ativo": ["Heterorhabditis bacteriophora", "Glyphosate"],
            "_dataset_id": ["embrapa.agroapi-agrofit", "embrapa.agroapi-agrofit"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "embrapa", "agroapi_agrofit", agrofit)

    shapefile = pa.table(
        {
            "objectid_1": ["595", "697"],
            "tip_situac": ["Desativada", "Desativada"],
            "sigla": ["Sem Info", "Sem Info"],
            "uf": ["SP", "SP"],
            "municipio": ["Leme", "Descalvado"],
            "_dataset_id": [
                "transportes.mtr-bit-malha-shapefile",
                "transportes.mtr-bit-malha-shapefile",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "transportes", "mtr_bit_malha_shapefile", shapefile)

    print(f"seeded wave 5 env/logistics silver under {lake_root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
