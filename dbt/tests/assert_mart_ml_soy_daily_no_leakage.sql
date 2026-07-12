-- Leakage guard (A7): fail if any feature as_of exceeds the label date.
select
    produto_slug,
    data,
    label_date,
    feature_as_of_max,
    as_of_cepea,
    as_of_ptax,
    as_of_b3,
    as_of_comex,
    as_of_frete,
    as_of_inmet
from {{ ref('mart_ml__soy_daily_features') }}
where feature_as_of_max > label_date
   or (as_of_cepea is not null and as_of_cepea > label_date)
   or (as_of_ptax is not null and as_of_ptax > label_date)
   or (as_of_b3 is not null and as_of_b3 > label_date)
   or (as_of_comex is not null and as_of_comex > label_date)
   or (as_of_frete is not null and as_of_frete > label_date)
   or (as_of_inmet is not null and as_of_inmet > label_date)
