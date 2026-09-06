## Source health daily probe

### Source health probe

- Executed at: 2026-09-06T07:43:51Z
- Run date: 2026-09-06
- Datasets probed: 131
- OK: 122 · Warning: 0 · Critical: 9
- Updated samples: 20
- Deprecated (2+ days): 9

### Source health alerts

- **antt.volume-trafego-pedagio** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (ckan package_show: Get "https://dados.antt.gov.br/api/3/action/package_show?id=volume-trafego-praca-pedagio": context deadline exceeded)
  - https://dados.antt.gov.br/dataset/volume-trafego-praca-pedagio
- **antt.receita-por-praca** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (ckan package_show: Get "https://dados.antt.gov.br/api/3/action/package_show?id=receita-por-praca": context deadline exceeded)
  - https://dados.antt.gov.br/dataset/receita-por-praca
- **antt.pracas-pedagio** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (ckan package_show: Get "https://dados.antt.gov.br/api/3/action/package_show?id=praca-de-pedagio": context deadline exceeded)
  - https://dados.antt.gov.br/dataset/praca-de-pedagio
- **dnit.condicoes-conservacao-rodovias** (critical, day 48): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 48): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **ibge.censo-agro-maquinario** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6880/n3/all/p/2017/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6880/n3/all/p/2017/v/all
  - https://censoagro2017.ibge.gov.br/
- **inmet.estacoes-convencionais** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://apitempo.inmet.gov.br/estacoes/M": EOF)
  - https://apitempo.inmet.gov.br/estacoes/M
  - https://portal.inmet.gov.br/paginas/catalogoman
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 48): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 48): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip": dial tcp 161.148.164.31:443: connect: connection refused)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
