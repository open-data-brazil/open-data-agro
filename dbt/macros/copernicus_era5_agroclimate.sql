{% macro copernicus_era5_agroclimate_fonte_oficial() %}
  https://cds.climate.copernicus.eu/
{% endmacro %}

{% macro copernicus_era5_agroclimate_columns() %}
    latitude, longitude, refdate, variable_slug, value, unit
{% endmacro %}

{% macro copernicus_era5_agroclimate_dedupe_keys() %}
    latitude, longitude, refdate, variable_slug
{% endmacro %}
