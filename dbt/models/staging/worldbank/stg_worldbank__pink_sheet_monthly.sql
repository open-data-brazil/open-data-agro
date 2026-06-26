with silver as (
    select * from {{ read_silver_parquet('pink_sheet_monthly', agency='worldbank') }}
),

renamed as (
    select
        refmonth,
        series_name,
        commodity_slug,
        unit,
        value,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ worldbank_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by refmonth, series_name
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
