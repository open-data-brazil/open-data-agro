{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__localidades_ufs/mart.parquet',
    format='parquet'
) }}

select
    codigo_uf,
    sigla_uf,
    nome,
    codigo_regiao,
    sigla_regiao,
    nome_regiao,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_ibge__localidades_ufs') }}
