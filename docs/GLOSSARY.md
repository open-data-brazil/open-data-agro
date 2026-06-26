# Domain glossary — Open Data Agro

> Ubiquitous language for this project. Code, APIs, docs, and agents MUST use these terms exactly.

---

## Open Data Agro

**Definition:** The open-source toolkit for embedding and querying Brazilian agricultural public data.
**Not the same as:** A government portal, SaaS analytics product, or farm management system.
**Code name:** `open-data-agro`

---

## Dataset

**Definition:** A versioned collection of records from one official source, embedded offline with metadata.
**Not the same as:** Live API client or scraped HTML page.
**Code name:** `Dataset`, `datasetId`

---

## Dataset metadata

**Definition:** Capture timestamp, source URL, endpoints used, row count, and license notes.
**Invariant:** Every embed ships with `metadata.json` conforming to project schema.
**Code name:** `DatasetMetadata`

---

## Fetch script

**Definition:** Maintainer-run script that downloads from official APIs/portals and writes embedded JSON/Parquet.
**Not the same as:** Runtime network call in library core.
**Code name:** `scripts/fetch-*.ts`

---

## Official source

**Definition:** Primary `.gov.br` document, API, or open data catalog entry cited for a field or code list.
**Not the same as:** Wikipedia, blog post, or unofficial GitHub mirror.
**Code name:** Referenced in [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md)

---

## IBGE municipality code (`cMun` / `codigo_ibge`)

**Definition:** Seven-digit IBGE code identifying a Brazilian municipality.
**Format:** `NNNNNNN` (e.g. `3550308` = São Paulo/SP).
**Code name:** `IbgeMunicipioCode`, `cMun`

---

## CONAB

**Definition:** Companhia Nacional de Abastecimento — national supply company publishing harvest estimates and stocks.
**Code name:** `conab` in dataset IDs and module paths.

---

## MAPA

**Definition:** Ministério da Agricultura e Pecuária — federal ministry; publishes policies, registers, and open data.
**Code name:** `mapa`

---

## INMET

**Definition:** Instituto Nacional de Meteorologia — weather and climate open data relevant to agriculture.
**Code name:** `inmet`

---

## CAR (Cadastro Ambiental Rural)

**Definition:** Rural environmental registry polygon identifiers — used for land-use context, not ownership proof in this project.
**Code name:** `car`, `CarCode`

---

## Crop code

**Definition:** Official code for a crop or product in a cited catalog (e.g. CONAB, IBGE PAM, MAPA classifications).
**Invariant:** Must map to a row in [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md).
**Code name:** `CropCode`

### PAM ↔ CONAB crop mapping (MVP)

| Crop (name) | IBGE PAM SIDRA code (`c81`/`c82`/`c782`) | CONAB grãos label |
|-------------|---------------------------------------------|-------------------|
| Soja | `2713` | SOJA |
| Milho | `2711` | MILHO |
| Trigo | `2716` | TRIGO |

### INMET missing values and timezone

| Policy | Value |
|--------|-------|
| Source timezone | UTC (store as-is in bronze; document BRT offset in analytics) |
| Missing sentinels | `9999`, `-9999`, `Null`, blank, `//` — normalized to empty in long-format `valor` |

---

## USDA marketing year

**Definition:** USDA FAS PSD reporting year for a commodity, which may span two calendar years (e.g. U.S. soy MY 2024/25). Stored as the four-digit **market year start** in bronze (`marketing_year` column).
**Not the same as:** Calendar year (`calendar_year`) or Brazilian safra label.
**Alignment:** Use `marketing_year` for cross-country PSD joins; use `calendar_year` + `month` for release timing. See [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md) USDA section.
**Code name:** `marketing_year`, `Market_Year`

---

## Golden vector

**Definition:** Test fixture with input/output pairs from official API samples or published tables.
**Code name:** `tests/vectors/*.official.json`

---

## Refresh bot

**Definition:** Scheduled job that runs fetch scripts, compares drift, and opens reports or PRs.
**Code name:** `data-refresh-bot` (planned)

---

## Normalize

**Definition:** Transform raw official records into canonical project types (codes, dates, units, CRS).
**Not the same as:** Business analytics or aggregation for dashboards.
**Code name:** `normalize*`, `toCanonical*`
