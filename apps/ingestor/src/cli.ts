import 'dotenv/config';

import { Command } from 'commander';

import { createAlertSink } from './alerts.js';
import { listCatalogEntries } from './catalog.js';
import { loadConfig, readPackageVersion } from './config.js';
import { JobRepository } from './db/job-repository.js';
import { runIngestJob } from './jobs/run-job.js';
import { exitWithError } from './lib/errors.js';
import { BronzeStorage } from './storage/r2-client.js';

async function withRepository<T>(fn: (repo: JobRepository) => Promise<T>): Promise<T> {
  const config = loadConfig();
  const repo = new JobRepository(config.databaseUrl);
  await repo.initialize();
  try {
    return await fn(repo);
  } finally {
    await repo.close();
  }
}

const program = new Command();

program
  .name('ingestor')
  .description('Open Data Agro — CONAB and public agro data ingestion CLI')
  .version(readPackageVersion());

const catalog = program.command('catalog').description('Dataset catalog commands');

catalog
  .command('list')
  .description('List registered datasets')
  .action(async () => {
    await withRepository(async (repo) => {
      await repo.syncRegistry(listCatalogEntries());
      const entries = listCatalogEntries();
      for (const entry of entries) {
        console.log(
          `${entry.datasetId}\t${entry.conabSection}\t${entry.format}\t${entry.schedule}\t${entry.sourceUrl}`,
        );
      }
    });
  });

program
  .command('run')
  .description('Run ingestion for a dataset')
  .argument('<dataset_id>', 'Dataset identifier (e.g. conab.estimativa-graos)')
  .option('--dry-run', 'Resolve source URL without download or upload')
  .action(async (datasetId: string, options: { dryRun?: boolean }) => {
    const config = loadConfig();
    const repo = new JobRepository(config.databaseUrl);
    const storage = new BronzeStorage(config);
    const alerts = createAlertSink(config.alertWebhookUrl);

    await repo.initialize();
    try {
      const result = await runIngestJob(config, repo, storage, alerts, {
        datasetId,
        ...(options.dryRun === true ? { dryRun: true } : {}),
      });
      console.log(JSON.stringify(result, null, 2));
      if (result.status === 'failed') {
        process.exitCode = 1;
      }
    } finally {
      await repo.close();
    }
  });

program
  .command('status')
  .description('Show recent ingest jobs')
  .option('--last <count>', 'Number of jobs to show', '10')
  .action(async (options: { last: string }) => {
    const limit = Number.parseInt(options.last, 10);
    await withRepository(async (repo) => {
      const jobs = await repo.listRecentJobs(Number.isNaN(limit) ? 10 : limit);
      for (const job of jobs) {
        console.log(
          `${job.startedAt.toISOString()}\t${job.datasetId}\t${job.status}\t${job.id}${
            job.errorMessage ? `\t${job.errorMessage}` : ''
          }`,
        );
      }
    });
  });

program.parseAsync(process.argv).catch(exitWithError);
