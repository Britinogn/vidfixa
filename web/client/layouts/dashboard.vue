<script setup lang="ts">
import { HugeiconsIcon } from '@hugeicons/vue'
import { Home01Icon, Download01Icon, CreditCardIcon, UserCircleIcon, Logout01Icon, PlayIcon } from '@hugeicons/core-free-icons'

const auth = useAuthStore()
const route = useRoute()
const isAdmin = computed(() => auth.user?.role === 'admin')
const links = computed(() => isAdmin.value
  ? [{ label: 'Overview', to: '/admin', icon: Home01Icon }]
  : [
      { label: 'Dashboard', to: '/dashboard', icon: Home01Icon },
      { label: 'Downloads', to: '/dashboard/downloads', icon: Download01Icon },
      { label: 'Subscription', to: '/dashboard/subscription', icon: CreditCardIcon },
      { label: 'Account', to: '/dashboard/account', icon: UserCircleIcon },
    ])

async function logout() {
  await auth.logout()
}
</script>

<template>
  <div class="min-h-screen bg-slate-50 md:flex">
    <aside class="w-full border-b border-slate-200 bg-white px-5 py-4 md:fixed md:inset-y-0 md:flex md:w-64 md:flex-col md:border-b-0 md:border-r md:px-6 md:py-7">
      <NuxtLink to="/" class="flex items-center gap-2 text-xl font-extrabold tracking-tight">
        <span class="grid size-9 place-items-center rounded-xl bg-brand-600 text-white"><HugeiconsIcon :icon="PlayIcon" :size="19" /></span>
        Vid<span class="text-brand-600">Fixa</span>
      </NuxtLink>
      <div class="mt-8 hidden text-xs font-semibold uppercase tracking-[0.12em] text-slate-400 md:block">Workspace</div>
      <nav class="mt-4 flex gap-2 overflow-x-auto md:flex-col">
        <NuxtLink v-for="link in links" :key="link.to" :to="link.to" class="flex shrink-0 items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition" :class="route.path === link.to ? 'bg-red-50 text-brand-600' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'">
          <HugeiconsIcon :icon="link.icon" :size="18" />{{ link.label }}
        </NuxtLink>
      </nav>
      <div class="mt-auto hidden border-t border-slate-100 pt-5 md:block">
        <div class="truncate text-sm font-semibold">{{ auth.user?.full_name }}</div>
        <div class="truncate text-xs text-slate-500">{{ auth.user?.email }}</div>
        <button class="mt-4 flex w-full items-center gap-2 rounded-lg px-3 py-2 text-left text-sm text-slate-600 hover:bg-slate-50" @click="logout"><HugeiconsIcon :icon="Logout01Icon" :size="17" /> Log out</button>
      </div>
    </aside>
    <main class="min-w-0 flex-1 md:ml-64">
      <header class="sticky top-0 z-10 flex h-[68px] items-center justify-between border-b border-slate-200 bg-white/95 px-5 backdrop-blur md:px-10">
        <div class="text-sm font-medium text-slate-500">{{ isAdmin ? 'Administration' : 'Your workspace' }}</div>
        <div class="flex items-center gap-3">
          <span class="hidden text-sm font-semibold sm:inline">{{ auth.user?.full_name }}</span>
          <span class="grid size-9 place-items-center rounded-full bg-slate-100 text-sm font-bold text-slate-700">{{ auth.user?.full_name?.slice(0, 1).toUpperCase() }}</span>
        </div>
      </header>
      <div class="mx-auto w-full max-w-7xl p-5 md:p-10"><slot /></div>
    </main>
  </div>
</template>
