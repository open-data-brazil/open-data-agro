with silver as (
    select * from {{ read_silver_parquet('agrofit_produtos_tecnicos', agency='mapa') }}
),

renamed as (
    select
        NUMERO_REGISTRO as numero_registro,
        PRODUTO_TECNICO_MARCA_COMERCIAL as produto_tecnico_marca_comercial,
        "INGREDIENTE_ATIVO(GRUPO_QUIMICI)(CONCENTRACAO)" as ingrediente_ativo,
        CLASSE as classe,
        TITULAR_REGISTRO as titular_registro,
        "EMPRESA_<PAIS>_TIPO" as empresa_pais_tipo,
        CLASSIFICACAO_TOXICOLOGICA as classificacao_toxicologica,
        CLASSIFICACAO_AMBIENTAL as classificacao_ambiental,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ mapa_agrofit_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by numero_registro, produto_tecnico_marca_comercial
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
