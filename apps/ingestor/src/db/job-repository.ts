import { readFile } from 'node:fs/promises';
import { join } from 'node:path';

import type { DatasetRegistryEntry } from '@open-data-agro/catalog';
import { loadRegistry } from '@open-data-agro/catalog';
import { Pool } from 'pg';

const SCHEMA_SQL = await readFile(
  join(import.meta.dirname, '../../../../infra/postgres/init/002_ingest_schema.sql'),
  'utf8',
);

export type JobStatus = 'running' | 'success' | 'failed' | 'skipped';

export type IngestJob = {
  id: string;
  datasetId: string;
  startedAt: Date;
  finishedAt: Date | null;
  status: JobStatus;
  errorMessage: string | null;
  dryRun: boolean;
};

export type IngestFile = {
  id: string;
  jobId: string;
  datasetId: string;
  sha256: string;
  rowCount: number | null;
  r2Key: string | null;
  contentType: string | null;
  lastModified: string | null;
  fileSizeBytes: number | null;
  discoveredAt: Date | null;
  createdAt: Date;
};

export class JobRepository {
  private readonly pool: Pool;

  constructor(databaseUrl: string) {
    this.pool = new Pool({ connectionString: databaseUrl });
  }

  async initialize(): Promise<void> {
    await this.pool.query(SCHEMA_SQL);
    await this.syncRegistry(loadRegistry());
  }

  async syncRegistry(entries: DatasetRegistryEntry[]): Promise<void> {
    for (const entry of entries) {
      await this.pool.query(
        `INSERT INTO catalog.dataset_registry (
          dataset_id, source_url, format, schedule, conab_section, portal_label, discovered_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (dataset_id) DO UPDATE SET
          source_url = EXCLUDED.source_url,
          format = EXCLUDED.format,
          schedule = EXCLUDED.schedule,
          conab_section = EXCLUDED.conab_section,
          portal_label = EXCLUDED.portal_label,
          discovered_at = EXCLUDED.discovered_at,
          updated_at = NOW()`,
        [
          entry.datasetId,
          entry.sourceUrl,
          entry.format,
          entry.schedule,
          entry.conabSection,
          entry.portalLabel,
          entry.discoveredAt,
        ],
      );
    }
  }

  async findIngestedSha256(datasetId: string, sha256: string): Promise<IngestFile | null> {
    const result = await this.pool.query<{
      id: string;
      job_id: string;
      dataset_id: string;
      sha256: string;
      row_count: number | null;
      r2_key: string | null;
      content_type: string | null;
      last_modified: string | null;
      file_size_bytes: string | null;
      discovered_at: Date | null;
      created_at: Date;
    }>(
      `SELECT id, job_id, dataset_id, sha256, row_count, r2_key, content_type,
              last_modified, file_size_bytes, discovered_at, created_at
       FROM catalog.ingest_files
       WHERE dataset_id = $1 AND sha256 = $2`,
      [datasetId, sha256],
    );

    const row = result.rows[0];
    if (!row) {
      return null;
    }

    return mapIngestFile(row);
  }

  async createJob(datasetId: string, dryRun: boolean): Promise<string> {
    const result = await this.pool.query<{ id: string }>(
      `INSERT INTO catalog.ingest_jobs (dataset_id, status, dry_run)
       VALUES ($1, 'running', $2)
       RETURNING id`,
      [datasetId, dryRun],
    );
    const row = result.rows[0];
    if (!row) {
      throw new Error('Failed to create ingest job');
    }
    return row.id;
  }

  async finishJob(
    jobId: string,
    status: JobStatus,
    errorMessage: string | null,
  ): Promise<void> {
    await this.pool.query(
      `UPDATE catalog.ingest_jobs
       SET status = $2, finished_at = NOW(), error_message = $3
       WHERE id = $1`,
      [jobId, status, errorMessage],
    );
  }

  async insertIngestFile(input: {
    jobId: string;
    datasetId: string;
    sha256: string;
    rowCount: number;
    r2Key: string;
    contentType: string | null;
    lastModified: string | null;
    fileSizeBytes: number;
    discoveredAt: string;
  }): Promise<void> {
    await this.pool.query(
      `INSERT INTO catalog.ingest_files (
        job_id, dataset_id, sha256, row_count, r2_key, content_type,
        last_modified, file_size_bytes, discovered_at
      ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
      [
        input.jobId,
        input.datasetId,
        input.sha256,
        input.rowCount,
        input.r2Key,
        input.contentType,
        input.lastModified,
        input.fileSizeBytes,
        input.discoveredAt,
      ],
    );
  }

  async listRecentJobs(limit: number): Promise<IngestJob[]> {
    const result = await this.pool.query<{
      id: string;
      dataset_id: string;
      started_at: Date;
      finished_at: Date | null;
      status: JobStatus;
      error_message: string | null;
      dry_run: boolean;
    }>(
      `SELECT id, dataset_id, started_at, finished_at, status, error_message, dry_run
       FROM catalog.ingest_jobs
       ORDER BY started_at DESC
       LIMIT $1`,
      [limit],
    );

    return result.rows.map((row) => ({
      id: row.id,
      datasetId: row.dataset_id,
      startedAt: row.started_at,
      finishedAt: row.finished_at,
      status: row.status,
      errorMessage: row.error_message,
      dryRun: row.dry_run,
    }));
  }

  async close(): Promise<void> {
    await this.pool.end();
  }
}

function mapIngestFile(row: {
  id: string;
  job_id: string;
  dataset_id: string;
  sha256: string;
  row_count: number | null;
  r2_key: string | null;
  content_type: string | null;
  last_modified: string | null;
  file_size_bytes: string | null;
  discovered_at: Date | null;
  created_at: Date;
}): IngestFile {
  return {
    id: row.id,
    jobId: row.job_id,
    datasetId: row.dataset_id,
    sha256: row.sha256,
    rowCount: row.row_count,
    r2Key: row.r2_key,
    contentType: row.content_type,
    lastModified: row.last_modified,
    fileSizeBytes: row.file_size_bytes ? Number(row.file_size_bytes) : null,
    discoveredAt: row.discovered_at,
    createdAt: row.created_at,
  };
}
