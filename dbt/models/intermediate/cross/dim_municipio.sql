{{ config(
    materialized='external',
    location=var('gold_root') ~ '/dim_municipio/mart.parquet',
    format='parquet'
) }}

select
    lpad(trim(codigo_ibge), 7, '0') as codigo_ibge,
    trim(nome) as nome_municipio,
    upper(trim(sigla_uf)) as uf,
    lpad(trim(codigo_uf), 2, '0') as codigo_uf,
    trim(codigo_regiao) as codigo_regiao,
    trim(nome_regiao) as nome_regiao,
    capturado_em,
    fonte_oficial,
    _dataset_id
from {{ read_gold_mart('mart_ibge__localidades_municipios') }}
where trim(codigo_ibge) is not null
  and trim(codigo_ibge) <> ''
