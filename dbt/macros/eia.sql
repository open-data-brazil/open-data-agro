{% macro eia_fonte_oficial() %}
  https://www.eia.gov/opendata/
{% endmacro %}

{% macro eia_petroleum_columns() %}
    series_id,
    series_name,
    commodity_slug,
    refdate,
    unit,
    value,
    frequency
{% endmacro %}
