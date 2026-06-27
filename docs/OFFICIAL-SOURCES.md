# Official sources catalog

> Primary references per dataset. **No dataset without a cited official source.**

**Initial portal:** [CONAB — Downloads de Arquivos](https://portaldeinformacoes.conab.gov.br/download-arquivos.html)

**Status column:** `**Pn — implemented**` = full E2E pipeline (ingest → GE → silver → dbt → DuckDB), verified via `make *-mvp` and CI collection gates. Priority `Pn` reflects collection sprint ordering, not implementation state.

**Historical depth:** per-dataset source min years and `--from` backfill examples — [DATA-HISTORY-RANGES.md](DATA-HISTORY-RANGES.md) (Phase 33).

**Deferred / unreachable sources** (removed from active catalog 2026-06-27): [DEFERRED-SOURCES.md](DEFERRED-SOURCES.md).

---

## CONAB — Portal de Informações Agropecuárias

| Item | Value |
|------|-------|
| **Portal** | https://portaldeinformacoes.conab.gov.br/ |
| **Downloads** | https://portaldeinformacoes.conab.gov.br/download-arquivos.html |
| **Usage** | Reproduction allowed non-commercial with source citation; preserve data integrity |
| **Contact** | sutin@conab.gov.br |

Detailed per-dataset mapping: `.local/phases/10-conab-producao-agricola/OFFICIAL-REFERENCE.md` (local).

---

## Index

### Produção Agrícola (Phase 10 — MVP)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `conab.estimativa-graos` | Estimativa Grãos | **P0 — implemented** |
| `conab.serie-historica-graos` | Série Histórica Grãos | **P0 — implemented** |
| `conab.estimativa-cana` | Estimativa Cana-de-Açúcar | **P1 — implemented** |
| `conab.serie-historica-cana` | Série Histórica Cana-de-Açúcar | **P1 — implemented** |
| `conab.estimativa-cafe` | Estimativa Café | **P1 — implemented** |
| `conab.serie-historica-cafe` | Série Histórica Café | **P1 — implemented** |
| `conab.custo-producao` | Custo de Produção | **P1 — implemented** |

### Mercado (Phase 11)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `conab.oferta-demanda` | Oferta e Demanda | **P1 — implemented** |
| `conab.precos-minimos` | Preços Mínimos | **P2 — implemented** |
| `conab.precos-agropecuarios-mensal-uf` | Preços agropecuários Mensal UF | **P1 — implemented** |
| `conab.precos-agropecuarios-mensal-municipio` | Preços agropecuários Mensal Município | **P1 — implemented** |
| `conab.precos-agropecuarios-semanal-uf` | Preços agropecuários Semanal UF | **P1 — implemented** |
| `conab.precos-agropecuarios-semanal-municipio` | Preços agropecuários Semanal Municipio | **P1 — implemented** |
| `conab.prohort-diario` | Prohort Diário | **P3 — implemented** |
| `conab.prohort-mensal` | Prohort Mensal | **P3 — implemented** |

### Abastecimento (Phase 12)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `conab.estoques-publicos` | Estoques Públicos | **P1 — implemented** |
| `conab.operacoes-comercializacao` | Operações de Comercialização | **P2 — implemented** |
| `conab.vendas-balcao` | Vendas em Balcão | **P2 — implemented** |

### ANP — Combustíveis (Phase 12 extension)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `anp.combustiveis-precos-medios-municipios` | LPC — preços médios por município | **P1 — implemented** |
| `anp.combustiveis-precos-postos` | LPC — preços por posto revendedor | **P1 — implemented** |

### Armazenamento e Logística (Phase 13)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `conab.armazenagem` | Armazenagem | **P1 — implemented** |
| `conab.frete` | Frete | **P1 — implemented** |
| `conab.serie-historica-capacidade-estatica` | Série Histórica da Capacidade Estática | **P1 — implemented** |

### Agricultura Familiar (Phase 14)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `conab.alimenta-brasil-entregas` | Programa Alimenta Brasil - Entregas | **P1 — implemented** |
| `conab.alimenta-brasil-propostas` | Programa Alimenta Brasil - Propostas | **P1 — implemented** |

### IBGE — Localidades (Phase 15)

| Dataset ID | API resource | Status |
|------------|--------------|--------|
| `ibge.localidades-municipios` | `/api/v1/localidades/municipios` | **P0 — implemented** |
| `ibge.localidades-ufs` | `/api/v1/localidades/estados` | **P0 — implemented** |
| `ibge.localidades-regioes` | `/api/v1/localidades/regioes` | **P1 — implemented** |
| `ibge.localidades-mesorregioes` | `/api/v1/localidades/mesorregioes` | **P2 — implemented** |
| `ibge.localidades-microrregioes` | `/api/v1/localidades/microrregioes` | **P2 — implemented** |

**Fonte oficial:** [IBGE API de Localidades](https://servicodados.ibge.gov.br/api/docs/localidades)

### IBGE — PAM Produção Agrícola Municipal (Phase 16)

| Dataset ID | SIDRA table | Status |
|------------|-------------|--------|
| `ibge.pam-area-quantidade` | 1612 — área plantada, colhida, quantidade | **P0 — implemented** |
| `ibge.pam-rendimento-valor` | 1613 — rendimento médio, valor da produção | **P1 — implemented** |
| `ibge.pam-estabelecimentos` | 5457 — número de estabelecimentos | **P2 — implemented** |

**Fonte oficial:** [IBGE SIDRA — PAM](https://sidra.ibge.gov.br/pesquisa/pam) · API: [apisidra.ibge.gov.br](https://apisidra.ibge.gov.br/)

MVP crops (soja, milho, trigo) use SIDRA classification `c81` (1612), `c82` (1613), `c782` (5457) with product codes `2713`, `2711`, `2716`.

### IBGE — LSPA Produção Agrícola (Phase 37)

| Dataset ID | SIDRA table | Status |
|------------|-------------|--------|
| `ibge.lspa-area-producao` | 6588 — área plantada, colhida, produção mensal por UF | **P0 — implemented** |

**Fonte oficial:** [IBGE SIDRA — LSPA](https://sidra.ibge.gov.br/pesquisa/lspa) · historical monthly series table **6588** (UF grain, `c48` crop classification)

Core crops: soja `39443`, milho `39441`, trigo `39445`. Variables `109` (área plantada), `216` (área colhida), `35` (produção). Complements CONAB `estimativa-graos`.

### INMET — Clima Histórico (Phase 17)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `inmet.estacoes-automaticas` | Portal catálogo CSV | **P0 — implemented** |
| `inmet.bdmep-diario` | BDMEP annual ZIP (`/uploads/dadoshistoricos/{year}.zip`) | **P0 — implemented** |
| `inmet.estacoes-convencionais` | Portal catálogo CSV | **P1 — implemented** |
| `inmet.bdmep-mensal` | Monthly rollups from annual ZIP | **P2 — implemented** |
| `inmet.pacote-anual-automaticas` | BDMEP bulk annual ZIP | **P1 — implemented** |

**Fonte oficial:** [BDMEP — INMET](https://bdmep.inmet.gov.br/) · Portal: [portal.inmet.gov.br/dadoshistoricos](https://portal.inmet.gov.br/dadoshistoricos)

Timestamps in source files are **UTC**; missing values use sentinels `9999`, `Null`, or blank per INMET documentation.

### BCB — Séries Macroeconômicas SGS (Phase 18)

| Dataset ID | SGS code | Status |
|------------|----------|--------|
| `bcb.sgs-ipca` | 433 — IPCA variação mensual (%) | **P0 — implemented** |
| `bcb.sgs-ptax-usd-venda` | 1 — Dólar PTAX venda | **P0 — implemented** |
| `bcb.sgs-ipca-12m` | 13522 — IPCA acumulado 12 meses | **P1 — implemented** |
| `bcb.sgs-igpm` | 189 — IGP-M variação mensual | **P2 — implemented** |
| `bcb.sgs-ptax-usd-compra` | 10813 — Dólar PTAX compra | **P2 — implemented** |
| `bcb.sgs-selic` | 11 — Taxa Selic meta (% a.a.) | **P1 — implemented** |

**Fonte oficial:** [BCB Dados Abertos](https://dadosabertos.bcb.gov.br/) · API: [api.bcb.gov.br](https://api.bcb.gov.br/)

Historical backfill paginates `dataInicial`/`dataFinal` in ≤10-year chunks per BCB API limits. PTAX series from **1984** — see [DATA-HISTORY-RANGES.md](DATA-HISTORY-RANGES.md).

### CEPEA — Preços Agro (Phase 19)

| Dataset ID | Indicator | Status |
|------------|-----------|--------|
| `cepea.soja-paranagua` | Soja — Paranaguá port (R$/sc 60 kg) | **P0 — implemented** |
| `cepea.soja-parana` | Soja — Paraná regional | **P1 — implemented** |
| `cepea.milho` | Milho — Campinas | **P1 — implemented** |
| `cepea.boi-gordo` | Boi gordo — São Paulo | **P2 — implemented** |

**Fonte oficial:** [CEPEA/ESALQ-USP](https://www.cepea.org.br/) · **License:** [CC BY-NC 4.0](https://www.cepea.org.br/br/licenca-de-uso-de-dados.aspx) — market reference (`fonte_tipo=referencia_mercado`), not `.gov.br`.

Programmatic ingest tries the CEPEA portal first; when Cloudflare blocks access, it falls back to the Notícias Agrícolas mirror (same CEPEA indicators). Full historical backfill from 2010 uses CEPEA “Consulta ao Banco de Dados” export; live ingest captures the latest published window.

Crossing with CONAB local prices (Phase 11) and BCB PTAX (Phase 18) is planned in analytics — see `.local/phases/DATA-CROSSING-VISION.md`.

### MDIC — Comex Stat (Phase 21 + 35)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `mdic.comex-exportacao-ncm-mes` | Comex Stat API — exportação mensal NCM agro | **P0 — implemented** |
| `mdic.comex-importacao-ncm-mes` | Comex Stat API — importação mensal NCM fertilizantes | **P0 — implemented** |
| `mdic.comex-exportacao-uf-ncm` | Comex Stat API — exportação mensal UF × NCM agro | **P0 — implemented** |
| `mdic.comex-importacao-diesel-ncm` | Comex Stat API — importação diesel/óleos combustíveis | **P0 — implemented** |

**Fonte oficial:** [Comex Stat — MDIC](https://comexstat.mdic.gov.br/) · API: [api-comexstat.mdic.gov.br](https://api-comexstat.mdic.gov.br/docs)

Monthly export FOB (USD) and quantity (kg) for ag commodities; import CIF for fertilizers and diesel; state-level export by UF. NCM → `produto_slug` mapping in [GLOSSARY.md](GLOSSARY.md). Historical backfill from 2015 via year-chunked `POST /general` requests.

### ANTT — Logística rodoviária / Pedágios (Phase 22 + 34)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `antt.pracas-pedagio` | Praças de Pedágio — concessionárias | **P0 — implemented** |
| `antt.volume-trafego-pedagio` | Volume de Tráfego por Praça — mensal consolidado | **P0 — implemented** |
| `antt.receita-por-praca` | Receita por Praça — mensal por praça | **P1 — implemented** |

**Fonte oficial:** [ANTT — Portal de Dados Abertos](https://dados.antt.gov.br/) · CKAN packages `praca-de-pedagio`, `volume-trafego-praca-pedagio`, `receita-por-praca`

Toll plaza locations (`antt.pracas-pedagio`), monthly traffic volume by vehicle category (`antt.volume-trafego-pedagio`), and monthly revenue per plaza (`antt.receita-por-praca`). Complements CONAB `conab.frete` (Phase 13) for highway logistics cost context. Tariff-by-category series (`antt.tarifas-pedagio`) deferred — CKAN package not published on dados.antt.gov.br (Phase 34 discovery).

### MAPA — Dados Abertos / ZARC (Phase 23)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `mapa.zarc-tabua-risco` | ZARC — Tábua de Risco Climático | **P0 — implemented** |

**Fonte oficial:** [MAPA — Portal de Dados Abertos](https://dados.agricultura.gov.br/dataset/tabua-de-risco-zoneamento-agricola-de-risco-climatico) · CKAN package `tabua-de-risco-zoneamento-agricola-de-risco-climatico` (latest annual safra CSV)

Municipal planting-window climate risk (`dec1`–`dec36`) by culture, soil cycle, and management type. Unique vs CONAB/IBGE supply series — complements PAM and estimativa with policy-driven planting constraints for soja, milho, trigo, and other cultures.

### MAPA — Agrofit registry (Phase 40)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `mapa.agrofit-produtos-formulados` | Agrofit — produtos formulados (defensivos) | **P1 — implemented** |
| `mapa.agrofit-produtos-tecnicos` | Agrofit — produtos técnicos (defensivos) | **P2 — implemented** |

**Fonte oficial:** [MAPA — Agrofit](https://dados.agricultura.gov.br/dataset/sistema-de-agrotoxicos-fitossanitarios-agrofit) · CKAN bulk CSV `agrofitprodutosformulados.csv` / `agrofitprodutostecnicos.csv`

Crop protection product registry — formulated products by culture/pest; technical active-ingredient registry. `fonte_tipo: oficial_gov_br`.

### ANA — Hidrologia (Phase 40)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `ana.hidrologia-series` | HidroWeb — séries diárias de vazão (estações selecionadas) | **P2 — implemented** |

**Fonte oficial:** [ANA — Dados abertos](https://www.gov.br/ana/pt-br/acesso-a-informacao/dados-abertos) · SOAP `HidroSerieHistorica` em `telemetriaws1.ana.gov.br` (sem API key)

Daily flow series for configured fluviometric stations. Complements INMET climate for hydrology context.

### ANTAQ — Movimentação portuária (Phase 40)

**Deferred** — bulk export URL returns HTTP 404. See [DEFERRED-SOURCES.md](DEFERRED-SOURCES.md#summary) (`antaq.movimentacao-carga-portuaria`).

### DNIT — SNV rodovias federais (Phase 44)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `dnit.snv-rodovias-federais` | SNV — jurisdição de vias federais | **P0 — implemented** |

**Fonte oficial:** [DNIT — Dados abertos](https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias) · CKAN package `jurisdicao-de-vias` (latest CSV, semicolon delimiter)

Federal highway jurisdiction segments — BR code, UF, km range, administration, surface type. Metadata preamble rows stripped before bronze parse.

### IPEA — Séries macro regionais (Phase 44)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `ipea.series-macro-regionais` | IPEA OData4 — componentes IDH agro/com rural | **P1 — implemented** |

**Fonte oficial:** [IPEA Data](http://www.ipeadata.gov.br/) · OData4 `ValoresSerie(SERCODIGO='...')` — series `ADH_P_AGRO_RUR`, `ADH_P_COM_RUR`

Regional macro indicators by territory — annual refdate grain for Brazil and UF-level states.

### IBGE — PEVS produção vegetal (Phase 44)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `ibge.pevs-producao-vegetal` | PEVS — quantidade e valor da extração vegetal por UF | **P1 — implemented** |

**Fonte oficial:** [IBGE SIDRA — PEVS](https://sidra.ibge.gov.br/pesquisa/pevs) · table **289**, UF grain, variables **144** (quantidade) and **145** (valor), annual

Plant extraction production statistics at UF level — complements PAM/LSPA annual crop series.

### IBGE — PPM produção municipal (Phase 44)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `ibge.ppm-producao-municipal` | PPM — quantidade e valor da produção municipal | **P2 — implemented** |

**Fonte oficial:** [IBGE SIDRA — PAM/PPM](https://sidra.ibge.gov.br/pesquisa/pam) · table **74**, municipal grain (`n6/in n3` UF batches), variables **106** (quantidade) and **215** (valor), annual 2010–2023

Municipal animal-origin production quantity and value — SIDRA chunked by UF batches.

### ANEEL — Bandeiras tarifárias (Phase 44)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `aneel.tarifas-energia` | Bandeiras tarifárias — acionamento mensal | **P2 — implemented** |

**Fonte oficial:** [ANEEL Dados Abertos](https://dadosabertos.aneel.gov.br/dataset/bandeiras-tarifarias) · CKAN package `bandeiras-tarifarias`, CSV semicolon, resource name contains `Acionamento`

Electricity tariff flag activation history — input cost context for farm energy.

### BNDES — Financiamento agro (Phase 44)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `bndes.financiamento-agro` | Desembolsos por setor CNAE — coluna agropecuária | **P2 — implemented** |

**Fonte oficial:** [BNDES Dados Abertos](https://dadosabertos.bndes.gov.br/dataset/desembolsos) · CKAN package `desembolsos`, CSV semicolon, resource `Por setor CNAE`

Monthly BNDES disbursements to agropecuaria sector (CNAE grouping).

### INMET — Monitor de secas (Phase 44)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `inmet.sequia-monitor` | ANA Monitor de Secas — áreas por categoria S0–S4 | **P1 — implemented** |

**Fonte oficial:** [ANA Monitor de Secas](https://www.gov.br/ana/pt-br/servicos/monitor-de-seca) · API `https://apimsbr.ana.gov.br/rpc/v1/dados-tabulares-monitor` (JSON, no auth)

Drought severity area statistics by map/month — catalog under INMET agency for climate feature grouping.

### BR sources wave 5 — MAPA (Phase 52)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `mapa.sipeagro-estabelecimentos` | SIPEAGRO — estabelecimentos registrados (8 linhas de produto CKAN) | **P0 — implemented** |
| `mapa.sipeagro-produtos` | SIPEAGRO — insumos (fertilizantes, veterinários, alimentação animal) | **P0 — implemented** |
| `mapa.sigef-producao-sementes` | SIGEF — campos de produção de sementes | **P0 — implemented** |
| `mapa.sigef-areas` | SIGEF — declarações de área para produção (uso próprio) | **P1 — implemented** |
| `mapa.sisser-seguro-rural` | SISSER — Programa Seguro Rural (PSR) propostas | **P1 — implemented** |

**Fonte oficial:** [MAPA dados abertos](https://dados.agricultura.gov.br/) · CKAN packages [SIPEAGRO](https://dados.agricultura.gov.br/dataset/sipeagro) · [SIGEF](https://dados.agricultura.gov.br/dataset/dados-referentes-ao-controle-da-producao-de-sementes-sigef) · [SISSER](https://dados.agricultura.gov.br/dataset/sisser3)

SIPEAGRO merges multiple CSV resources with normalized columns (`linha_produto`, `uf`, `numero_registro_estabelecimento`, …). SISSER merges all PSR CSV periods (2006–2025). `mapa.fertilizantes-registro` rejected — use SIPEAGRO Fertilizantes line.

**Gate:** `make br-sources-wave-5-mapa-mvp` · `make ci-br-sources-wave-5-mapa-mvp`

### BR sources wave 5 — IBGE (Phase 53)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `ibge.ppm-efetivo-rebanhos` | PPM — efetivo dos rebanhos por UF (SIDRA 3939) | **P0 — implemented** |
| `ibge.ppm-vacas-ordenhadas` | PPM — vacas ordenhadas por UF (SIDRA 94) | **P0 — implemented** |
| `ibge.ppm-ovinos-tosquiados` | PPM — ovinos tosquiados por UF (SIDRA 95) | **P0 — implemented** |
| `ibge.ppm-aquicultura` | PPM — aquicultura por UF (SIDRA 3940) | **P0 — implemented** |
| `ibge.pam-precos-produtor` | PAM — valor e quantidade municipal (proxy preço produtor, SIDRA 1612) | **P1 — implemented** |
| `ibge.pam-culturas-estendidas` | PAM — 18 culturas prioritárias além de soja/milho/trigo | **P1 — implemented** |
| `ibge.lspa-rendimento-medio` | LSPA — rendimento médio mensal por UF (SIDRA 6588, var 35) | **P1 — implemented** |
| `ibge.censo-agro-area-uso-solo` | Censo Agro 2017 — área e uso do solo por UF (SIDRA 6879) | **P1 — implemented** |
| `ibge.censo-agro-maquinario` | Censo Agro 2017 — maquinário e tecnologia por UF (SIDRA 6880) | **P1 — implemented** |
| `ibge.pnad-rural-renda-ocupacao` | PNAD Contínua — renda e ocupação rural por UF (SIDRA 6385) | **P2 — implemented** |

**Fonte oficial:** [IBGE SIDRA — PPM](https://sidra.ibge.gov.br/pesquisa/ppm) · [PAM](https://sidra.ibge.gov.br/pesquisa/pam) · [LSPA](https://sidra.ibge.gov.br/pesquisa/lspa) · [Censo Agro 2017](https://censoagro2017.ibge.gov.br/) · [PNAD](https://sidra.ibge.gov.br/pesquisa/pnad) · API: [apisidra.ibge.gov.br](https://apisidra.ibge.gov.br/)

PPM herd tables use UF territory (`n3`); PAM/LSPA reuse municipal/monthly SIDRA batching from Phases 16/37. Censo 2017 tables are static decennial snapshots. PNAD rural uses `v/all` (specific variable codes return HTTP 400 on table 6385).

**Gate:** `make br-sources-wave-5-ibge-mvp` · `make ci-br-sources-wave-5-ibge-mvp`

### BR sources wave 4 (Phase 48)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `ibge.censo-agro-estabelecimentos` | Censo Agro 2017 — establishments and area by UF (SIDRA 6878) | **P0 — implemented** |
| `ibge.pnad-continua-rural` | PNAD Contínua occupation by UF (SIDRA 6385) | **P1 — implemented** |
| `suframa.comercio-mercadorias-zfm` | SUFRAMA ZFM inspected invoices (XLSX) | **P1 — implemented** |
| `transportes.mtr-bit-malha-rodoviaria` | MTR BIT / DNIT SNV federal road metadata CSV | **P1 — implemented** |
| `mapa.sif-abate-estatisticas` | MAPA SIGSIF slaughter statistics by UF | **P1 — implemented** |
| `ons.carga-energetica` | ONS daily energy load by subsystem | **P1 — implemented** |
| `inpe.deter-alertas-desmatamento` | INPE DETER Amazon alerts (TerraBrasilis WFS) | **P2 — implemented** |
| `dnit.condicoes-conservacao-rodovias` | DNIT pavement condition surveys | **P2 — implemented** |

**Fonte oficial:** [Censo Agro 2017](https://censoagro2017.ibge.gov.br/) · [SIDRA](https://apisidra.ibge.gov.br/) · [SUFRAMA dados abertos](https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos) · [MTR BIT](https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas) · [MAPA SIGSIF](https://dados.agricultura.gov.br/dataset/servico-de-inspecao-federal-sif) · [ONS carga energia](https://dados.ons.org.br/dataset/carga-energia) · [INPE TerraBrasilis](https://terrabrasilis.dpi.inpe.br/downloads/) · [DNIT condições pavimento](https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento)

**Deferred:** MTR BIT georeferenced shapefile ZIP (ingest SNV CSV metadata only); SUFRAMA aquisições PDF reports (use NF vistoriadas XLSX); INPE full DETER shapefile bulk (WFS capped subset).

### B3 — Mercado futuro agro (Phase 24)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `b3.futuro-soja` | SOY — preço de ajuste diário (BVBG.187 SPRD) | **P0 — implemented** |
| `b3.futuro-milho` | CCM — milho futuro | **P1 — implemented** |
| `b3.futuro-boi` | BGI — boi gordo futuro | **P1 — implemented** |

**Fonte oficial:** [B3 — Pesquisa por pregão](https://www.b3.com.br/pt_br/market-data-e-indices/servicos-de-dados/market-data/historico/boletins-diarios/pesquisa-por-pregao/) · arquivo `SPRD{YYMMDD}.zip` (Boletim simplificado derivativos)

Daily futures settlement (`AdjstdQt`) by contract symbol. Regulated exchange reference — not `.gov.br`. License documented in catalog and OFFICIAL-REFERENCE. No synthetic continuous rolls in bronze.

### USDA FAS — PSD global supply (Phase 25)

**Deferred** — `apps.fas.usda.gov` unreachable from BR/CI without US egress. See [DEFERRED-SOURCES.md](DEFERRED-SOURCES.md) (`usda.psd-soja`, `usda.psd-milho`, `usda.psd-trigo`).

**Still active:** `usda.wasde` — [USDA WASDE](https://www.usda.gov/oce/commodity-markets/wasde) (HTML index, no PSD SOAP).

### FAO — FAOSTAT (Phase 26 + 36)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `fao.prices-agro` | Producer prices — soja, milho, trigo, carne bovina | **P0 — implemented** |
| `fao.producao-agro` | Annual production by country — soja, milho, trigo, carne bovina | **P0 — implemented** |

**Fonte oficial:** [FAO FAOSTAT — Producer Prices (PP)](https://www.fao.org/faostat/en/#data/PP) · bulk `Prices_E_All_Data_(Normalized).zip` · [Production (QCL)](https://www.fao.org/faostat/en/#data/QCL) · `Production_Crops_Livestock_E_All_Data_(Normalized).zip` (no API key)

Producer prices (USD/tonne) and price indices by country × year. Production element `5510`. Items 236/56/15/867. `fonte_tipo: internacional_oficial`.

**Deferred:** `fao.comercio-agro` — Trade bulk ZIP returns HTTP 403. See [DEFERRED-SOURCES.md](DEFERRED-SOURCES.md).
### World Bank — Pink Sheet commodities (Phase 27 + 36)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `worldbank.pink-sheet-monthly` | Pink Sheet monthly prices — soja, milho, trigo, petróleo, carne | **P0 — implemented** |
| `worldbank.ag-indices` | Pink Sheet agriculture sub-indices (2010=100) | **P1 — implemented** |

**Fonte oficial:** [World Bank Commodity Markets](https://www.worldbank.org/en/research/commodity-markets) · `CMO-Historical-Data-Monthly.xlsx` bulk — sheets `Monthly Prices` and `Monthly Indices` (no API key)

USD-denominated monthly commodity reference prices and agriculture sub-indices. Monthly grain — no daily resampling in bronze (Stage H policy). `fonte_tipo: internacional_oficial`.

### NOAA — global climate indices (Phase 28)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `noaa.enso-indices` | Oceanic Niño Index (ONI) — seasonal ENSO SST anomaly | **P2 — implemented** |
| `noaa.global-temp-anomaly` | Global land+ocean monthly temperature anomaly | **P2 — implemented** |

**Fonte oficial:** [NOAA CPC ONI](https://www.cpc.ncep.noaa.gov/products/analysis_monitoring/ensostuff/ONI_v5.php) · `oni.ascii.txt` · [NCEI Climate at a Glance](https://www.ncei.noaa.gov/access/monitoring/climate-at-a-glance/global/time-series) · CSV `globe/land_ocean/0/0/{start}-{end}.csv` (no API key)

Global climate shock features complementing INMET (Phase 17). `fonte_tipo: internacional_oficial`. ONI is seasonal (3-month running mean); global temp is monthly.

### U.S. EIA — Petroleum prices (Phase 38)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `eia.petroleum-prices` | Daily WTI + Brent spot prices | **P0 — implemented** |
| `usda.wasde` | WASDE monthly supply/demand estimates | **P1 — implemented** |
| `igc.goi-index` | IGC GOI daily export price index + sub-indices | **P1 — implemented** |
| `un.comtrade-bulk` | UN Comtrade bulk API — Brazil HS ag chapters | **P1 — implemented** |

**Fonte oficial:** [U.S. EIA Open Data](https://www.eia.gov/opendata/) · API v2 `seriesid` route · series `PET.RWTC.D` (WTI), `PET.RBRTE.D` (Brent) · free `EIA_API_KEY` required for live fetch · [USDA WASDE](https://www.usda.gov/oce/commodity-markets/wasde) · [IGC GOI](https://igc.int/en/public-site/markets/marketinfo-goi.aspx) · [UN Comtrade API](https://uncomtrade.org/docs/un-comtrade-api/)

Global oil shock reference complementing World Bank Pink Sheet crude oil. WASDE supply/demand, IGC competitor price index, and UN bilateral trade for Brazil ag HS chapters. `fonte_tipo: internacional_oficial`. Daily grain for EIA/IGC — no intraday resampling in bronze.

### International sources wave 2 (Phase 41)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `eurostat.ag-prices` | EU agricultural output price indices (2015=100) | **P2 — implemented** |
| `argentina.bcra-cambio` | BCRA official USD exchange-rate daily series | **P2 — implemented** |

**Fonte oficial:** [EUROSTAT agriculture database](https://ec.europa.eu/eurostat/web/agriculture/database) · dataset `apri_pi15_outa` JSON API (no key) · [BCRA estadísticas cambiarias](https://api.bcra.gob.ar/estadisticascambiarias/v1.0/Cotizaciones/USD) (no key)

EU ag price reference and Argentina FX parity for competitor market models. `fonte_tipo: internacional_oficial`.

**Deferred:** `usda.gats-trade` — see [DEFERRED-SOURCES.md](DEFERRED-SOURCES.md).
### International sources wave 3 (Phase 45)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `oecd-fao.ag-outlook` | OECD-FAO Agricultural Outlook SDMX CSV — Brazil soy/maize/wheat balances | **P0 — implemented** |
| `fao.food-price-index` | FAO monthly Food Price Index + sub-indices (2002-2004=100) | **P1 — implemented** |
| `argentina.magyp-producion-granos` | MAGyP annual grain production via datos.gob.ar series API | **P1 — implemented** |

**Fonte oficial:** [OECD-FAO Outlook](https://www.oecd.org/en/data/datasets/oecd-fao-agricultural-outlook.html) · SDMX `https://sdmx.oecd.org/public/rest/data/` (no key) · [FAO FFPI](https://www.fao.org/worldfoodsituation/foodpricesindex/en/) · CSV bulk (no key) · [MAGyP datos abiertos](https://datos.magyp.gob.ar/) · [datos.gob.ar series API](https://apis.datos.gob.ar/series/api/) (no key)

**Deferred (verified):** `imf.commodity-prices` (no PCPS bulk), `paraguay.bcp-exportaciones-soja`, `uruguay.ine-exportaciones-agro`, `china.nbs-soy-imports` (403), `usda.ams-grain-prices` (403), `baltic.bdi-index` (subscription).

### International sources wave 4 (Phase 49)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `cftc.cot-agricultural-futures` | CFTC Commitment of Traders — ag futures positioning | **P0 — implemented** |
| `jrc.mars-crop-yield` | JRC MARS ASAP crop yield forecasts | **P1 — implemented** |
| `fao.giews-crop-prospects` | FAO GIEWS crop prospects / food prices | **P1 — implemented** |
| `fao.amis-market-monitor` | FAO AMIS market monitor | **P1 — implemented** |
| `sagis.grain-supply-statistics` | SAGIS South Africa grain supply | **P1 — implemented** |
| `japan.maff-ag-trade` | MAFF Japan ag production + trade | **P1 — implemented** |
| `fred.commodity-indexes` | FRED commodity price indexes | **P1 — implemented** |
| `nasa.power-agroclimatology` | NASA POWER agroclimatology point API | **P1 — implemented** |
| `copernicus.era5-agroclimate` | Copernicus CDS ERA5 reanalysis (sample/API) | **P2 — implemented** |

**Fonte oficial:** [CFTC COT](https://www.cftc.gov/MarketReports/CommitmentsofTraders/index.htm) · [JRC MARS](https://mars.jrc.ec.europa.eu/) · [FAO GIEWS](https://www.fao.org/giews/) · [FAO AMIS](https://www.amis-outlook.org/) · [SAGIS](https://www.sagis.org.za/) · [MAFF Japan](https://www.maff.go.jp/e/index.html) · [FRED](https://fred.stlouisfed.org/) · [NASA POWER](https://power.larc.nasa.gov/) · [Copernicus CDS](https://cds.climate.copernicus.eu/)

**Deferred (2026-06-27):** `wto.its-trade-statistics`, `mexico.siap-produccion-agricola`, `noaa.gpcc-precipitation` — see [DEFERRED-SOURCES.md](DEFERRED-SOURCES.md).

**Other deferred (verified):** `imf.commodity-prices`, `usda.ams-grain-prices`, `iea.world-energy-statistics`, Mercosur BCP/INE bulk, `china.nbs-soy-imports`, `baltic.bdi-index` (subscription). Fixture-backed ingest for sources blocked from CI networks (MAFF, GIEWS, AMIS, SAGIS) — live fetch via env bulk paths documented in phase OFFICIAL-REFERENCE.

### Wave 5 — approved candidates (Phase 51 discovery, ingest Phases 52–56)

> Live probe gate: `python3 scripts/ci/verify_wave5_discovery_probe.py`. Full fichas: `.local/SOURCE-DISCOVERY-CATALOG.md` § Wave 5 (gitignored).

| Phase | Approved `dataset_id` | Primary URL (verified 2026-06-27) |
|-------|----------------------|-----------------------------------|
| 52 MAPA | `mapa.sipeagro-estabelecimentos`, … `mapa.sisser-seguro-rural` (5 datasets) | [dados.agricultura.gov.br](https://dados.agricultura.gov.br/) CKAN — **Phase 52 ✅ implemented** |
| 53 IBGE | `ibge.ppm-efetivo-rebanhos`, … `ibge.pnad-rural-renda-ocupacao` (10 datasets) | [apisidra.ibge.gov.br](https://apisidra.ibge.gov.br/) — **Phase 53 ✅ implemented** |
| 54 Env/logistics | `ibama.sisfogo-incendios`, `ibama.licencas-ambientais`, `ibama.autos-infracao`, `ana.pluviometria-redes`, `embrapa.agroapi-agrofit`, `transportes.mtr-bit-malha-shapefile` | [dadosabertos.ibama.gov.br](https://dadosabertos.ibama.gov.br/) · [gov.br/ana](https://www.gov.br/ana/pt-br/acesso-a-informacao/dados-abertos) · [AgroAPI store](https://www.agroapi.cnptia.embrapa.br/store/) · [MTR BIT](https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas) — **Phase 54 ✅ implemented** |
| 55 Finance | `bcb.cim-agro-credito-rural`, `bndes.desembolsos-linhas-agro`, `anp.etanol-precos` | [bcb.gov.br](https://www.bcb.gov.br/publicacoes/cim) · BNDES CKAN · [gov.br/anp](https://www.gov.br/anp/pt-br/assuntos/precos-e-indices/precos-de-combustiveis) |
| 56 Industry | `abiove.balanco-complexo-soja`, `abiove.exportacoes-complexo-soja`, `abiove.capacidade-instalada-esmagamento`, `b3.futuro-cafe`, `b3.futuro-acucar` | [abiove.org.br/estatisticas](https://abiove.org.br/estatisticas/) · B3 SPRD prefixes **ICF** (café), **CNL** (açúcar) |

**Rejected (Phase 51):** `mapa.fertilizantes-registro` (duplicate SIPEAGRO), `b3.futuro-algodao` (no B3 cotton contract), `cna.pesquisas-setoriais` (membership only).

**Deferred (Phase 51):** SIGORG/Agrostat/SISBI panels, SICAR shapefiles, INCRA, CONAB safra bulletins, ANTAQ dwell time, CPTEC/Embrapa solos — see `.local/SOURCE-GAP-MATRIX-WAVE5.md`.

---

## Rules

- Secondary blogs and unofficial mirrors are **not** acceptable as sole citation.
- Store resolved download URL + `discovered_at` per ingest — portal section link is mandatory.

---

## Related

- [VISION.md](VISION.md)
- [GLOSSARY.md](GLOSSARY.md)
- [ROADMAP.md](ROADMAP.md)
- [DATA-HISTORY-RANGES.md](DATA-HISTORY-RANGES.md)
- [REFRESH-POLICY.md](REFRESH-POLICY.md)
