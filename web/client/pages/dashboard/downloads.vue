<script setup lang="ts">
import { readDownloadIDs } from '~/utils/download-history'

definePageMeta({ layout: 'dashboard' })

const auth = useAuthStore()
const { ids, useDownloadHistoryPages } = useDownloads()

// Page size 10: small enough to prove cursor pagination works.
const history = useDownloadHistoryPages(10)

onMounted(() => {
  if (ids.value.length === 0) ids.value = readDownloadIDs()
})

// Logged-in → server history (survives logout, cross-device).
// Anonymous → device-local list.
const serverIDs = computed(() =>
  (history.data.value?.pages ?? []).flat().map((d) => d.id),
)
const displayIDs = computed(() =>
  auth.isAuthenticated ? serverIDs.value : ids.value,
)
const loadingHistory = computed(
  () => auth.isAuthenticated && history.isPending.value,
)
const historyFailed = computed(
  () => auth.isAuthenticated && history.isError.value,
)
const countLabel = computed(() => {
  const n = displayIDs.value.length
  if (n === 0) return null
  return n === 1 ? '1 download' : `${n} downloads`
})
</script>

<template>
  <div class="mx-auto max-w-5xl">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <PageTitle
          title="Downloads"
          description="Status of your recent video downloads."
        />
        <p
          v-if="countLabel && !loadingHistory"
          class="mt-1 text-sm font-medium text-slate-500"
        >
          {{ countLabel }}
          <span v-if="!auth.isAuthenticated" class="font-normal text-slate-400">
            · this device
          </span>
        </p>
      </div>

      <NuxtLink
        to="/"
        class="inline-flex h-10 shrink-0 items-center justify-center rounded-[10px] bg-brand-600 px-4 text-sm font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100"
      >
        New download
      </NuxtLink>
    </div>

    <!-- Loading -->
    <div
      v-if="loadingHistory"
      class="mt-8 space-y-3"
      aria-busy="true"
      aria-live="polite"
    >
      <span class="sr-only">Loading download history…</span>
      <div
        v-for="n in 3"
        :key="n"
        class="h-16 animate-pulse rounded-2xl bg-slate-100"
      />
    </div>

    <!-- Error -->
    <div
      v-else-if="historyFailed"
      role="alert"
      class="mt-8 flex flex-col items-start gap-4 rounded-2xl border border-red-100 bg-red-50 p-5 sm:flex-row sm:items-center sm:justify-between"
    >
      <div>
        <p class="text-sm font-semibold text-red-800">
          Couldn’t load your download history
        </p>
        <p class="mt-1 text-sm text-red-700/80">
          Check your connection and try again.
        </p>
      </div>
      <button
        type="button"
        class="shrink-0 rounded-[10px] bg-red-600 px-4 py-2 text-sm font-semibold text-white outline-none transition hover:bg-red-700 focus-visible:ring-2 focus-visible:ring-red-600 focus-visible:ring-offset-2"
        @click="history.refetch()"
      >
        Try again
      </button>
    </div>

    <!-- List -->
    <template v-else-if="displayIDs.length">
      <div
        class="mt-8 space-y-3"
        role="list"
        aria-label="Download history"
      >
        <DownloadStatusCard
          v-for="id in displayIDs"
          :key="id"
          :id="id"
        />
      </div>

      <div v-if="auth.isAuthenticated" class="mt-6 text-center">
        <button
          v-if="history.hasNextPage.value"
          type="button"
          :disabled="history.isFetchingNextPage.value"
          class="inline-flex items-center justify-center rounded-[10px] border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold text-slate-700 outline-none transition hover:bg-slate-50 focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60"
          @click="history.fetchNextPage()"
        >
          {{ history.isFetchingNextPage.value ? 'Loading…' : 'Load more' }}
        </button>
        <p v-else class="text-xs text-slate-400">
          You’ve reached the end of your history.
        </p>
      </div>
    </template>

    <!-- Empty -->
    <div
      v-else
      class="mt-8 rounded-2xl border border-dashed border-slate-200 bg-slate-50/80 px-6 py-14 text-center"
    >
      <h2 class="font-semibold text-slate-900">No downloads yet</h2>
      <p class="mx-auto mt-1.5 max-w-sm text-sm leading-6 text-slate-500">
        {{
          auth.isAuthenticated
            ? 'Downloads tied to your account will show up here.'
            : 'Downloads started on this device will show up here.'
        }}
      </p>
      <NuxtLink
        to="/"
        class="mt-5 inline-flex rounded-[10px] bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100"
      >
        Start a download
      </NuxtLink>
    </div>
  </div>
</template>