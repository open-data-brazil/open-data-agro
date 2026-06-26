{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__localidades_regioes/mart.parquet',
    format='parquet'
) }}

select
    codigo_regiao,
    sigla_regiao,
    nome,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_ibge__localidades_regioes') }}
