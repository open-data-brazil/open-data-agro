{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__pevs_producao_vegetal/mart.parquet',
    format='parquet'
) }}

select
    sidra_tabela,
    codigo_uf,
    uf,
    ano,
    variavel_codigo,
    variavel,
    produto_codigo,
    produto,
    valor,
    unidade_codigo,
    unidade,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_ibge__pevs_producao_vegetal') }}
