import { useInfiniteQuery, useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import type { Download, DownloadHistoryItem, FileTicketResponse } from '~/types/api'
import { getApiErrorMessage } from '~/utils/api-error'
import { forgetDownload, readDownloadIDs, rememberDownload } from '~/utils/download-history'

const TERMINAL_STATUSES = ['completed', 'failed', 'cancelled']

export function useDownloads() {
  const api = useApi()
  const auth = useAuthStore()
  const queryClient = useQueryClient()
  const ids = useState<string[]>('download-history', () => [])

  // Every cache key is scoped to whoever the data belongs to. Without this a
  // second account on the same tab renders the first one's cached dashboard,
  // subscription and download statuses before any request goes out.
  const scope = computed(() => (auth.user ? `user:${auth.user.id}` : 'anonymous'))

  const createDownload = useMutation({
    mutationFn: (url: string) => api<Download>('/downloads/', { method: 'POST', body: { url } }),
    onSuccess: (download) => {
      rememberDownload(download.id)
      ids.value = [download.id, ...ids.value.filter((id) => id !== download.id)].slice(0, 20)
      // Logged-in views render from the server list, so refetch it at once.
      // Otherwise the new card — and its status pill — stays invisible until
      // some unrelated refresh happens.
      queryClient.invalidateQueries({ queryKey: ['downloads'] })
    },
  })

  function useDownload(id: string) {
    return useQuery({
      queryKey: computed(() => ['download', scope.value, id]),
      queryFn: () => api<Download>(`/downloads/${id}`),
      // Wait for the stored token to load. A signed-in user's first poll would
      // otherwise go out anonymously and be rejected as someone else's file.
      enabled: computed(() => auth.ready && Boolean(id)),
      refetchInterval: (query) => {
        // A rejected poll has no status to read, so treat it as final. Left
        // out, an expired token or a 500 would retry every 2s indefinitely.
        if (query.state.error) return false
        const status = query.state.data?.status
        return status && TERMINAL_STATUSES.includes(status) ? false : 2_000
      },
    })
  }

  function useDownloadHistory(limit = 20) {
    return useQuery({
      queryKey: computed(() => ['downloads', 'history', scope.value, limit]),
      queryFn: () => api<DownloadHistoryItem[]>('/downloads/', { params: { limit } }),
      // Users-only endpoint: anonymous histories stay device-local, since an
      // anon cookie is not an account whose history could follow it.
      enabled: computed(() => auth.ready && auth.isAuthenticated),
      refetchInterval: (query) => {
        if (query.state.error) return false
        const items = query.state.data
        if (!items) return 5_000
        return items.some((d) => !TERMINAL_STATUSES.includes(d.status)) ? 5_000 : false
      },
    })
  }

    type HistoryCursor = { before: string; before_id: string } | null

  // Paged variant of useDownloadHistory for the full history view. Each
  // page is fetched with the last item of the previous page as the cursor,
  // which is what exercises the server's keyset pagination.
  function useDownloadHistoryPages(limit = 10) {
    return useInfiniteQuery({
      queryKey: computed(() => ['downloads', 'history-pages', scope.value, limit]),
      queryFn: ({ pageParam }: { pageParam: HistoryCursor }) =>
        api<DownloadHistoryItem[]>('/downloads/', {
          params: {
            limit,
            ...(pageParam ? { before: pageParam.before, before_id: pageParam.before_id } : {}),
          },
        }),
      initialPageParam: null as HistoryCursor,
      getNextPageParam: (lastPage) => {
        if (lastPage.length < limit) return undefined
        const last = lastPage[lastPage.length - 1]
        if (!last) return undefined
        return { before: last.created_at, before_id: last.id }
      },
      enabled: computed(() => auth.ready && auth.isAuthenticated),
      refetchInterval: (query) => {
        if (query.state.error) return false
        const items = query.state.data?.pages.flat()
        if (!items) return 5_000
        return items.some((d) => !TERMINAL_STATUSES.includes(d.status)) ? 5_000 : false
      },
    })
  }

  function useDownloadFile(id: string) {    const config = useRuntimeConfig()
    const previewUrl = ref('')
    const previewError = ref('')
    const saveError = ref('')
    const opening = ref(false)
    const saving = ref(false)

    async function fileUrl(inline: boolean) {
      const { path } = await api<FileTicketResponse>(`/downloads/${id}/file-ticket${inline ? '?inline=true' : ''}`)
      return `${config.public.apiBaseUrl}${path}`
    }

    async function openPreview() {
      if (opening.value || previewUrl.value) return
      opening.value = true
      previewError.value = ''
      try {
        previewUrl.value = await fileUrl(true)
      } catch (error) {
        previewError.value = getApiErrorMessage(error, 'Could not load the preview. Try saving the file instead.')
      } finally {
        opening.value = false
      }
    }

    function closePreview() {
      previewUrl.value = ''
    }

    async function saveFile() {
      if (saving.value) return
      saving.value = true
      saveError.value = ''
      try {
        // Hand the file to the browser's own download manager rather than
        // fetching it into a blob. The blob path held the entire video in
        // memory, which is what made large files fail on phones.
        window.location.assign(await fileUrl(false))
      } catch (error) {
        saveError.value = getApiErrorMessage(error, 'Could not start the download. Please try again.')
      } finally {
        saving.value = false
      }
    }

    return { previewUrl, previewError, saveError, opening, saving, openPreview, closePreview, saveFile }
  }

  return { ids, createDownload, useDownload, useDownloadHistory, useDownloadHistoryPages, useDownloadFile, readDownloadIDs, forgetDownload }
}
