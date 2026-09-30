<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  AlertCircleIcon,
  Cancel01Icon,
  CheckmarkCircle01Icon,
  Copy01Icon,
  Download01Icon,
  Loading03Icon,
  PlayIcon,
  RefreshIcon,
} from '@hugeicons/core-free-icons'
import { useDownloads } from '~/composables/useDownloads'

const props = defineProps<{ id: string }>()

const { ids, useDownload, useDownloadFile, forgetDownload } = useDownloads()
const query = useDownload(props.id)

/* Poll while non-terminal; stop once the status can no longer change.
   Local to this component because useDownload(id) doesn't accept options. */
let pollTimer: number | undefined

function stopPolling() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = undefined
  }
}

function startPolling() {
  stopPolling()
  pollTimer = window.setInterval(() => {
    query.refetch()
  }, 2500)
}

watch(
  () => query.data.value?.status,
  (s) => {
    if (s === 'queued' || s === 'processing') startPolling()
    else stopPolling()
  },
  { immediate: true },
)

onBeforeUnmount(stopPolling)


// Poll while non-terminal, stop once complete/failed/cancelled.
// const query = useDownload(props.id, {
//   refetchInterval: (q) => {
//     const status = q.state.data?.status
//     if (status === 'queued' || status === 'processing') return 2500
//     return false
//   },
// })

const {
  previewUrl,
  previewError,
  saveError,
  opening,
  saving,
  openPreview,
  closePreview,
  saveFile,
} = useDownloadFile(props.id)

const removing = ref(false)
const copied = ref(false)

const status = computed(() => query.data.value?.status ?? 'loading')

const statusConfig = computed(() => {
  switch (status.value) {
    case 'completed':
      return {
        label: 'Ready',
        icon: CheckmarkCircle01Icon,
        classes: 'bg-green-50 text-green-700',
      }
    case 'processing':
      return {
        label: 'Processing',
        icon: Loading03Icon,
        classes: 'bg-amber-50 text-amber-700',
        spin: true,
      }
    case 'queued':
      return {
        label: 'Queued',
        icon: Loading03Icon,
        classes: 'bg-slate-100 text-slate-600',
        spin: true,
      }
    case 'failed':
      return {
        label: 'Failed',
        icon: AlertCircleIcon,
        classes: 'bg-red-50 text-red-700',
      }
    case 'cancelled':
      return {
        label: 'Cancelled',
        icon: Cancel01Icon,
        classes: 'bg-slate-100 text-slate-500',
      }
    default:
      return {
        label: 'Loading…',
        icon: Loading03Icon,
        classes: 'bg-slate-100 text-slate-500',
        spin: true,
      }
  }
})

const shortUrl = computed(() => {
  const raw = query.data.value?.url
  if (!raw) return `Download ${props.id.slice(0, 8)}`
  try {
    const u = new URL(raw)
    return `${u.hostname.replace(/^www\./, '')}${u.pathname.length > 1 ? u.pathname : ''}`
  } catch {
    return raw
  }
})

