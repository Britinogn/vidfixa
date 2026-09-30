<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  AlertCircleIcon,
  CheckmarkCircle01Icon,
  ClipboardIcon,
  Download01Icon,
  Link01Icon,
  Loading03Icon,
} from '@hugeicons/core-free-icons'
import { useDownloads } from '~/composables/useDownloads'
import { getApiErrorMessage } from '~/utils/api-error'

const url = ref('')
const message = ref('')
const isMessageError = ref(false)
const { createDownload } = useDownloads()

const isPending = computed(() => createDownload.isPending.value)

// Known hosts we advertise support for — used only to give the user a
// gentle hint, never to block a submit.
const SUPPORTED_HOSTS = [
  // Instagram: Handles www., m., and shortcode paths
  { match: /(?:www\.|m\.)?instagram\.com|instagr\.am/i, name: 'Instagram' },
  
  // Facebook: Handles www., m., fb.watch, and share paths
  { match: /(?:www\.|m\.)?facebook\.com|fb\.watch/i, name: 'Facebook' },
  
  // X (Twitter): Handles www., mobile, and status paths
  { match: /(?:www\.|mobile\.)?(?:twitter|x)\.com/i, name: 'X' },
  
  // LinkedIn: Handles www. and various feed paths
  { match: /(?:www\.)?linkedin\.com/i, name: 'LinkedIn' },
  
  // TikTok: Handles www., m., vm., and vt. subdomains
  { match: /(?:www\.|m\.|vm\.|vt\.)?tiktok\.com/i, name: 'TikTok' },
] as const

const detectedPlatform = computed(() => {
  if (!url.value) return null
  return SUPPORTED_HOSTS.find((h) => h.match.test(url.value))?.name ?? null
})

const isValidUrl = computed(() => {
  if (!url.value) return false
  try {
    const u = new URL(url.value)
    return u.protocol === 'http:' || u.protocol === 'https:'
  } catch {
    return false
  }
})

const canSubmit = computed(() => isValidUrl.value && !isPending.value)

let dismissTimer: number | undefined

function clearMessage() {
  message.value = ''
  if (dismissTimer) {
    clearTimeout(dismissTimer)
    dismissTimer = undefined
  }
}

// Auto-dismiss success messages after a few seconds so they don't linger
// next to a fresh empty input. Errors stay until the user acts.
function scheduleAutoDismiss() {
  if (dismissTimer) clearTimeout(dismissTimer)
  dismissTimer = window.setTimeout(() => {
    if (!isMessageError.value) message.value = ''
  }, 6000)
}

async function pasteFromClipboard() {
  try {
    const text = await navigator.clipboard.readText()
    if (text) {
      url.value = text.trim()
      clearMessage()
    } else {
      message.value = 'Your clipboard is empty. Copy a video link first.'
      isMessageError.value = true
    }
  } catch {
    message.value =
      'Could not access your clipboard. Paste the link manually instead.'
    isMessageError.value = true
  }
}

function clearInput() {
  url.value = ''
  clearMessage()
}

async function submit() {
  if (!canSubmit.value) return
  clearMessage()

  try {
    await createDownload.mutateAsync(url.value.trim())
    url.value = ''
    message.value = 'Your download is queued. We’ll keep its status in your download history.'
    isMessageError.value = false
    scheduleAutoDismiss()
  } catch (error) {
    message.value = getApiErrorMessage(
      error,
      'Could not start the download. Check the link and try again.',
    )
    isMessageError.value = true
  }
}

// Clear any stale "success" message the moment the user starts typing again.
watch(url, () => {
  if (message.value && !isMessageError.value) clearMessage()
})

onBeforeUnmount(() => {
  if (dismissTimer) clearTimeout(dismissTimer)
})
</script>

