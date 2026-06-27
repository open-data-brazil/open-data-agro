with silver as (
    select * from {{ read_silver_parquet('futuro_acucar', agency='b3') }}
),

renamed as (
    select
        refdate,
        symbol,
        commodity,
        maturity_code,
        previous_price,
        price,
        currency,
        price_change,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ b3_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by refdate, symbol
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
