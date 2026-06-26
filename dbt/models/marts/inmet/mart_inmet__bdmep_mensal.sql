{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_inmet__bdmep_mensal/mart.parquet',
    format='parquet'
) }}

select
    cd_estacao,
    mes,
    variavel,
    valor,
    uf,
    ano,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_inmet__bdmep_mensal') }}
