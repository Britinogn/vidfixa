<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { Link01Icon, Settings01Icon, Download01Icon, ArrowRight01Icon } from '@hugeicons/core-free-icons'

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
        <h2 class="text-2xl font-bold tracking-tight sm:text-3xl md:text-[32px]">How it works</h2>
        <p class="mt-3 text-sm text-slate-500 sm:text-base">Get your video in three simple steps.</p>
      </div>

      <div class="mt-10 grid gap-10 sm:mt-14 md:grid-cols-3 md:gap-6">
        <div v-for="(step, index) in steps" :key="step.title" class="relative flex flex-col items-center text-center">
          <span class="mb-5 grid size-12 place-items-center rounded-full bg-red-50 text-sm font-bold text-brand-600">0{{ index + 1 }}</span>
          <HugeiconsIcon :icon="step.icon" :size="30" class="text-slate-800" />
          <h3 class="mt-4 text-lg font-semibold">{{ step.title }}</h3>
          <p class="mt-2 max-w-[250px] text-sm leading-6 text-slate-500">{{ step.detail }}</p>
          <HugeiconsIcon v-if="index < 2" :icon="ArrowRight01Icon" :size="23" class="absolute -right-4 top-16 hidden text-slate-400 md:block" />
        </div>
      </div>
    </div>
  </section>
</template>