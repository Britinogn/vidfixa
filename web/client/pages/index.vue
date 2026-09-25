<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { ZapIcon, Shield01Icon, CircleIcon, InfinityIcon, LaptopIcon, LockKeyIcon, Link01Icon, Settings01Icon, Download01Icon, ArrowRight01Icon, CheckmarkCircle01Icon } from '@hugeicons/core-free-icons'
import { readDownloadIDs } from '~/utils/download-history'

const auth = useAuthStore()
const { ids } = useDownloads()
const features = [
  { title: 'Super fast', detail: 'Download videos in seconds, not minutes.', icon: ZapIcon },
  { title: 'High quality', detail: 'Get the best available quality from the source.', icon: Shield01Icon },
  { title: 'No watermark', detail: 'Your videos, clean and ready to use.', icon: CircleIcon },
  { title: 'Download history', detail: 'Keep track of your recent downloads.', icon: InfinityIcon },
  { title: 'Works everywhere', detail: 'Use VidFixa on desktop, tablet, or mobile.', icon: LaptopIcon },
  { title: 'Secure and private', detail: 'Your account stays protected.', icon: LockKeyIcon },
]
const steps = [
  { title: 'Paste the URL', detail: 'Copy and paste a video link from a supported platform.', icon: Link01Icon },
  { title: 'We process it', detail: 'VidFixa prepares your video for download.', icon: Settings01Icon },
  { title: 'Download and enjoy', detail: 'Save the finished video to your device.', icon: Download01Icon },
]

const primaryLink = computed(() => auth.user?.role === 'admin' ? '/admin' : auth.user ? '/dashboard' : '/register')
const primaryLabel = computed(() => auth.user ? 'Open dashboard' : 'Get started')

onMounted(() => {
  if (ids.value.length === 0) ids.value = readDownloadIDs()
})
</script>

