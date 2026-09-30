<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import type {
  AdminDownload,
  AdminOverview,
  AdminPayment,
  AdminSubscription,
  AdminUser,
} from '~/types/api'

definePageMeta({ layout: 'dashboard' })

const auth = useAuthStore()
const api = useApi()
const enabled = computed(() => auth.isAuthenticated && auth.isAdmin)

// Keys carry the admin's id so one admin session can never paint another's
// cached data in the same tab.
const key = computed(() => ['admin', auth.user?.id])

const overview = useQuery({ queryKey: computed(() => [...key.value, 'overview']), queryFn: () => api<AdminOverview>('/admin/overview'), enabled })
const users = useQuery({ queryKey: computed(() => [...key.value, 'users']), queryFn: () => api<AdminUser[]>('/admin/users'), enabled })
const payments = useQuery({ queryKey: computed(() => [...key.value, 'payments']), queryFn: () => api<AdminPayment[]>('/admin/payments'), enabled })
const downloads = useQuery({ queryKey: computed(() => [...key.value, 'downloads']), queryFn: () => api<AdminDownload[]>('/admin/downloads'), enabled })
const subscriptions = useQuery({ queryKey: computed(() => [...key.value, 'subscriptions']), queryFn: () => api<AdminSubscription[]>('/admin/subscriptions'), enabled })

const rowLimit = 25
const metrics = computed(() => overview.data.value)
const refreshing = ref(false)

const money = (amount: string, currency: string) =>
  `${currency} ${Number(amount).toLocaleString(undefined, { minimumFractionDigits: 2 })}`

async function refreshAll() {
  if (refreshing.value) return
  refreshing.value = true
  try {
    await Promise.all([
      overview.refetch(),
      users.refetch(),
      payments.refetch(),
      downloads.refetch(),
      subscriptions.refetch(),
    ])
  } finally {
    refreshing.value = false
  }
}

/* ── Tabs ─────────────────────────────────────────────────────── */
const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'payments', label: 'Payments' },
  { id: 'downloads', label: 'Downloads' },
  { id: 'users', label: 'Users' },
  { id: 'subscriptions', label: 'Subscriptions' },
] as const

type TabId = (typeof tabs)[number]['id']
const activeTab = ref<TabId>('overview')
</script>

