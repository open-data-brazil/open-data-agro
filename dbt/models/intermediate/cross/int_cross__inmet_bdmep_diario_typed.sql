{{ config(materialized='view') }}

select
    trim(cd_estacao) as cd_estacao,
    {{ to_date('data') }} as data,
    trim(variavel) as variavel,
    {{ ascii_to_double('valor') }} as valor,
    upper(trim(uf)) as uf,
    try_cast(ano as integer) as ano,
    ({{ to_date('data') }} + interval 2 day)::date as as_of,
    capturado_em,
    fonte_oficial,
    _dataset_id
from {{ read_gold_mart('mart_inmet__bdmep_diario') }}
where {{ to_date('data') }} is not null
