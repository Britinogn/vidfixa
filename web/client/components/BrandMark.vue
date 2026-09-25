<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { PlayIcon } from '@hugeicons/core-free-icons'

withDefaults(defineProps<{ light?: boolean }>(), { light: false })
</script>

<template>
  <span class="inline-flex items-center gap-2 text-[22px] font-extrabold tracking-tight" :class="light ? 'text-white' : 'text-slate-900'">
    <span class="grid size-9 place-items-center rounded-xl bg-brand-600 text-white"><HugeiconsIcon :icon="PlayIcon" :size="19" /></span>
    Vid<span class="-ml-2 text-brand-600">Fixa</span>
  </span>
</template>
