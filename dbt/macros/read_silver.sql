{% macro read_silver_parquet(table_name) %}
  read_parquet('{{ var('lake_root') }}/silver/conab/{{ table_name }}/**/*.parquet')
{% endmacro %}
