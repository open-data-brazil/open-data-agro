{% macro antt_fonte_oficial() %}
  https://dados.antt.gov.br/dataset/praca-de-pedagio
{% endmacro %}

{% macro antt_pracas_pedagio_columns() %}
    concessionaria,
    praca_de_pedagio,
    ano_do_pnv_snv,
    rodovia,
    uf,
    km_m,
    municipio,
    tipo_de_pista,
    sentido,
    situacao,
    data_da_inativacao,
    latitude,
    longitude
{% endmacro %}
