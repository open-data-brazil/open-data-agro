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
