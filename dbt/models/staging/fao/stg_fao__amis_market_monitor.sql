with silver as (
    select * from {{ read_silver_parquet('amis_market_monitor', agency='fao') }}
),

renamed as (
    select
        {{ fao_amis_market_monitor_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ fao_amis_market_monitor_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by commodity_slug, refmonth, indicator_slug
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
