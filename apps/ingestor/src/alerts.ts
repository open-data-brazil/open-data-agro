export type AlertEvent = {
  type: 'ingest_failed' | 'partial_upload_rollback';
  datasetId: string;
  jobId: string;
  message: string;
};

export type AlertSink = {
  emit: (event: AlertEvent) => void;
};

export function createAlertSink(webhookUrl: string | null): AlertSink {
  return {
    emit(event: AlertEvent): void {
      console.error(`[ALERT] ${event.type} dataset=${event.datasetId} job=${event.jobId}: ${event.message}`);
      if (webhookUrl) {
        console.error(`[ALERT] webhook stub — would POST to ${webhookUrl}`);
      }
    },
  };
}
