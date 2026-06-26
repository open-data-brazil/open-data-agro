{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__localidades_municipios/mart.parquet',
    format='parquet'
) }}

select
    codigo_ibge,
    nome,
    sigla_uf,
    codigo_uf,
    codigo_regiao,
    nome_regiao,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_ibge__localidades_municipios') }}