<template>
  <form class="w-full" novalidate @submit.prevent="submit">
    <div
      class="flex flex-col gap-2 rounded-2xl border bg-white p-2 shadow-card transition-colors sm:flex-row sm:items-center"
      :class="
        url && !isValidUrl
          ? 'border-red-200'
          : 'border-slate-200 focus-within:border-slate-300'
      "
    >
      <HugeiconsIcon
        :icon="Link01Icon"
        :size="19"
        class="ml-3 hidden shrink-0 text-slate-400 sm:block"
        aria-hidden="true"
      />

      <div class="flex min-w-0 flex-1 items-center gap-1">
        <input
          v-model="url"
          type="url"
          inputmode="url"
          autocapitalize="off"
          autocorrect="off"
          spellcheck="false"
          :placeholder="
            detectedPlatform
              ? `Detected ${detectedPlatform} link`
              : 'Paste a video URL (Instagram, Facebook, X, LinkedIn)'
          "
          class="h-12 min-w-0 flex-1 border-0 bg-transparent px-3 text-sm outline-none placeholder:text-slate-400 focus:ring-0 disabled:opacity-60 sm:px-1"
          :disabled="isPending"
          aria-label="Video URL"
          :aria-invalid="Boolean(url) && !isValidUrl"
          aria-describedby="download-form-message"
        >

        <!-- Clear button: appears once there's content -->
        <button
          v-if="url && !isPending"
          type="button"
          class="grid size-9 shrink-0 place-items-center rounded-lg text-slate-400 outline-none transition hover:bg-slate-50 hover:text-slate-600 focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2"
          aria-label="Clear input"
          @click="clearInput"
        >
          <svg viewBox="0 0 16 16" fill="none" class="size-4">
            <path
              d="M4 4l8 8M12 4l-8 8"
              stroke="currentColor"
              stroke-width="1.6"
              stroke-linecap="round"
            />
          </svg>
        </button>

        <!-- Clipboard paste — visible on all viewports now -->
        <button
          type="button"
          class="grid size-10 shrink-0 place-items-center rounded-lg text-slate-400 outline-none transition hover:bg-slate-50 hover:text-slate-600 focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 sm:size-11"
          aria-label="Paste from clipboard"
          :disabled="isPending"
          @click="pasteFromClipboard"
        >
          <HugeiconsIcon :icon="ClipboardIcon" :size="19" aria-hidden="true" />
        </button>
      </div>

      <button
        type="submit"
        :disabled="!canSubmit"
        class="inline-flex h-12 w-full shrink-0 items-center justify-center gap-2 whitespace-nowrap rounded-[10px] bg-brand-600 px-6 text-sm font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100 disabled:cursor-not-allowed disabled:opacity-60 sm:w-auto"
      >
        <HugeiconsIcon
          :icon="isPending ? Loading03Icon : Download01Icon"
          :size="18"
          :class="isPending && 'animate-spin motion-reduce:animate-none'"
          aria-hidden="true"
        />
        {{ isPending ? 'Starting…' : 'Download' }}
      </button>
    </div>

    <!-- Inline hint when the URL looks wrong -->
    <p
      v-if="url && !isValidUrl"
      class="mt-2 flex items-center gap-1.5 text-xs text-slate-500"
    >
      <HugeiconsIcon :icon="AlertCircleIcon" :size="13" aria-hidden="true" />
      That doesn’t look like a full link. Try pasting the whole URL from the
      address bar.
    </p>

    <!-- Status message -->
    <p
      v-if="message"
      id="download-form-message"
      :role="isMessageError ? 'alert' : 'status'"
      aria-live="polite"
      class="mt-3 flex items-start gap-2 text-sm"
      :class="isMessageError ? 'text-red-700' : 'text-green-700'"
    >
      <HugeiconsIcon
        :icon="isMessageError ? AlertCircleIcon : CheckmarkCircle01Icon"
        :size="16"
        class="mt-0.5 shrink-0"
        aria-hidden="true"
      />
      <span>{{ message }}</span>
    </p>
  </form>
</template>