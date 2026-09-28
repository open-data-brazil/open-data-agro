## Source health daily probe

### Source health probe

- Executed at: 2026-09-28T09:28:34Z
- Run date: 2026-09-28
- Datasets probed: 131
- OK: 122 · Warning: 3 · Critical: 6
- Updated samples: 13
- Deprecated (2+ days): 6

### Source health alerts

- **dnit.condicoes-conservacao-rodovias** (critical, day 70): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 70): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **ibge.localidades-mesorregioes** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://servicodados.ibge.gov.br/api/v1/localidades/mesorregioes?orderBy=nome": dial tcp 170.84.40.205:443: i/o timeout)
  - https://servicodados.ibge.gov.br/api/v1/localidades/mesorregioes?orderBy=nome
  - https://servicodados.ibge.gov.br/api/docs/localidades
- **inmet.estacoes-automaticas** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apitempo.inmet.gov.br/estacoes/T": EOF)
  - https://apitempo.inmet.gov.br/estacoes/T
  - https://portal.inmet.gov.br/paginas/catalogoaut
- **inmet.bdmep-diario** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://portal.inmet.gov.br/uploads/dadoshistoricos/2025.zip": read tcp 10.1.0.149:42148->201.57.198.188:443: read: connection reset by peer)
  - https://portal.inmet.gov.br/uploads/dadoshistoricos/2025.zip
  - https://bdmep.inmet.gov.br/
- **jrc.mars-crop-yield** (critical, day 13): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv)
  - https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv
  - https://mars.jrc.ec.europa.eu/dataset
- **suframa.comercio-mercadorias-zfm** (critical, day 4): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx": dial tcp 161.148.164.31:443: i/o timeout)
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 70): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 70): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
