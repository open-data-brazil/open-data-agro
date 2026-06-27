{% macro sagis_grain_supply_statistics_fonte_oficial() %}
  https://www.sagis.org.za/
{% endmacro %}

{% macro sagis_grain_supply_statistics_columns() %}
    commodity_slug, marketing_year, supply_t, demand_t, opening_stocks_t, closing_stocks_t
{% endmacro %}

{% macro sagis_grain_supply_statistics_dedupe_keys() %}
    commodity_slug, marketing_year
{% endmacro %}
