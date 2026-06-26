{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__capacidade_estatica/mart.parquet',
    format='parquet'
) }}

select
    ano,
    uf,
    quantidade_mil_t,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__capacidade_estatica') }}
