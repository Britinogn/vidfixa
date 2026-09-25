import { useMutation, useQuery } from '@tanstack/vue-query'
import type { Download } from '~/types/api'
import { rememberDownload } from '~/utils/download-history'

export function useDownloads() {
  const api = useApi()
  const auth = useAuthStore()
  const ids = useState<string[]>('download-history', () => [])

  const createDownload = useMutation({
    mutationFn: (url: string) => api<Download>('/downloads/', { method: 'POST', body: { url } }),
    onSuccess: (download) => {
      rememberDownload(download.id)
      ids.value = [download.id, ...ids.value.filter((id) => id !== download.id)].slice(0, 20)
    },
  })

  function useDownload(id: string) {
    return useQuery({
      queryKey: ['download', id],
      queryFn: () => api<Download>(`/downloads/${id}`),
      enabled: computed(() => Boolean(id)),
      refetchInterval: (query) => {
        const status = query.state.data?.status
        return status === 'completed' || status === 'failed' || status === 'cancelled' ? false : 2_000
      },
    })
  }

  async function downloadFile(id: string) {
    const config = useRuntimeConfig()
    const response = await fetch(`${config.public.apiBaseUrl}/downloads/${id}/file`, {
      headers: auth.token ? { Authorization: `Bearer ${auth.token}` } : undefined,
      credentials: 'include',
    })
    if (!response.ok) throw new Error('This video is not ready to download yet.')
    const blob = await response.blob()
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `vidfixa-${id}.mp4`
    link.click()
    window.setTimeout(() => URL.revokeObjectURL(url), 1_000)
  }

  return { ids, createDownload, useDownload, downloadFile }
}
