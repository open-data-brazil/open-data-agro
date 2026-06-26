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
    gold_tables = [
        "mart_mapa__agrofit_produtos_formulados",
        "mart_mapa__agrofit_produtos_tecnicos",
        "mart_ana__hidrologia_series",
        "mart_antaq__movimentacao_carga_portuaria",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    agrofit_formulados = pa.table(
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
    write_table(lake_root, "mapa", "agrofit_produtos_formulados", agrofit_formulados)

    agrofit_tecnicos = pa.table(
        {
            "NUMERO_REGISTRO": ["4215", "10116"],
            "PRODUTO_TECNICO_MARCA_COMERCIAL": ["2,4 D Técnico Mol", "2,4-D Acid Tecnico"],
            "INGREDIENTE_ATIVO(GRUPO_QUIMICI)(CONCENTRACAO)": [
                "2,4-D (ácido ariloxialcanóico) (973 g/kg)",
                "2,4-D (ácido ariloxialcanóico) (980 g/kg)",
            ],
            "CLASSE": ["Herbicida", "Herbicida"],
            "TITULAR_REGISTRO": ["Meghmani Organics", "Sharda do Brasil"],
            "EMPRESA_<PAIS>_TIPO": ["ÍNDIA", "ÍNDIA"],
            "CLASSIFICACAO_TOXICOLOGICA": ["I", "II"],
            "CLASSIFICACAO_AMBIENTAL": ["III", "III"],
            "_dataset_id": ["mapa.agrofit-produtos-tecnicos", "mapa.agrofit-produtos-tecnicos"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "mapa", "agrofit_produtos_tecnicos", agrofit_tecnicos)

    ana_hidro = pa.table(
        {
            "station_code": ["15400000", "15400000"],
            "consistency_level": ["2", "2"],
            "data_type": ["3", "3"],
            "observed_at": ["2024-01-01 00:00:00", "2024-01-02 00:00:00"],
            "daily_mean": ["13225.42", "13100.00"],
            "acquisition_method": ["1", "1"],
            "max_value": ["14000.00", "13800.00"],
            "min_value": ["12500.00", "12400.00"],
            "mean_value": ["13225.42", "13100.00"],
            "_dataset_id": ["ana.hidrologia-series", "ana.hidrologia-series"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ana", "hidrologia_series", ana_hidro)

    antaq_carga = pa.table(
        {
            "Ano": ["2024", "2024"],
            "Mes": ["01", "01"],
            "CodigoInstalacaoPortuaria": ["BRSP001", "BRSP001"],
            "NomeInstalacaoPortuaria": ["Santos", "Santos"],
            "TipoMovimentacao": ["Carga", "Carga"],
            "TipoNavegacao": ["Longo Curso", "Longo Curso"],
            "Sentido": ["Entrada", "Saída"],
            "NaturezaCarga": ["Granéis Sólidos", "Contêineres"],
            "TipoOperacao": ["Descarga", "Embarque"],
            "PesoToneladas": ["12500000,50", "8200000,00"],
            "_dataset_id": ["antaq.movimentacao-carga-portuaria", "antaq.movimentacao-carga-portuaria"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "antaq", "movimentacao_carga_portuaria", antaq_carga)

    print(f"seeded BR sources wave 2 silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
