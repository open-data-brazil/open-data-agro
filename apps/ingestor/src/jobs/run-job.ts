import { mkdtemp, readFile, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { randomUUID } from 'node:crypto';

import { requireDataset } from '@open-data-agro/catalog';
import { downloadToTemp, resolveDownloadUrl } from '@open-data-agro/conab-client';

import type { AlertSink } from '../alerts.js';
import { convertSourceToParquet } from '../convert/index.js';
import type { IngestorConfig } from '../config.js';
import { JobRepository } from '../db/job-repository.js';
import { BronzeStorage } from '../storage/r2-client.js';

export type RunJobOptions = {
  datasetId: string;
  dryRun?: boolean;
};

export type RunJobResult = {
  jobId: string;
  status: 'success' | 'failed' | 'skipped';
  sha256?: string;
  r2Key?: string;
  rowCount?: number;
  message: string;
};

export async function runIngestJob(
  config: IngestorConfig,
  repo: JobRepository,
  storage: BronzeStorage,
  alerts: AlertSink,
  options: RunJobOptions,
): Promise<RunJobResult> {
  const entry = requireDataset(options.datasetId);
  const dryRun = options.dryRun ?? false;
  const jobId = await repo.createJob(entry.datasetId, dryRun);
  let tempDir: string | null = null;
  let uploadedKey: string | null = null;

  try {
    const resolved = resolveDownloadUrl(entry.datasetId);

    if (dryRun) {
      await repo.finishJob(jobId, 'success', null);
      return {
        jobId,
        status: 'success',
        message: `Dry run — resolved ${resolved.url}`,
      };
    }

    tempDir = await mkdtemp(join(tmpdir(), 'open-data-agro-ingest-'));
    const download = await downloadToTemp({
      datasetId: entry.datasetId,
      tempDir,
    });

    const existing = await repo.findIngestedSha256(entry.datasetId, download.sha256);
    if (existing) {
      await repo.finishJob(jobId, 'skipped', null);
      const skipped: RunJobResult = {
        jobId,
        status: 'skipped',
        sha256: download.sha256,
        message: `Skipped — sha256 already ingested (${existing.r2Key ?? 'no key'})`,
      };
      if (existing.r2Key) {
        skipped.r2Key = existing.r2Key;
      }
      return skipped;
    }

    const parquetPath = join(tempDir, 'bronze.parquet');
    const converted = await convertSourceToParquet({
      entry,
      inputPath: download.filePath,
      outputPath: parquetPath,
    });

    const ingestDate = new Date().toISOString().slice(0, 10);
    const partId = randomUUID();
    const r2Key = storage.buildBronzeKey(entry.datasetId, ingestDate, partId);

    const upload = await storage.uploadParquet(r2Key, parquetPath);
    uploadedKey = upload.key;

    try {
      await repo.insertIngestFile({
        jobId,
        datasetId: entry.datasetId,
        sha256: download.sha256,
        rowCount: converted.rowCount,
        r2Key: upload.key,
        contentType: download.contentType,
        lastModified: download.lastModified,
        fileSizeBytes: download.contentLength,
        discoveredAt: entry.discoveredAt,
      });
    } catch (error) {
      await storage.deleteObject(upload.key);
      uploadedKey = null;
      throw error;
    }

    await repo.finishJob(jobId, 'success', null);

    return {
      jobId,
      status: 'success',
      sha256: download.sha256,
      r2Key: upload.key,
      rowCount: converted.rowCount,
      message: `Ingested ${String(converted.rowCount)} rows to ${upload.storage}:${upload.key}`,
    };
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);

    if (uploadedKey) {
      try {
        await storage.deleteObject(uploadedKey);
        alerts.emit({
          type: 'partial_upload_rollback',
          datasetId: entry.datasetId,
          jobId,
          message: `Rolled back orphan object ${uploadedKey}`,
        });
      } catch (rollbackError) {
        const rollbackMessage =
          rollbackError instanceof Error ? rollbackError.message : String(rollbackError);
        alerts.emit({
          type: 'partial_upload_rollback',
          datasetId: entry.datasetId,
          jobId,
          message: `Failed to rollback ${uploadedKey}: ${rollbackMessage}`,
        });
      }
    }

    await repo.finishJob(jobId, 'failed', message);
    alerts.emit({
      type: 'ingest_failed',
      datasetId: entry.datasetId,
      jobId,
      message,
    });

    return {
      jobId,
      status: 'failed',
      message,
    };
  } finally {
    if (tempDir) {
      await rm(tempDir, { recursive: true, force: true });
    }
  }
}

export async function downloadConabRaw(datasetId: string, tempDir: string): Promise<Buffer> {
  const download = await downloadToTemp({ datasetId, tempDir });
  return readFile(download.filePath);
}
