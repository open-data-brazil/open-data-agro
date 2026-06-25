import { DeleteObjectCommand, PutObjectCommand, S3Client } from '@aws-sdk/client-s3';
import { mkdir, readFile, unlink, writeFile } from 'node:fs/promises';
import { dirname, join } from 'node:path';

import type { IngestorConfig } from '../config.js';
import { isR2Configured } from '../config.js';

export type BronzeUploadResult = {
  key: string;
  storage: 'r2' | 'local';
};

export class BronzeStorage {
  private readonly s3: S3Client | null;
  private readonly bucket: string;
  private readonly localRoot: string;

  constructor(config: IngestorConfig) {
    this.bucket = config.r2Bucket;
    this.localRoot = config.lakeLocalRoot;

    if (isR2Configured(config)) {
      this.s3 = new S3Client({
        region: 'auto',
        endpoint: config.r2Endpoint as string,
        credentials: {
          accessKeyId: config.r2AccessKeyId as string,
          secretAccessKey: config.r2SecretAccessKey as string,
        },
      });
    } else {
      this.s3 = null;
    }
  }

  buildBronzeKey(datasetId: string, ingestDate: string, partId: string): string {
    const slug = datasetSlug(datasetId);
    return `bronze/conab/${slug}/ingest_date=${ingestDate}/part-${partId}.parquet`;
  }

  async uploadParquet(key: string, filePath: string): Promise<BronzeUploadResult> {
    const body = await readFile(filePath);

    if (this.s3) {
      await this.s3.send(
        new PutObjectCommand({
          Bucket: this.bucket,
          Key: key,
          Body: body,
          ContentType: 'application/octet-stream',
        }),
      );
      return { key, storage: 'r2' };
    }

    const localPath = join(this.localRoot, key);
    await mkdir(dirname(localPath), { recursive: true });
    await writeFile(localPath, body);
    return { key, storage: 'local' };
  }

  async deleteObject(key: string): Promise<void> {
    if (this.s3) {
      await this.s3.send(
        new DeleteObjectCommand({
          Bucket: this.bucket,
          Key: key,
        }),
      );
      return;
    }

    const localPath = join(this.localRoot, key);
    await unlink(localPath);
  }
}

export function datasetSlug(datasetId: string): string {
  const parts = datasetId.split('.');
  return parts.length > 1 ? parts.slice(1).join('.') : datasetId;
}
