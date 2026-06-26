# Quality scripts — post-dbt cross-mart checks

Complements **Great Expectations** bronze checkpoints (`expectations/suites/bronze/`). See `.local/phases/06-quality-great-expectations/README.md` for the full GE vs post-dbt split.

## `validate_codigo_ibge.py`

Validates IBGE municipality codes in **gold** Parquet marts against `mart_ibge__localidades_municipios`.

| Makefile target | Purpose |
|-----------------|---------|
| `make validate-codigo-ibge` | Run against `LAKE_LOCAL_ROOT` (default `./lake`) |
| `make validate-codigo-ibge-lake` | Explicit `./lake` after live localidades ingest |
| `make ci-validate-codigo-ibge` | CI gate: seed + dbt + validate on `/tmp/cod-ibge-ci-lake` |

Marts checked (when present under `lake/gold/`):

- CONAB: municipal prices, frete, armazenagem, estoques, prohort, custo produção, PAA propostas
- IBGE PAM: area-quantidade, rendimento-valor, estabelecimentos

Sentinel codes (`9999999`) are skipped. Invalid codes fail with sample rows per column.

```bash
python3 scripts/quality/validate_codigo_ibge.py --lake-root ./lake
python3 scripts/quality/validate_codigo_ibge.py --lake-root ./lake --json
```

**Not covered here (GE bronze):** column existence, nullability, row count — those run at `processor promote` via GE checkpoint.
