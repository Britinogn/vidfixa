<script setup lang="ts">
import { computed } from 'vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import { CheckmarkCircle01Icon } from '@hugeicons/core-free-icons'

const auth = useAuthStore()

const primaryLink = computed(() =>
  auth.user?.role === 'admin' ? '/admin' : auth.user ? '/dashboard' : '/register',
)

/* Logged-out users on the Free plan should land on the homepage — the plan
   itself promises "no account required to start," so pushing them to
   /register contradicts the feature list. Logged-in users still go to the
   dashboard regardless of which plan they click. */
const freeLink = computed(() =>
  auth.user?.role === 'admin'
    ? '/admin'
    : auth.user
      ? '/dashboard'
      : '/',
)

const plans = computed(() => [
  {
    name: 'Free',
    tagline: 'For occasional downloads',
    // Naira price for the future NGN cutover — do not delete.
    // price: '₦0',
    price: '$0',
    limit: '10 downloads / month',
    features: [
      'No account required to start',
      'Works instantly, right from the homepage',
      'Create an account to track usage and manage your plan',
    ],
    cta: 'Start free',
    to: freeLink.value,
    highlighted: false,
  },
  {
    name: 'Plus',
    tagline: 'For regular creators',
    // Naira price for the future NGN cutover — do not delete.
    // price: '₦2,500',
    price: '$2',
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
    // Naira price for the future NGN cutover — do not delete.
    // price: '₦5,000',
    price: '$4',
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
])
</script>

<template>
  <section
    id="pricing"
    aria-labelledby="pricing-heading"
    class="bg-slate-50 py-16 sm:py-20 md:py-24"
  >
    <div class="mx-auto w-full max-w-5xl px-4 sm:px-5 md:max-w-6xl md:px-8">
      <div class="text-center">
        <span
          class="inline-flex rounded-full bg-red-50 px-3 py-1.5 text-xs font-semibold text-brand-600"
        >
          Simple pricing
        </span>
        <h2
          id="pricing-heading"
          class="mt-5 text-2xl font-bold tracking-tight sm:text-3xl md:text-4xl"
        >
          Start free. Upgrade when ready.
        </h2>
        <p
          class="mt-3 text-sm leading-6 text-slate-500 sm:text-base md:mx-auto md:max-w-2xl"
        >
          Every plan gets the same fast, no-watermark downloads — the only
          difference is how many you can make each month.
        </p>
      </div>

      <div
        class="mx-auto mt-8 grid w-full max-w-3xl gap-4 sm:mt-10 sm:gap-5 md:max-w-none md:grid-cols-3 md:items-stretch md:gap-6 lg:gap-8"
      >
        <article
          v-for="plan in plans"
          :key="plan.name"
          class="relative flex min-w-0 flex-col rounded-2xl border bg-white p-5 sm:p-6 md:p-8"
          :class="
            plan.highlighted
              ? 'border-brand-600 shadow-[0_12px_32px_rgba(225,29,72,.12)] md:z-10 md:scale-[1.04] md:shadow-[0_20px_48px_rgba(225,29,72,.16)]'
              : 'border-slate-200 md:shadow-sm'
          "
        >
          <span
            v-if="plan.highlighted"
            class="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-brand-600 px-3 py-1 text-[11px] font-semibold text-white"
          >
            Most popular
          </span>

          <h3 class="text-lg font-semibold">{{ plan.name }}</h3>
          <p class="mt-2 text-sm text-slate-500">{{ plan.tagline }}</p>

          <div class="mt-5 flex items-baseline gap-1 md:mt-6">
            <span class="text-3xl font-bold md:text-4xl">{{ plan.price }}</span>
            <span class="text-sm font-medium text-slate-500">/ month</span>
          </div>
          <p class="mt-2 text-sm font-semibold text-slate-700">
            {{ plan.limit }}
          </p>

          <ul class="mt-5 space-y-2.5">
            <li
              v-for="feature in plan.features"
              :key="feature"
              class="flex items-start gap-2 text-sm text-slate-600"
            >
              <HugeiconsIcon
                :icon="CheckmarkCircle01Icon"
                :size="16"
                class="mt-0.5 shrink-0 text-brand-600"
                aria-hidden="true"
              />
              <span>{{ feature }}</span>
            </li>
          </ul>

          <!-- Wrapper: mt-auto pins the CTA to the bottom, pt-8 guarantees
               a minimum gap above it so it never crowds the feature list. -->
          <div class="mt-auto pt-8">
            <NuxtLink
              :to="plan.to"
              class="block rounded-[10px] py-3 text-center text-sm font-semibold outline-none transition focus-visible:ring-4 focus-visible:ring-red-100"
              :class="
                plan.highlighted
                  ? 'bg-brand-600 text-white hover:bg-brand-700'
                  : 'border border-slate-200 hover:bg-slate-50'
              "
            >
              {{ plan.cta }}
            </NuxtLink>
          </div>
        </article>
      </div>

      <div
        class="mx-auto mt-8 max-w-2xl rounded-2xl border border-slate-200 bg-white p-5 sm:p-6 md:max-w-3xl md:p-8"
      >
        <h3 class="text-sm font-semibold text-slate-800 md:text-base">
          Good to know
        </h3>
        <ul
          class="mt-3 list-disc space-y-2 pl-4 text-xs leading-5 text-slate-500 marker:text-slate-300 sm:text-sm md:grid md:grid-cols-2 md:gap-x-8 md:gap-y-2.5 md:space-y-0"
        >
          <li>
            Your download limit resets every month — unused downloads don't
            carry over.
          </li>
          <li>
            You can start downloading immediately without an account; sign up
            anytime to keep your history and unlock Plus or Pro.
          </li>
          <!-- Naira billing line for the future NGN cutover — do not delete. -->
          <!-- <li>Plus and Pro are billed monthly in Nigerian Naira (₦). Payments are processed securely — we never see or store your card details.</li> -->
          <li>
            Plus and Pro are billed monthly in US Dollars ($). Payments are
            processed securely — we never see or store your card details.
          </li>
          <li>
            Download history is kept in this browser. Clearing your browser
            data clears it too.
          </li>
          <li>Cancel anytime — no long-term contract.</li>
        </ul>
      </div>
    </div>
  </section>
</template>