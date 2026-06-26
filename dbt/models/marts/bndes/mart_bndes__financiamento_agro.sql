{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_bndes__financiamento_agro/mart.parquet',
    format='parquet'
) }}

select
    ano,
    mes,
    agropecuaria,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_bndes__financiamento_agro') }}
