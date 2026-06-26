with silver as (
    select * from {{ read_silver_parquet('precos_minimos') }}
),

renamed as (
    select
        {{ conab_precos_minimos_columns() }},
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
            partition by
                produto,
                id_produto,
                uf,
                regionalizacao,
                ano_inicio_vigencia,
                mes_inicio_vigencia,
                ano_termino_vigencia,
                mes_termino_vigencia
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
