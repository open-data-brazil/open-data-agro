{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__lspa_area_producao/mart.parquet',
    format='parquet'
) }}

select
    sidra_tabela,
    codigo_uf,
    uf,
    mes,
    variavel_codigo,
    variavel,
    produto_codigo,
    produto,
    produto_slug,
    valor,
    unidade_codigo,
    unidade,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_ibge__lspa_area_producao') }}
