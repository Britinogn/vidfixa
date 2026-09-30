<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { HugeiconsIcon } from '@hugeicons/vue'
import { ArrowUpRight01Icon } from '@hugeicons/core-free-icons'

const route = useRoute()

const currentYear = new Date().getFullYear()

const productLinks = [
  { label: 'Features', to: '/#features' },
  { label: 'Pricing', to: '/#pricing' },
  { label: 'How it works', to: '/#how-it-works' },
] as const

const companyLinks = [
  { label: 'About', to: '/about' },
] as const

/* Hash link: smooth-scroll if we're already on '/', otherwise navigate. */
function handleHashClick(hash: string, event: MouseEvent) {
  if (route.path !== '/') return
  const el = document.querySelector(hash)
  if (!el) return
  event.preventDefault()
  el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}
</script>

<template>
  <footer class="border-t border-slate-200 bg-white">
    <div class="mx-auto w-full max-w-7xl px-4 py-10 sm:px-5 sm:py-12 md:px-8">
      <div class="flex flex-col gap-10 md:flex-row md:items-start md:justify-between md:gap-12">
        <!-- Brand -->
        <div class="max-w-xs">
          <NuxtLink
            to="/"
            aria-label="VidFixa home"
            class="inline-flex items-center gap-2 rounded-lg outline-none focus-visible:ring-4 focus-visible:ring-red-100"
          >
            <img
              src="/logo.png"
              alt=""
              class="size-12 rounded-xl object-cover"
              width="48"
              height="48"
            />
            <span class="text-[22px] font-extrabold tracking-tight text-slate-900">
              Vid<span class="text-brand-600">Fixa</span>
            </span>
          </NuxtLink>

          <p class="mt-3 text-sm leading-6 text-slate-500">
            A simple way to save supported videos from the web.
          </p>
        </div>

        <!-- Link columns -->
        <div class="grid grid-cols-2 gap-x-10 gap-y-8 sm:grid-cols-3 sm:gap-x-12">
          <!-- Product -->
          <nav aria-labelledby="footer-product">
            <h2
              id="footer-product"
              class="text-xs font-semibold uppercase tracking-wider text-slate-900"
            >
              Product
            </h2>

            <ul class="mt-3 flex flex-col gap-2.5 text-sm">
              <li v-for="link in productLinks" :key="link.to">
                <NuxtLink
                  :to="link.to"
                  class="rounded text-slate-500 outline-none transition hover:text-slate-900 focus-visible:ring-4 focus-visible:ring-red-100"
                  @click="handleHashClick(link.to.replace('/', ''), $event)"
                >
                  {{ link.label }}
                </NuxtLink>
              </li>
            </ul>
          </nav>

          <!-- Company -->
          <nav aria-labelledby="footer-company">
            <h2
              id="footer-company"
              class="text-xs font-semibold uppercase tracking-wider text-slate-900"
            >
              Company
            </h2>

            <ul class="mt-3 flex flex-col gap-2.5 text-sm">
              <li v-for="link in companyLinks" :key="link.to">
                <NuxtLink
                  :to="link.to"
                  class="rounded text-slate-500 outline-none transition hover:text-slate-900 focus-visible:ring-4 focus-visible:ring-red-100"
                >
                  {{ link.label }}
                </NuxtLink>
              </li>

              <li>
                <a
                  href="mailto:brightonwuemeri@gmail.com"
                  class="rounded text-slate-500 outline-none transition hover:text-slate-900 focus-visible:ring-4 focus-visible:ring-red-100"
                >
                  Contact
                </a>
              </li>
            </ul>
          </nav>

          <!-- Builder -->
          <nav aria-labelledby="footer-builder">
            <h2
              id="footer-builder"
              class="text-xs font-semibold uppercase tracking-wider text-slate-900"
            >
              Builder
            </h2>

            <ul class="mt-3 flex flex-col gap-2.5 text-sm">
              <li>
                <NuxtLink
                  to="/about"
                  class="rounded text-slate-500 outline-none transition hover:text-slate-900 focus-visible:ring-4 focus-visible:ring-red-100"
                >
                  Tino Ctemz
                </NuxtLink>
              </li>

              <li>
                <a
                  href="https://x.com/Britinogn"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="group inline-flex items-center gap-1 rounded text-slate-500 outline-none transition hover:text-slate-900 focus-visible:ring-4 focus-visible:ring-red-100"
                >
                  <span>X / @Britinogn</span>
                  <HugeiconsIcon
                    :icon="ArrowUpRight01Icon"
                    :size="12"
                    class="text-slate-400 transition group-hover:text-slate-700"
                    aria-hidden="true"
                  />
                  <span class="sr-only">(opens in a new tab)</span>
                </a>
              </li>
            </ul>
          </nav>
        </div>
      </div>

      <!-- Bottom -->
      <div
        class="mt-10 flex flex-col gap-2 border-t border-slate-100 pt-6 text-xs text-slate-400 sm:flex-row sm:items-center sm:justify-between"
      >
        <span>© {{ currentYear }} VidFixa. All rights reserved.</span>

        <span>
          Built by
          <NuxtLink
            to="/about"
            class="rounded font-medium text-slate-600 outline-none transition hover:text-brand-600 focus-visible:ring-4 focus-visible:ring-red-100"
          >
            Tino Ctemz
          </NuxtLink>
        </span>
      </div>
    </div>
  </footer>
</template>