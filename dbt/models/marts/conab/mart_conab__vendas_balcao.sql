{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__vendas_balcao/mart.parquet',
    format='parquet'
) }}

select
    ano,
    mes,
    municipio_armazem_venda,
    uf,
    qtd_produto_kg,
    valor_comercializado,
    numero_atendimentos,
    clientes_atendidos,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__vendas_balcao') }}
