with silver as (
    select * from {{ read_silver_parquet('ppm_aquicultura', agency='ibge') }}
),

renamed as (
    select
        {{ ibge_ppm_uf_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ ibge_ppm_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by codigo_uf, ano, variavel_codigo, categoria_codigo
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
