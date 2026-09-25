import type { NitroFetchOptions } from 'nitropack'

export function useApi() {
  const config = useRuntimeConfig()
  const auth = useAuthStore()

  return function apiFetch<T>(request: string, options: NitroFetchOptions<string> = {}) {
    const headers = new Headers(options.headers as HeadersInit | undefined)
    if (auth.token) headers.set('Authorization', `Bearer ${auth.token}`)
    return $fetch<T>(request, {
      ...options,
      baseURL: config.public.apiBaseUrl,
      headers,
      credentials: 'include',
    })
  }
}
