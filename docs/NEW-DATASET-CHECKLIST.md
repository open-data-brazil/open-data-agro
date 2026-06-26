# New dataset checklist

Copy this checklist when adding a **new** catalog dataset (new `dataset_id` row in `docs/OFFICIAL-SOURCES.md`).  
Reference implementation: **`conab.estimativa-graos`** (Phases 10 + collection patterns).

---

## 1. Sample fixtures

- [ ] Add official sample under `internal/<agency>/testdata/` (or agency-equivalent path)
- [ ] Regenerate reference vectors if CONAB: `make conab-reference` (or agency equivalent)

**Example:** `internal/conab/testdata/estimativa-graos/`

---

## 2. Golden ingest test

- [ ] Add `internal/ingest/*_golden_test.go` covering TXT/CSV/XLS → Parquet for the dataset slug

**Example:** `internal/ingest/conab_estimativa_graos_golden_test.go`

---

## 3. Great Expectations bronze suite

- [ ] Add `expectations/suites/bronze/<agency>/<name>.json`
- [ ] Map suite in `internal/processor/quality.go`

**Example:** `expectations/suites/bronze/conab/estimativa_graos.json`

---

## 4. dbt staging and marts

- [ ] Macro in `dbt/macros/` (if column normalization needed)
- [ ] `stg_<agency>__<name>.sql` + `mart_<agency>__<name>.sql`
- [ ] Source/mart entries in `dbt/models/staging/<agency>/_<agency>__sources.yml` and `dbt/models/marts/`

**Example:** `stg_conab__estimativa_graos.sql`, `mart_conab__estimativa_graos.sql`

---

## 5. CI silver seed

- [ ] Add or extend `scripts/ci/seed_*_silver.py` with minimal Delta rows for offline dbt

**Example:** `scripts/ci/seed_dbt_silver.py`

---

## 6. Makefile MVP gate

- [ ] `make <agency>-mvp` or extend existing MVP target to include the dataset
- [ ] Optional isolated-lake CI mirror: `make ci-<agency>-mvp`

**Example:** `make conab-mvp`, `make ci-p1-collection-mvp`

---

## 7. Official sources catalog

- [ ] Update `docs/OFFICIAL-SOURCES.md` → status **`Pn — implemented`**
- [ ] Run `python3 scripts/ci/check_official_sources_status.py`

---

## 8. Phase reference doc

- [ ] Column mapping in `.local/phases/<phase>/OFFICIAL-REFERENCE.md`

---

## 9. Changelog

- [ ] Entry under `[Unreleased]` in `CHANGELOG.md`

---

## 10. Validation commands

```bash
go test ./internal/ingest -run Golden -count=1
make <your-mvp-target>
python3 scripts/ci/check_official_sources_status.py
```

---

## Related

- [NEW-PROJECT-CHECKLIST.md](NEW-PROJECT-CHECKLIST.md) — project bootstrap
- [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md) — catalog index
- [AGENTS.md](../AGENTS.md) — agent protocol
