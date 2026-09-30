<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: [
    function () {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) return navigateTo('/login')
    },
  ],
})

const auth = useAuthStore()

const initial = computed(
  () => auth.user?.full_name?.trim().charAt(0).toUpperCase() || '?',
)

// Computed so it doesn't re-evaluate on every render, and so we can control
// the format — toLocaleDateString() alone varies by user locale and can cause
// SSR/CSR hydration mismatches.
const memberSince = computed(() => {
  const iso = auth.user?.created_at
  if (!iso) return '—'
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return '—'
  return d.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
})

async function handleLogout() {
  await auth.logout()
  await navigateTo('/login')
}
</script>

<template>
  <div>
    <PageTitle title="Account" description="Your VidFixa profile." />

    <section
      class="max-w-2xl rounded-2xl border border-slate-200 bg-white p-6 shadow-card"
    >
      <!-- Identity -->
      <div class="flex items-center gap-4">
        <span
          aria-hidden="true"
          class="grid size-14 shrink-0 place-items-center rounded-full bg-red-50 text-xl font-bold text-brand-600"
        >
          {{ initial }}
        </span>
        <div class="min-w-0">
          <h2 class="truncate font-semibold">
            {{ auth.user?.full_name || 'Unnamed account' }}
          </h2>
          <p class="mt-1 truncate text-sm text-slate-500">
            {{ auth.user?.email || '—' }}
          </p>
        </div>
      </div>

      <!-- Details -->
      <dl
        class="mt-7 grid gap-4 border-t border-slate-100 pt-6 sm:grid-cols-2"
      >
        <div>
          <dt class="text-xs font-medium uppercase tracking-wide text-slate-400">
            Account role
          </dt>
          <dd class="mt-1 text-sm font-medium capitalize">
            {{ auth.user?.role || '—' }}
          </dd>
        </div>
        <div>
          <dt class="text-xs font-medium uppercase tracking-wide text-slate-400">
            Member since
          </dt>
          <dd class="mt-1 text-sm font-medium">
            {{ memberSince }}
          </dd>
        </div>
      </dl>

      <!-- Actions -->
      <div class="mt-7 border-t border-slate-100 pt-6">
        <button
          type="button"
          class="rounded-[10px] border border-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-700 outline-none transition hover:bg-slate-50 focus-visible:ring-4 focus-visible:ring-red-100"
          @click="handleLogout"
        >
          Log out
        </button>
      </div>
    </section>
  </div>
</template>