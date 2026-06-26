with silver as (
    select * from {{ read_silver_parquet('gats_trade', agency='usda') }}
),

renamed as (
    select
        commodity_code,
        commodity_name,
        partner_code,
        partner_name,
        flow,
        year,
        value,
        unit,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ usda_gats_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by commodity_code, partner_code, flow, year
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