const createdAt = computed(() => {
  const iso = query.data.value?.created_at
  if (!iso) return ''
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleString('en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
})

const isTerminal = computed(() =>
  ['completed', 'failed', 'cancelled'].includes(status.value),
)

const canPreview = computed(
  () => status.value === 'completed' && !previewUrl.value && !previewError.value,
)

function remove() {
  if (removing.value) return
  removing.value = true
  // Let the leave transition finish before we mutate the list.
  // The Transition wrapper in the template handles the actual DOM removal.
}

function onAfterLeave() {
  forgetDownload(props.id)
  ids.value = ids.value.filter((id) => id !== props.id)
}

async function copyUrl() {
  const raw = query.data.value?.url
  if (!raw) return
  try {
    await navigator.clipboard.writeText(raw)
    copied.value = true
    window.setTimeout(() => (copied.value = false), 1800)
  } catch {
    // Silent — the URL is visible in the title attribute anyway.
  }
}

function retry() {
  query.refetch()
}

onBeforeUnmount(() => {
  // Nothing to clean up now that we removed the setTimeout in remove().
})
</script>

<template>
  <Transition
    enter-active-class="transition ease-out duration-200 motion-reduce:transition-none"
    enter-from-class="opacity-0 -translate-y-1"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition ease-in duration-150 motion-reduce:transition-none"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 scale-95"
    @after-leave="onAfterLeave"
  >
    <article
      v-if="!removing"
      class="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white p-4 transition-shadow hover:shadow-card"
    >
      <!-- <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between"> -->
      <div class="flex flex-col gap-4 @[600px]:flex-row @[600px]:items-start @[600px]:justify-between">
        <!-- Identity -->
        <div class="min-w-0 flex-1 md:max-w-[55%] lg:max-w-none">
          <div class="flex items-center gap-2">
            <a
              v-if="query.data.value?.url"
              :href="query.data.value.url"
              target="_blank"
              rel="noopener noreferrer"
              :title="query.data.value.url"
              class="truncate text-sm font-semibold text-slate-800 outline-none transition hover:text-brand-600 focus-visible:underline focus-visible:underline-offset-2"
            >
              {{ shortUrl }}
            </a>
            <span v-else class="truncate text-sm font-semibold text-slate-800">
              {{ shortUrl }}
            </span>

            <button
              v-if="query.data.value?.url"
              type="button"
              class="grid size-6 shrink-0 place-items-center rounded-md text-slate-400 outline-none transition hover:bg-slate-100 hover:text-slate-700 focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2"
              :aria-label="copied ? 'URL copied' : 'Copy URL'"
              @click="copyUrl"
            >
              <HugeiconsIcon
                :icon="copied ? CheckmarkCircle01Icon : Copy01Icon"
                :size="13"
                aria-hidden="true"
              />
            </button>
          </div>

          <div class="mt-1 text-xs text-slate-500">
            <template v-if="createdAt">{{ createdAt }}</template>
            <template v-else-if="query.isPending.value">Checking status…</template>
            <template v-else-if="query.isError.value">Status unavailable</template>
          </div>

          <!-- Failure reason -->
          <p
            v-if="query.data.value?.error"
            class="mt-2 flex items-start gap-1.5 rounded-lg bg-red-50 px-2.5 py-1.5 text-xs text-red-700"
          >
            <HugeiconsIcon
              :icon="AlertCircleIcon"
              :size="13"
              class="mt-0.5 shrink-0"
              aria-hidden="true"
            />
            <span>{{ query.data.value.error }}</span>
          </p>
        </div>

        <!-- Actions -->
        <div class="flex flex-wrap items-center justify-end gap-2 md:justify-end lg:flex-nowrap lg:shrink-0">
          <span
            class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-semibold"
            :class="statusConfig.classes"
            aria-live="polite"
          >
            <HugeiconsIcon
              :icon="statusConfig.icon"
              :size="12"
              :class="statusConfig.spin && 'animate-spin motion-reduce:animate-none'"
              aria-hidden="true"
            />
            {{ statusConfig.label }}
          </span>

          <button
            v-if="canPreview"
            type="button"
            class="inline-flex items-center gap-2 rounded-[10px] border border-slate-200 px-3 py-2 text-xs font-semibold text-slate-700 outline-none transition hover:bg-slate-50 focus-visible:ring-4 focus-visible:ring-red-100 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="opening"
            @click="openPreview"
          >
            <HugeiconsIcon
              :icon="opening ? Loading03Icon : PlayIcon"
              :size="14"
              :class="opening && 'animate-spin motion-reduce:animate-none'"
              aria-hidden="true"
            />
            {{ opening ? 'Loading…' : 'Preview' }}
          </button>

          <button
            v-if="status === 'completed'"
            type="button"
            :disabled="saving"
            class="inline-flex items-center gap-2 rounded-[10px] bg-brand-600 px-3 py-2 text-xs font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100 disabled:cursor-not-allowed disabled:opacity-60"
            @click="saveFile"
          >
            <HugeiconsIcon
              :icon="saving ? Loading03Icon : Download01Icon"
              :size="14"
              :class="saving && 'animate-spin motion-reduce:animate-none'"
              aria-hidden="true"
            />
            {{ saving ? 'Starting…' : 'Save file' }}
          </button>

          <button
            v-if="status === 'failed'"
            type="button"
            class="inline-flex items-center gap-2 rounded-[10px] border border-slate-200 px-3 py-2 text-xs font-semibold text-slate-700 outline-none transition hover:bg-slate-50 focus-visible:ring-4 focus-visible:ring-red-100"
            @click="retry"
          >
            <HugeiconsIcon :icon="RefreshIcon" :size="14" aria-hidden="true" />
            Retry
          </button>

          <button
            type="button"
            class="grid size-9 shrink-0 place-items-center rounded-lg text-slate-400 outline-none transition hover:bg-slate-100 hover:text-slate-700 focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2"
            aria-label="Remove from history"
            :disabled="removing"
            @click="remove"
          >
            <HugeiconsIcon :icon="Cancel01Icon" :size="17" aria-hidden="true" />
          </button>
        </div>
      </div>

      <!-- Inline errors from preview/save hooks -->
      <p
        v-if="previewError || saveError"
        role="alert"
        class="flex items-start gap-1.5 text-xs text-red-700"
      >
        <HugeiconsIcon
          :icon="AlertCircleIcon"
          :size="13"
          class="mt-0.5 shrink-0"
          aria-hidden="true"
        />
        <span>{{ previewError || saveError }}</span>
      </p>

      <!-- Preview -->
      <div v-if="previewUrl" class="relative">
        <video
          :src="previewUrl"
          controls
          playsinline
          preload="metadata"
          class="w-full rounded-xl bg-black"
        />
        <button
          type="button"
          class="absolute right-2 top-2 grid size-8 place-items-center rounded-lg bg-black/70 text-white outline-none transition hover:bg-black/90 focus-visible:ring-2 focus-visible:ring-white"
          aria-label="Close preview"
          @click="closePreview"
        >
          <HugeiconsIcon :icon="Cancel01Icon" :size="16" aria-hidden="true" />
        </button>
      </div>
    </article>
  </Transition>
</template>