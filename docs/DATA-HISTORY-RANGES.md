# Data history ranges — backfill reference

> Per-dataset **source history** and **catalog default** for the Open Data Agro ingestor.  
> Phase 33 (collection hardening). **Not IA** — documents how to load all available years.

**Related:** [REFRESH-POLICY.md](REFRESH-POLICY.md) · [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md)

---

## How to read this table

| Column | Meaning |
|--------|---------|
| **Source min** | Earliest year/date available from the official primary source |
| **Catalog default** | `period_start` / `start_date` in `configs/catalog/` when no `--from` flag |
| **Backfill command** | Example full-history ingest (live network) |
| **MVP / CI** | Offline pipeline uses fixtures + seeds — not full history |

CONAB portal files are **full snapshots** (no `--from`); each download replaces the complete published series.

---

## CONAB (Phases 10–14)

| `dataset_id` | Source min | Catalog default | Backfill |
|--------------|------------|-----------------|----------|
| `conab.estimativa-graos` | ~1990 (safra) | full snapshot | `go run ./cmd/ingestor run conab.estimativa-graos` |
| `conab.serie-historica-graos` | ~1976 | full snapshot | same pattern |
| `conab.estimativa-cana` | ~1990 | full snapshot | |
| `conab.serie-historica-cana` | ~1976 | full snapshot | |
| `conab.estimativa-cafe` | ~1990 | full snapshot | |
| `conab.serie-historica-cafe` | ~1976 | full snapshot | |
| `conab.custo-producao` | ~2006 | full snapshot | |
| `conab.oferta-demanda` | ~2005 | full snapshot | |
| `conab.precos-*` (mercado) | ~2004–2010 | full snapshot | See Phase 11 OFFICIAL-REFERENCE |
| `conab.prohort-*` | ~2018 | full snapshot | |
| `conab.estoques-publicos` | ~2008 | full snapshot | |
| `conab.operacoes-comercializacao` | ~2008 | full snapshot | |
| `conab.vendas-balcao` | ~2010 | full snapshot | |
| `conab.frete` | ~2005 | full snapshot | |
| `conab.armazenagem` | ~2008 | full snapshot | |
| `conab.serie-historica-capacidade-estatica` | ~2000 | full snapshot | |
| `conab.alimenta-brasil-*` | ~2020 | full snapshot | |

Portal audit: all listed sections mapped — see `.local/phases/33-collection-hardening/CONAB-PORTAL-AUDIT.md`.

---

## ANP (Phase 12)

| `dataset_id` | Source min | Catalog default | Backfill |
|--------------|------------|-----------------|----------|
| `anp.combustiveis-precos-medios-municipios` | 2004 | full snapshot | `go run ./cmd/ingestor run anp.combustiveis-precos-medios-municipios` |
| `anp.combustiveis-precos-postos` | 2004 | full snapshot | |

---

## IBGE (Phases 15–16)

| `dataset_id` | Source min | Catalog default | Backfill |
|--------------|------------|-----------------|----------|
| `ibge.localidades-*` | current | full API pull | `make ibge-localidades-mvp` |
| `ibge.pam-area-quantidade` | **1974** | 2010 (MVP window) | `go run ./cmd/ingestor run ibge.pam-area-quantidade --from 1974 --to 2024` |
| `ibge.pam-rendimento-valor` | **1974** | 2010 | `--from 1974` |
| `ibge.pam-estabelecimentos` | **2007** | 2010 | `--from 2007` |

Use `--uf` and `--crop` to chunk large SIDRA pulls.

---

## INMET (Phase 17)

| `dataset_id` | Source min | Catalog default | Backfill |
|--------------|------------|-----------------|----------|
| `inmet.bdmep-diario` | **1961** | 2000 | `go run ./cmd/ingestor run inmet.bdmep-diario --year 1961` … per year |
| `inmet.bdmep-mensal` | 1961 (derived) | 2000 | `--year` + `--uf` filters |
| `inmet.pacote-anual-automaticas` | 2000 | 2000 | Annual ZIP per year |
| `inmet.estacoes-*` | current | catalog | Station catalogs |

National full backfill can exceed 100 GB — prefer `--uf` + `--year` increments (see `configs/catalog/inmet/clima.yaml`).

---

## BCB SGS (Phase 18)

