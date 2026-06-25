export const CONAB_USER_AGENT = 'open-data-agro-ingestor/0.1.0 (+https://github.com/open-data-agro)';

export const DEFAULT_TIMEOUT_MS = 60_000;
export const DEFAULT_RETRIES = 3;

export type HttpResponseMeta = {
  contentType: string | null;
  lastModified: string | null;
  contentLength: number | null;
};

export type DownloadResult = {
  filePath: string;
  bytes: Buffer;
  sha256: string;
  meta: HttpResponseMeta;
};

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

function isRetryableStatus(status: number): boolean {
  return status === 408 || status === 429 || status >= 500;
}

export async function fetchWithRetry(
  url: string,
  options: { timeoutMs?: number; retries?: number } = {},
): Promise<Response> {
  const timeoutMs = options.timeoutMs ?? DEFAULT_TIMEOUT_MS;
  const retries = options.retries ?? DEFAULT_RETRIES;
  let lastError: Error | undefined;

  for (let attempt = 0; attempt < retries; attempt += 1) {
    const controller = new AbortController();
    const timer = setTimeout(() => {
      controller.abort();
    }, timeoutMs);

    try {
      const response = await fetch(url, {
        signal: controller.signal,
        headers: {
          'User-Agent': CONAB_USER_AGENT,
          Accept: '*/*',
        },
      });

      if (!response.ok) {
        if (isRetryableStatus(response.status) && attempt < retries - 1) {
          await sleep(2 ** attempt * 500);
          continue;
        }
        throw new Error(`HTTP ${String(response.status)} for ${url}`);
      }

      return response;
    } catch (error) {
      lastError = error instanceof Error ? error : new Error(String(error));
      if (attempt < retries - 1) {
        await sleep(2 ** attempt * 500);
        continue;
      }
    } finally {
      clearTimeout(timer);
    }
  }

  throw lastError ?? new Error(`Failed to fetch ${url}`);
}
