## Source health daily probe

### Source health probe

- Executed at: 2026-10-04T09:26:03Z
- Run date: 2026-10-04
- Datasets probed: 131
- OK: 119 · Warning: 0 · Critical: 12
- Updated samples: 15
- Deprecated (2+ days): 12

### Source health alerts

- **bcb.cim-agro-credito-rural** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://api.bcb.gov.br/dados/serie/bcdata.sgs.21087/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json": dial tcp: lookup api.bcb.gov.br on 127.0.0.53:53: no such host)
  - https://api.bcb.gov.br/dados/serie/bcdata.sgs.21087/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json
  - https://www.bcb.gov.br/publicacoes/cim
- **bcb.sgs-ipca** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://api.bcb.gov.br/dados/serie/bcdata.sgs.433/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json": dial tcp: lookup api.bcb.gov.br on 127.0.0.53:53: no such host)
  - https://api.bcb.gov.br/dados/serie/bcdata.sgs.433/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json
  - https://dadosabertos.bcb.gov.br/
- **bcb.sgs-ipca-12m** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://api.bcb.gov.br/dados/serie/bcdata.sgs.13522/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json": dial tcp: lookup api.bcb.gov.br on 127.0.0.53:53: no such host)
  - https://api.bcb.gov.br/dados/serie/bcdata.sgs.13522/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json
  - https://dadosabertos.bcb.gov.br/
- **bcb.sgs-ptax-usd-venda** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://api.bcb.gov.br/dados/serie/bcdata.sgs.1/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json": dial tcp: lookup api.bcb.gov.br on 127.0.0.53:53: no such host)
  - https://api.bcb.gov.br/dados/serie/bcdata.sgs.1/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json
  - https://dadosabertos.bcb.gov.br/
- **bcb.sgs-ptax-usd-compra** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://api.bcb.gov.br/dados/serie/bcdata.sgs.10813/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json": dial tcp: lookup api.bcb.gov.br on 127.0.0.53:53: no such host)
  - https://api.bcb.gov.br/dados/serie/bcdata.sgs.10813/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json
  - https://dadosabertos.bcb.gov.br/
- **bcb.sgs-igpm** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://api.bcb.gov.br/dados/serie/bcdata.sgs.189/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json": dial tcp: lookup api.bcb.gov.br on 127.0.0.53:53: no such host)
  - https://api.bcb.gov.br/dados/serie/bcdata.sgs.189/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json
  - https://dadosabertos.bcb.gov.br/
- **bcb.sgs-selic** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://api.bcb.gov.br/dados/serie/bcdata.sgs.11/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json": dial tcp: lookup api.bcb.gov.br on 127.0.0.53:53: no such host)
  - https://api.bcb.gov.br/dados/serie/bcdata.sgs.11/dados?dataFinal=04%2F10%2F2026&dataInicial=04%2F07%2F2026&formato=json
  - https://dadosabertos.bcb.gov.br/
- **dnit.condicoes-conservacao-rodovias** (critical, day 76): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 76): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **jrc.mars-crop-yield** (critical, day 19): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv)
  - https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv
  - https://mars.jrc.ec.europa.eu/dataset
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 76): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 76): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
