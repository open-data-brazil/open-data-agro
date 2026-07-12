{{ config(materialized='view') }}

select
    lower(trim(produto)) as produto_slug,
    trim(praca) as praca,
    {{ to_date('data') }} as data,
    {{ ascii_to_double('preco_rs_sc') }} as preco_rs_sc,
    {{ ascii_to_double('variacao_dia_pct') }} as variacao_dia_pct,
    {{ ascii_to_double('preco_usd_sc') }} as preco_usd_sc,
    try_cast(ano as integer) as ano,
    {{ to_date('data') }} as as_of,
    capturado_em,
    fonte_oficial,
    _dataset_id
from {{ read_gold_mart('mart_cepea__soja_paranagua') }}
where {{ to_date('data') }} is not null
