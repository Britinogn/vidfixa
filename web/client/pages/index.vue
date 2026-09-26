<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'

import {
  ZapIcon,
  Shield01Icon,
  CircleIcon,
  InfinityIcon,
  LaptopIcon,
  LockKeyIcon,
  Link01Icon,
  Settings01Icon,
  Download01Icon,
  ArrowRight01Icon,
  CheckmarkCircle01Icon,
} from '@hugeicons/core-free-icons'

import { readDownloadIDs } from '~/utils/download-history'
//

// useSeoMeta({
//   title: 'Download videos. Simple. Fast.',
//   description:
//     'Paste a link. Download in seconds. Instagram, Facebook, X, LinkedIn — high quality, no watermark.',
// })

// useHead({
//   script: [
//     {
//       type: 'application/ld+json',
//       children: JSON.stringify({
//         '@context': 'https://schema.org',
//         '@type': 'WebApplication',
//         name: 'VidFixa',
//         url: 'https://vidfixa.onrender.com',
//         description:
//           'Download videos from Instagram, Facebook, X, and LinkedIn. Fast, secure, no watermark.',
//         applicationCategory: 'MultimediaApplication',
//         operatingSystem: 'Any',
//         offers: {
//           '@type': 'Offer',
//           price: '0',
//           priceCurrency: 'NGN',
//         },
//       }),
//     },
//   ],
// })

//
const auth = useAuthStore()
const { ids } = useDownloads()

const features = [
  {
    title: 'Super fast',
    detail: 'Download videos in seconds, not minutes.',
    icon: ZapIcon,
  },
  {
    title: 'High quality',
    detail: 'Get the best available quality from the source.',
    icon: Shield01Icon,
  },
  {
    title: 'No watermark',
    detail: 'Your videos, clean and ready to use.',
    icon: CircleIcon,
  },
  {
    title: 'Download history',
    detail: 'Keep track of your recent downloads.',
    icon: InfinityIcon,
  },
  {
    title: 'Works everywhere',
    detail: 'Use VidFixa on desktop, tablet, or mobile.',
    icon: LaptopIcon,
  },
  {
    title: 'Secure and private',
    detail: 'Your account stays protected.',
    icon: LockKeyIcon,
  },
]

const steps = [
  {
    title: 'Paste the URL',
    detail: 'Copy and paste a video link from a supported platform.',
    icon: Link01Icon,
  },
  {
    title: 'We process it',
    detail: 'VidFixa prepares your video for download.',
    icon: Settings01Icon,
  },
  {
    title: 'Download and enjoy',
    detail: 'Save the finished video to your device.',
    icon: Download01Icon,
  },
]

const primaryLink = computed(() =>
  auth.user?.role === 'admin'
    ? '/admin'
    : auth.user
      ? '/dashboard'
      : '/register',
)

onMounted(() => {
  if (ids.value.length === 0) {
    ids.value = readDownloadIDs()
  }
})
</script>

