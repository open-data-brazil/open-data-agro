export function exitWithError(reason: Error | string | number | boolean | object): never {
  let message: string;
  if (reason instanceof Error) {
    message = reason.message;
  } else if (typeof reason === 'string') {
    message = reason;
  } else if (typeof reason === 'number' || typeof reason === 'boolean') {
    message = String(reason);
  } else {
    message = JSON.stringify(reason);
  }
  console.error(message);
  process.exit(1);
}
