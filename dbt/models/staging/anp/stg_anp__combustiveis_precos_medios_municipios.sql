with silver as (
    select * from {{ read_silver_parquet('combustiveis_precos_medios_municipios', 'anp') }}
),

renamed as (
    select
        {{ anp_combustiveis_precos_medios_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ anp_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by data_inicial, data_final, estado, municipio, produto
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
