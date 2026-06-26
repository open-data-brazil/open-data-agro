{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__prohort_diario/mart.parquet',
    format='parquet'
) }}

select
    municipio_ceasa,
    cod_ibge_municipio,
    uf_ceasa,
    ceasa,
    produto,
    unidade_medida,
    data_preco,
    preco_diario,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__prohort_diario') }}
