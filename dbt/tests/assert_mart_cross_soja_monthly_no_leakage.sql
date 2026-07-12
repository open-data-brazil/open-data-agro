-- Leakage guard: feature as_of must not exceed municipal label_date (CONAB publish as_of).
select
    produto_slug,
    cod_ibge,
    refmonth,
    label_date,
    feature_as_of_max,
    as_of_cepea,
    as_of_ptax,
    as_of_comex,
    as_of_frete
from {{ ref('mart_cross__soja_features_monthly') }}
where feature_as_of_max > label_date
   or (as_of_cepea is not null and as_of_cepea > label_date)
   or (as_of_ptax is not null and as_of_ptax > label_date)
   or (as_of_comex is not null and as_of_comex > label_date)
   or (as_of_frete is not null and as_of_frete > label_date)
