with silver as (
    select * from {{ read_silver_parquet('its_trade_statistics', agency='wto') }}
),

renamed as (
    select
        {{ wto_its_trade_statistics_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ wto_its_trade_statistics_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by reporter_code, indicator_code, period, flow_code
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
