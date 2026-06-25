import { parse } from 'csv-parse/sync';
import { readFile } from 'node:fs/promises';
import parquet from 'parquetjs';

export type CsvToParquetOptions = {
  inputPath: string;
  outputPath: string;
  delimiter?: string;
};

export type ConvertResult = {
  rowCount: number;
  columnNames: string[];
};

function normalizeColumnName(name: string): string {
  return name.trim();
}

function buildSchema(columns: string[]): parquet.ParquetSchema {
  const fields: Record<string, { type: string; optional: boolean }> = {};
  for (const column of columns) {
    fields[column] = { type: 'UTF8', optional: true };
  }
  return new parquet.ParquetSchema(fields);
}

export async function csvToParquet(options: CsvToParquetOptions): Promise<ConvertResult> {
  const raw = await readFile(options.inputPath, 'utf8');
  const records = parse(raw, {
    columns: true,
    delimiter: options.delimiter ?? ',',
    relax_column_count: true,
    skip_empty_lines: true,
    trim: true,
  }) as Record<string, string>[];

  if (records.length === 0) {
    throw new Error(`No rows found in ${options.inputPath}`);
  }

  const firstRecord = records[0];
  if (!firstRecord) {
    throw new Error(`No rows found in ${options.inputPath}`);
  }

  const columnNames = Object.keys(firstRecord).map(normalizeColumnName);
  const schema = buildSchema(columnNames);
  const writer = await parquet.ParquetWriter.openFile(schema, options.outputPath);

  try {
    for (const record of records) {
      const row: Record<string, string> = {};
      for (const column of columnNames) {
        const cell = record[column];
        row[column] = cell === undefined ? '' : cell;
      }
      await writer.appendRow(row);
    }
  } finally {
    await writer.close();
  }

  return { rowCount: records.length, columnNames };
}
