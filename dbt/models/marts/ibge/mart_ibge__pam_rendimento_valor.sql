{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__pam_rendimento_valor/mart.parquet',
    format='parquet'
) }}

select
    sidra_tabela,
    codigo_ibge,
    municipio,
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
from {{ ref('stg_ibge__pam_rendimento_valor') }}
