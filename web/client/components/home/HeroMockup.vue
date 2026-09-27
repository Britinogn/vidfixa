<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { ArrowRight01Icon, CheckmarkCircle01Icon } from '@hugeicons/core-free-icons'
</script>

<template>
  <div class="relative mx-auto mt-2 w-full max-w-[570px] min-w-0 py-2 sm:mt-4 sm:py-4 md:mt-0 md:py-0">
    <div class="absolute -inset-x-3 -inset-y-4 rounded-[35%] bg-red-50/70 blur-2xl sm:-inset-x-6 sm:-inset-y-6 md:-inset-x-8 md:-inset-y-8" />

    <div class="relative min-w-0 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-[0_20px_55px_rgba(15,23,42,.09)] sm:rounded-2xl sm:shadow-[0_24px_70px_rgba(15,23,42,.09)]">
      <div class="flex h-9 items-center gap-1.5 border-b border-slate-100 px-3 sm:h-10 sm:px-4">
        <i class="size-2 rounded-full bg-red-300 sm:size-2.5" />
        <i class="size-2 rounded-full bg-amber-300 sm:size-2.5" />
        <i class="size-2 rounded-full bg-green-300 sm:size-2.5" />
        <span class="ml-2 truncate text-[10px] font-semibold text-slate-600 sm:ml-3 sm:text-xs">VidFixa</span>
      </div>

      <div class="min-w-0 p-3 sm:p-4 md:p-5">
        <div class="mb-3 flex min-w-0 gap-1.5 overflow-x-auto pb-0.5 text-[9px] font-medium text-slate-500 sm:hidden">
          <div class="shrink-0 rounded-lg bg-red-50 px-2.5 py-2 text-brand-600">Home</div>
          <div class="shrink-0 rounded-lg px-2.5 py-2">Downloads</div>
          <div class="shrink-0 rounded-lg px-2.5 py-2">Subscription</div>
          <div class="shrink-0 rounded-lg px-2.5 py-2">Account</div>
        </div>

        <div class="grid min-w-0 grid-cols-1 gap-3 sm:grid-cols-[92px_minmax(0,1fr)] sm:gap-4 md:grid-cols-[108px_minmax(0,1fr)]">
          <div class="hidden space-y-1.5 text-[10px] font-medium text-slate-500 sm:block md:space-y-2 md:text-[11px]">
            <div class="rounded-lg bg-red-50 px-2 py-2 text-brand-600 md:px-2.5">Home</div>
            <div class="rounded-lg px-2 py-2 md:px-2.5">Downloads</div>
            <div class="rounded-lg px-2 py-2 md:px-2.5">Subscription</div>
            <div class="rounded-lg px-2 py-2 md:px-2.5">Account</div>
          </div>

          <div class="grid min-w-0 grid-cols-1 gap-3 sm:grid-cols-[minmax(0,1fr)_52px] md:grid-cols-[minmax(0,1fr)_60px]">
            <div class="media-preview flex aspect-[1.65] min-h-0 items-center justify-center rounded-lg border border-slate-200 sm:rounded-xl">
              <span class="grid size-9 place-items-center rounded-full border border-white/60 bg-slate-900/70 text-white sm:size-11">
                <HugeiconsIcon :icon="ArrowRight01Icon" :size="18" class="sm:hidden" />
                <HugeiconsIcon :icon="ArrowRight01Icon" :size="21" class="hidden sm:block" />
              </span>
            </div>

            <div class="grid grid-cols-3 gap-2 sm:grid-cols-1 sm:space-y-2 sm:gap-0">
              <div class="media-thumb h-12 rounded-lg sm:h-[44px] md:h-[52px]" />
              <div class="media-thumb media-thumb-two h-12 rounded-lg sm:h-[44px] md:h-[52px]" />
              <div class="media-thumb media-thumb-three h-12 rounded-lg sm:h-[44px] md:h-[52px]" />
            </div>

            <div class="col-span-1 mt-0 min-w-0 rounded-xl border border-slate-200 bg-white p-2.5 sm:col-span-2 sm:mt-1 sm:p-3">
              <div class="flex min-w-0 items-center gap-2.5 sm:gap-3">
                <div class="media-thumb h-10 w-14 shrink-0 rounded-lg sm:h-12 sm:w-16" />
                <div class="min-w-0 flex-1">
                  <div class="truncate text-[11px] font-bold sm:text-xs">Beautiful landscape.mp4</div>
                  <div class="mt-1 text-[9px] text-slate-500 sm:text-[10px]">1080p · Processing</div>
                  <div class="mt-1.5 h-1.5 overflow-hidden rounded-full bg-slate-100 sm:mt-2">
                    <div class="h-full w-2/3 rounded-full bg-brand-600" />
                  </div>
                  <div class="mt-1 text-right text-[9px] text-slate-500 sm:text-[10px]">65%</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="relative ml-auto mt-2 flex w-full max-w-[260px] items-center gap-2.5 rounded-xl border border-slate-200 bg-white px-3 py-2.5 shadow-card sm:mt-3 sm:gap-3 sm:px-4 sm:py-3">
      <span class="grid size-8 shrink-0 place-items-center rounded-lg bg-red-50 text-brand-600 sm:size-9">
        <HugeiconsIcon :icon="CheckmarkCircle01Icon" :size="19" />
      </span>
      <div class="min-w-0">
        <div class="text-[11px] font-bold sm:text-xs">Download ready</div>
        <div class="mt-0.5 text-[10px] leading-4 text-slate-500 sm:text-[11px]">Your video is ready to save.</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.media-preview { background: linear-gradient(150deg, #bae6fd, #38bdf8 42%, #0f766e 43%, #164e63 67%, #cbd5e1 68%); }
.media-thumb { background: linear-gradient(150deg, #bae6fd, #38bdf8 45%, #0f766e 46%, #164e63); }
.media-thumb-two { background: linear-gradient(150deg, #fde68a, #fb923c 50%, #7c2d12 51%, #1e293b); }
.media-thumb-three { background: linear-gradient(150deg, #cbd5e1, #64748b 50%, #334155 51%, #0f172a); }
</style>