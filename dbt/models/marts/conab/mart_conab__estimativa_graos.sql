{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__estimativa_graos/mart.parquet',
    format='parquet'
) }}

select
    produto,
    uf,
    safra,
    regiao,
    producao_mil_t,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__estimativa_graos') }}
