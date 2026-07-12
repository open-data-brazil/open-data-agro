{{ config(materialized='view') }}

-- UF-pair frete monthly aggregates (corridor proxies for ML features).
select
    uf_origem,
    uf_destino,
    refmonth,
    ano,
    mes,
    count(*) as n_contratos,
    avg(distancia_km) as distancia_km_avg,
    avg(valor_frete_tonelada) as valor_frete_ton_avg,
    avg(valor_tonelada_km) as valor_tonelada_km_avg,
    max(as_of) as as_of
from {{ ref('int_cross__conab_frete_typed') }}
group by 1, 2, 3, 4, 5
