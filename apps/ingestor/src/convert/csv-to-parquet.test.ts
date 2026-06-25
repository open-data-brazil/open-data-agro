import { mkdir, readFile, rm, writeFile } from 'node:fs/promises';
import { join } from 'node:path';
import { tmpdir } from 'node:os';
import { randomUUID } from 'node:crypto';
import { describe, expect, it } from 'vitest';

import { csvToParquet } from './csv-to-parquet.js';

describe('csvToParquet', () => {
  it('preserves raw column names', async () => {
    const dir = join(tmpdir(), `open-data-agro-test-${randomUUID()}`);
    await mkdir(dir, { recursive: true });

    const inputPath = join(dir, 'sample.txt');
    const outputPath = join(dir, 'sample.parquet');
    await writeFile(inputPath, 'ano_agricola;uf;produto\n2024/25;SP;SOJA\n');

    const result = await csvToParquet({
      inputPath,
      outputPath,
      delimiter: ';',
    });

    expect(result.rowCount).toBe(1);
    expect(result.columnNames).toEqual(['ano_agricola', 'uf', 'produto']);
    const parquetBytes = await readFile(outputPath);
    expect(parquetBytes.length).toBeGreaterThan(0);

    await rm(dir, { recursive: true, force: true });
  });
});
