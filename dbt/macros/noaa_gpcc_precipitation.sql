{% macro noaa_gpcc_precipitation_fonte_oficial() %}
  https://www.ncei.noaa.gov/data/gpcc-monthly/
{% endmacro %}

{% macro noaa_gpcc_precipitation_columns() %}
    refmonth, latitude, longitude, precip_mm, grid_resolution
{% endmacro %}

{% macro noaa_gpcc_precipitation_dedupe_keys() %}
    refmonth, latitude, longitude
{% endmacro %}
