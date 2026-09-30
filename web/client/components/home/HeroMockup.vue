<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  ArrowRight01Icon,
  CheckmarkCircle01Icon,
  Download01Icon,
} from '@hugeicons/core-free-icons'

/* ---------------------------------------------------------------- nav --- */
const navItems = [
  { id: 'home', label: 'Home' },
  { id: 'downloads', label: 'Downloads' },
  { id: 'subscription', label: 'Subscription' },
  { id: 'account', label: 'Account' },
] as const

type NavId = (typeof navItems)[number]['id']
const activeNav = ref<NavId>('home')

/* ------------------------------------------------------------ preview --- */
const scenes = [
  { id: 'scene-1', label: 'Scene 1 — mountain lake', thumb: 'media-thumb', preview: 'media-preview' },
  { id: 'scene-2', label: 'Scene 2 — desert sunset', thumb: 'media-thumb media-thumb-two', preview: 'media-preview media-preview-two' },
  { id: 'scene-3', label: 'Scene 3 — night sky', thumb: 'media-thumb media-thumb-three', preview: 'media-preview media-preview-three' },
] as const

const activeScene = ref(0)
const currentScene = computed(() => scenes[activeScene.value] ?? scenes[0])

/* ----------------------------------------------------------- progress --- */
const fileName = 'Beautiful landscape.mp4'
const progress = ref(45)
const isComplete = computed(() => progress.value >= 100)

let timer: number | undefined

onMounted(() => {
  const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  if (reduceMotion) {
    progress.value = 100
    return
  }

  timer = window.setInterval(() => {
    progress.value = Math.min(100, progress.value + 5)
    if (progress.value >= 100 && timer) {
      clearInterval(timer)
      timer = undefined
    }
  }, 220)
})

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})

function handleDownload() {
  // wire up the real download here
}
</script>

