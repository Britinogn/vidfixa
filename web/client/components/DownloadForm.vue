<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { Link01Icon, Download01Icon, ClipboardIcon } from '@hugeicons/core-free-icons'
import { useDownloads } from '~/composables/useDownloads'

const url = ref('')
const message = ref('')
const { createDownload } = useDownloads()

function clearMessage() {
  message.value = ''
}

async function pasteFromClipboard() {
  try {
    const text = await navigator.clipboard.readText()
    if (text) {
      url.value = text.trim()
      clearMessage()
    }
  } catch {
    // Clipboard permission denied or unsupported — fail silently,
    // the user can still paste manually into the field.
  }
}

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

      <div class="flex flex-1 items-center gap-1">
        <input
          v-model="url"
          required
          type="url"
          inputmode="url"
          autocapitalize="off"
          autocorrect="off"
          placeholder="Paste a video URL (Instagram, Facebook, X, LinkedIn)"
          class="h-12 min-w-0 flex-1 border-0 bg-transparent px-3 text-sm outline-none placeholder:text-slate-400 focus:ring-0 disabled:opacity-60 sm:px-1"
          :disabled="createDownload.isPending.value"
          @input="clearMessage"
        >
        <button
          type="button"
          class="grid size-11 shrink-0 place-items-center rounded-lg text-slate-400 hover:bg-slate-50 hover:text-slate-600 active:bg-slate-100 sm:hidden"
          aria-label="Paste from clipboard"
          :disabled="createDownload.isPending.value"
          @click="pasteFromClipboard"
        >
          <HugeiconsIcon :icon="ClipboardIcon" :size="19" />
        </button>
      </div>

      <button
        type="submit"
        :disabled="createDownload.isPending.value"
        class="inline-flex h-12 w-full items-center justify-center gap-2 whitespace-nowrap rounded-[10px] bg-brand-600 px-6 text-sm font-semibold text-white transition hover:bg-brand-700 active:bg-brand-800 disabled:cursor-not-allowed disabled:opacity-60 sm:w-auto"
      >
        <HugeiconsIcon :icon="Download01Icon" :size="18" />{{ createDownload.isPending.value ? 'Starting…' : 'Download' }}
      </button>
    </div>
    <p v-if="message" class="mt-3 text-sm" :class="createDownload.isError.value ? 'text-red-700' : 'text-green-700'" role="status">{{ message }}</p>
  </form>
</template>