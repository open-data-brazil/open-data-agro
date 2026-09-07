## Source health daily probe

### Source health probe

- Executed at: 2026-09-07T08:01:47Z
- Run date: 2026-09-07
- Datasets probed: 131
- OK: 124 · Warning: 2 · Critical: 5
- Updated samples: 12
- Deprecated (2+ days): 5

### Source health alerts

- **dnit.condicoes-conservacao-rodovias** (critical, day 49): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 49): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **ibge.censo-agro-maquinario** (critical, day 3): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6880/n3/all/p/2017/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6880/n3/all/p/2017/v/all
  - https://censoagro2017.ibge.gov.br/
- **ibge.pnad-rural-renda-ocupacao** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6385/n3/11,12,13,14,15,16,17/p/last%201/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6385/n3/11,12,13,14,15,16,17/p/last%201/v/all
  - https://sidra.ibge.gov.br/pesquisa/pnad
- **suframa.comercio-mercadorias-zfm** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx": dial tcp 161.148.164.31:443: i/o timeout)
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 49): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 49): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip": dial tcp 161.148.164.31:443: i/o timeout)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
