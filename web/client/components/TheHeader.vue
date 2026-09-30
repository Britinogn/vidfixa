<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { HugeiconsIcon } from '@hugeicons/vue'
import { MenuIcon, CancelIcon } from '@hugeicons/core-free-icons'
import { useAuthStore } from '~/stores/auth'

type NavLink = {
  label: string
  to: string
  hash?: string
}

const auth = useAuthStore()
const route = useRoute()

const isMenuOpen = ref(false)
const menuId = 'primary-mobile-menu'
const toggleButton = ref<HTMLButtonElement | null>(null)

const isLoggedIn = computed(() => Boolean(auth.user))
const primaryLink = computed(() =>
  auth.user?.role === 'admin' ? '/admin' : isLoggedIn.value ? '/dashboard' : '/register'
)
const primaryLabel = computed(() => (isLoggedIn.value ? 'Open dashboard' : 'Get started'))

const navLinks: NavLink[] =  [
  { label: 'Features', to: '/#features', hash: '#features' },
  { label: 'How it works', to: '/#how-it-works', hash: '#how-it-works' },
  { label: 'Pricing', to: '/#pricing', hash: '#pricing' },
  { label: 'About', to: '/about' },
]

function isActive(to: string, hash?: string) {
  if (hash) return route.path === '/' && route.hash === hash
  return route.path === to || route.path.startsWith(`${to}/`)
}

function closeMenu() {
  isMenuOpen.value = false
}

function toggleMenu() {
  isMenuOpen.value = !isMenuOpen.value
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && isMenuOpen.value) {
    closeMenu()
    toggleButton.value?.focus()
  }
}

watch(isMenuOpen, (open) => {
  if (typeof document === 'undefined') return
  document.body.style.overflow = open ? 'hidden' : ''
})

onMounted(() => document.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeydown)
  if (typeof document !== 'undefined') document.body.style.overflow = ''
})

function handleHashClick(hash: string, event: MouseEvent) {
  if (route.path !== '/') return
  const el = document.querySelector(hash)
  if (!el) return
  event.preventDefault()
  el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  closeMenu()
}

watch(
  () => route.fullPath,
  () => closeMenu()
)
</script>

<template>
  <header class="fixed inset-x-0 top-0 z-50">
    <!-- Bar -->
    <div
      class="border-b border-slate-200/80 bg-[var(--color-background)]/90 backdrop-blur-md supports-[backdrop-filter]:bg-[var(--color-background)]/75"
    >
      <div class="mx-auto flex h-16 w-full max-w-7xl items-center justify-between px-5 lg:h-[4.5rem] lg:px-8">
        <!-- Logo -->
        <NuxtLink
          to="/"
          aria-label="VidFixa home"
          class="flex items-center gap-2.5 rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2"
          @click="closeMenu"
        >
          <img
            src="/logo.png"
            alt=""
            class="size-10 rounded-xl object-cover lg:size-11"
            width="44"
            height="44"
          />
          <span class="text-xl font-extrabold tracking-tight text-slate-900 lg:text-[22px]">
            Vid<span class="text-brand-600">Fixa</span>
          </span>
        </NuxtLink>

        <!-- Desktop nav -->
        <nav
          aria-label="Primary"
          class="hidden items-center gap-1 text-sm font-medium text-slate-600 lg:flex"
        >
          <NuxtLink
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            class="rounded-lg px-3 py-2 transition-colors hover:bg-slate-100 hover:text-slate-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600"
            :class="isActive(link.to, link.hash) ? 'bg-slate-100 text-slate-900' : ''"
            @click="link.hash ? handleHashClick(link.hash, $event) : undefined"
          >
            {{ link.label }}
          </NuxtLink>
        </nav>

        <!-- Desktop actions -->
        <div class="hidden items-center gap-2 lg:flex">
          <NuxtLink
            v-if="!isLoggedIn"
            to="/login"
            class="rounded-lg px-3 py-2 text-sm font-semibold text-slate-700 transition-colors hover:bg-slate-100 hover:text-brand-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600"
          >
            Login
          </NuxtLink>
          <NuxtLink
            :to="primaryLink"
            class="rounded-[10px] bg-brand-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-brand-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2"
          >
            {{ primaryLabel }}
          </NuxtLink>
        </div>

        <!-- Mobile toggle -->
        <button
          ref="toggleButton"
          type="button"
          class="grid size-11 place-items-center rounded-lg text-slate-700 transition-colors hover:bg-slate-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2 lg:hidden"
          :aria-expanded="isMenuOpen"
          :aria-controls="menuId"
          :aria-label="isMenuOpen ? 'Close menu' : 'Open menu'"
          @click="toggleMenu"
        >
          <HugeiconsIcon
            :icon="isMenuOpen ? CancelIcon : MenuIcon"
            :size="24"
            aria-hidden="true"
          />
        </button>
      </div>
    </div>

    <!-- Backdrop -->
    <Transition
      enter-active-class="transition-opacity duration-150 ease-out motion-reduce:transition-none"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-100 ease-in motion-reduce:transition-none"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isMenuOpen"
        class="fixed inset-0 top-16 z-30 bg-slate-900/25 lg:hidden"
        aria-hidden="true"
        @click="closeMenu"
      />
    </Transition>

    <!-- Mobile panel -->
    <Transition
      enter-active-class="transition duration-150 ease-out motion-reduce:transition-none"
      enter-from-class="opacity-0 -translate-y-1"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-100 ease-in motion-reduce:transition-none"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-1"
    >
      <div
        v-if="isMenuOpen"
        :id="menuId"
        class="absolute inset-x-4 top-[4.5rem] z-40 overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-[0_16px_40px_rgba(15,23,42,0.12)] lg:hidden md:inset-x-6"
        role="dialog"
        aria-modal="true"
        aria-label="Navigation menu"
      >
        <nav aria-label="Mobile" class="flex flex-col p-2 text-sm font-medium text-slate-600">
          <NuxtLink
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            class="rounded-xl px-3 py-3 transition-colors hover:bg-slate-50 hover:text-slate-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-inset"
            :class="isActive(link.to, link.hash) ? 'bg-slate-50 text-slate-900' : ''"
            @click="link.hash ? handleHashClick(link.hash, $event) : closeMenu()"
          >
            {{ link.label }}
          </NuxtLink>
        </nav>

        <div class="flex flex-col gap-2 border-t border-slate-100 p-3 sm:flex-row">
          <NuxtLink
            v-if="!isLoggedIn"
            to="/login"
            class="rounded-[10px] border border-slate-200 px-3 py-2.5 text-center text-sm font-semibold text-slate-700 transition-colors hover:bg-slate-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600 sm:flex-1"
            @click="closeMenu"
          >
            Login
          </NuxtLink>
          <NuxtLink
            :to="primaryLink"
            class="rounded-[10px] bg-brand-600 px-3 py-2.5 text-center text-sm font-semibold text-white transition hover:bg-brand-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-600 focus-visible:ring-offset-2 sm:flex-1"
            @click="closeMenu"
          >
            {{ primaryLabel }}
          </NuxtLink>
        </div>
      </div>
    </Transition>
  </header>
</template>