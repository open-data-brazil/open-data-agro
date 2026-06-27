{% macro jrc_mars_crop_yield_fonte_oficial() %}
  https://mars.jrc.ec.europa.eu/dataset
{% endmacro %}

{% macro jrc_mars_crop_yield_columns() %}
    country, crop, crop_slug, forecast_yield_kg_ha, five_yr_avg_kg_ha, harvest_year, forecast_timing, region_name
{% endmacro %}

{% macro jrc_mars_crop_yield_dedupe_keys() %}
    country, crop_slug, harvest_year, forecast_timing
{% endmacro %}
