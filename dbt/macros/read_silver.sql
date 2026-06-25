{% macro read_silver_parquet(table_name, agency='conab') %}
  read_parquet('{{ var('lake_root') }}/silver/{{ agency }}/{{ table_name }}/**/*.parquet')
{% endmacro %}
