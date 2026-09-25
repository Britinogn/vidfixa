<script setup lang="ts">
defineProps<{ title: string; description?: string }>()
</script>

<template>
  <div class="mb-8">
    <h1 class="text-2xl font-bold tracking-tight text-slate-900 md:text-[30px]">{{ title }}</h1>
    <p v-if="description" class="mt-2 text-sm text-slate-500 md:text-base">{{ description }}</p>
  </div>
</template>
