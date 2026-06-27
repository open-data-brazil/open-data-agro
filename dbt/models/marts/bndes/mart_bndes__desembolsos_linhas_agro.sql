{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_bndes__desembolsos_linhas_agro/mart.parquet',
    format='parquet'
) }}

select
    ano,
    mes,
    bndes_finem,
    bndes_exim,
    bndes_mercado_de_capitais,
    bndes_nao_reembolsavel,
    bndes_microcredito,
    bndes_prestacao_de_garantia,
    bndes_finame,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_bndes__desembolsos_linhas_agro') }}
