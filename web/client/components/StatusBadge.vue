<script setup lang="ts">
import type { DownloadStatus } from '~/types/api'
defineProps<{ status: DownloadStatus | string }>()
const styles: Record<string, string> = {
  completed: 'bg-green-50 text-green-700',
  processing: 'bg-amber-50 text-amber-700',
  queued: 'bg-slate-100 text-slate-600',
  failed: 'bg-red-50 text-red-700',
  cancelled: 'bg-slate-100 text-slate-500',
  active: 'bg-green-50 text-green-700',
  SUCCEEDED: 'bg-green-50 text-green-700',
  pending: 'bg-amber-50 text-amber-700',
}
</script>

<template><span class="inline-flex rounded-full px-2.5 py-1 text-xs font-semibold capitalize" :class="styles[status] || 'bg-slate-100 text-slate-600'">{{ status.replaceAll('_', ' ').toLowerCase() }}</span></template>
