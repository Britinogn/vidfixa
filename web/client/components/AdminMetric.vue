<script setup lang="ts">
defineProps<{ label: string; value: number | string; detail: string }>()
</script>

<template>
  <article class="rounded-2xl border border-slate-200 bg-white p-5 shadow-card">
    <div class="text-sm font-medium text-slate-500">{{ label }}</div>
    <div class="mt-3 text-3xl font-bold tracking-tight text-slate-900">{{ value }}</div>
    <div class="mt-2 text-xs text-slate-400">{{ detail }}</div>
  </article>
</template>
