## Source health daily probe

### Source health probe

- Executed at: 2026-11-18T11:13:05Z
- Run date: 2026-11-18
- Datasets probed: 131
- OK: 119 · Warning: 3 · Critical: 9
- Updated samples: 20
- Deprecated (2+ days): 9

### Source health alerts

- **aneel.tarifas-energia** (critical, day 6): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (ckan package_show: Get "https://dadosabertos.aneel.gov.br/api/3/action/package_show?id=bandeiras-tarifarias": dial tcp 200.198.220.169:443: i/o timeout)
  - https://dadosabertos.aneel.gov.br/dataset/bandeiras-tarifarias
- **antt.volume-trafego-pedagio** (warning, day 1): Possible link deprecation — official source unreachable after retries. (parse ckan response: invalid character '<' looking for beginning of value)
  - https://dados.antt.gov.br/dataset/volume-trafego-praca-pedagio
- **antt.receita-por-praca** (warning, day 1): Possible link deprecation — official source unreachable after retries. (parse ckan response: invalid character '<' looking for beginning of value)
  - https://dados.antt.gov.br/dataset/receita-por-praca
- **antt.pracas-pedagio** (warning, day 1): Possible link deprecation — official source unreachable after retries. (parse ckan response: invalid character '<' looking for beginning of value)
  - https://dados.antt.gov.br/dataset/praca-de-pedagio
- **cepea.soja-paranagua** (critical, day 6): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/soja/soja-indicador-cepea-esalq-porto-paranagua)
  - https://www.noticiasagricolas.com.br/cotacoes/soja/soja-indicador-cepea-esalq-porto-paranagua
  - https://www.cepea.org.br/br/indicador/soja.aspx
- **cepea.soja-parana** (critical, day 6): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/soja/indicador-cepea-esalq-soja-parana)
  - https://www.noticiasagricolas.com.br/cotacoes/soja/indicador-cepea-esalq-soja-parana
  - https://www.cepea.org.br/br/indicador/soja.aspx
- **cepea.milho** (critical, day 6): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/milho/indicador-cepea-esalq-milho)
  - https://www.noticiasagricolas.com.br/cotacoes/milho/indicador-cepea-esalq-milho
  - https://www.cepea.org.br/br/indicador/milho.aspx
- **cepea.boi-gordo** (critical, day 6): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/boi-gordo/boi-gordo-indicador-esalq-bmf)
  - https://www.noticiasagricolas.com.br/cotacoes/boi-gordo/boi-gordo-indicador-esalq-bmf
  - https://www.cepea.org.br/br/indicador/boi-gordo.aspx
- **dnit.condicoes-conservacao-rodovias** (critical, day 14): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 14): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **jrc.mars-crop-yield** (critical, day 14): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv)
  - https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv
  - https://mars.jrc.ec.europa.eu/dataset
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 14): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
