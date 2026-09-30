<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  Link01Icon,
  Settings01Icon,
  Download01Icon,
} from '@hugeicons/core-free-icons'

const steps = [
  { title: 'Paste the URL', detail: 'Copy and paste a video link from a supported platform.', icon: Link01Icon },
  { title: 'We process it', detail: 'VidFixa prepares your video for download.', icon: Settings01Icon },
  { title: 'Download and enjoy', detail: 'Save the finished video to your device.', icon: Download01Icon },
]
</script>

<template>
  <section id="how-it-works" class="bg-slate-50 py-16 sm:py-20 md:py-24">
    <div class="mx-auto w-full max-w-7xl px-4 sm:px-5 md:px-8">
      <div class="text-center">
        <h2 class="text-2xl font-bold tracking-tight sm:text-3xl md:text-[32px]">
          How it works
        </h2>
        <p class="mt-3 text-sm text-slate-500 sm:text-base">
          Get your video in three simple steps.
        </p>
      </div>

      <div class="relative mt-12 sm:mt-16">
        <!-- Dashed connector behind the icon row (md and up) -->
        <div
          class="pointer-events-none absolute left-[16%] right-[16%] top-7 hidden border-t border-dashed border-slate-300 md:block"
          aria-hidden="true"
        />

        <ol class="relative grid gap-10 md:grid-cols-3 md:gap-6">
          <li
            v-for="(step, index) in steps"
            :key="step.title"
            class="flex flex-col items-center text-center"
          >
            <!-- Icon tile with number badge -->
            <div class="relative">
              <span
                class="grid size-14 place-items-center rounded-2xl bg-red-50 text-brand-600 ring-1 ring-red-100"
              >
                <HugeiconsIcon :icon="step.icon" :size="24" aria-hidden="true" />
              </span>
              <span
                class="absolute -right-2 -top-2 grid size-6 place-items-center rounded-full bg-white text-[10px] font-bold text-brand-600 ring-1 ring-red-100"
                aria-hidden="true"
              >
                {{ index + 1 }}
              </span>
            </div>

            <h3 class="mt-5 text-lg font-semibold">{{ step.title }}</h3>
            <p class="mt-2 max-w-[260px] text-sm leading-6 text-slate-500">
              {{ step.detail }}
            </p>
          </li>
        </ol>
      </div>
    </div>
  </section>
</template>