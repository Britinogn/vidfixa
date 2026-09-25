<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { Download01Icon, Cancel01Icon } from '@hugeicons/core-free-icons'
import { useDownloads } from '~/composables/useDownloads'
import { forgetDownload } from '~/utils/download-history'

const props = defineProps<{ id: string }>()
const { ids, useDownload, downloadFile } = useDownloads()
const query = useDownload(props.id)
const removing = ref(false)
const statusStyle = computed(() => ({
  completed: 'bg-green-50 text-green-700',
  processing: 'bg-amber-50 text-amber-700',
  queued: 'bg-slate-100 text-slate-600',
  failed: 'bg-red-50 text-red-700',
  cancelled: 'bg-slate-100 text-slate-500',
}[query.data.value?.status || 'queued']))

function remove() {
  removing.value = true
  forgetDownload(props.id)
  ids.value = ids.value.filter((id) => id !== props.id)
}
</script>

<template>
  <article v-if="!removing" class="flex flex-col gap-4 rounded-2xl border border-slate-200 bg-white p-4 sm:flex-row sm:items-center sm:justify-between">
    <div class="min-w-0">
      <div class="truncate text-sm font-semibold text-slate-800">{{ query.data.value?.url || `Download ${id.slice(0, 8)}` }}</div>
      <div class="mt-1 text-xs text-slate-500">{{ query.data.value ? new Date(query.data.value.created_at).toLocaleString() : query.isPending.value ? 'Checking status…' : 'Status unavailable' }}</div>
      <p v-if="query.data.value?.error" class="mt-2 text-xs text-red-700">{{ query.data.value.error }}</p>
    </div>
    <div class="flex shrink-0 items-center gap-2">
      <span class="rounded-full px-3 py-1 text-xs font-semibold capitalize" :class="statusStyle">{{ query.data.value?.status || 'loading' }}</span>
      <button v-if="query.data.value?.status === 'completed'" class="inline-flex items-center gap-2 rounded-[10px] bg-brand-600 px-3 py-2 text-xs font-semibold text-white hover:bg-brand-700" @click="downloadFile(id)"><HugeiconsIcon :icon="Download01Icon" :size="15" /> Save file</button>
      <button class="rounded-lg p-2 text-slate-400 hover:bg-slate-100 hover:text-slate-700" aria-label="Remove from history" @click="remove"><HugeiconsIcon :icon="Cancel01Icon" :size="17" /></button>
    </div>
  </article>
</template>
