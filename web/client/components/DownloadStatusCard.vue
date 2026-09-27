<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { Download01Icon, Cancel01Icon, PlayIcon, Loading03Icon } from '@hugeicons/core-free-icons'
import { useDownloads } from '~/composables/useDownloads'

const props = defineProps<{ id: string }>()
const { ids, useDownload, useDownloadFile, forgetDownload } = useDownloads()
const query = useDownload(props.id)
const { previewUrl, previewError, saveError, opening, saving, openPreview, closePreview, saveFile } = useDownloadFile(props.id)
const removing = ref(false)

const statusStyle = computed(() => ({
  completed: 'bg-green-50 text-green-700',
  processing: 'bg-amber-50 text-amber-700',
  queued: 'bg-slate-100 text-slate-600',
  failed: 'bg-red-50 text-red-700',
  cancelled: 'bg-slate-100 text-slate-500',
}[query.data.value?.status || 'queued']))

function remove() {
  if (removing.value) return
  removing.value = true
  setTimeout(() => {
    forgetDownload(props.id)
    ids.value = ids.value.filter((id) => id !== props.id)
  }, 150)
}
</script>

<template>
  <Transition
    enter-active-class="transition ease-out duration-200"
    enter-from-class="opacity-0 -translate-y-1"
    leave-active-class="transition ease-in duration-150"
    leave-to-class="opacity-0 scale-95"
  >
    <article v-if="!removing" class="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white p-4">
      <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div class="min-w-0">
          <div class="truncate text-sm font-semibold text-slate-800" :title="query.data.value?.url">
            {{ query.data.value?.url || `Download ${id.slice(0, 8)}` }}
          </div>
          <div class="mt-1 text-xs text-slate-500">
            {{ query.data.value ? new Date(query.data.value.created_at).toLocaleString() : query.isPending.value ? 'Checking status…' : query.isError.value ? 'Status unavailable' : '' }}
          </div>
          <p v-if="query.data.value?.error" class="mt-2 text-xs text-red-700">{{ query.data.value.error }}</p>
        </div>

        <div class="flex flex-wrap items-center justify-end gap-2 sm:flex-nowrap sm:shrink-0">
          <span class="rounded-full px-3 py-1 text-xs font-semibold capitalize" :class="statusStyle" aria-live="polite">
            {{ query.data.value?.status || 'loading' }}
          </span>

          <button
            v-if="query.data.value?.status === 'completed' && !previewUrl && !previewError"
            class="inline-flex items-center gap-2 rounded-[10px] border border-slate-200 px-3 py-2.5 text-xs font-semibold text-slate-700 hover:bg-slate-50"
            :disabled="opening"
            @click="openPreview"
          >
            <HugeiconsIcon :icon="opening ? Loading03Icon : PlayIcon" :size="15" :class="opening && 'animate-spin'" />
            {{ opening ? 'Loading…' : 'Preview' }}
          </button>

          <button
            v-if="query.data.value?.status === 'completed'"
            :disabled="saving"
            class="inline-flex items-center gap-2 rounded-[10px] bg-brand-600 px-3 py-2.5 text-xs font-semibold text-white hover:bg-brand-700 active:bg-brand-800"
            @click="saveFile"
          >
            <HugeiconsIcon :icon="saving ? Loading03Icon : Download01Icon" :size="15" :class="saving && 'animate-spin'" />
            {{ saving ? 'Starting…' : 'Save file' }}
          </button>

          <button
            class="grid size-9 shrink-0 place-items-center rounded-lg text-slate-400 hover:bg-slate-100 hover:text-slate-700 active:bg-slate-200"
            aria-label="Remove from history"
            :disabled="removing"
            @click="remove"
          >
            <HugeiconsIcon :icon="Cancel01Icon" :size="17" />
          </button>
        </div>
      </div>

      <p v-if="previewError" class="text-xs text-red-700">{{ previewError }}</p>
      <p v-if="saveError" class="text-xs text-red-700">{{ saveError }}</p>

      <div v-if="previewUrl" class="relative">
        <video
          :src="previewUrl"
          controls
          playsinline
          preload="metadata"
          class="w-full rounded-xl bg-black"
        />
        <button
          class="absolute right-2 top-2 grid size-8 place-items-center rounded-lg bg-black/60 text-white hover:bg-black/80"
          aria-label="Close preview"
          @click="closePreview"
        >
          <HugeiconsIcon :icon="Cancel01Icon" :size="16" />
        </button>
      </div>
    </article>
  </Transition>
</template>
