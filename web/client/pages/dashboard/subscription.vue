<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { CheckmarkCircle01Icon } from '@hugeicons/core-free-icons'
import { useSubscription } from '~/composables/useSubscription'
import { getApiErrorMessage } from '~/utils/api-error'
import type { PlanTier } from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: [
    function () {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) return navigateTo('/login')
    },
  ],
})

type PaidPlan = Extract<PlanTier, 'plus' | 'pro'>

const route = useRoute()
const { subscription, checkout } = useSubscription()

const pendingTier = ref<PaidPlan | null>(null)

const currentPlan = computed(() => subscription.data.value?.plan ?? 'free')

// A checkout left open elsewhere (another tab, or a return from the
// provider's cancel page) still belongs to this account. Surface it as
// "resume" instead of letting the server's 409 read as a dead end.
const pendingCheckout = computed(() => subscription.data.value?.pending ?? null)
const hasPendingCheckout = computed(() => pendingCheckout.value !== null)

const cancelledNotice = computed(() => route.query.checkout === 'cancelled')

const pendingExpiry = computed(() => {
  const expiresAt = pendingCheckout.value?.expires_at
  if (!expiresAt) return ''
  const d = new Date(expiresAt)
  if (Number.isNaN(d.getTime())) return ''
  // Pinned locale — server (Node on Render) and browser locales differ,
  // which causes a hydration mismatch warning.
  return d.toLocaleString('en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
})

function isCheckingOut(tier: PlanTier) {
  return checkout.isPending.value && pendingTier.value === tier
}

function choosePlan(tier: PaidPlan) {
  if (checkout.isPending.value) return
  pendingTier.value = tier
  checkout.mutate(tier, {
    onSettled: () => {
      pendingTier.value = null
    },
  })
}

function resumePending() {
  const tier = pendingCheckout.value?.tier
  if (tier === 'plus' || tier === 'pro') choosePlan(tier)
}

// Reset transient spinner state on navigation, so a failed attempt doesn't
// leave "Opening checkout…" stuck if the user leaves and comes back.
onBeforeUnmount(() => {
  pendingTier.value = null
})

// Returning from the provider means our cached plan state is stale by
// definition — refresh it instead of showing last visit's data.
onMounted(() => {
  if (route.query.checkout) subscription.refetch()
})

const plans: {
  tier: PlanTier
  price: string
  limit: string
  description: string
}[] = [
  // Naira prices for the future NGN cutover — do not delete.
  // { tier: 'free', price: '₦0', limit: '10 downloads / month', description: 'For occasional downloads.' },
  // { tier: 'plus', price: '₦2,500', limit: '50 downloads / cycle', description: 'For regular creators.' },
  // { tier: 'pro', price: '₦5,000', limit: '200 downloads / cycle', description: 'For power users.' },
  { tier: 'free', price: '$0', limit: '10 downloads / month', description: 'For occasional downloads.' },
  { tier: 'plus', price: '$2', limit: '50 downloads / cycle', description: 'For regular creators.' },
  { tier: 'pro', price: '$4', limit: '200 downloads / cycle', description: 'For power users.' },
]

const errorMessage = computed(() => {
  if (!checkout.error.value) return ''
  return getApiErrorMessage(
    checkout.error.value,
    'Could not start checkout. Please try again.',
  )
})
</script>

<template>
  <div>
    <PageTitle
      title="Subscription"
      description="Review your current plan or choose a plan that fits your needs."
    />

    <!-- Cancelled notice -->
    <p
      v-if="cancelledNotice"
      role="status"
      class="mb-5 rounded-xl border border-amber-100 bg-amber-50 p-4 text-sm text-amber-800"
    >
      That checkout was cancelled — nothing was charged. You can start again
      whenever you're ready.
    </p>

    <!-- Pending checkout banner -->
    <div
      v-if="pendingCheckout"
      class="mb-5 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-brand-600 bg-red-50 p-5"
    >
      <div class="min-w-0">
        <p class="text-sm font-bold capitalize">
          Your {{ pendingCheckout.tier }} checkout is still open
        </p>
        <p class="mt-1 text-sm text-slate-600">
          <span v-if="pendingExpiry">It expires {{ pendingExpiry }}. </span>
          Finish it to activate your plan, or wait for it to expire before
          starting a different one.
        </p>
      </div>
      <button
        type="button"
        :disabled="checkout.isPending.value"
        class="rounded-[10px] bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100 disabled:cursor-not-allowed disabled:opacity-60"
        @click="resumePending"
      >
        {{ checkout.isPending.value ? 'Opening…' : 'Resume checkout' }}
      </button>
    </div>

    <!-- ── Loading skeleton ─────────────────────────────────────────── -->
    <div
      v-if="subscription.isPending.value"
      class="grid gap-5 lg:grid-cols-3"
      aria-busy="true"
      aria-live="polite"
    >
      <span class="sr-only">Loading your plan…</span>

      <article
        v-for="n in 3"
        :key="n"
        class="flex flex-col rounded-2xl border border-slate-200 bg-white p-6 shadow-card"
      >
        <!-- Header row: title + badge -->
        <div class="flex items-center justify-between">
          <div class="h-5 w-16 animate-pulse rounded bg-slate-200" />
          <div class="h-6 w-20 animate-pulse rounded-full bg-slate-100" />
        </div>

        <!-- Description -->
        <div class="mt-3 space-y-2">
          <div class="h-3.5 w-full animate-pulse rounded bg-slate-100" />
          <div class="h-3.5 w-3/4 animate-pulse rounded bg-slate-100" />
        </div>

        <!-- Price -->
        <div class="mt-6 flex items-baseline gap-2">
          <div class="h-8 w-16 animate-pulse rounded bg-slate-200" />
          <div class="h-3.5 w-16 animate-pulse rounded bg-slate-100" />
        </div>

        <!-- Feature line -->
        <div class="mt-5 flex items-center gap-2 border-t border-slate-100 pt-5">
          <div class="size-4.5 animate-pulse rounded-full bg-slate-100" />
          <div class="h-3.5 w-40 animate-pulse rounded bg-slate-100" />
        </div>

        <!-- CTA button -->
        <div class="mt-auto pt-7">
          <div class="h-11 w-full animate-pulse rounded-[10px] bg-slate-100" />
        </div>
      </article>
    </div>

    <!-- ── Plans ────────────────────────────────────────────────────── -->
    <div v-else class="grid gap-5 lg:grid-cols-3">
      <article
        v-for="plan in plans"
        :key="plan.tier"
        class="flex flex-col rounded-2xl border bg-white p-6 shadow-card"
        :class="
          currentPlan === plan.tier
            ? 'border-brand-600 ring-1 ring-brand-600'
            : 'border-slate-200'
        "
        :aria-current="currentPlan === plan.tier ? 'true' : undefined"
      >
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-bold capitalize">{{ plan.tier }}</h2>
          <span
            v-if="currentPlan === plan.tier"
            class="rounded-full bg-red-50 px-2.5 py-1 text-xs font-semibold text-brand-600"
          >
            Current plan
          </span>
        </div>

        <p class="mt-2 text-sm text-slate-500">{{ plan.description }}</p>

        <div class="mt-6 text-3xl font-bold">
          {{ plan.price }}
          <span class="text-sm font-medium text-slate-500"> / month</span>
        </div>

        <div
          class="mt-5 flex items-center gap-2 border-t border-slate-100 pt-5 text-sm text-slate-600"
        >
          <HugeiconsIcon
            :icon="CheckmarkCircle01Icon"
            :size="18"
            class="text-green-600"
            aria-hidden="true"
          />
          {{ plan.limit }}
        </div>

        <!-- Paid plan CTA -->
        <button
          v-if="plan.tier !== 'free'"
          type="button"
          :disabled="
            currentPlan === plan.tier ||
            checkout.isPending.value ||
            hasPendingCheckout
          "
          class="mt-auto pt-7 text-left disabled:cursor-not-allowed"
          @click="choosePlan(plan.tier as PaidPlan)"
        >
          <span
            class="block rounded-[10px] px-4 py-3 text-center text-sm font-semibold outline-none transition focus-visible:ring-4 focus-visible:ring-red-100"
            :class="
              currentPlan === plan.tier
                ? 'bg-slate-100 text-slate-500'
                : 'bg-brand-600 text-white hover:bg-brand-700'
            "
          >
            {{
              currentPlan === plan.tier
                ? 'Current plan'
                : isCheckingOut(plan.tier)
                  ? 'Opening checkout…'
                  : `Choose ${plan.tier}`
            }}
          </span>
        </button>

        <!-- Free plan (no action, just label) -->
        <div v-else class="mt-auto pt-7">
          <span
            v-if="currentPlan === 'free'"
            class="block rounded-[10px] bg-slate-100 px-4 py-3 text-center text-sm font-semibold text-slate-500"
          >
            Current plan
          </span>
          <span
            v-else
            class="block rounded-[10px] border border-slate-200 px-4 py-3 text-center text-sm font-semibold text-slate-600"
          >
            Free tier
          </span>
        </div>
      </article>
    </div>

    <!-- Checkout error -->
    <p
      v-if="errorMessage"
      role="alert"
      class="mt-5 rounded-xl border border-red-100 bg-red-50 p-4 text-sm text-red-700"
    >
      {{ errorMessage }}
    </p>

    <p class="mt-5 text-xs leading-5 text-slate-400">
      A completed payment is confirmed by VidFixa after the payment provider
      notifies the server.
    </p>
  </div>
</template>