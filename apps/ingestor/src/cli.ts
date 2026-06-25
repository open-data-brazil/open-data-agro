import { Command } from 'commander';

const program = new Command();

program
  .name('ingestor')
  .description('Open Data Agro — CONAB and public agro data ingestion CLI')
  .version('0.0.0');

program.parse();
