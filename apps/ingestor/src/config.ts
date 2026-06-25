import { readFileSync } from 'node:fs';
import { join } from 'node:path';

export type IngestorConfig = {
  databaseUrl: string;
  r2AccountId: string | null;
  r2AccessKeyId: string | null;
  r2SecretAccessKey: string | null;
  r2Bucket: string;
  r2Endpoint: string | null;
  lakeLocalRoot: string;
  alertWebhookUrl: string | null;
};

function requireEnv(name: string): string {
  const value = process.env[name];
  if (!value) {
    throw new Error(`Missing required environment variable: ${name}`);
  }
  return value;
}

function optionalEnv(name: string): string | null {
  const value = process.env[name];
  return value && value.length > 0 ? value : null;
}

export function loadConfig(): IngestorConfig {
  const databaseUrl = requireEnv('DATABASE_URL');
  const r2AccountId = optionalEnv('R2_ACCOUNT_ID');
  const r2AccessKeyId = optionalEnv('R2_ACCESS_KEY_ID');
  const r2SecretAccessKey = optionalEnv('R2_SECRET_ACCESS_KEY');
  const r2Bucket = process.env.R2_BUCKET ?? 'open-data-agro';
  const endpointTemplate = optionalEnv('R2_ENDPOINT');
  const r2Endpoint =
    endpointTemplate && r2AccountId
      ? endpointTemplate.replace('{account_id}', r2AccountId)
      : r2AccountId
        ? `https://${r2AccountId}.r2.cloudflarestorage.com`
        : null;

  return {
    databaseUrl,
    r2AccountId,
    r2AccessKeyId,
    r2SecretAccessKey,
    r2Bucket,
    r2Endpoint,
    lakeLocalRoot: process.env.LAKE_LOCAL_ROOT ?? './lake',
    alertWebhookUrl: optionalEnv('ALERT_WEBHOOK_URL'),
  };
}

export function readPackageVersion(): string {
  const packageJsonPath = join(import.meta.dirname, '../package.json');
  const raw = readFileSync(packageJsonPath, 'utf8');
  const parsed = JSON.parse(raw) as { version?: string };
  return parsed.version ?? '0.0.0';
}

export function isR2Configured(config: IngestorConfig): boolean {
  return Boolean(config.r2AccessKeyId && config.r2SecretAccessKey && config.r2Endpoint);
}
