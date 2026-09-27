<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { CheckmarkCircle01Icon } from '@hugeicons/core-free-icons'

const auth = useAuthStore()
const primaryLink = computed(() => auth.user?.role === 'admin' ? '/admin' : auth.user ? '/dashboard' : '/register')

const plans = [
  {
    name: 'Free',
    tagline: 'For occasional downloads',
    price: '₦0',
    limit: '10 downloads / month',
    features: [
      'No account required to start',
      'Works instantly, right from the homepage',
      'Create an account to track usage and manage your plan',
    ],
    cta: 'Get started',
    to: undefined, // uses primaryLink
    highlighted: false,
  },
  {
    name: 'Plus',
    tagline: 'For regular creators',
    price: '₦2,500',
    limit: '50 downloads / month',
    features: [
      'Everything in Free',
      '5x the monthly download limit',
      'Download history saved on this device',
    ],
    cta: 'Choose Plus',
    to: '/register',
    highlighted: true,
  },
  {
    name: 'Pro',
    tagline: 'For power users',
    price: '₦5,000',
    limit: '200 downloads / month',
    features: [
      'Everything in Plus',
      '20x the Free monthly limit',
      'Best value per download for heavy use',
    ],
    cta: 'Choose Pro',
    to: '/register',
    highlighted: false,
  },
]
</script>

<template>
  <section id="pricing" class="bg-slate-50 py-16 sm:py-20 md:py-24">
    <div class="mx-auto w-full max-w-5xl px-4 sm:px-5 md:px-8">
      <div class="text-center">
        <span class="inline-flex rounded-full bg-red-50 px-3 py-1.5 text-xs font-semibold text-brand-600">Simple pricing</span>
        <h2 class="mt-5 text-2xl font-bold tracking-tight sm:text-3xl">Start free. Upgrade when ready.</h2>
        <p class="mt-3 text-sm leading-6 text-slate-500 sm:text-base">
          Every plan gets the same fast, no-watermark downloads — the only difference is how many you can make each month.
        </p>
      </div>

      <div class="mx-auto mt-8 grid w-full max-w-3xl gap-4 sm:mt-10 sm:gap-5 md:grid-cols-3">
        <article
          v-for="plan in plans"
          :key="plan.name"
          class="relative flex min-w-0 flex-col rounded-2xl border bg-white p-5 sm:p-6"
          :class="plan.highlighted ? 'border-brand-600 shadow-[0_12px_32px_rgba(225,29,72,.12)]' : 'border-slate-200'"
        >
          <span
            v-if="plan.highlighted"
            class="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-brand-600 px-3 py-1 text-[11px] font-semibold text-white"
          >
            Most popular
          </span>

          <h3 class="font-semibold">{{ plan.name }}</h3>
          <p class="mt-2 text-sm text-slate-500">{{ plan.tagline }}</p>

          <div class="mt-5 text-3xl font-bold">
            {{ plan.price }}<span class="text-sm font-medium text-slate-500"> / month</span>
          </div>
          <p class="mt-2 text-sm font-semibold text-slate-700">{{ plan.limit }}</p>

          <ul class="mt-4 space-y-2.5">
            <li v-for="feature in plan.features" :key="feature" class="flex items-start gap-2 text-sm text-slate-600">
              <HugeiconsIcon :icon="CheckmarkCircle01Icon" :size="16" class="mt-0.5 shrink-0 text-brand-600" />
              <span>{{ feature }}</span>
            </li>
          </ul>

          <NuxtLink
            :to="plan.to || primaryLink"
            class="mt-6 block rounded-[10px] py-3 text-center text-sm font-semibold transition"
            :class="plan.highlighted
              ? 'bg-brand-600 text-white hover:bg-brand-700'
              : 'border border-slate-200 hover:bg-slate-50'"
          >
            {{ plan.cta }}
          </NuxtLink>
        </article>
      </div>

      <div class="mx-auto mt-8 max-w-2xl rounded-2xl border border-slate-200 bg-white p-5 sm:p-6">
        <h3 class="text-sm font-semibold text-slate-800">Good to know</h3>
        <ul class="mt-3 space-y-2 text-xs leading-5 text-slate-500 sm:text-sm">
          <li>• Your download limit resets every month — unused downloads don't carry over.</li>
          <li>• You can start downloading immediately without an account; sign up anytime to keep your history and unlock Plus or Pro.</li>
          <li>• Plus and Pro are billed monthly in Nigerian Naira (₦). Payments are processed securely — we never see or store your card details.</li>
          <li>• Download history is kept in this browser. Clearing your browser data clears it too.</li>
          <li>• Cancel anytime — no long-term contract.</li>
        </ul>
      </div>
    </div>
  </section>
</template>