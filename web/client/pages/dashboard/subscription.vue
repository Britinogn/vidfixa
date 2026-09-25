<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { CheckmarkCircle01Icon } from '@hugeicons/core-free-icons'
import { useSubscription } from '~/composables/useSubscription'
import type { PlanTier } from '~/types/api'

definePageMeta({ layout: 'dashboard' })
const { subscription, checkout } = useSubscription()
const plans: { tier: PlanTier; price: string; limit: string; description: string }[] = [
  { tier: 'free', price: '₦0', limit: '10 downloads / month', description: 'For occasional downloads.' },
  { tier: 'plus', price: '₦2,500', limit: '50 downloads / cycle', description: 'For regular creators.' },
  { tier: 'pro', price: '₦5,000', limit: '200 downloads / cycle', description: 'For power users.' },
]
const errorMessage = computed(() => {
  const error = checkout.error.value as { data?: string; message?: string } | null
  return error?.data || error?.message || ''
})
</script>

<template><div><PageTitle title="Subscription" description="Review your current plan or choose a plan that fits your needs."/><div v-if="subscription.isPending.value" class="text-sm text-slate-500">Loading your plan…</div><div v-else class="grid gap-5 lg:grid-cols-3">
  <article v-for="plan in plans" :key="plan.tier" class="flex flex-col rounded-2xl border bg-white p-6 shadow-card" :class="subscription.data.value?.plan === plan.tier ? 'border-brand-600 ring-1 ring-brand-600' : 'border-slate-200'"><div class="flex items-center justify-between"><h2 class="text-lg font-bold capitalize">{{ plan.tier }}</h2><span v-if="subscription.data.value?.plan === plan.tier" class="rounded-full bg-red-50 px-2.5 py-1 text-xs font-semibold text-brand-600">Current plan</span></div><p class="mt-2 text-sm text-slate-500">{{ plan.description }}</p><div class="mt-6 text-3xl font-bold">{{ plan.price }}<span class="text-sm font-medium text-slate-500"> / month</span></div><div class="mt-5 flex items-center gap-2 border-t border-slate-100 pt-5 text-sm text-slate-600"><HugeiconsIcon :icon="CheckmarkCircle01Icon" :size="18" class="text-green-600"/>{{ plan.limit }}</div>
    <button v-if="plan.tier !== 'free'" :disabled="subscription.data.value?.plan === plan.tier || checkout.isPending.value" class="mt-auto pt-7"><span class="block rounded-[10px] px-4 py-3 text-center text-sm font-semibold" :class="subscription.data.value?.plan === plan.tier ? 'bg-slate-100 text-slate-500' : 'bg-brand-600 text-white hover:bg-brand-700'" @click="checkout.mutate(plan.tier as 'plus' | 'pro')">{{ subscription.data.value?.plan === plan.tier ? 'Current plan' : checkout.isPending.value ? 'Opening checkout…' : `Choose ${plan.tier}` }}</span></button><div v-else class="mt-auto pt-7"><span v-if="subscription.data.value?.plan === 'free'" class="block rounded-[10px] bg-slate-100 px-4 py-3 text-center text-sm font-semibold text-slate-500">Current plan</span><span v-else class="block rounded-[10px] border border-slate-200 px-4 py-3 text-center text-sm font-semibold text-slate-600">Free tier</span></div>
  </article></div><p v-if="errorMessage" class="mt-5 rounded-xl bg-red-50 p-4 text-sm text-red-700">{{ errorMessage }}</p><p class="mt-5 text-xs leading-5 text-slate-400">A completed payment is confirmed by VidFixa after the payment provider notifies the server.</p></div></template>
