{{ config(materialized='view') }}

select
    try_cast(sgs_codigo as integer) as sgs_codigo,
    {{ to_date('data') }} as data,
    {{ ascii_to_double('valor') }} as ptax_usd_venda,
    try_cast(ano as integer) as ano,
    {{ to_date('data') }} as as_of,
    capturado_em,
    fonte_oficial,
    _dataset_id
from {{ read_gold_mart('mart_bcb__sgs_ptax_usd_venda') }}
where {{ to_date('data') }} is not null
