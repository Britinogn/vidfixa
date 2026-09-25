<script setup lang="ts">
import type { Usage } from '~/types/api'
const props = defineProps<{ usage: Usage; plan: string }>()
const percent = computed(() => props.usage.limit > 0 ? Math.min(100, Math.round(props.usage.used / props.usage.limit * 100)) : 0)
</script>

<template>
  <section class="rounded-2xl border border-slate-200 bg-white p-6 shadow-card">
    <div class="flex items-start justify-between gap-4">
      <div><div class="text-sm font-medium text-slate-500">Monthly usage</div><div class="mt-2 text-3xl font-bold tracking-tight">{{ usage.used }} <span class="text-base font-medium text-slate-400">/ {{ usage.limit }}</span></div></div>
      <span class="rounded-full bg-red-50 px-3 py-1.5 text-xs font-semibold capitalize text-brand-600">{{ plan }} plan</span>
    </div>
    <div class="mt-5 h-2 overflow-hidden rounded-full bg-slate-100"><div class="h-full rounded-full bg-brand-600 transition-all" :style="{ width: `${percent}%` }" /></div>
    <div class="mt-3 flex justify-between text-xs text-slate-500"><span>{{ usage.remaining }} downloads remaining</span><span>{{ percent }}%</span></div>
  </section>
</template>
