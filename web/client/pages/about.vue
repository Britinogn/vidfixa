<script setup lang="ts">
import { computed } from 'vue'

useSeoMeta({
  title: 'About VidFixa',
  description:
    'Learn about VidFixa, a simple and focused video downloader built to make saving supported videos easier.',
})

const auth = useAuthStore()

/* Same routing logic as the header, footer, and pricing page. Logged-in
   users shouldn't land on /register — that's a dead-end for them. */
const primaryLink = computed(() =>
  auth.user?.role === 'admin'
    ? '/admin'
    : auth.user
      ? '/dashboard'
      : '/register',
)

const primaryLabel = computed(() =>
  auth.user ? 'Open dashboard' : 'Start using VidFixa',
)

/* Smooth-scroll to in-page anchors, matching the header's behavior so the
   scroll feels the same everywhere on the site. */
function scrollTo(id: string, event: MouseEvent) {
  const el = document.querySelector(id)
  if (!el) return
  event.preventDefault()
  el.scrollIntoView({ behavior: 'smooth', block: 'start' })
}
</script>

<template>
  <main class="bg-white">
    <!-- ── Hero ───────────────────────────────────────────────────── -->
    <section class="border-b border-slate-100">
      <div class="mx-auto max-w-5xl px-5 py-16 sm:px-8 sm:py-24">
        <div class="max-w-2xl">
          <p class="text-sm font-semibold tracking-wide text-brand-600">
            About VidFixa
          </p>

          <h1
            class="mt-3 text-4xl font-bold tracking-tight text-slate-900 sm:text-5xl"
          >
            Built to make saving videos simpler.
          </h1>

          <p class="mt-5 text-lg leading-8 text-slate-500">
            VidFixa is a focused video downloader for people who want a
            straightforward way to save supported videos from the web — no
            clutter, no extra steps.
          </p>

          <div class="mt-8 flex flex-wrap gap-3">
            <NuxtLink
              :to="primaryLink"
              class="inline-flex h-12 items-center justify-center rounded-[10px] bg-brand-600 px-6 text-sm font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100"
            >
              {{ primaryLabel }}
            </NuxtLink>

            <a
              href="#why"
              class="inline-flex h-12 items-center justify-center rounded-[10px] border border-slate-200 bg-white px-6 text-sm font-semibold text-slate-700 outline-none transition hover:border-slate-300 hover:text-slate-900 focus-visible:ring-4 focus-visible:ring-red-100"
              @click="scrollTo('#why', $event)"
            >
              Why it exists
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- ── Why + Values ───────────────────────────────────────────── -->
    <section id="why" class="scroll-mt-20">
      <div class="mx-auto max-w-5xl px-5 py-16 sm:px-8 sm:py-20">
        <div
          class="grid gap-14 lg:grid-cols-[1.15fr_0.85fr] lg:items-start lg:gap-16"
        >
          <!-- Story -->
          <div>
            <h2
              class="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl"
            >
              Why VidFixa exists
            </h2>

            <div class="mt-5 space-y-4 text-base leading-7 text-slate-600">
              <p>
                Saving a video shouldn't require a complicated process.
                VidFixa was built around that simple idea.
              </p>
              <p>
                Paste a supported link, get your download. No extra steps, no
                noisy UI — just a clean path from link to file.
              </p>
              <p>
                The product stays focused on what matters: a clear workflow
                that improves as more platforms and features are added.
              </p>
            </div>
          </div>

          <!-- Values -->
          <div class="rounded-2xl border border-slate-200 bg-slate-50 p-6 sm:p-7">
            <h2 class="text-lg font-bold text-slate-900">
              What we care about
            </h2>

            <ul class="mt-6 space-y-5">
              <li class="flex gap-3">
                <span
                  aria-hidden="true"
                  class="mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-lg bg-white text-xs font-semibold text-brand-600 ring-1 ring-slate-200"
                >
                  1
                </span>
                <div>
                  <h3 class="font-semibold text-slate-900">Simplicity</h3>
                  <p class="mt-0.5 text-sm leading-6 text-slate-500">
                    Clear, easy to understand, no clutter.
                  </p>
                </div>
              </li>

              <li class="flex gap-3">
                <span
                  aria-hidden="true"
                  class="mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-lg bg-white text-xs font-semibold text-brand-600 ring-1 ring-slate-200"
                >
                  2
                </span>
                <div>
                  <h3 class="font-semibold text-slate-900">Speed</h3>
                  <p class="mt-0.5 text-sm leading-6 text-slate-500">
                    From supported link to download, fast.
                  </p>
                </div>
              </li>

              <li class="flex gap-3">
                <span
                  aria-hidden="true"
                  class="mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-lg bg-white text-xs font-semibold text-brand-600 ring-1 ring-slate-200"
                >
                  3
                </span>
                <div>
                  <h3 class="font-semibold text-slate-900">Focus</h3>
                  <p class="mt-0.5 text-sm leading-6 text-slate-500">
                    Built around the core experience only.
                  </p>
                </div>
              </li>

              <li class="flex gap-3">
                <span
                  aria-hidden="true"
                  class="mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-lg bg-white text-xs font-semibold text-brand-600 ring-1 ring-slate-200"
                >
                  4
                </span>
                <div>
                  <h3 class="font-semibold text-slate-900">Improvement</h3>
                  <p class="mt-0.5 text-sm leading-6 text-slate-500">
                    Continuously refined as VidFixa grows.
                  </p>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- ── Builder ────────────────────────────────────────────────── -->
    <section class="border-y border-slate-100 bg-slate-50">
      <div class="mx-auto max-w-5xl px-5 py-16 sm:px-8 sm:py-20">
        <div class="max-w-2xl">
          <p class="text-sm font-semibold tracking-wide text-brand-600">
            The builder
          </p>

          <h2 class="mt-2 text-3xl font-bold tracking-tight text-slate-900">
            Built independently by Tino Ctemz.
          </h2>

          <div class="mt-5 space-y-4 text-base leading-7 text-slate-600">
            <p>
              VidFixa is an independent project by Tino Ctemz, a software
              developer focused on practical web products and useful tools.
            </p>
            <p>
              It's treated as a real product — improved through development,
              testing, and feedback.
            </p>
          </div>

          <a
            href="https://x.com/Britinogn"
            target="_blank"
            rel="noopener noreferrer"
            class="mt-7 inline-flex items-center gap-2 rounded-[10px] border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold text-slate-700 outline-none transition hover:border-slate-300 hover:text-slate-900 focus-visible:ring-4 focus-visible:ring-red-100"
          >
            <svg
              class="size-4"
              viewBox="0 0 24 24"
              fill="currentColor"
              aria-hidden="true"
            >
              <path
                d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"
              />
            </svg>
            Follow Tino on X
            <span class="sr-only">(opens in a new tab)</span>
          </a>
        </div>
      </div>
    </section>

    <!-- ── Final CTA ──────────────────────────────────────────────── -->
    <section>
      <div class="mx-auto max-w-5xl px-5 py-16 text-center sm:px-8 sm:py-20">
        <h2
          class="text-3xl font-bold tracking-tight text-slate-900 sm:text-4xl"
        >
          Simple tools are often the most useful.
        </h2>

        <p class="mx-auto mt-4 max-w-xl text-base leading-7 text-slate-500">
          Keep the experience simple, make it useful, and keep building.
          That's the idea behind VidFixa.
        </p>

        <NuxtLink
          :to="primaryLink"
          class="mt-8 inline-flex h-12 items-center justify-center rounded-[10px] bg-brand-600 px-7 text-sm font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100"
        >
          {{ primaryLabel }}
        </NuxtLink>
      </div>
    </section>
  </main>
</template>