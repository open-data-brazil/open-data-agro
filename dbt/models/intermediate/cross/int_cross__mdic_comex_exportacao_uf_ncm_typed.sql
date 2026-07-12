{{ config(materialized='view') }}

select
    trim(co_ncm) as co_ncm,
    trim(ncm_descricao) as ncm_descricao,
    lower(trim(produto_slug)) as produto_slug,
    upper(trim(uf)) as uf,
    {{ to_date('data') }} as refmonth,
    {{ ascii_to_double('valor_fob_usd') }} as valor_fob_usd,
    {{ ascii_to_double('quantidade_kg') }} as quantidade_kg,
    try_cast(ano as integer) as ano,
    (last_day({{ to_date('data') }}) + interval 45 day)::date as as_of,
    capturado_em,
    fonte_oficial,
    _dataset_id
from {{ read_gold_mart('mart_mdic__comex_exportacao_uf_ncm') }}
where {{ to_date('data') }} is not null