<template>
  <div>
    <header class="mx-auto flex h-[76px] max-w-7xl items-center justify-between px-5 md:px-8">
      <NuxtLink to="/" aria-label="VidFixa home"><BrandMark /></NuxtLink>
      <nav class="hidden items-center gap-10 text-sm font-medium text-slate-600 md:flex"><a href="#features" class="hover:text-slate-900">Features</a><a href="#how-it-works" class="hover:text-slate-900">How it works</a><a href="#pricing" class="hover:text-slate-900">Pricing</a></nav>
      <div class="flex items-center gap-3"><NuxtLink to="/login" class="px-3 py-2 text-sm font-semibold text-slate-700 hover:text-brand-600">Login</NuxtLink><NuxtLink :to="primaryLink" class="rounded-[10px] bg-brand-600 px-5 py-3 text-sm font-semibold text-white shadow-sm hover:bg-brand-700">{{ primaryLabel }}</NuxtLink></div>
    </header>

    <main>
      <section class="relative overflow-hidden">
        <div class="mx-auto grid max-w-7xl items-center gap-12 px-5 pb-20 pt-12 md:min-h-[500px] md:grid-cols-[1.05fr_.95fr] md:px-8 md:pb-24 md:pt-10">
          <div class="relative z-[1]">
            <div class="mb-5 inline-flex items-center gap-2 rounded-full bg-red-50 px-3 py-1.5 text-xs font-semibold text-brand-600"><span class="size-1.5 rounded-full bg-brand-600" />Fast <span>·</span> Secure <span>·</span> No watermark</div>
            <h1 class="max-w-2xl text-[40px] font-extrabold leading-[1.08] tracking-[-0.045em] text-slate-900 sm:text-5xl md:text-[56px]">Download videos.<br>Simple. <span class="text-brand-600">Fast.</span></h1>
            <p class="mt-5 max-w-xl text-base leading-7 text-slate-500 md:text-lg">Save videos from your favorite platforms in just a few seconds. No hassle, no watermark.</p>
            <div class="mt-8 max-w-[610px]"><DownloadForm /></div>
            <div class="mt-7 flex flex-wrap items-center gap-x-7 gap-y-3 text-sm font-medium text-slate-500"><span class="flex items-center gap-2"><i class="platform-dot instagram" />Instagram</span><span class="flex items-center gap-2"><i class="platform-dot facebook" />Facebook</span><span class="flex items-center gap-2"><i class="platform-dot x-platform">X</i>X</span><span class="flex items-center gap-2"><i class="platform-dot linkedin">in</i>LinkedIn</span></div>
            <section v-if="ids.length" class="mt-8 max-w-[610px]" aria-live="polite">
              <div class="mb-3 flex items-center justify-between"><h2 class="text-sm font-semibold text-slate-800">Recent downloads</h2><NuxtLink v-if="auth.user && auth.user.role === 'user'" to="/dashboard/downloads" class="text-xs font-semibold text-brand-600 hover:text-brand-700">Open dashboard</NuxtLink></div>
              <div class="space-y-3"><DownloadStatusCard v-for="id in ids.slice(0, 3)" :key="id" :id="id" /></div>
            </section>
          </div>
          <div class="relative mx-auto w-full max-w-[570px] py-4 md:py-0">
            <div class="absolute -inset-x-8 -inset-y-8 rounded-[45%] bg-red-50/70 blur-2xl" />
            <div class="relative rounded-2xl border border-slate-200 bg-white shadow-[0_24px_70px_rgba(15,23,42,.09)]">
              <div class="flex h-10 items-center gap-1.5 border-b border-slate-100 px-4"><i class="size-2.5 rounded-full bg-red-300"/><i class="size-2.5 rounded-full bg-amber-300"/><i class="size-2.5 rounded-full bg-green-300"/><span class="ml-3 text-xs font-semibold text-slate-600">VidFixa</span></div>
              <div class="grid grid-cols-[108px_1fr] gap-4 p-4 sm:grid-cols-[125px_1fr] sm:p-5">
                <div class="space-y-2 text-[11px] font-medium text-slate-500"><div class="rounded-lg bg-red-50 px-2.5 py-2 text-brand-600">⌂ &nbsp;Home</div><div class="px-2.5 py-2">↓ &nbsp;Downloads</div><div class="px-2.5 py-2">▣ &nbsp;Subscription</div><div class="px-2.5 py-2">◉ &nbsp;Account</div></div>
                <div class="grid grid-cols-[1fr_60px] gap-3"><div class="media-preview flex aspect-[1.65] items-center justify-center rounded-xl border border-slate-200"><span class="grid size-11 place-items-center rounded-full border border-white/60 bg-slate-900/70 text-white"><HugeiconsIcon :icon="ArrowRight01Icon" :size="21" /></span></div><div class="space-y-2"><div class="media-thumb h-[52px] rounded-lg"/><div class="media-thumb media-thumb-two h-[52px] rounded-lg"/><div class="media-thumb media-thumb-three h-[52px] rounded-lg"/></div>
                  <div class="col-span-2 mt-1 rounded-xl border border-slate-200 bg-white p-3"><div class="flex items-center gap-3"><div class="media-thumb h-12 w-16 shrink-0 rounded-lg"/><div class="min-w-0 flex-1"><div class="truncate text-xs font-bold">Beautiful landscape.mp4</div><div class="mt-1 text-[10px] text-slate-500">1080p · Processing</div><div class="mt-2 h-1.5 overflow-hidden rounded-full bg-slate-100"><div class="h-full w-2/3 rounded-full bg-brand-600"/></div><div class="mt-1 text-right text-[10px] text-slate-500">65%</div></div></div></div>
                </div>
              </div>
            </div>
            <div class="relative ml-auto mt-3 flex max-w-[260px] items-center gap-3 rounded-xl border border-slate-200 bg-white px-4 py-3 shadow-card"><span class="grid size-9 place-items-center rounded-lg bg-red-50 text-brand-600"><HugeiconsIcon :icon="CheckmarkCircle01Icon" :size="20"/></span><div><div class="text-xs font-bold">Download ready</div><div class="mt-0.5 text-[11px] text-slate-500">Your video is ready to save.</div></div></div>
          </div>
        </div>
      </section>

      <section id="how-it-works" class="bg-slate-50 py-20 md:py-24">
        <div class="mx-auto max-w-7xl px-5 md:px-8"><div class="text-center"><h2 class="text-3xl font-bold tracking-tight md:text-[32px]">How it works</h2><p class="mt-3 text-slate-500">Get your video in three simple steps.</p></div><div class="mt-14 grid gap-10 md:grid-cols-3 md:gap-6">
          <div v-for="(step, index) in steps" :key="step.title" class="relative flex flex-col items-center text-center"><span class="mb-5 grid size-12 place-items-center rounded-full bg-red-50 text-sm font-bold text-brand-600">0{{ index + 1 }}</span><HugeiconsIcon :icon="step.icon" :size="30" class="text-slate-800"/><h3 class="mt-4 text-lg font-semibold">{{ step.title }}</h3><p class="mt-2 max-w-[250px] text-sm leading-6 text-slate-500">{{ step.detail }}</p><HugeiconsIcon v-if="index < 2" :icon="ArrowRight01Icon" :size="23" class="absolute -right-4 top-16 hidden text-slate-400 md:block"/></div>
        </div></div>
      </section>

      <section id="features" class="py-20 md:py-24"><div class="mx-auto max-w-7xl px-5 md:px-8"><div class="mx-auto max-w-2xl text-center"><span class="rounded-full bg-red-50 px-3 py-1.5 text-xs font-semibold text-brand-600">Why choose VidFixa</span><h2 class="mt-5 text-3xl font-bold tracking-tight md:text-[32px]">Everything you need, in one place.</h2><p class="mt-3 text-slate-500">A straightforward video download experience, built around your workflow.</p></div><div class="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <article v-for="feature in features" :key="feature.title" class="flex gap-4 rounded-2xl border border-slate-200 bg-white p-5"><span class="grid size-11 shrink-0 place-items-center rounded-xl bg-red-50 text-brand-600"><HugeiconsIcon :icon="feature.icon" :size="21"/></span><div><h3 class="text-sm font-semibold">{{ feature.title }}</h3><p class="mt-1.5 text-xs leading-5 text-slate-500">{{ feature.detail }}</p></div></article>
        </div></div></section>

      <section id="pricing" class="bg-slate-50 py-20 md:py-24"><div class="mx-auto max-w-5xl px-5 md:px-8"><div class="text-center"><span class="rounded-full bg-red-50 px-3 py-1.5 text-xs font-semibold text-brand-600">Simple pricing</span><h2 class="mt-5 text-3xl font-bold tracking-tight">Start free. Upgrade when ready.</h2><p class="mt-3 text-slate-500">Choose a plan that fits how you download.</p></div><div class="mx-auto mt-10 grid max-w-3xl gap-5 md:grid-cols-3">
          <article class="rounded-2xl border border-slate-200 bg-white p-6"><h3 class="font-semibold">Free</h3><p class="mt-2 text-sm text-slate-500">For occasional downloads</p><div class="mt-5 text-3xl font-bold">₦0<span class="text-sm font-medium text-slate-500"> / month</span></div><p class="mt-2 text-sm text-slate-500">10 downloads each month</p><NuxtLink :to="primaryLink" class="mt-6 block rounded-[10px] border border-slate-200 py-3 text-center text-sm font-semibold hover:bg-slate-50">Get started</NuxtLink></article>
          <article class="rounded-2xl border border-slate-200 bg-white p-6"><h3 class="font-semibold">Plus</h3><p class="mt-2 text-sm text-slate-500">For regular creators</p><div class="mt-5 text-3xl font-bold">₦2,500<span class="text-sm font-medium text-slate-500"> / month</span></div><p class="mt-2 text-sm text-slate-500">50 downloads per subscription</p><NuxtLink to="/register" class="mt-6 block rounded-[10px] bg-brand-600 py-3 text-center text-sm font-semibold text-white hover:bg-brand-700">Choose Plus</NuxtLink></article>
          <article class="rounded-2xl border border-slate-200 bg-white p-6"><h3 class="font-semibold">Pro</h3><p class="mt-2 text-sm text-slate-500">For power users</p><div class="mt-5 text-3xl font-bold">₦5,000<span class="text-sm font-medium text-slate-500"> / month</span></div><p class="mt-2 text-sm text-slate-500">200 downloads per subscription</p><NuxtLink to="/register" class="mt-6 block rounded-[10px] bg-brand-600 py-3 text-center text-sm font-semibold text-white hover:bg-brand-700">Choose Pro</NuxtLink></article>
        </div><p class="mt-5 text-center text-xs text-slate-400">Plus and Pro checkout require an account.</p></div></section>
    </main>
    <footer class="border-t border-slate-200 bg-white py-8"><div class="mx-auto flex max-w-7xl flex-col items-center justify-between gap-4 px-5 text-sm text-slate-500 sm:flex-row md:px-8"><BrandMark/><span>© {{ new Date().getFullYear() }} VidFixa. Save videos simply.</span><a href="#features" class="hover:text-slate-900">Features</a></div></footer>
  </div>
</template>

<style scoped>
.platform-dot { display: inline-grid; width: 19px; height: 19px; place-items: center; border-radius: 5px; color: white; font-size: 10px; font-style: normal; font-weight: 700; }
.instagram { background: linear-gradient(135deg, #7c3aed, #f43f5e, #f59e0b); }
.facebook { background: #1877f2; border-radius: 50%; }
.x-platform { background: #0f172a; }
.linkedin { background: #0a66c2; border-radius: 3px; font-size: 11px; }
.media-preview { background: linear-gradient(150deg, #bae6fd, #38bdf8 42%, #0f766e 43%, #164e63 67%, #cbd5e1 68%); }
.media-thumb { background: linear-gradient(150deg, #bae6fd, #38bdf8 45%, #0f766e 46%, #164e63); }
.media-thumb-two { background: linear-gradient(150deg, #fde68a, #fb923c 50%, #7c2d12 51%, #1e293b); }
.media-thumb-three { background: linear-gradient(150deg, #cbd5e1, #64748b 50%, #334155 51%, #0f172a); }
</style>
