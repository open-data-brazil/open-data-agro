{% macro read_gold_mart(mart_name) %}
  read_parquet('{{ var('gold_root') }}/{{ mart_name }}/mart.parquet')
{% endmacro %}

{% macro br_to_double(expr) %}
  try_cast(replace(replace(cast({{ expr }} as varchar), '.', ''), ',', '.') as double)
{% endmacro %}

{% macro ascii_to_double(expr) %}
  try_cast(replace(cast({{ expr }} as varchar), ',', '.') as double)
{% endmacro %}

{% macro to_date(expr) %}
  try_cast({{ expr }} as date)
{% endmacro %}
