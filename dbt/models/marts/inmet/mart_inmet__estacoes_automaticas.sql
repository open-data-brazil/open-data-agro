{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_inmet__estacoes_automaticas/mart.parquet',
    format='parquet'
) }}

select
    cd_estacao,
    nome,
    latitude,
    longitude,
    uf,
    situacao,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_inmet__estacoes_automaticas') }}
