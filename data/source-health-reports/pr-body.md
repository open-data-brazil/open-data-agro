## Source health daily probe

### Source health probe

- Executed at: 2026-07-22T11:17:30Z
- Run date: 2026-07-22
- Datasets probed: 131
- OK: 120 · Warning: 3 · Critical: 8
- Updated samples: 25
- Deprecated (2+ days): 8

### Source health alerts

- **dnit.condicoes-conservacao-rodovias** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **eurostat.ag-prices** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://ec.europa.eu/eurostat/api/dissemination/statistics/1.0/data/apri_pi15_outa?format=JSON&geo=EU27_2020&lang=en&product=010000&product=011000&product=015000&sinceTimePeriod=2020)
  - https://ec.europa.eu/eurostat/api/dissemination/statistics/1.0/data/apri_pi15_outa?format=JSON&geo=EU27_2020&lang=en&product=010000&product=011000&product=015000&sinceTimePeriod=2020
  - https://ec.europa.eu/eurostat/web/agriculture/database
- **mdic.comex-exportacao-ncm-mes** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 429 for https://api-comexstat.mdic.gov.br/general)
  - https://api-comexstat.mdic.gov.br/general
  - https://comexstat.mdic.gov.br/
- **mdic.comex-importacao-ncm-mes** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 429 for https://api-comexstat.mdic.gov.br/general)
  - https://api-comexstat.mdic.gov.br/general
  - https://comexstat.mdic.gov.br/
- **mdic.comex-exportacao-uf-ncm** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 429 for https://api-comexstat.mdic.gov.br/general)
  - https://api-comexstat.mdic.gov.br/general
  - https://comexstat.mdic.gov.br/
- **mdic.comex-importacao-diesel-ncm** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 429 for https://api-comexstat.mdic.gov.br/general)
  - https://api-comexstat.mdic.gov.br/general
  - https://comexstat.mdic.gov.br/
- **noaa.global-temp-anomaly** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://www.ncei.noaa.gov/access/monitoring/climate-at-a-glance/global/time-series/globe/land_ocean/0/0/2010-2010.csv": dial tcp 205.167.25.178:443: connect: connection refused)
  - https://www.ncei.noaa.gov/access/monitoring/climate-at-a-glance/global/time-series/globe/land_ocean/0/0/2010-2010.csv
  - https://www.ncei.noaa.gov/access/monitoring/climate-at-a-glance/global/time-series
- **suframa.comercio-mercadorias-zfm** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx": dial tcp 161.148.164.31:443: i/o timeout)
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
