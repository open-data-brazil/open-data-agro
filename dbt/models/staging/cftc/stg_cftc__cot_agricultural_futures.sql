with silver as (
    select * from {{ read_silver_parquet('cot_agricultural_futures', agency='cftc') }}
),

renamed as (
    select
        {{ cftc_cot_agricultural_futures_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ cftc_cot_agricultural_futures_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by commodity_slug, report_date
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
