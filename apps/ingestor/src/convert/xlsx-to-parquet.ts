import { readFile } from 'node:fs/promises';
import parquet from 'parquetjs';
import * as XLSX from 'xlsx';

export type XlsxToParquetOptions = {
  inputPath: string;
  outputPath: string;
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

export async function xlsxToParquet(options: XlsxToParquetOptions): Promise<ConvertResult> {
  const buffer = await readFile(options.inputPath);
  const workbook = XLSX.read(buffer, { type: 'buffer' });
  const sheetName = workbook.SheetNames[0];
  if (!sheetName) {
    throw new Error(`No sheets found in ${options.inputPath}`);
  }

  const sheet = workbook.Sheets[sheetName];
  if (!sheet) {
    throw new Error(`Missing sheet ${sheetName} in ${options.inputPath}`);
  }

  const records = XLSX.utils.sheet_to_json<Record<string, string | number | boolean | null>>(sheet, {
    defval: '',
    raw: false,
  });

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
        const value = record[column];
        row[column] = value === undefined || value === null ? '' : String(value);
      }
      await writer.appendRow(row);
    }
  } finally {
    await writer.close();
  }

  return { rowCount: records.length, columnNames };
}
