{% macro fred_commodity_indexes_fonte_oficial() %}
  https://fred.stlouisfed.org/
{% endmacro %}

{% macro fred_commodity_indexes_columns() %}
    series_id, refmonth, value
{% endmacro %}

{% macro fred_commodity_indexes_dedupe_keys() %}
    series_id, refmonth
{% endmacro %}
