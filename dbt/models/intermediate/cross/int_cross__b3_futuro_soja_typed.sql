{{ config(materialized='view') }}

select
    {{ to_date('refdate') }} as refdate,
    trim(symbol) as symbol,
    upper(trim(commodity)) as commodity,
    trim(maturity_code) as maturity_code,
    {{ ascii_to_double('previous_price') }} as previous_price,
    {{ ascii_to_double('price') }} as price,
    trim(currency) as currency,
    {{ ascii_to_double('price_change') }} as price_change,
    {{ to_date('refdate') }} as as_of,
    capturado_em,
    fonte_oficial,
    _dataset_id
from {{ read_gold_mart('mart_b3__futuro_soja') }}
where {{ to_date('refdate') }} is not null
