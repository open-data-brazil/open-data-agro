{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_inmet__estacoes_convencionais/mart.parquet',
    format='parquet'
) }}

select
    cd_estacao,
    nome,
    latitude,
    longitude,
    uf,
    situacao,
    regiao,
    altitude,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_inmet__estacoes_convencionais') }}
