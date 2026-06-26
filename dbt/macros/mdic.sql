{% macro mdic_fonte_oficial() %}
  https://comexstat.mdic.gov.br/
{% endmacro %}

{% macro mdic_comex_columns() %}
    co_ncm,
    ncm_descricao,
    produto_slug,
    data,
    valor_fob_usd,
    quantidade_kg,
    ano
{% endmacro %}

{% macro mdic_comex_import_columns() %}
    co_ncm,
    ncm_descricao,
    produto_slug,
    data,
    valor_cif_usd,
    quantidade_kg,
    valor_frete_usd,
    valor_seguro_usd,
    ano
{% endmacro %}

{% macro mdic_comex_export_uf_columns() %}
    co_ncm,
    ncm_descricao,
    produto_slug,
    uf,
    data,
    valor_fob_usd,
    quantidade_kg,
    ano
{% endmacro %}
