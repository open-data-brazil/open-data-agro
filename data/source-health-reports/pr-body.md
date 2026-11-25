## Source health daily probe

### Source health probe

- Executed at: 2026-11-25T12:12:46Z
- Run date: 2026-11-25
- Datasets probed: 131
- OK: 111 · Warning: 6 · Critical: 14
- Updated samples: 23
- Deprecated (2+ days): 14

### Source health alerts

- **aneel.tarifas-energia** (critical, day 13): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (ckan package_show: Get "https://dadosabertos.aneel.gov.br/api/3/action/package_show?id=bandeiras-tarifarias": dial tcp 200.198.220.169:443: i/o timeout)
  - https://dadosabertos.aneel.gov.br/dataset/bandeiras-tarifarias
- **anp.combustiveis-precos-medios-municipios** (critical, day 3): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (fetch LPC listing: Get "https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas": net/http: TLS handshake timeout)
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas
- **anp.combustiveis-precos-postos** (critical, day 3): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (fetch LPC listing: Get "https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas": net/http: TLS handshake timeout)
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas
- **anp.etanol-precos** (critical, day 3): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (fetch LPC listing: Get "https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas": net/http: TLS handshake timeout)
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-indices/precos-de-combustiveis
- **antt.volume-trafego-pedagio** (warning, day 1): Possible link deprecation — official source unreachable after retries. (ckan package_show: Get "https://dados.antt.gov.br/api/3/action/package_show?id=volume-trafego-praca-pedagio": context deadline exceeded)
  - https://dados.antt.gov.br/dataset/volume-trafego-praca-pedagio
- **antt.receita-por-praca** (warning, day 1): Possible link deprecation — official source unreachable after retries. (ckan package_show: Get "https://dados.antt.gov.br/api/3/action/package_show?id=receita-por-praca": context deadline exceeded)
  - https://dados.antt.gov.br/dataset/receita-por-praca
- **antt.pracas-pedagio** (warning, day 1): Possible link deprecation — official source unreachable after retries. (ckan package_show: Get "https://dados.antt.gov.br/api/3/action/package_show?id=praca-de-pedagio": context deadline exceeded)
  - https://dados.antt.gov.br/dataset/praca-de-pedagio
- **cepea.soja-paranagua** (critical, day 13): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/soja/soja-indicador-cepea-esalq-porto-paranagua)
  - https://www.noticiasagricolas.com.br/cotacoes/soja/soja-indicador-cepea-esalq-porto-paranagua
  - https://www.cepea.org.br/br/indicador/soja.aspx
- **cepea.soja-parana** (critical, day 13): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/soja/indicador-cepea-esalq-soja-parana)
  - https://www.noticiasagricolas.com.br/cotacoes/soja/indicador-cepea-esalq-soja-parana
  - https://www.cepea.org.br/br/indicador/soja.aspx
- **cepea.milho** (critical, day 13): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/milho/indicador-cepea-esalq-milho)
  - https://www.noticiasagricolas.com.br/cotacoes/milho/indicador-cepea-esalq-milho
  - https://www.cepea.org.br/br/indicador/milho.aspx
- **cepea.boi-gordo** (critical, day 13): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://www.noticiasagricolas.com.br/cotacoes/boi-gordo/boi-gordo-indicador-esalq-bmf)
  - https://www.noticiasagricolas.com.br/cotacoes/boi-gordo/boi-gordo-indicador-esalq-bmf
  - https://www.cepea.org.br/br/indicador/boi-gordo.aspx
- **dnit.condicoes-conservacao-rodovias** (critical, day 21): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 21): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **inmet.estacoes-automaticas** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apitempo.inmet.gov.br/estacoes/T": context deadline exceeded (Client.Timeout exceeded while awaiting headers))
  - https://apitempo.inmet.gov.br/estacoes/T
  - https://portal.inmet.gov.br/paginas/catalogoaut
- **inmet.estacoes-convencionais** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apitempo.inmet.gov.br/estacoes/M": read tcp 10.1.0.243:41230->201.57.198.188:443: read: connection reset by peer)
  - https://apitempo.inmet.gov.br/estacoes/M
  - https://portal.inmet.gov.br/paginas/catalogoman
- **inmet.pacote-anual-automaticas** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://portal.inmet.gov.br/uploads/dadoshistoricos/2025.zip": read tcp 10.1.0.243:56944->201.57.198.188:443: read: connection reset by peer)
  - https://portal.inmet.gov.br/uploads/dadoshistoricos/2025.zip
  - https://portal.inmet.gov.br/dadoshistoricos
- **jrc.mars-crop-yield** (critical, day 21): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv)
  - https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv
  - https://mars.jrc.ec.europa.eu/dataset
- **suframa.comercio-mercadorias-zfm** (critical, day 3): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx": net/http: TLS handshake timeout)
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 21): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 3): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip": net/http: TLS handshake timeout)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