<template>
  <main class="min-w-0 overflow-x-hidden">
    <!-- Hero -->
    <section class="relative overflow-hidden">
      <div
        class="mx-auto grid w-full max-w-7xl items-center gap-10 px-4 pb-16 pt-8 sm:px-5 sm:pb-20 sm:pt-10 md:grid-cols-[1.05fr_.95fr] md:gap-12 md:px-8 md:pb-24 md:pt-10"
      >
        <!-- Hero content -->
        <div class="relative z-[1] min-w-0">
          <div
            class="mb-5 inline-flex max-w-full items-center gap-1.5 rounded-full bg-red-50 px-3 py-1.5 text-[11px] font-semibold text-brand-600 sm:gap-2 sm:text-xs"
          >
            <span class="size-1.5 shrink-0 rounded-full bg-brand-600" />
            <span>Fast</span>
            <span>·</span>
            <span>Secure</span>
            <span>·</span>
            <span>No watermark</span>
          </div>

          <h1
            class="max-w-2xl text-[36px] font-extrabold leading-[1.08] tracking-[-0.045em] text-slate-900 min-[375px]:text-[40px] sm:text-5xl md:text-[56px]"
          >
            Download videos.<br />
            Simple. <span class="text-brand-600">Fast.</span>
          </h1>

          <p
            class="mt-4 max-w-xl text-sm leading-6 text-slate-500 sm:mt-5 sm:text-base sm:leading-7 md:text-lg"
          >
            Save videos from your favorite platforms in just a few seconds.
            No hassle, no watermark.
          </p>

          <!-- Download form -->
          <div class="mt-6 w-full max-w-[610px] sm:mt-8">
            <DownloadForm />
          </div>

          <!-- Platforms -->
          <div
            class="mt-6 flex flex-wrap items-center gap-x-4 gap-y-3 text-xs font-medium text-slate-500 sm:mt-7 sm:gap-x-7 sm:text-sm"
          >
            <span class="flex items-center gap-2">
              <i class="platform-dot instagram"> </i>
              <span>Instagram</span>
            </span>

            <span class="flex items-center gap-2">
              <i class="platform-dot facebook"> </i>
              <span>Facebook</span>
            </span>

            <span class="flex items-center gap-2">
              <i class="platform-dot x-platform">X</i>
              <span>X</span>
            </span>

            <span class="flex items-center gap-2">
              <i class="platform-dot linkedin">in</i>
              <span>LinkedIn</span>
            </span>
          </div>

          <!-- Recent downloads -->
          <section
            v-if="ids.length"
            class="mt-7 w-full max-w-[610px] sm:mt-8"
            aria-live="polite"
          >
            <div class="mb-3 flex min-w-0 items-center justify-between gap-3">
              <h2 class="text-sm font-semibold text-slate-800">
                Recent downloads
              </h2>

              <NuxtLink
                v-if="auth.user && auth.user.role === 'user'"
                to="/dashboard/downloads"
                class="shrink-0 text-xs font-semibold text-brand-600 hover:text-brand-700"
              >
                Open dashboard
              </NuxtLink>
            </div>

            <div class="space-y-3">
              <DownloadStatusCard
                v-for="id in ids.slice(0, 3)"
                :key="id"
                :id="id"
              />
            </div>
          </section>
        </div>

        <!-- Hero mockup -->
        <div
          class="relative mx-auto mt-2 w-full max-w-[570px] min-w-0 py-2 sm:mt-4 sm:py-4 md:mt-0 md:py-0"
        >
          <!-- Glow -->
          <div
            class="absolute -inset-x-3 -inset-y-4 rounded-[35%] bg-red-50/70 blur-2xl sm:-inset-x-6 sm:-inset-y-6 md:-inset-x-8 md:-inset-y-8"
          />

          <!-- Browser -->
          <div
            class="relative min-w-0 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-[0_20px_55px_rgba(15,23,42,.09)] sm:rounded-2xl sm:shadow-[0_24px_70px_rgba(15,23,42,.09)]"
          >
            <!-- Browser header -->
            <div
              class="flex h-9 items-center gap-1.5 border-b border-slate-100 px-3 sm:h-10 sm:px-4"
            >
              <i class="size-2 rounded-full bg-red-300 sm:size-2.5" />
              <i class="size-2 rounded-full bg-amber-300 sm:size-2.5" />
              <i class="size-2 rounded-full bg-green-300 sm:size-2.5" />

              <span
                class="ml-2 truncate text-[10px] font-semibold text-slate-600 sm:ml-3 sm:text-xs"
              >
                VidFixa
              </span>
            </div>

            <!-- Mockup body -->
            <div class="min-w-0 p-3 sm:p-4 md:p-5">
              <!-- Mobile sidebar -->
              <div
                class="mb-3 flex min-w-0 gap-1.5 overflow-x-auto pb-0.5 text-[9px] font-medium text-slate-500 sm:hidden"
              >
                <div
                  class="shrink-0 rounded-lg bg-red-50 px-2.5 py-2 text-brand-600"
                >
                  Home
                </div>

                <div class="shrink-0 rounded-lg px-2.5 py-2">
                  Downloads
                </div>

                <div class="shrink-0 rounded-lg px-2.5 py-2">
                  Subscription
                </div>

                <div class="shrink-0 rounded-lg px-2.5 py-2">
                  Account
                </div>
              </div>

              <!-- Desktop/tablet body -->
              <div class="grid min-w-0 grid-cols-1 gap-3 sm:grid-cols-[92px_minmax(0,1fr)] sm:gap-4 md:grid-cols-[108px_minmax(0,1fr)]">
                <!-- Sidebar -->
                <div
                  class="hidden space-y-1.5 text-[10px] font-medium text-slate-500 sm:block md:space-y-2 md:text-[11px]"
                >
                  <div
                    class="rounded-lg bg-red-50 px-2 py-2 text-brand-600 md:px-2.5"
                  >
                    Home
                  </div>

                  <div class="rounded-lg px-2 py-2 md:px-2.5">
                    Downloads
                  </div>

                  <div class="rounded-lg px-2 py-2 md:px-2.5">
                    Subscription
                  </div>

                  <div class="rounded-lg px-2 py-2 md:px-2.5">
                    Account
                  </div>
                </div>

                <!-- Main mockup -->
                <div class="grid min-w-0 grid-cols-1 gap-3 sm:grid-cols-[minmax(0,1fr)_52px] md:grid-cols-[minmax(0,1fr)_60px]">
                  <!-- Main preview -->
                  <div
                    class="media-preview flex aspect-[1.65] min-h-0 items-center justify-center rounded-lg border border-slate-200 sm:rounded-xl"
                  >
                    <span
                      class="grid size-9 place-items-center rounded-full border border-white/60 bg-slate-900/70 text-white sm:size-11"
                    >
                      <HugeiconsIcon
                        :icon="ArrowRight01Icon"
                        :size="18"
                        class="sm:hidden"
                      />

                      <HugeiconsIcon
                        :icon="ArrowRight01Icon"
                        :size="21"
                        class="hidden sm:block"
                      />
                    </span>
                  </div>

                  <!-- Thumbnails -->
                  <div class="grid grid-cols-3 gap-2 sm:grid-cols-1 sm:space-y-2 sm:gap-0">
                    <div
                      class="media-thumb h-12 rounded-lg sm:h-[44px] md:h-[52px]"
                    />

                    <div
                      class="media-thumb media-thumb-two h-12 rounded-lg sm:h-[44px] md:h-[52px]"
                    />

                    <div
                      class="media-thumb media-thumb-three h-12 rounded-lg sm:h-[44px] md:h-[52px]"
                    />
                  </div>

                  <!-- Processing card -->
                  <div
                    class="col-span-1 mt-0 min-w-0 rounded-xl border border-slate-200 bg-white p-2.5 sm:col-span-2 sm:mt-1 sm:p-3"
                  >
                    <div class="flex min-w-0 items-center gap-2.5 sm:gap-3">
                      <div
                        class="media-thumb h-10 w-14 shrink-0 rounded-lg sm:h-12 sm:w-16"
                      />

                      <div class="min-w-0 flex-1">
                        <div
                          class="truncate text-[11px] font-bold sm:text-xs"
                        >
                          Beautiful landscape.mp4
                        </div>

                        <div class="mt-1 text-[9px] text-slate-500 sm:text-[10px]">
                          1080p · Processing
                        </div>

                        <div
                          class="mt-1.5 h-1.5 overflow-hidden rounded-full bg-slate-100 sm:mt-2"
                        >
                          <div
                            class="h-full w-2/3 rounded-full bg-brand-600"
                          />
                        </div>

                        <div
                          class="mt-1 text-right text-[9px] text-slate-500 sm:text-[10px]"
                        >
                          65%
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Download ready -->
          <div
            class="relative ml-auto mt-2 flex w-full max-w-[260px] items-center gap-2.5 rounded-xl border border-slate-200 bg-white px-3 py-2.5 shadow-card sm:mt-3 sm:gap-3 sm:px-4 sm:py-3"
          >
            <span
              class="grid size-8 shrink-0 place-items-center rounded-lg bg-red-50 text-brand-600 sm:size-9"
            >
              <HugeiconsIcon :icon="CheckmarkCircle01Icon" :size="19" />
            </span>

            <div class="min-w-0">
              <div class="text-[11px] font-bold sm:text-xs">
                Download ready
              </div>

              <div class="mt-0.5 text-[10px] leading-4 text-slate-500 sm:text-[11px]">
                Your video is ready to save.
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section
      id="how-it-works"
      class="bg-slate-50 py-16 sm:py-20 md:py-24"
    >
      <div class="mx-auto w-full max-w-7xl px-4 sm:px-5 md:px-8">
        <div class="text-center">
          <h2
            class="text-2xl font-bold tracking-tight sm:text-3xl md:text-[32px]"
          >
            How it works
          </h2>

          <p class="mt-3 text-sm text-slate-500 sm:text-base">
            Get your video in three simple steps.
          </p>
        </div>

        <div class="mt-10 grid gap-10 sm:mt-14 md:grid-cols-3 md:gap-6">
          <div
            v-for="(step, index) in steps"
            :key="step.title"
            class="relative flex flex-col items-center text-center"
          >
            <span
              class="mb-5 grid size-12 place-items-center rounded-full bg-red-50 text-sm font-bold text-brand-600"
            >
              0{{ index + 1 }}
            </span>

            <HugeiconsIcon
              :icon="step.icon"
              :size="30"
              class="text-slate-800"
            />

            <h3 class="mt-4 text-lg font-semibold">
              {{ step.title }}
            </h3>

            <p
              class="mt-2 max-w-[250px] text-sm leading-6 text-slate-500"
            >
              {{ step.detail }}
            </p>

            <HugeiconsIcon
              v-if="index < 2"
              :icon="ArrowRight01Icon"
              :size="23"
              class="absolute -right-4 top-16 hidden text-slate-400 md:block"
            />
          </div>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section id="features" class="py-16 sm:py-20 md:py-24">
      <div class="mx-auto w-full max-w-7xl px-4 sm:px-5 md:px-8">
        <div class="mx-auto max-w-2xl text-center">
          <span
            class="inline-flex rounded-full bg-red-50 px-3 py-1.5 text-xs font-semibold text-brand-600"
          >
            Why choose VidFixa
          </span>

          <h2
            class="mt-5 text-2xl font-bold tracking-tight sm:text-3xl md:text-[32px]"
          >
            Everything you need, in one place.
          </h2>

          <p class="mt-3 text-sm leading-6 text-slate-500 sm:text-base">
            A straightforward video download experience, built around your
            workflow.
          </p>
        </div>

        <div class="mt-8 grid gap-3 sm:mt-10 sm:grid-cols-2 sm:gap-4 lg:grid-cols-3">
          <article
            v-for="feature in features"
            :key="feature.title"
            class="flex min-w-0 gap-3.5 rounded-2xl border border-slate-200 bg-white p-4 sm:gap-4 sm:p-5"
          >
            <span
              class="grid size-10 shrink-0 place-items-center rounded-xl bg-red-50 text-brand-600 sm:size-11"
            >
              <HugeiconsIcon :icon="feature.icon" :size="21" />
            </span>

            <div class="min-w-0">
              <h3 class="text-sm font-semibold">
                {{ feature.title }}
              </h3>

              <p class="mt-1.5 text-xs leading-5 text-slate-500">
                {{ feature.detail }}
              </p>
            </div>
          </article>
        </div>
      </div>
    </section>

    <!-- Pricing -->
    <section
      id="pricing"
      class="bg-slate-50 py-16 sm:py-20 md:py-24"
    >
      <div class="mx-auto w-full max-w-5xl px-4 sm:px-5 md:px-8">
        <div class="text-center">
          <span
            class="inline-flex rounded-full bg-red-50 px-3 py-1.5 text-xs font-semibold text-brand-600"
          >
            Simple pricing
          </span>

          <h2
            class="mt-5 text-2xl font-bold tracking-tight sm:text-3xl"
          >
            Start free. Upgrade when ready.
          </h2>

          <p class="mt-3 text-sm leading-6 text-slate-500 sm:text-base">
            Choose a plan that fits how you download.
          </p>
        </div>

        <div
          class="mx-auto mt-8 grid w-full max-w-3xl gap-4 sm:mt-10 sm:gap-5 md:grid-cols-3"
        >
          <!-- Free -->
          <article
            class="flex min-w-0 flex-col rounded-2xl border border-slate-200 bg-white p-5 sm:p-6"
          >
            <h3 class="font-semibold">Free</h3>

            <p class="mt-2 text-sm text-slate-500">
              For occasional downloads
            </p>

            <div class="mt-5 text-3xl font-bold">
              ₦0
              <span class="text-sm font-medium text-slate-500">
                / month
              </span>
            </div>

            <p class="mt-2 text-sm text-slate-500">
              10 downloads each month
            </p>

            <NuxtLink
              :to="primaryLink"
              class="mt-6 block rounded-[10px] border border-slate-200 py-3 text-center text-sm font-semibold transition hover:bg-slate-50"
            >
              Get started
            </NuxtLink>
          </article>

          <!-- Plus -->
          <article
            class="flex min-w-0 flex-col rounded-2xl border border-slate-200 bg-white p-5 sm:p-6"
          >
            <h3 class="font-semibold">Plus</h3>

            <p class="mt-2 text-sm text-slate-500">
              For regular creators
            </p>

            <div class="mt-5 text-3xl font-bold">
              ₦2,500
              <span class="text-sm font-medium text-slate-500">
                / month
              </span>
            </div>

            <p class="mt-2 text-sm text-slate-500">
              50 downloads per subscription
            </p>

            <NuxtLink
              to="/register"
              class="mt-6 block rounded-[10px] bg-brand-600 py-3 text-center text-sm font-semibold text-white transition hover:bg-brand-700"
            >
              Choose Plus
            </NuxtLink>
          </article>

          <!-- Pro -->
          <article
            class="flex min-w-0 flex-col rounded-2xl border border-slate-200 bg-white p-5 sm:p-6"
          >
            <h3 class="font-semibold">Pro</h3>

            <p class="mt-2 text-sm text-slate-500">
              For power users
            </p>

            <div class="mt-5 text-3xl font-bold">
              ₦5,000
              <span class="text-sm font-medium text-slate-500">
                / month
              </span>
            </div>

            <p class="mt-2 text-sm text-slate-500">
              200 downloads per subscription
            </p>

            <NuxtLink
              to="/register"
              class="mt-6 block rounded-[10px] bg-brand-600 py-3 text-center text-sm font-semibold text-white transition hover:bg-brand-700"
            >
              Choose Pro
            </NuxtLink>
          </article>
        </div>

        <p class="mt-5 text-center text-xs leading-5 text-slate-400">
          Plus and Pro checkout require an account.
        </p>
      </div>
    </section>
  </main>

  <!-- Footer -->
  <footer class="border-t border-slate-200 bg-white py-7 sm:py-8">
    <div
      class="mx-auto flex w-full max-w-7xl flex-col items-center gap-4 px-4 text-center text-sm text-slate-500 sm:px-5 md:flex-row md:justify-between md:px-8 md:text-left"
    >
      <BrandMark />

      <span class="max-w-full">
        © {{ new Date().getFullYear() }} VidFixa. Save videos simply.
      </span>

      <a
        href="#features"
        class="font-medium hover:text-slate-900"
      >
        Features
      </a>
    </div>
  </footer>
