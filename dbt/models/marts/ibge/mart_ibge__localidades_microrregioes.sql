{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__localidades_microrregioes/mart.parquet',
    format='parquet'
) }}

select
    codigo_microrregiao,
    nome,
    codigo_mesorregiao,
    nome_mesorregiao,
    codigo_uf,
    sigla_uf,
    nome_uf,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_ibge__localidades_microrregioes') }}
