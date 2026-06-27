with silver as (
    select * from {{ read_silver_parquet('maff_ag_trade', agency='japan') }}
),

renamed as (
    select
        {{ japan_maff_ag_trade_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ japan_maff_ag_trade_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by commodity_slug, refyear, flow_code
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
