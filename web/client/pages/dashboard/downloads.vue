<script setup lang="ts">
import { readDownloadIDs } from '~/utils/download-history'
definePageMeta({ layout: 'dashboard' })
const { ids } = useDownloads()
onMounted(() => { if (ids.value.length === 0) ids.value = readDownloadIDs() })
</script>

<template><div><PageTitle title="Downloads" description="Check the status of your recent video downloads."/><div v-if="ids.length" class="space-y-3"><DownloadStatusCard v-for="id in ids" :key="id" :id="id"/></div><div v-else class="rounded-2xl border border-dashed border-slate-300 bg-white px-6 py-14 text-center"><h2 class="font-semibold">No downloads to show</h2><p class="mt-2 text-sm text-slate-500">Your recent downloads will appear here on this device.</p><NuxtLink to="/" class="mt-5 inline-flex rounded-[10px] bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white">Start a download</NuxtLink></div></div></template>
