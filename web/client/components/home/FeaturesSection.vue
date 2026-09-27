<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { ZapIcon, Shield01Icon, CircleIcon, InfinityIcon, LaptopIcon, LockKeyIcon } from '@hugeicons/core-free-icons'

const features = [
  { title: 'Super fast', detail: 'Download videos in seconds, not minutes.', icon: ZapIcon },
  { title: 'High quality', detail: 'Get the best available quality from the source.', icon: Shield01Icon },
  { title: 'No watermark', detail: 'Your videos, clean and ready to use.', icon: CircleIcon },
  { title: 'Download history', detail: 'Keep track of your recent downloads.', icon: InfinityIcon },
  { title: 'Works everywhere', detail: 'Use VidFixa on desktop, tablet, or mobile.', icon: LaptopIcon },
  { title: 'Secure and private', detail: 'Your account stays protected.', icon: LockKeyIcon },
]
</script>

<template>
  <section id="features" class="py-16 sm:py-20 md:py-24">
    <div class="mx-auto w-full max-w-7xl px-4 sm:px-5 md:px-8">
      <div class="mx-auto max-w-2xl text-center">
        <span class="inline-flex rounded-full bg-red-50 px-3 py-1.5 text-xs font-semibold text-brand-600">Why choose VidFixa</span>
        <h2 class="mt-5 text-2xl font-bold tracking-tight sm:text-3xl md:text-[32px]">Everything you need, in one place.</h2>
        <p class="mt-3 text-sm leading-6 text-slate-500 sm:text-base">A straightforward video download experience, built around your workflow.</p>
      </div>

      <div class="mt-8 grid gap-3 sm:mt-10 sm:grid-cols-2 sm:gap-4 lg:grid-cols-3">
        <article v-for="feature in features" :key="feature.title" class="flex min-w-0 gap-3.5 rounded-2xl border border-slate-200 bg-white p-4 sm:gap-4 sm:p-5">
          <span class="grid size-10 shrink-0 place-items-center rounded-xl bg-red-50 text-brand-600 sm:size-11">
            <HugeiconsIcon :icon="feature.icon" :size="21" />
          </span>
          <div class="min-w-0">
            <h3 class="text-sm font-semibold">{{ feature.title }}</h3>
            <p class="mt-1.5 text-xs leading-5 text-slate-500">{{ feature.detail }}</p>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>