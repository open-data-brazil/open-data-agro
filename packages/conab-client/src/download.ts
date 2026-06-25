import { createHash } from 'node:crypto';
import { mkdir, writeFile } from 'node:fs/promises';
import { join } from 'node:path';

import { fetchWithRetry } from './http-client.js';
import { resolveDownloadUrl } from './resolve-url.js';

export type DownloadToTempOptions = {
  datasetId: string;
  tempDir: string;
  timeoutMs?: number;
};

export type DownloadToTempResult = {
  filePath: string;
  bytes: Buffer;
  sha256: string;
  sourceUrl: string;
  contentType: string | null;
  lastModified: string | null;
  contentLength: number;
};

export async function downloadToTemp(
  options: DownloadToTempOptions,
): Promise<DownloadToTempResult> {
  const resolved = resolveDownloadUrl(options.datasetId);
  const response = await fetchWithRetry(
    resolved.url,
    options.timeoutMs === undefined ? {} : { timeoutMs: options.timeoutMs },
  );
  const bytes = Buffer.from(await response.arrayBuffer());
  const sha256 = createHash('sha256').update(bytes).digest('hex');

  await mkdir(options.tempDir, { recursive: true });
  const extension = resolved.url.split('.').pop() ?? 'bin';
  const filePath = join(options.tempDir, `raw.${extension}`);
  await writeFile(filePath, bytes);

  return {
    filePath,
    bytes,
    sha256,
    sourceUrl: resolved.url,
    contentType: response.headers.get('content-type'),
    lastModified: response.headers.get('last-modified'),
    contentLength: bytes.length,
  };
}
