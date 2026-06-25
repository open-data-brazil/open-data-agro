# Cloudflare R2 (production bronze)

Optional production backend when `STORAGE_MODE=r2`. Local development uses `STORAGE_MODE=local` by default.

## Provisioning

1. Create a bucket in the [Cloudflare R2 dashboard](https://developers.cloudflare.com/r2/buckets/) (default name: `open-data-agro`).
2. Create an API token with **Object Read & Write** on that bucket.
3. Copy account ID, access key, and secret into `.env`:

```text
STORAGE_MODE=r2
R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET=open-data-agro
```

`R2_ENDPOINT` is derived from `R2_ACCOUNT_ID` when unset (`https://{account_id}.r2.cloudflarestorage.com`).

## S3 client settings

| Setting | Value |
|---------|-------|
| Region | `auto` |
| Path style | `false` |
| Keys | Same as local `lake/` layout |

## CORS (optional)

Enable CORS on the bucket only if browser clients need direct object access. The ingestor and DuckDB CLI use server-side credentials.

## Terraform

Infrastructure-as-code for R2 buckets is out of scope for the MVP; use the dashboard or add Terraform under `infra/r2/terraform/` in a later phase.
