with silver as (
    select * from {{ read_silver_parquet('mars_crop_yield', agency='jrc') }}
),

renamed as (
    select
        {{ jrc_mars_crop_yield_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ jrc_mars_crop_yield_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by country, crop_slug, harvest_year, forecast_timing
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
