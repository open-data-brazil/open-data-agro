with silver as (
    select * from {{ read_silver_parquet('pam_estabelecimentos', agency='ibge') }}
),

renamed as (
    select
        {{ ibge_pam_core_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ ibge_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by codigo_ibge, ano, variavel_codigo, produto_codigo
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
