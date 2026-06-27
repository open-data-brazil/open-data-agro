with silver as (
    select * from {{ read_silver_parquet('grain_supply_statistics', agency='sagis') }}
),

renamed as (
    select
        {{ sagis_grain_supply_statistics_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ sagis_grain_supply_statistics_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by commodity_slug, marketing_year
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
