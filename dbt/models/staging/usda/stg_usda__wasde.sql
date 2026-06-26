with silver as (
    select * from {{ read_silver_parquet('wasde', agency='usda') }}
),

renamed as (
    select
        {{ usda_wasde_columns() }},
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
            partition by report_month, commodity, market_year, attribute
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
