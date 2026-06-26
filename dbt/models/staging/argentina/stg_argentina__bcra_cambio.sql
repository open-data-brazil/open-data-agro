with silver as (
    select * from {{ read_silver_parquet('bcra_cambio', agency='argentina') }}
),

renamed as (
    select
        currency_code,
        currency_name,
        refdate,
        exchange_rate,
        rate_type,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ argentina_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by currency_code, refdate, rate_type
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
