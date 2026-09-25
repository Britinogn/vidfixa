<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { Link01Icon, Download01Icon } from '@hugeicons/core-free-icons'
import { useDownloads } from '~/composables/useDownloads'

const url = ref('')
const message = ref('')
const { createDownload } = useDownloads()

async function submit() {
  message.value = ''
  try {
    await createDownload.mutateAsync(url.value.trim())
    url.value = ''
    message.value = 'Your download is queued. We’ll keep its status in your download history.'
  } catch (error) {
    const response = error as { data?: string | { message?: string }; message?: string }
    message.value = typeof response.data === 'string' ? response.data : response.data?.message || response.message || 'Could not start the download. Check the link and try again.'
  }
}
</script>

<template>
  <form class="w-full" @submit.prevent="submit">
    <div class="flex flex-col gap-2 rounded-2xl border border-slate-200 bg-white p-2 shadow-card sm:flex-row sm:items-center">
      <HugeiconsIcon :icon="Link01Icon" :size="19" class="ml-3 hidden shrink-0 text-slate-400 sm:block" />
      <input v-model="url" required type="url" placeholder="Paste a video URL (Instagram, Facebook, X, LinkedIn)" class="h-12 min-w-0 flex-1 border-0 bg-transparent px-3 text-sm outline-none placeholder:text-slate-400 focus:ring-0 sm:px-1">
      <button type="submit" :disabled="createDownload.isPending.value" class="inline-flex h-12 items-center justify-center gap-2 rounded-[10px] bg-brand-600 px-6 text-sm font-semibold text-white transition hover:bg-brand-700 disabled:cursor-not-allowed disabled:opacity-60">
        <HugeiconsIcon :icon="Download01Icon" :size="18" />{{ createDownload.isPending.value ? 'Starting…' : 'Download' }}
      </button>
    </div>
    <p v-if="message" class="mt-3 text-sm" :class="createDownload.isError.value ? 'text-red-700' : 'text-green-700'" role="status">{{ message }}</p>
  </form>
</template>
