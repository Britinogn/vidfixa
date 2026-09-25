const storageKey = 'vidfixa.downloads.v1'

export function readDownloadIDs(): string[] {
  if (!import.meta.client) return []
  try {
    const value: unknown = JSON.parse(localStorage.getItem(storageKey) || '[]')
    return Array.isArray(value) ? value.filter((item): item is string => typeof item === 'string').slice(0, 20) : []
  } catch {
    return []
  }
}

export function rememberDownload(id: string): void {
  if (!import.meta.client) return
  const next = [id, ...readDownloadIDs().filter((savedID) => savedID !== id)].slice(0, 20)
  localStorage.setItem(storageKey, JSON.stringify(next))
}

export function forgetDownload(id: string): void {
  if (!import.meta.client) return
  localStorage.setItem(storageKey, JSON.stringify(readDownloadIDs().filter((savedID) => savedID !== id)))
}
