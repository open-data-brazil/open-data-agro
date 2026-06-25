# ADR 003 — Delta promotion runtime

**Status:** Accepted  
**Date:** 2026-06-25  
**Context:** Phase 3 — bronze Parquet → silver Delta

## Decision

`cmd/processor promote` shells out to **`scripts/delta/promote.py`** using the pinned **`deltalake`** Python package (delta-rs 1.6.0 bindings).

DuckDB `delta` extension is used for **read/query** smoke tests and analytics (Phase 4/8), not for writes.

## Alternatives considered

| Option | Outcome |
|--------|---------|
| **delta-rs Rust CLI** (`deltalake-cli`) | No stable parquet→delta write command; rejected |
| **DuckDB `COPY … FORMAT DELTA`** | Requires extension install per machine; deferred to Phase 4 orchestration |
| **Spark** | Violates local-first / no-cluster policy |

## Consequences

- Developers need Python 3.12 + `pip install -r toolchain/python-requirements.txt` for promotion.
- Go remains the CLI entrypoint; promotion logic stays in one script for schema/metadata columns.
- `STORAGE_MODE=minio|r2` reuses the same script with S3 `storage_options`.

## References

- [delta-rs writing](https://delta-io.github.io/delta-rs/usage/writing/)
- [Phase 3 OFFICIAL-REFERENCE](../../.local/phases/03-lakehouse-delta/OFFICIAL-REFERENCE.md)
