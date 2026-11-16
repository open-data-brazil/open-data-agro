## Source health daily probe

### Source health probe

- Executed at: 2026-11-16T10:46:47Z
- Run date: 2026-11-16
- Datasets probed: 131
- OK: 112 · Warning: 3 · Critical: 16
- Updated samples: 22
- Deprecated (2+ days): 16

### Source health alerts

- **ana.hidrologia-series** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL": context deadline exceeded (Client.Timeout exceeded while awaiting headers))
  - https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL
  - https://www.gov.br/ana/pt-br/acesso-a-informacao/dados-abertos
- **ana.pluviometria-redes** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL": context deadline exceeded (Client.Timeout exceeded while awaiting headers))
  - https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL
  - https://www.gov.br/ana/pt-br/acesso-a-informacao/dados-abertos
- **aneel.tarifas-energia** (critical, day 4): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (ckan package_show: Get "https://dadosabertos.aneel.gov.br/api/3/action/package_show?id=bandeiras-tarifarias": dial tcp 200.198.220.169:443: i/o timeout)
  - https://dadosabertos.aneel.gov.br/dataset/bandeiras-tarifarias
- **bndes.desembolsos-linhas-agro** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (ckan package_show: Get "https://dadosabertos.bndes.gov.br/api/3/action/package_show?id=desembolsos": dial tcp 200.225.103.94:443: connect: connection refused)
  - https://dadosabertos.bndes.gov.br/dataset/desembolsos
- **bndes.financiamento-agro** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (ckan package_show: Get "https://dadosabertos.bndes.gov.br/api/3/action/package_show?id=desembolsos": dial tcp 200.225.103.94:443: connect: connection refused)
  - https://dadosabertos.bndes.gov.br/dataset/desembolsos
- **cepea.soja-paranagua** (critical, day 4): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/soja/soja-indicador-cepea-esalq-porto-paranagua)
  - https://www.noticiasagricolas.com.br/cotacoes/soja/soja-indicador-cepea-esalq-porto-paranagua
  - https://www.cepea.org.br/br/indicador/soja.aspx
- **cepea.soja-parana** (critical, day 4): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/soja/indicador-cepea-esalq-soja-parana)
  - https://www.noticiasagricolas.com.br/cotacoes/soja/indicador-cepea-esalq-soja-parana
  - https://www.cepea.org.br/br/indicador/soja.aspx
- **cepea.milho** (critical, day 4): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/milho/indicador-cepea-esalq-milho)
  - https://www.noticiasagricolas.com.br/cotacoes/milho/indicador-cepea-esalq-milho
  - https://www.cepea.org.br/br/indicador/milho.aspx
- **cepea.boi-gordo** (critical, day 4): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/boi-gordo/boi-gordo-indicador-esalq-bmf)
  - https://www.noticiasagricolas.com.br/cotacoes/boi-gordo/boi-gordo-indicador-esalq-bmf
  - https://www.cepea.org.br/br/indicador/boi-gordo.aspx
- **dnit.condicoes-conservacao-rodovias** (critical, day 12): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 12): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **ibge.localidades-municipios** (critical, day 12): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicodados.ibge.gov.br/api/v1/localidades/municipios?orderBy=nome": dial tcp 170.84.40.178:443: i/o timeout)
  - https://servicodados.ibge.gov.br/api/v1/localidades/municipios?orderBy=nome
  - https://servicodados.ibge.gov.br/api/docs/localidades
- **ibge.localidades-ufs** (critical, day 12): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicodados.ibge.gov.br/api/v1/localidades/estados?orderBy=nome": dial tcp 170.84.40.178:443: i/o timeout)
  - https://servicodados.ibge.gov.br/api/v1/localidades/estados?orderBy=nome
  - https://servicodados.ibge.gov.br/api/docs/localidades
- **ibge.localidades-regioes** (critical, day 12): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicodados.ibge.gov.br/api/v1/localidades/regioes?orderBy=nome": dial tcp 170.84.40.178:443: i/o timeout)
  - https://servicodados.ibge.gov.br/api/v1/localidades/regioes?orderBy=nome
  - https://servicodados.ibge.gov.br/api/docs/localidades
- **ibge.localidades-mesorregioes** (critical, day 12): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicodados.ibge.gov.br/api/v1/localidades/mesorregioes?orderBy=nome": dial tcp 170.84.40.178:443: i/o timeout)
  - https://servicodados.ibge.gov.br/api/v1/localidades/mesorregioes?orderBy=nome
  - https://servicodados.ibge.gov.br/api/docs/localidades
- **ibge.localidades-microrregioes** (critical, day 12): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicodados.ibge.gov.br/api/v1/localidades/microrregioes?orderBy=nome": dial tcp 170.84.40.178:443: i/o timeout)
  - https://servicodados.ibge.gov.br/api/v1/localidades/microrregioes?orderBy=nome
  - https://servicodados.ibge.gov.br/api/docs/localidades
- **jrc.mars-crop-yield** (critical, day 12): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv)
  - https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv
  - https://mars.jrc.ec.europa.eu/dataset
- **suframa.comercio-mercadorias-zfm** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx": dial tcp 161.148.164.31:443: connect: connection refused)
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 12): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
