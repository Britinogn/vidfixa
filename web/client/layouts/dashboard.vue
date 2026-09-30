<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  Home01Icon,
  Download01Icon,
  CreditCardIcon,
  UserCircleIcon,
  Logout01Icon,
} from '@hugeicons/core-free-icons'

useSeoMeta({ robots: 'noindex, nofollow' })

const auth = useAuthStore()
const route = useRoute()

const isAdmin = computed(() => auth.user?.role === 'admin')

const links = computed(() =>
  isAdmin.value
    ? [{ label: 'Overview', to: '/admin', icon: Home01Icon }]
    : [
        { label: 'Dashboard', to: '/dashboard', icon: Home01Icon },
        { label: 'Downloads', to: '/dashboard/downloads', icon: Download01Icon },
        { label: 'Subscription', to: '/dashboard/subscription', icon: CreditCardIcon },
        { label: 'Account', to: '/dashboard/account', icon: UserCircleIcon },
      ],
)

const initial = computed(
  () => auth.user?.full_name?.trim().charAt(0).toUpperCase() || '?',
)

function isActive(to: string) {
  if (to === '/dashboard' || to === '/admin') return route.path === to
  return route.path === to || route.path.startsWith(`${to}/`)
}

const isUserMenuOpen = ref(false)

async function logout() {
  isUserMenuOpen.value = false
  await auth.logout()
  await navigateTo('/login')
}