<template>
  <div>
    <!-- Header + refresh -->
    <div class="flex flex-wrap items-start justify-between gap-4">
      <PageTitle
        title="Admin overview"
        description="Monitor VidFixa accounts, subscriptions, payments, and downloads."
      />
      <button
        type="button"
        class="inline-flex items-center gap-2 rounded-[10px] border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold text-slate-700 outline-none transition hover:bg-slate-50 focus-visible:ring-4 focus-visible:ring-red-100 disabled:cursor-not-allowed disabled:opacity-60"
        :disabled="refreshing"
        @click="refreshAll"
      >
        <span
          v-if="refreshing"
          class="size-4 animate-spin rounded-full border-2 border-slate-300 border-t-slate-600"
        />
        {{ refreshing ? 'Refreshing…' : 'Refresh' }}
      </button>
    </div>

    <!-- Tabs -->
    <div
      role="tablist"
      aria-label="Admin sections"
      class="mt-6 flex gap-1 overflow-x-auto border-b border-slate-200 pb-px"
    >
      <button
        v-for="tab in tabs"
        :key="tab.id"
        type="button"
        role="tab"
        :id="`admin-tab-${tab.id}`"
        :aria-selected="activeTab === tab.id"
        :aria-controls="`admin-panel-${tab.id}`"
        class="shrink-0 border-b-2 px-4 py-2.5 text-sm font-medium outline-none transition focus-visible:ring-4 focus-visible:ring-red-100"
        :class="
          activeTab === tab.id
            ? 'border-brand-600 text-brand-600'
            : 'border-transparent text-slate-500 hover:text-slate-900'
        "
        @click="activeTab = tab.id"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Global states -->
    <div
      v-if="overview.isPending.value"
      class="mt-6 rounded-2xl border border-slate-200 bg-white p-8 text-sm text-slate-500"
    >
      Loading admin data…
    </div>

    <div
      v-else-if="overview.isError.value"
      class="mt-6 rounded-2xl bg-red-50 p-5 text-sm text-red-700"
    >
      Could not load the overview. Use Refresh to try again.
    </div>

    <template v-else-if="metrics">
      <!-- ── Overview tab ─────────────────────────────────────────── -->
      <div
        v-if="activeTab === 'overview'"
        role="tabpanel"
        :id="`admin-panel-overview`"
        :aria-labelledby="`admin-tab-overview`"
        class="mt-6"
      >
        <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          <AdminMetric
            label="Total users"
            :value="metrics.total_users.toLocaleString()"
            detail="Registered accounts"
          />
          <AdminMetric
            label="Active subscriptions"
            :value="metrics.active_subscriptions.toLocaleString()"
            detail="Active paid access"
          />
          <AdminMetric
            label="Successful payments"
            :value="metrics.successful_payments.toLocaleString()"
            detail="Confirmed by provider webhook"
          />
          <AdminMetric
            label="Failed downloads"
            :value="metrics.failed_downloads.toLocaleString()"
            :detail="`${metrics.total_downloads.toLocaleString()} total downloads`"
          />
        </div>

        <section class="mt-9 grid gap-8 xl:grid-cols-2">
          <div>
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-lg font-bold">Recent payments</h2>
              <span class="text-xs text-slate-400">
                Read only · {{ payments.data.value?.length || 0 }} total
              </span>
            </div>

            <div
              v-if="payments.isError.value"
              class="rounded-2xl bg-red-50 p-4 text-sm text-red-700"
            >
              Payments could not be loaded.
              <button
                class="font-semibold underline"
                @click="payments.refetch()"
              >
                Retry
              </button>
            </div>

            <div
              v-else-if="payments.isPending.value"
              class="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500"
            >
              Loading payments…
            </div>

            <div
              v-else
              class="overflow-hidden rounded-2xl border border-slate-200 bg-white"
            >
              <div
                v-if="payments.data.value?.length"
                class="divide-y divide-slate-100"
              >
                <div
                  v-for="payment in payments.data.value.slice(0, 6)"
                  :key="payment.id"
                  class="flex items-center justify-between gap-3 px-4 py-3.5"
                >
                  <div class="min-w-0">
                    <div class="truncate text-sm font-semibold">
                      {{ payment.user_email }}
                    </div>
                    <div class="mt-1 text-xs text-slate-500">
                      {{ new Date(payment.created_at).toLocaleDateString() }} ·
                      {{ payment.provider_reference }}
                    </div>
                  </div>
                  <div class="shrink-0 text-right">
                    <div class="text-sm font-semibold">
                      {{ money(payment.amount, payment.currency) }}
                    </div>
                    <div class="mt-1">
                      <StatusBadge :status="payment.status" />
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="p-8 text-center text-sm text-slate-500">
                No payment records yet.
              </div>
            </div>
          </div>

          <div>
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-lg font-bold">Recent downloads</h2>
              <span class="text-xs text-slate-400">
                {{ metrics.total_downloads.toLocaleString() }} total
              </span>
            </div>

            <div
              v-if="downloads.isError.value"
              class="rounded-2xl bg-red-50 p-4 text-sm text-red-700"
            >
              Downloads could not be loaded.
              <button
                class="font-semibold underline"
                @click="downloads.refetch()"
              >
                Retry
              </button>
            </div>

            <div
              v-else-if="downloads.isPending.value"
              class="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500"
            >
              Loading downloads…
            </div>

            <div
              v-else
              class="overflow-hidden rounded-2xl border border-slate-200 bg-white"
            >
              <div
                v-if="downloads.data.value?.length"
                class="divide-y divide-slate-100"
              >
                <div
                  v-for="item in downloads.data.value.slice(0, 6)"
                  :key="item.id"
                  class="flex items-center justify-between gap-3 px-4 py-3.5"
                >
                  <div class="min-w-0">
                    <div class="truncate text-sm font-semibold">
                      {{ item.user_email || 'Anonymous visitor' }}
                    </div>
                    <div class="mt-1 truncate text-xs text-slate-500">
                      {{ item.platform }} ·
                      {{ new Date(item.created_at).toLocaleString() }}
                    </div>
                  </div>
                  <StatusBadge :status="item.status" />
                </div>
              </div>
              <div v-else class="p-8 text-center text-sm text-slate-500">
                No downloads yet.
              </div>
            </div>
          </div>
        </section>
      </div>

      <!-- ── Payments tab ─────────────────────────────────────────── -->
      <div
        v-else-if="activeTab === 'payments'"
        role="tabpanel"
        :id="`admin-panel-payments`"
        :aria-labelledby="`admin-tab-payments`"
        class="mt-6"
      >
        <div class="mb-4 flex flex-wrap items-center justify-between gap-2">
          <h2 class="text-lg font-bold">Payments</h2>
          <span class="text-xs text-slate-400">
            Read only · {{ payments.data.value?.length || 0 }} total
          </span>
        </div>

        <div
          v-if="payments.isError.value"
          class="rounded-2xl bg-red-50 p-4 text-sm text-red-700"
        >
          Payments could not be loaded.
          <button class="font-semibold underline" @click="payments.refetch()">
            Retry
          </button>
        </div>

        <div
          v-else-if="payments.isPending.value"
          class="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500"
        >
          Loading payments…
        </div>

        <div
          v-else
          class="overflow-x-auto rounded-2xl border border-slate-200 bg-white"
        >
          <table class="w-full min-w-[620px] text-left text-sm">
            <thead
              class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500"
            >
              <tr>
                <th class="px-4 py-3 font-semibold">User</th>
                <th class="px-4 py-3 font-semibold">Reference</th>
                <th class="px-4 py-3 font-semibold">Amount</th>
                <th class="px-4 py-3 font-semibold">Status</th>
                <th class="px-4 py-3 font-semibold">Date</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr
                v-for="payment in payments.data.value?.slice(0, rowLimit)"
                :key="payment.id"
              >
                <td class="px-4 py-3 font-medium">{{ payment.user_email }}</td>
                <td class="px-4 py-3 text-xs text-slate-500">
                  {{ payment.provider_reference }}
                </td>
                <td class="px-4 py-3 font-semibold">
                  {{ money(payment.amount, payment.currency) }}
                </td>
                <td class="px-4 py-3">
                  <StatusBadge :status="payment.status" />
                </td>
                <td class="px-4 py-3 text-slate-500">
                  {{ new Date(payment.created_at).toLocaleDateString() }}
                </td>
              </tr>
              <tr v-if="!payments.data.value?.length">
                <td colspan="5" class="px-4 py-8 text-center text-slate-500">
                  No payment records yet.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ── Downloads tab ────────────────────────────────────────── -->
      <div
        v-else-if="activeTab === 'downloads'"
        role="tabpanel"
        :id="`admin-panel-downloads`"
        :aria-labelledby="`admin-tab-downloads`"
        class="mt-6"
      >
        <div class="mb-4 flex flex-wrap items-center justify-between gap-2">
          <h2 class="text-lg font-bold">Downloads</h2>
          <span class="text-xs text-slate-400">
            {{ metrics.total_downloads.toLocaleString() }} total
          </span>
        </div>

        <div
          v-if="downloads.isError.value"
          class="rounded-2xl bg-red-50 p-4 text-sm text-red-700"
        >
          Downloads could not be loaded.
          <button class="font-semibold underline" @click="downloads.refetch()">
            Retry
          </button>
        </div>

        <div
          v-else-if="downloads.isPending.value"
          class="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500"
        >
          Loading downloads…
        </div>

        <div
          v-else
          class="overflow-x-auto rounded-2xl border border-slate-200 bg-white"
        >
          <table class="w-full min-w-[620px] text-left text-sm">
            <thead
              class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500"
            >
              <tr>
                <th class="px-4 py-3 font-semibold">User</th>
                <th class="px-4 py-3 font-semibold">Platform</th>
                <th class="px-4 py-3 font-semibold">Status</th>
                <th class="px-4 py-3 font-semibold">Date</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr
                v-for="item in downloads.data.value?.slice(0, rowLimit)"
                :key="item.id"
              >
                <td class="px-4 py-3 font-medium">
                  {{ item.user_email || 'Anonymous visitor' }}
                </td>
                <td class="px-4 py-3 capitalize text-slate-600">
                  {{ item.platform }}
                </td>
                <td class="px-4 py-3">
                  <StatusBadge :status="item.status" />
                </td>
                <td class="px-4 py-3 text-slate-500">
                  {{ new Date(item.created_at).toLocaleString() }}
                </td>
              </tr>
              <tr v-if="!downloads.data.value?.length">
                <td colspan="4" class="px-4 py-8 text-center text-slate-500">
                  No downloads yet.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ── Users tab ────────────────────────────────────────────── -->
      <div
        v-else-if="activeTab === 'users'"
        role="tabpanel"
        :id="`admin-panel-users`"
        :aria-labelledby="`admin-tab-users`"
        class="mt-6"
      >
        <div class="mb-4 flex flex-wrap items-center justify-between gap-2">
          <h2 class="text-lg font-bold">Users</h2>
          <span class="text-xs text-slate-400">
            {{ users.data.value?.length || 0 }} accounts<span
              v-if="(users.data.value?.length || 0) > rowLimit"
            >
              · showing first {{ rowLimit }}</span
            >
          </span>
        </div>

        <div
          v-if="users.isError.value"
          class="rounded-2xl bg-red-50 p-4 text-sm text-red-700"
        >
          Users could not be loaded.
          <button class="font-semibold underline" @click="users.refetch()">
            Retry
          </button>
        </div>

        <div
          v-else-if="users.isPending.value"
          class="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500"
        >
          Loading users…
        </div>

        <div
          v-else
          class="overflow-x-auto rounded-2xl border border-slate-200 bg-white"
        >
          <table class="w-full min-w-[540px] text-left text-sm">
            <thead
              class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500"
            >
              <tr>
                <th class="px-4 py-3 font-semibold">User</th>
                <th class="px-4 py-3 font-semibold">Role</th>
                <th class="px-4 py-3 font-semibold">Joined</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr
                v-for="user in users.data.value?.slice(0, rowLimit)"
                :key="user.id"
              >
                <td class="px-4 py-3">
                  <div class="font-semibold">{{ user.full_name }}</div>
                  <div class="mt-0.5 text-xs text-slate-500">
                    {{ user.email }}
                  </div>
                </td>
                <td class="px-4 py-3 capitalize text-slate-600">
                  {{ user.role }}
                </td>
                <td class="px-4 py-3 text-slate-500">
                  {{ new Date(user.created_at).toLocaleDateString() }}
                </td>
              </tr>
              <tr v-if="!users.data.value?.length">
                <td colspan="3" class="px-4 py-8 text-center text-slate-500">
                  No user records yet.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ── Subscriptions tab ────────────────────────────────────── -->
      <div
        v-else-if="activeTab === 'subscriptions'"
        role="tabpanel"
        :id="`admin-panel-subscriptions`"
        :aria-labelledby="`admin-tab-subscriptions`"
        class="mt-6"
      >
        <div class="mb-4 flex flex-wrap items-center justify-between gap-2">
          <h2 class="text-lg font-bold">Subscriptions</h2>
          <span class="text-xs text-slate-400">
            Read only · {{ subscriptions.data.value?.length || 0 }} total<span
              v-if="(subscriptions.data.value?.length || 0) > rowLimit"
            >
              · showing first {{ rowLimit }}</span
            >
          </span>
        </div>

        <div
          v-if="subscriptions.isError.value"
          class="rounded-2xl bg-red-50 p-4 text-sm text-red-700"
        >
          Subscriptions could not be loaded.
          <button
            class="font-semibold underline"
            @click="subscriptions.refetch()"
          >
            Retry
          </button>
        </div>

        <div
          v-else-if="subscriptions.isPending.value"
          class="rounded-2xl border border-slate-200 bg-white p-6 text-sm text-slate-500"
        >
          Loading subscriptions…
        </div>

        <div
          v-else
          class="overflow-x-auto rounded-2xl border border-slate-200 bg-white"
        >
          <table class="w-full min-w-[620px] text-left text-sm">
            <thead
              class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500"
            >
              <tr>
                <th class="px-4 py-3 font-semibold">Account</th>
                <th class="px-4 py-3 font-semibold">Plan</th>
                <th class="px-4 py-3 font-semibold">Status</th>
                <th class="px-4 py-3 font-semibold">Expires</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr
                v-for="subscription in subscriptions.data.value?.slice(
                  0,
                  rowLimit,
                )"
                :key="subscription.id"
              >
                <td class="px-4 py-3 font-medium">
                  {{ subscription.user_email }}
                </td>
                <td class="px-4 py-3 capitalize">
                  {{ subscription.plan }}
                </td>
                <td class="px-4 py-3">
                  <StatusBadge :status="subscription.status" />
                </td>
                <td class="px-4 py-3 text-slate-500">
                  {{
                    subscription.expires_at
                      ? new Date(subscription.expires_at).toLocaleDateString()
                      : '—'
                  }}
                </td>
              </tr>
              <tr v-if="!subscriptions.data.value?.length">
                <td colspan="4" class="px-4 py-8 text-center text-slate-500">
                  No subscriptions yet.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>