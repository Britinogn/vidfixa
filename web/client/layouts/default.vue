<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { MenuIcon, CancelIcon } from '@hugeicons/core-free-icons'

const auth = useAuthStore()
const isMenuOpen = ref(false)

const primaryLink = computed(() => auth.user?.role === 'admin' ? '/admin' : auth.user ? '/dashboard' : '/register')
const primaryLabel = computed(() => auth.user ? 'Open dashboard' : 'Get started')

function closeMenu() {
  isMenuOpen.value = false
}
</script>

<template>
  <div class="min-h-screen bg-white text-slate-900">
    <header class="relative mx-auto flex h-[76px] max-w-7xl items-center justify-between px-5 md:px-8">
      <NuxtLink to="/" aria-label="VidFixa home" @click="closeMenu"><BrandMark /></NuxtLink>

      <nav class="hidden items-center gap-10 text-sm font-medium text-slate-600 md:flex">
        <a href="#features" class="hover:text-slate-900">Features</a>
        <a href="#how-it-works" class="hover:text-slate-900">How it works</a>
        <a href="#pricing" class="hover:text-slate-900">Pricing</a>
      </nav>

      <div class="hidden items-center gap-3 md:flex">
        <NuxtLink to="/login" class="px-3 py-2 text-sm font-semibold text-slate-700 hover:text-brand-600">Login</NuxtLink>
        <NuxtLink :to="primaryLink" class="rounded-[10px] bg-brand-600 px-5 py-3 text-sm font-semibold text-white shadow-sm hover:bg-brand-700">{{ primaryLabel }}</NuxtLink>
      </div>

      <button
        type="button"
        class="grid size-10 place-items-center rounded-lg text-slate-700 hover:bg-slate-100 md:hidden"
        :aria-expanded="isMenuOpen"
        aria-label="Toggle menu"
        @click="isMenuOpen = !isMenuOpen"
      >
        <HugeiconsIcon :icon="MenuIcon" :altIcon="CancelIcon" :showAlt="isMenuOpen" :size="24" />
      </button>

      <Transition
        enter-active-class="transition ease-out duration-150"
        enter-from-class="opacity-0 -translate-y-2"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition ease-in duration-100"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 -translate-y-2"
      >
        <div
          v-if="isMenuOpen"
          class="absolute inset-x-4 top-[68px] z-20 rounded-2xl border border-slate-200 bg-white p-4 shadow-[0_16px_40px_rgba(15,23,42,.12)] md:hidden"
        >
          <nav class="flex flex-col text-sm font-medium text-slate-600">
            <a href="#features" class="rounded-lg px-3 py-2.5 hover:bg-slate-50 hover:text-slate-900" @click="closeMenu">Features</a>
            <a href="#how-it-works" class="rounded-lg px-3 py-2.5 hover:bg-slate-50 hover:text-slate-900" @click="closeMenu">How it works</a>
            <a href="#pricing" class="rounded-lg px-3 py-2.5 hover:bg-slate-50 hover:text-slate-900" @click="closeMenu">Pricing</a>
          </nav>
          <div class="mt-3 flex flex-col gap-2 border-t border-slate-100 pt-3">
            <NuxtLink to="/login" class="rounded-[10px] px-3 py-2.5 text-center text-sm font-semibold text-slate-700 hover:bg-slate-50" @click="closeMenu">Login</NuxtLink>
            <NuxtLink :to="primaryLink" class="rounded-[10px] bg-brand-600 px-3 py-2.5 text-center text-sm font-semibold text-white hover:bg-brand-700" @click="closeMenu">{{ primaryLabel }}</NuxtLink>
          </div>
        </div>
      </Transition>
    </header>

    <slot />
  </div>
</template>