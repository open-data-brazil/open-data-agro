with silver as (
    select * from {{ read_silver_parquet('psd_soja', agency='usda') }}
),

renamed as (
    select
        commodity_code,
        commodity_name,
        commodity_slug,
        country_code,
        country_name,
        marketing_year,
        calendar_year,
        month,
        attribute_id,
        attribute_name,
        unit_id,
        unit_description,
        value,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ usda_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by country_code, marketing_year, attribute_id, calendar_year, month
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
