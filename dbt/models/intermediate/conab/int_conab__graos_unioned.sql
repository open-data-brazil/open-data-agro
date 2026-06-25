select
    'estimativa' as origem,
    produto,
    uf,
    safra as periodo,
    producao_mil_t,
    capturado_em,
    fonte_oficial
from {{ ref('stg_conab__estimativa_graos') }}

union all

select
    'serie_historica' as origem,
    produto,
    uf,
    cast(ano as varchar) as periodo,
    producao_mil_t,
    capturado_em,
    fonte_oficial
from {{ ref('stg_conab__serie_historica_graos') }}
