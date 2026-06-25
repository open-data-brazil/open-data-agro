with silver as (
    select * from {{ read_silver_parquet('estoques_publicos') }}
),

renamed as (
    select
        {{ conab_estoques_publicos_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ conab_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by produto, municipio, cod_ibge, ano, mes, conta_operacional
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
