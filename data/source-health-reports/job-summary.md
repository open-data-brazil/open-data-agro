### Source health probe

- Executed at: 2026-11-19T11:11:56Z
- Run date: 2026-11-19
- Datasets probed: 131
- OK: 120 · Warning: 2 · Critical: 9
- Updated samples: 24
- Deprecated (2+ days): 9

### Source health alerts

- **ana.hidrologia-series** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL": context deadline exceeded (Client.Timeout exceeded while awaiting headers))
  - https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL
  - https://www.gov.br/ana/pt-br/acesso-a-informacao/dados-abertos
- **ana.pluviometria-redes** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL": context deadline exceeded (Client.Timeout exceeded while awaiting headers))
  - https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL
  - https://www.gov.br/ana/pt-br/acesso-a-informacao/dados-abertos
- **aneel.tarifas-energia** (critical, day 7): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (ckan package_show: Get "https://dadosabertos.aneel.gov.br/api/3/action/package_show?id=bandeiras-tarifarias": dial tcp 200.198.220.169:443: i/o timeout)
  - https://dadosabertos.aneel.gov.br/dataset/bandeiras-tarifarias
- **cepea.soja-paranagua** (critical, day 7): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/soja/soja-indicador-cepea-esalq-porto-paranagua)
  - https://www.noticiasagricolas.com.br/cotacoes/soja/soja-indicador-cepea-esalq-porto-paranagua
  - https://www.cepea.org.br/br/indicador/soja.aspx
- **cepea.soja-parana** (critical, day 7): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/soja/indicador-cepea-esalq-soja-parana)
  - https://www.noticiasagricolas.com.br/cotacoes/soja/indicador-cepea-esalq-soja-parana
  - https://www.cepea.org.br/br/indicador/soja.aspx
- **cepea.milho** (critical, day 7): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/milho/indicador-cepea-esalq-milho)
  - https://www.noticiasagricolas.com.br/cotacoes/milho/indicador-cepea-esalq-milho
  - https://www.cepea.org.br/br/indicador/milho.aspx
- **cepea.boi-gordo** (critical, day 7): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/boi-gordo/boi-gordo-indicador-esalq-bmf)
  - https://www.noticiasagricolas.com.br/cotacoes/boi-gordo/boi-gordo-indicador-esalq-bmf
  - https://www.cepea.org.br/br/indicador/boi-gordo.aspx
- **dnit.condicoes-conservacao-rodovias** (critical, day 15): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 15): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **jrc.mars-crop-yield** (critical, day 15): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv)
  - https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv
  - https://mars.jrc.ec.europa.eu/dataset
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 15): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