<template>
  <div class="relative mx-auto mt-2 w-full max-w-[570px] min-w-0 py-2 sm:mt-4 sm:py-4 md:mt-0 md:py-0">
    <!-- decorative glow -->
    <div
      aria-hidden="true"
      class="absolute -inset-x-3 -inset-y-4 rounded-[35%] bg-red-50/70 blur-2xl sm:-inset-x-6 sm:-inset-y-6 md:-inset-x-8 md:-inset-y-8"
    />

    <div
      class="relative min-w-0 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-[0_20px_55px_rgba(15,23,42,.09)] sm:rounded-2xl sm:shadow-[0_24px_70px_rgba(15,23,42,.09)]"
    >
      <!-- window chrome -->
      <div class="flex h-9 items-center gap-1.5 border-b border-slate-100 px-3 sm:h-10 sm:px-4">
        <span aria-hidden="true" class="flex items-center gap-1.5">
          <i class="size-2 rounded-full bg-red-300 sm:size-2.5" />
          <i class="size-2 rounded-full bg-amber-300 sm:size-2.5" />
          <i class="size-2 rounded-full bg-green-300 sm:size-2.5" />
        </span>
        <span class="ml-2 truncate text-[11px] font-semibold text-slate-700 sm:ml-3 sm:text-xs">
          VidFixa
        </span>
      </div>

      <div class="min-w-0 p-3 sm:p-4 md:p-5">
        <div
          class="grid min-w-0 grid-cols-1 gap-3 sm:grid-cols-[92px_minmax(0,1fr)] sm:gap-4 md:grid-cols-[108px_minmax(0,1fr)]"
        >
          <!-- nav: chips on mobile, sidebar on desktop -->
          <nav
            aria-label="Primary"
            class="flex min-w-0 gap-1.5 overflow-x-auto pb-0.5 text-[10px] font-medium sm:flex-col sm:overflow-visible sm:pb-0 sm:text-[11px]"
          >
            <button
              v-for="item in navItems"
              :key="item.id"
              type="button"
              :aria-current="activeNav === item.id ? 'page' : undefined"
              class="shrink-0 rounded-lg px-2.5 py-2 text-left transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2 sm:w-full"
              :class="
                activeNav === item.id
                  ? 'bg-red-50 font-semibold text-brand-600'
                  : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
              "
              @click="activeNav = item.id"
            >
              {{ item.label }}
            </button>
          </nav>

          <div
            class="grid min-w-0 grid-cols-1 gap-3 sm:grid-cols-[minmax(0,1fr)_52px] md:grid-cols-[minmax(0,1fr)_60px]"
          >
            <!-- preview -->
            <div
              class="flex aspect-[1.65] min-h-0 items-center justify-center rounded-lg border border-slate-200 sm:rounded-xl"
              :class="currentScene.preview"
            >
              <button
                type="button"
                aria-label="Play preview"
                class="grid size-10 place-items-center rounded-full border border-white/60 bg-slate-900/70 text-white transition duration-200 hover:scale-105 hover:bg-slate-900/85 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900/70 motion-reduce:transition-none motion-reduce:hover:scale-100 sm:size-12"
              >
                <HugeiconsIcon :icon="ArrowRight01Icon" :size="20" aria-hidden="true" />
              </button>
            </div>

            <!-- scene picker -->
            <div
              role="group"
              aria-label="Choose a scene"
              class="grid grid-cols-3 gap-2 sm:grid-cols-1"
            >
              <button
                v-for="(scene, index) in scenes"
                :key="scene.id"
                type="button"
                :aria-label="scene.label"
                :aria-pressed="activeScene === index"
                class="h-12 rounded-lg transition duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2 motion-reduce:transition-none sm:h-[44px] md:h-[52px]"
                :class="[
                  scene.thumb,
                  activeScene === index
                    ? 'ring-2 ring-brand-600 ring-offset-2'
                    : 'opacity-75 hover:opacity-100',
                ]"
                @click="activeScene = index"
              />
            </div>

            <!-- file card -->
            <div
              class="col-span-1 mt-0 min-w-0 rounded-xl border border-slate-200 bg-white p-2.5 sm:col-span-2 sm:mt-1 sm:p-3"
            >
              <div class="flex min-w-0 items-center gap-2.5 sm:gap-3">
                <div
                  aria-hidden="true"
                  class="media-thumb h-10 w-14 shrink-0 rounded-lg sm:h-12 sm:w-16"
                />
                <div class="min-w-0 flex-1">
                  <div class="truncate text-[11px] font-bold sm:text-xs" :title="fileName">
                    {{ fileName }}
                  </div>
                  <div class="mt-1 text-[10px] text-slate-600 sm:text-[11px]">
                    {{ isComplete ? '1080p · Ready' : '1080p · Processing' }}
                  </div>

                  <div
                    class="mt-1.5 h-1.5 overflow-hidden rounded-full bg-slate-100 sm:mt-2"
                    role="progressbar"
                    :aria-valuenow="progress"
                    aria-valuemin="0"
                    aria-valuemax="100"
                    :aria-label="`Processing ${fileName}`"
                  >
                    <div
                      class="h-full rounded-full bg-brand-600 transition-[width] duration-300 ease-out motion-reduce:transition-none"
                      :style="{ width: `${progress}%` }"
                    />
                  </div>

                  <div
                    aria-hidden="true"
                    class="mt-1 text-right text-[10px] text-slate-600 sm:text-[11px]"
                  >
                    {{ progress }}%
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- status toast (persistent live region so completion is announced) -->
    <div
      role="status"
      aria-live="polite"
      class="relative ml-auto mt-2 w-full max-w-[290px] sm:mt-3"
    >
      <Transition
        enter-active-class="transition duration-300 ease-out motion-reduce:transition-none"
        enter-from-class="translate-y-1 opacity-0"
        leave-active-class="transition duration-200 ease-in motion-reduce:transition-none"
        leave-to-class="translate-y-1 opacity-0"
      >
        <div
          v-if="isComplete"
          class="flex items-center gap-2.5 rounded-xl border border-slate-200 bg-white px-3 py-2.5 shadow-card sm:gap-3 sm:px-4 sm:py-3"
        >
          <span
            aria-hidden="true"
            class="grid size-8 shrink-0 place-items-center rounded-lg bg-red-50 text-brand-600 sm:size-9"
          >
            <HugeiconsIcon :icon="CheckmarkCircle01Icon" :size="19" />
          </span>

          <div class="min-w-0 flex-1">
            <div class="text-[11px] font-bold sm:text-xs">Download ready</div>
            <div class="mt-0.5 text-[10px] leading-4 text-slate-600 sm:text-[11px]">
              Your video is ready to save.
            </div>
          </div>

          <button
            type="button"
            class="inline-flex shrink-0 items-center gap-1.5 rounded-lg bg-brand-600 px-2.5 py-1.5 text-[11px] font-semibold text-white transition hover:brightness-110 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2"
            @click="handleDownload"
          >
            <HugeiconsIcon :icon="Download01Icon" :size="14" aria-hidden="true" />
            Save
          </button>
        </div>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.media-preview {
  background: linear-gradient(150deg, #bae6fd, #38bdf8 42%, #0f766e 43%, #164e63 67%, #cbd5e1 68%);
}
.media-preview-two {
  background: linear-gradient(150deg, #fde68a, #fb923c 50%, #7c2d12 51%, #1e293b);
}
.media-preview-three {
  background: linear-gradient(150deg, #cbd5e1, #64748b 50%, #334155 51%, #0f172a);
}
.media-thumb {
  background: linear-gradient(150deg, #bae6fd, #38bdf8 45%, #0f766e 46%, #164e63);
}
.media-thumb-two {
  background: linear-gradient(150deg, #fde68a, #fb923c 50%, #7c2d12 51%, #1e293b);
}
.media-thumb-three {
  background: linear-gradient(150deg, #cbd5e1, #64748b 50%, #334155 51%, #0f172a);
}
</style>