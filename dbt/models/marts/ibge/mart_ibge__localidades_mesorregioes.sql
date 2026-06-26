{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__localidades_mesorregioes/mart.parquet',
    format='parquet'
) }}

select
    codigo_mesorregiao,
    nome,
    codigo_uf,
    sigla_uf,
    nome_uf,
    codigo_regiao,
    sigla_regiao,
    nome_regiao,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_ibge__localidades_mesorregioes') }}
