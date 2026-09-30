<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import type { DashboardResponse } from '~/types/api'
import { readDownloadIDs } from '~/utils/download-history'

definePageMeta({
  layout: 'dashboard',
  middleware: [
    function () {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) return navigateTo('/login')
    },
  ],
})

const auth = useAuthStore()
const api = useApi()
const { ids, useDownloadHistory } = useDownloads()
const history = useDownloadHistory(5)

if (ids.value.length === 0) {
  ids.value = readDownloadIDs()
}

const dashboard = useQuery({
  queryKey: ['dashboard'] as const,
  queryFn: () => api<DashboardResponse>('/dashboard/'),
  enabled: computed(() => auth.isAuthenticated),
  staleTime: 60_000,
})

function firstNameOf(name?: string | null) {
  return name?.trim().split(/\s+/)[0] || ''
}

const greetingName = computed(
  () =>
    firstNameOf(dashboard.data.value?.user.full_name) ||
    firstNameOf(auth.user?.full_name) ||
    'there',
)

const plan = computed(() => dashboard.data.value?.plan ?? 'free')
const isFree = computed(() => plan.value === 'free')
// Logged-in users read recent downloads from the server so the list
// survives logout; anonymous users keep the device-local list.
const recentIds = computed(() =>
  auth.isAuthenticated && history.data.value
    ? history.data.value.map((d) => d.id)
    : ids.value.slice(0, 5),
)
</script>

<template>
  <div class="mx-auto max-w-5xl">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <PageTitle
        :title="`Welcome back, ${greetingName}`"
        description="Your usage, plan, and recent downloads."
      />
      <NuxtLink
        to="/"
        class="inline-flex h-10 shrink-0 items-center justify-center rounded-[10px] bg-brand-600 px-4 text-sm font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100"
      >
        New download
      </NuxtLink>
    </div>

    <!-- Loading -->
    <div
      v-if="dashboard.isPending.value"
      class="mt-8 space-y-8"
      aria-busy="true"
      aria-live="polite"
    >
      <span class="sr-only">Loading your dashboard…</span>

      <div class="grid gap-4 sm:gap-5 lg:grid-cols-[1.4fr_1fr]">
        <div class="h-40 animate-pulse rounded-2xl bg-slate-100 sm:h-44" />
        <div class="h-40 animate-pulse rounded-2xl bg-slate-100 sm:h-44" />
      </div>

      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <div class="h-6 w-36 animate-pulse rounded bg-slate-100" />
          <div class="h-4 w-16 animate-pulse rounded bg-slate-100" />
        </div>
        <div
          v-for="n in 3"
          :key="n"
          class="h-16 animate-pulse rounded-2xl bg-slate-100"
        />
      </div>
    </div>

    <!-- Error -->
    <div
      v-else-if="dashboard.isError.value"
      role="alert"
      class="mt-8 flex flex-col items-start gap-4 rounded-2xl border border-red-100 bg-red-50 p-5 sm:flex-row sm:items-center sm:justify-between"
    >
      <div>
        <p class="text-sm font-semibold text-red-800">Couldn’t load your dashboard</p>
        <p class="mt-1 text-sm text-red-700/80">
          Check your connection and try again.
        </p>
      </div>
      <button
        type="button"
        class="shrink-0 rounded-[10px] bg-red-600 px-4 py-2 text-sm font-semibold text-white outline-none transition hover:bg-red-700 focus-visible:ring-2 focus-visible:ring-red-600 focus-visible:ring-offset-2"
        @click="dashboard.refetch()"
      >
        Try again
      </button>
    </div>

    <template v-else-if="dashboard.data.value">
      <!-- Overview -->
      <div class="mt-8 grid gap-4 sm:gap-5 lg:grid-cols-[1.4fr_1fr] lg:items-stretch">
        <UsageCard :usage="dashboard.data.value.usage" :plan="plan" />

        <section
          class="flex flex-col justify-between rounded-2xl border border-slate-200 bg-white p-5 shadow-card sm:p-6"
        >
          <div>
            <p class="text-sm font-medium text-slate-500">Current plan</p>
            <div class="mt-2 flex items-center gap-2">
              <p class="text-2xl font-bold capitalize tracking-tight text-slate-900">
                {{ plan }}
              </p>
              <span
                class="rounded-full px-2 py-0.5 text-xs font-semibold"
                :class="
                  isFree
                    ? 'bg-slate-100 text-slate-600'
                    : 'bg-brand-50 text-brand-700'
                "
              >
                {{ isFree ? 'Limited' : 'Active' }}
              </span>
            </div>
            <p class="mt-2 text-sm leading-6 text-slate-500">
              {{
                isFree
                  ? 'Upgrade for more monthly downloads and higher limits.'
                  : 'Manage billing or change your plan anytime.'
              }}
            </p>
          </div>

          <NuxtLink
            to="/dashboard/subscription"
            class="mt-5 inline-flex w-full items-center justify-center rounded-[10px] px-4 py-2.5 text-sm font-semibold outline-none transition focus-visible:ring-4 focus-visible:ring-red-100 sm:w-fit"
            :class="
              isFree
                ? 'bg-brand-600 text-white hover:bg-brand-700'
                : 'border border-slate-200 bg-white text-slate-800 hover:border-slate-300 hover:bg-slate-50'
            "
          >
            {{ isFree ? 'Explore plans' : 'Manage plan' }}
          </NuxtLink>
        </section>
      </div>

      <!-- Recent downloads -->
      <section class="mt-10" aria-labelledby="recent-downloads-heading">
        <div class="mb-4 flex items-center justify-between gap-3">
          <h2
            id="recent-downloads-heading"
            class="text-lg font-bold tracking-tight text-slate-900 sm:text-xl"
          >
            Recent downloads
          </h2>
          <NuxtLink
            v-if="recentIds.length"
            to="/dashboard/downloads"
            class="rounded text-sm font-semibold text-brand-600 outline-none transition hover:text-brand-700 focus-visible:ring-4 focus-visible:ring-red-100"
          >
            View all
          </NuxtLink>
        </div>

        <div v-if="recentIds.length" class="space-y-3">
          <DownloadStatusCard
            v-for="id in recentIds"
            :key="id"
            :id="id"
          />
        </div>

        <div
          v-else
          class="rounded-2xl border border-dashed border-slate-200 bg-slate-50/80 px-6 py-14 text-center"
        >
          <div
            class="mx-auto flex size-11 items-center justify-center rounded-xl bg-white text-slate-400 shadow-sm ring-1 ring-slate-200"
            aria-hidden="true"
          >
            <svg class="size-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.75">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
            </svg>
          </div>
          <p class="mt-4 font-semibold text-slate-900">No downloads yet</p>
          <p class="mx-auto mt-1 max-w-xs text-sm leading-6 text-slate-500">
            Paste a supported video URL on the home page to create your first download.
          </p>
          <NuxtLink
            to="/"
            class="mt-5 inline-flex rounded-[10px] bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100"
          >
            Create a download
          </NuxtLink>
        </div>
      </section>
    </template>
  </div>
</template>