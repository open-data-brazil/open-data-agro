{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__prohort_mensal/mart.parquet',
    format='parquet'
) }}

select
    ano,
    mes,
    municipio_origem,
    cod_ibge_municipio_origem,
    uf_origem,
    ceasa,
    uf_ceasa,
    municipio_ceasa,
    cod_ibge_municipio_ceasa,
    produto,
    qtd_comercializada_kg,
    valor_comercializado,
    pais_origem,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__prohort_mensal') }}
