{% macro cepea_fonte_oficial() %}
  https://www.cepea.org.br/
{% endmacro %}

{% macro cepea_indicador_columns() %}
  trim(produto) as produto,
  trim(praca) as praca,
  trim(data) as data,
  cast(preco_rs_sc as varchar) as preco_rs_sc,
  cast(variacao_dia_pct as varchar) as variacao_dia_pct,
  cast(preco_usd_sc as varchar) as preco_usd_sc,
  cast(ano as varchar) as ano
{% endmacro %}
