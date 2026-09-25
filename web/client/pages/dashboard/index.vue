<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import type { DashboardResponse } from '~/types/api'
import { readDownloadIDs } from '~/utils/download-history'

definePageMeta({ layout: 'dashboard' })
const auth = useAuthStore()
const api = useApi()
const dashboard = useQuery({ queryKey: ['dashboard'], queryFn: () => api<DashboardResponse>('/dashboard/'), enabled: computed(() => auth.isAuthenticated) })
const { ids } = useDownloads()
const greetingName = computed(() => dashboard.data.value?.user.full_name.split(' ')[0] || auth.user?.full_name.split(' ')[0] || 'there')
onMounted(() => { if (ids.value.length === 0) ids.value = readDownloadIDs() })
</script>

<template>
  <div><PageTitle :title="`Welcome back, ${greetingName}`" description="Your downloads and plan at a glance." />
    <div v-if="dashboard.isPending.value" class="rounded-2xl border border-slate-200 bg-white p-8 text-sm text-slate-500">Loading your dashboard…</div>
    <div v-else-if="dashboard.isError.value" class="rounded-2xl bg-red-50 p-5 text-sm text-red-700">Could not load your dashboard. Please refresh to try again.</div>
    <template v-else-if="dashboard.data.value">
      <div class="grid gap-5 lg:grid-cols-[1.4fr_1fr]"><UsageCard :usage="dashboard.data.value.usage" :plan="dashboard.data.value.plan"/><section class="flex flex-col justify-between rounded-2xl border border-slate-200 bg-white p-6 shadow-card"><div><div class="text-sm font-medium text-slate-500">Current plan</div><div class="mt-2 text-2xl font-bold capitalize">{{ dashboard.data.value.plan }}</div><p class="mt-2 text-sm leading-6 text-slate-500">{{ dashboard.data.value.plan === 'free' ? 'Upgrade for more monthly downloads.' : 'Manage your plan or start a new billing cycle.' }}</p></div><NuxtLink to="/dashboard/subscription" class="mt-5 inline-flex w-fit items-center rounded-[10px] bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700">{{ dashboard.data.value.plan === 'free' ? 'Explore plans' : 'Manage plan' }}</NuxtLink></section></div>
      <section class="mt-9"><div class="mb-4 flex items-center justify-between"><h2 class="text-xl font-bold">Recent downloads</h2><NuxtLink to="/dashboard/downloads" class="text-sm font-semibold text-brand-600 hover:text-brand-700">View all</NuxtLink></div><div v-if="ids.length" class="space-y-3"><DownloadStatusCard v-for="id in ids.slice(0, 5)" :key="id" :id="id"/></div><div v-else class="rounded-2xl border border-dashed border-slate-300 bg-white px-6 py-12 text-center"><p class="font-semibold">No downloads yet</p><p class="mt-1 text-sm text-slate-500">Paste a video URL to get started.</p><NuxtLink to="/" class="mt-4 inline-flex rounded-[10px] bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-brand-700">Create a download</NuxtLink></div></section>
    </template>
  </div>
</template>
