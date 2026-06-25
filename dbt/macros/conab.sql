{% macro conab_silver_path(table_name) %}
  {{ var('lake_root') }}/silver/conab/{{ table_name }}
{% endmacro %}

{% macro conab_fonte_oficial() %}
  {{ return(var('conab_fonte_oficial')) }}
{% endmacro %}

{# Map CONAB LevantamentoGraos headers to glossary snake_case. #}
{% macro conab_graos_core_columns() %}
  "Produto" as produto,
  "UF" as uf,
  "Safra" as safra,
  "Região" as regiao,
  "Produção (mil t)" as producao_mil_t
{% endmacro %}

{% macro conab_serie_core_columns() %}
  "Produto" as produto,
  "UF" as uf,
  "Ano" as ano,
  "Produção (mil t)" as producao_mil_t
{% endmacro %}
