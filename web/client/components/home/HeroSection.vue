<script setup lang="ts">
import { readDownloadIDs } from '~/utils/download-history'

const auth = useAuthStore()
const { ids } = useDownloads()

onMounted(() => {
  if (ids.value.length === 0) ids.value = readDownloadIDs()
})
</script>

<template>
  <section class="relative overflow-hidden">
    <div class="mx-auto grid w-full max-w-7xl items-center gap-10 px-4 pb-16 pt-8 sm:px-5 sm:pb-20 sm:pt-10 md:grid-cols-[1.05fr_.95fr] md:gap-12 md:px-8 md:pb-24 md:pt-10">
      <div class="relative z-[1] min-w-0">
        <div class="mb-5 inline-flex max-w-full items-center gap-1.5 rounded-full bg-red-50 px-3 py-1.5 text-[11px] font-semibold text-brand-600 sm:gap-2 sm:text-xs">
          <span class="size-1.5 shrink-0 rounded-full bg-brand-600" />
          <span>Fast</span><span>·</span><span>Secure</span><span>·</span><span>No watermark</span>
        </div>

        <h1 class="max-w-2xl text-[36px] font-extrabold leading-[1.08] tracking-[-0.045em] text-slate-900 min-[375px]:text-[40px] sm:text-5xl md:text-[56px]">
          Download videos.<br>Simple. <span class="text-brand-600">Fast.</span>
        </h1>

        <p class="mt-4 max-w-xl text-sm leading-6 text-slate-500 sm:mt-5 sm:text-base sm:leading-7 md:text-lg">
          Save videos from your favorite platforms in just a few seconds. No hassle, no watermark.
        </p>

        <div class="mt-6 w-full max-w-[610px] sm:mt-8">
          <DownloadForm />
        </div>

        <HomePlatformBadges class="mt-6 sm:mt-7" />

        <section v-if="ids.length" class="mt-7 w-full max-w-[610px] sm:mt-8" aria-live="polite">
          <div class="mb-3 flex min-w-0 items-center justify-between gap-3">
            <h2 class="text-sm font-semibold text-slate-800">Recent downloads</h2>
            <NuxtLink v-if="auth.user && auth.user.role === 'user'" to="/dashboard/downloads" class="shrink-0 text-xs font-semibold text-brand-600 hover:text-brand-700">
              Open dashboard
            </NuxtLink>
          </div>
          <div class="space-y-3">
            <DownloadStatusCard v-for="id in ids.slice(0, 3)" :key="id" :id="id" />
          </div>
        </section>
      </div>

      <HomeHeroMockup />
    </div>
  </section>
</template>