import { FetchError } from 'ofetch'

export function getApiErrorMessage(error: unknown, fallback: string) {
  if (!(error instanceof FetchError)) {
    return fallback
  }

  // No HTTP response means the request could not reach the backend.
  if (!error.response) {
    return 'Unable to connect to VidFixa. Please check your connection and try again.'
  }

  // The backend currently returns plain-text error messages.
  const data = error.response._data

  if (typeof data === 'string' && data.trim()) {
    return data
  }

  return fallback
}