function onDocumentClick(e: MouseEvent) {
  if (!isUserMenuOpen.value) return
  const target = e.target as HTMLElement
  if (!target.closest('[data-user-menu]')) {
    isUserMenuOpen.value = false
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') isUserMenuOpen.value = false
}

onMounted(() => {
  document.addEventListener('click', onDocumentClick)
  document.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocumentClick)
  document.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <div class="min-h-screen bg-slate-50 md:flex">
    <!-- ── Sidebar ────────────────────────────────────────────────── -->
    <aside
      class="w-full border-b border-slate-200 bg-white px-5 py-4 md:fixed md:inset-y-0 md:flex md:w-64 md:flex-col md:border-b-0 md:border-r md:px-6 md:py-7"
    >
      <!-- Top row: logo + mobile avatar -->
      <div class="flex items-center justify-between gap-3 md:block">
        <NuxtLink
          to="/"
          aria-label="VidFixa home"
          class="inline-flex items-center gap-2 rounded-lg text-xl font-extrabold tracking-tight outline-none focus-visible:ring-4 focus-visible:ring-red-100"
        >
          <img
            src="/logo.png"
            alt=""
            class="size-12 rounded-xl object-cover"
            width="48"
            height="48"
          />
          <span>
            Vid<span class="text-brand-600">Fixa</span>
          </span>
        </NuxtLink>

        <!-- Mobile user avatar (same row as logo) -->
        <div data-user-menu class="relative md:hidden">
          <button
            type="button"
            class="grid size-10 place-items-center rounded-full bg-slate-100 text-sm font-bold text-slate-700 outline-none transition focus-visible:ring-4 focus-visible:ring-red-100"
            :aria-expanded="isUserMenuOpen"
            aria-haspopup="menu"
            :aria-label="`Account menu for ${auth.user?.full_name || 'user'}`"
            @click="isUserMenuOpen = !isUserMenuOpen"
          >
            {{ initial }}
          </button>

          <Transition
            enter-active-class="transition ease-out duration-150 motion-reduce:transition-none"
            enter-from-class="opacity-0 -translate-y-1"
            enter-to-class="opacity-100 translate-y-0"
            leave-active-class="transition ease-in duration-100 motion-reduce:transition-none"
            leave-from-class="opacity-100 translate-y-0"
            leave-to-class="opacity-0 -translate-y-1"
          >
            <div
              v-if="isUserMenuOpen"
              role="menu"
              class="absolute right-0 top-full z-30 mt-2 w-64 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-[0_16px_40px_rgba(15,23,42,.12)]"
            >
              <div class="border-b border-slate-100 px-4 py-3">
                <div class="truncate text-sm font-semibold">
                  {{ auth.user?.full_name || 'Unnamed account' }}
                </div>
                <div class="truncate text-xs text-slate-500">
                  {{ auth.user?.email || '—' }}
                </div>
              </div>

              <div class="p-1.5">
                <NuxtLink
                  to="/dashboard/account"
                  role="menuitem"
                  class="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-left text-sm text-slate-700 outline-none transition hover:bg-slate-50 focus-visible:ring-4 focus-visible:ring-red-100"
                  @click="isUserMenuOpen = false"
                >
                  <HugeiconsIcon
                    :icon="UserCircleIcon"
                    :size="17"
                    aria-hidden="true"
                  />
                  Account settings
                </NuxtLink>

                <button
                  type="button"
                  role="menuitem"
                  class="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-left text-sm text-red-600 outline-none transition hover:bg-red-50 focus-visible:ring-4 focus-visible:ring-red-100"
                  @click="logout"
                >
                  <HugeiconsIcon
                    :icon="Logout01Icon"
                    :size="17"
                    aria-hidden="true"
                  />
                  Log out
                </button>
              </div>
            </div>
          </Transition>
        </div>
      </div>

      <!-- Workspace label (desktop only) -->
      <div
        class="mt-8 hidden text-xs font-semibold uppercase tracking-[0.12em] text-slate-400 md:block"
      >
        Workspace
      </div>

      <!-- Nav -->
      <nav
        aria-label="Workspace"
        class="mt-4 flex gap-2 overflow-x-auto pb-1 md:mt-3 md:flex-col md:overflow-visible md:pb-0"
      >
        <NuxtLink
          v-for="link in links"
          :key="link.to"
          :to="link.to"
          :aria-current="isActive(link.to) ? 'page' : undefined"
          class="flex shrink-0 items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium outline-none transition focus-visible:ring-4 focus-visible:ring-red-100"
          :class="
            isActive(link.to)
              ? 'bg-red-50 text-brand-600'
              : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
          "
        >
          <HugeiconsIcon :icon="link.icon" :size="18" aria-hidden="true" />
          {{ link.label }}
        </NuxtLink>
      </nav>

      <!-- Desktop user card — pinned to bottom -->
      <div
        class="mt-auto hidden rounded-2xl border border-slate-100 bg-slate-50/60 p-3 md:block"
      >
        <div class="flex items-center gap-3">
          <span
            aria-hidden="true"
            class="grid size-10 shrink-0 place-items-center rounded-full bg-white text-sm font-bold text-slate-700 ring-1 ring-slate-200"
          >
            {{ initial }}
          </span>
          <div class="min-w-0 flex-1">
            <div class="truncate text-sm font-semibold text-slate-900">
              {{ auth.user?.full_name || 'Unnamed account' }}
            </div>
            <div class="truncate text-xs text-slate-500">
              {{ auth.user?.email || '—' }}
            </div>
          </div>
        </div>

        <button
          type="button"
          class="mt-3 flex w-full items-center gap-2 rounded-lg border border-slate-200 bg-white px-3 py-2 text-left text-sm font-medium text-slate-600 outline-none transition hover:bg-slate-50 hover:text-slate-900 focus-visible:ring-4 focus-visible:ring-red-100"
          @click="logout"
        >
          <HugeiconsIcon :icon="Logout01Icon" :size="17" aria-hidden="true" />
          Log out
        </button>
      </div>
    </aside>

    <!-- ── Main ───────────────────────────────────────────────────── -->
    <main class="min-w-0 flex-1 md:ml-64">
      <header
        class="sticky top-0 z-10 flex h-[68px] items-center justify-between border-b border-slate-200 bg-white/95 px-5 backdrop-blur md:px-10"
      >
        <div class="hidden text-sm font-medium text-slate-500 sm:block">
          {{ isAdmin ? 'Administration' : 'Your workspace' }}
        </div>

        <!-- Desktop user menu -->
        <div data-user-menu class="relative ml-auto hidden md:block">
          <button
            type="button"
            class="flex items-center gap-3 rounded-full outline-none transition hover:bg-slate-50 focus-visible:ring-4 focus-visible:ring-red-100"
            :aria-expanded="isUserMenuOpen"
            aria-haspopup="menu"
            :aria-label="`Account menu for ${auth.user?.full_name || 'user'}`"
            @click="isUserMenuOpen = !isUserMenuOpen"
          >
            <span class="truncate text-sm font-semibold">
              {{ auth.user?.full_name || 'Account' }}
            </span>
            <span
              aria-hidden="true"
              class="grid size-9 shrink-0 place-items-center rounded-full bg-slate-100 text-sm font-bold text-slate-700"
            >
              {{ initial }}
            </span>
          </button>

          <Transition
            enter-active-class="transition ease-out duration-150 motion-reduce:transition-none"
            enter-from-class="opacity-0 -translate-y-1"
            enter-to-class="opacity-100 translate-y-0"
            leave-active-class="transition ease-in duration-100 motion-reduce:transition-none"
            leave-from-class="opacity-100 translate-y-0"
            leave-to-class="opacity-0 -translate-y-1"
          >
            <div
              v-if="isUserMenuOpen"
              role="menu"
              class="absolute right-0 top-full z-20 mt-2 w-64 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-[0_16px_40px_rgba(15,23,42,.12)]"
            >
              <div class="border-b border-slate-100 px-4 py-3">
                <div class="truncate text-sm font-semibold">
                  {{ auth.user?.full_name || 'Unnamed account' }}
                </div>
                <div class="truncate text-xs text-slate-500">
                  {{ auth.user?.email || '—' }}
                </div>
              </div>

              <div class="p-1.5">
                <NuxtLink
                  to="/dashboard/account"
                  role="menuitem"
                  class="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-left text-sm text-slate-700 outline-none transition hover:bg-slate-50 focus-visible:ring-4 focus-visible:ring-red-100"
                  @click="isUserMenuOpen = false"
                >
                  <HugeiconsIcon
                    :icon="UserCircleIcon"
                    :size="17"
                    aria-hidden="true"
                  />
                  Account settings
                </NuxtLink>

                <button
                  type="button"
                  role="menuitem"
                  class="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-left text-sm text-red-600 outline-none transition hover:bg-red-50 focus-visible:ring-4 focus-visible:ring-red-100"
                  @click="logout"
                >
                  <HugeiconsIcon
                    :icon="Logout01Icon"
                    :size="17"
                    aria-hidden="true"
                  />
                  Log out
                </button>
              </div>
            </div>
          </Transition>
        </div>
      </header>

      <div class="mx-auto w-full max-w-7xl p-5 md:p-10">
        <slot />
      </div>
    </main>
  </div>
</template>