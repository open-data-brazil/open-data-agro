with silver as (
    select * from {{ read_silver_parquet('pracas_pedagio', agency='antt') }}
),

renamed as (
    select
        concessionaria,
        praca_de_pedagio,
        ano_do_pnv_snv,
        rodovia,
        uf,
        km_m,
        municipal as municipio,
        tipo_de_pista,
        sentido,
        situacao,
        data_da_inativacao,
        latitude,
        longitude,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ antt_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by concessionaria, praca_de_pedagio, rodovia, uf, km_m
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
