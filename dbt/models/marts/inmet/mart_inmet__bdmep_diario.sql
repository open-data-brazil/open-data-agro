{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_inmet__bdmep_diario/mart.parquet',
    format='parquet'
) }}

select
    cd_estacao,
    data,
    variavel,
    valor,
    uf,
    ano,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_inmet__bdmep_diario') }}
