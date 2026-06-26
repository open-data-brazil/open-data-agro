{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__precos_mensal_municipio/mart.parquet',
    format='parquet'
) }}

select
    produto,
    classificacao_produto,
    id_produto,
    municipio,
    cod_ibge,
    uf,
    regiao,
    ano,
    mes,
    nivel_comercializacao,
    valor_produto_kg,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__precos_mensal_municipio') }}
