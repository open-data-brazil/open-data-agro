import type { DatasetRegistryEntry } from '@open-data-agro/catalog';

import { csvToParquet } from './csv-to-parquet.js';
import { xlsxToParquet } from './xlsx-to-parquet.js';

export type ConvertSourceOptions = {
  entry: DatasetRegistryEntry;
  inputPath: string;
  outputPath: string;
};

export type ConvertSourceResult = {
  rowCount: number;
  columnNames: string[];
};

export async function convertSourceToParquet(
  options: ConvertSourceOptions,
): Promise<ConvertSourceResult> {
  const { entry, inputPath, outputPath } = options;

  if (entry.format === 'xlsx') {
    return xlsxToParquet({ inputPath, outputPath });
  }

  return csvToParquet({
    inputPath,
    outputPath,
    delimiter: entry.delimiter ?? (entry.format === 'txt' ? ';' : ','),
  });
}