| `dataset_id` | SGS | Source min | Catalog `start_date` | Backfill |
|--------------|-----|------------|----------------------|----------|
| `bcb.sgs-ipca` | 433 | 1980 | 01/01/1995 | `--from 1995-01-01` or earlier if API returns |
| `bcb.sgs-ipca-12m` | 13522 | 1995 | 01/01/1995 | |
| `bcb.sgs-ptax-usd-venda` | 1 | **1984** | 01/01/1984 | `--from 1984-01-01` |
| `bcb.sgs-ptax-usd-compra` | 10813 | **1984** | 01/01/1984 | `--from 1984-01-01` |
| `bcb.sgs-igpm` | 189 | 1989 | 01/01/1995 | |

API paginates ≤10 years per request (`internal/bcb/sgs.go`).

---

## CEPEA (Phase 19)

| `dataset_id` | Source min | Catalog default | Backfill |
|--------------|------------|-----------------|----------|
| `cepea.soja-paranagua` | **2004** (daily) | 2010-01-01 | `--from 2010-01-01` or `2004-01-01` for full |
| `cepea.soja-parana` | 2004 | 2010-01-01 | |
| `cepea.milho` | 2004 | 2010-01-01 | |
| `cepea.boi-gordo` | 2004 | 2010-01-01 | |

Historical export via CEPEA “Consulta ao Banco de Dados”; live window uses portal or Notícias Agrícolas mirror.

---

## MDIC Comex (Phase 21 + 35)

| `dataset_id` | Source min | Catalog default | Backfill |
|--------------|------------|-----------------|----------|
| `mdic.comex-exportacao-ncm-mes` | **1997** (API) | 2015 (ag NCM spot-check) | `--from 2015-01-01` or `1997-01-01` |
| `mdic.comex-importacao-ncm-mes` | **1997** (API) | 2015 (fertilizer NCMs) | `make mdic-comex-extended-mvp` |
| `mdic.comex-exportacao-uf-ncm` | **1997** (API) | 2015 (ag NCM × UF) | `make mdic-comex-extended-mvp` |
| `mdic.comex-importacao-diesel-ncm` | **1997** (API) | 2015 (diesel NCMs) | `make mdic-comex-extended-mvp` |

API chunks one calendar year per request. Ag NCM series may start later than 1997 depending on code.

---

## ANTT, MAPA, B3 (Phases 22–24, 34)

| `dataset_id` | Source min | Catalog default | Backfill |
|--------------|------------|-----------------|----------|
| `antt.pracas-pedagio` | ~2010 | full CKAN CSV | `make dnit-antt-logistica-mvp` |
| `antt.volume-trafego-pedagio` | 2010 | latest mensal consolidado CSV | `make br-logistica-extended-mvp` |
| `antt.receita-por-praca` | 2010 | latest annual CSV (`Receita por Praça - YYYY`) | `make br-logistica-extended-mvp` |
| `mapa.zarc-tabua-risco` | ~2010 (safra) | latest CKAN CSV | `make mapa-dados-mvp` |
| `b3.futuro-soja` | ~2000 (SOY) | 2024 (MVP window) | `--from 2024-01-01`; extend for history |

---

## International (Phases 25–28)

| `dataset_id` | Source min | Catalog default | Backfill |
|--------------|------------|-----------------|----------|
| `usda.psd-*` | ~1960 | 2010 | SOAP per marketing year |
| `fao.prices-agro` | ~1991 | 2010 | FAOSTAT bulk zip |
| `worldbank.pink-sheet-monthly` | ~1960 | 2010 | CMO monthly xlsx |
| `noaa.enso-indices` | 1950 | 2010 | ONI ascii |
| `noaa.global-temp-anomaly` | 1880 | 2010 | NCEI CSV |

---

## Spot-check policy (Phase 33 acceptance)

For each major agency, verify **3 representative years** after a live backfill:

| Agency | Years to spot-check | Command |
|--------|---------------------|---------|
| CONAB | 2010, 2018, 2024 | Row count vs portal file size order-of-magnitude |
| IBGE PAM | 2010, 2015, 2020 | `--from` year in bronze partition |
| INMET | 2010, 2018, 2024 | `--year` + `--uf MT` sample |
| BCB | 1995, 2010, 2024 | Min/max `data` in mart |

CI uses offline seeds — spot-checks are **manual** after production backfill.

---

## Related

- [REFRESH-POLICY.md](REFRESH-POLICY.md)
- [NEW-DATASET-CHECKLIST.md](NEW-DATASET-CHECKLIST.md)
- `.local/phases/33-collection-hardening/`