</template>

<style scoped>
.platform-dot {
  display: inline-grid;
  width: 19px;
  height: 19px;
  flex-shrink: 0;
  place-items: center;
  border-radius: 5px;
  color: white;
  font-size: 10px;
  font-style: normal;
  font-weight: 700;
}

.instagram {
  background: linear-gradient(
    135deg,
    #7c3aed,
    #f43f5e,
    #f59e0b
  );
}

.facebook {
  background: #1877f2;
  border-radius: 50%;
}

.x-platform {
  background: #0f172a;
}

.linkedin {
  background: #0a66c2;
  border-radius: 3px;
  font-size: 11px;
}

.media-preview {
  background: linear-gradient(
    150deg,
    #bae6fd,
    #38bdf8 42%,
    #0f766e 43%,
    #164e63 67%,
    #cbd5e1 68%
  );
}

.media-thumb {
  background: linear-gradient(
    150deg,
    #bae6fd,
    #38bdf8 45%,
    #0f766e 46%,
    #164e63
  );
}

.media-thumb-two {
  background: linear-gradient(
    150deg,
    #fde68a,
    #fb923c 50%,
    #7c2d12 51%,
    #1e293b
  );
}

.media-thumb-three {
  background: linear-gradient(
    150deg,
    #cbd5e1,
    #64748b 50%,
    #334155 51%,
    #0f172a
  );
}
</style>