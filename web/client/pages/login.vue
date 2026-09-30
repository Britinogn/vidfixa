<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  AlertCircleIcon,
  ArrowLeft01Icon,
  Loading03Icon,
  ViewIcon,
  ViewOffIcon,
} from '@hugeicons/core-free-icons'
import { getApiErrorMessage } from '~/utils/api-error'

useSeoMeta({ robots: 'noindex, nofollow' })

definePageMeta({ layout: false })

const auth = useAuthStore()
const route = useRoute()

const formError = ref('')
const showPassword = ref(false)
const capsLockOn = ref(false)
const emailInput = ref<HTMLInputElement | null>(null)

const schema = toTypedSchema(
  z.object({
    email: z
      .string()
      .min(1, 'Email is required.')
      .email('Enter a valid email.'),
    password: z.string().min(8, 'Password must be at least 8 characters.'),
  }),
)

const { defineField, handleSubmit, errors, isSubmitting, validateField, setFieldError } =
  useForm({
    validationSchema: schema,
    initialValues: { email: '', password: '' },
    // validate on blur, re-validate on change once the field has an error
    validateOnMount: false,
  })

const [email, emailAttrs] = defineField('email', {
  validateOnModelUpdate: false,
})
const [password, passwordAttrs] = defineField('password', {
  validateOnModelUpdate: false,
})

/* Clear the API error banner as soon as the user edits anything. */
watch([email, password], () => {
  if (formError.value) formError.value = ''
})

/* Detect Caps Lock while typing in the password field. */
function onPasswordKey(e: KeyboardEvent) {
  if (typeof e.getModifierState === 'function') {
    capsLockOn.value = e.getModifierState('CapsLock')
  }
}

const submit = handleSubmit(async (values) => {
  formError.value = ''

  try {
    const user = await auth.login(values.email, values.password)

    const target =
      typeof route.query.redirect === 'string' ? route.query.redirect : ''

    if (target.startsWith('/') && !target.startsWith('//')) {
      await navigateTo(target)
      return
    }

    await navigateTo(user.role === 'admin' ? '/admin' : '/dashboard')
  } catch (error) {
    const message = getApiErrorMessage(
      error,
      'Could not sign in. Check your details and try again.',
    )
    formError.value = message

    // Attach 401-type errors to the password field for a tighter loop.
    if (/password|credentials|invalid/i.test(message)) {
      setFieldError('password', 'Incorrect email or password.')
    }
  }
})

onMounted(() => {
  // Focus the first empty field on load.
  if (!email.value) emailInput.value?.focus()
})
</script>

<template>
  <main class="grid min-h-screen lg:grid-cols-[1fr_1fr]">
    <!-- ─── Login form ─────────────────────────────────────────────── -->
    <section class="flex items-center justify-center px-5 py-10 sm:px-8 sm:py-12">
      <div class="w-full max-w-[420px]">
        <!-- Brand (logo + wordmark, matching header) -->
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

        <!-- Heading -->
        <div class="mt-10 sm:mt-12">
          <h1 class="text-3xl font-bold tracking-tight text-slate-900">
            Welcome back
          </h1>
          <p class="mt-2 text-sm leading-6 text-slate-500 sm:text-base">
            Sign in to continue to VidFixa.
          </p>
        </div>

        <form class="mt-8 space-y-5" novalidate @submit.prevent="submit">
          <!-- Email -->
          <div>
            <label
              for="login-email"
              class="block text-sm font-medium text-slate-700"
            >
              Email
            </label>

            <input
              id="login-email"
              ref="emailInput"
              v-model="email"
              v-bind="emailAttrs"
              type="email"
              inputmode="email"
              autocomplete="email"
              autocapitalize="none"
              spellcheck="false"
              :disabled="isSubmitting"
              :aria-invalid="Boolean(errors.email)"
              :aria-describedby="errors.email ? 'login-email-error' : undefined"
              class="mt-2 h-12 w-full rounded-[10px] border bg-white px-4 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:ring-4 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400"
              :class="
                errors.email
                  ? 'border-red-300 focus:border-red-500 focus:ring-red-50'
                  : 'border-slate-200 focus:border-brand-600 focus:ring-red-50'
              "
              placeholder="you@example.com"
              @blur="validateField('email')"
            >

            <p
              v-if="errors.email"
              id="login-email-error"
              class="mt-1.5 flex items-center gap-1.5 text-xs text-red-600"
            >
              <HugeiconsIcon :icon="AlertCircleIcon" :size="14" aria-hidden="true" />
              {{ errors.email }}
            </p>
          </div>

          <!-- Password -->
          <div>
            <div class="flex items-center justify-between">
              <label
                for="login-password"
                class="block text-sm font-medium text-slate-700"
              >
                Password
              </label>
              <NuxtLink
                to="/forgot-password"
                class="rounded text-xs font-medium text-slate-500 outline-none transition hover:text-brand-600 focus-visible:ring-4 focus-visible:ring-red-100"
              >
                Forgot password?
              </NuxtLink>
            </div>

            <div class="relative mt-2">
              <input
                id="login-password"
                v-model="password"
                v-bind="passwordAttrs"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                :disabled="isSubmitting"
                :aria-invalid="Boolean(errors.password)"
                :aria-describedby="
                  [errors.password ? 'login-password-error' : '', capsLockOn ? 'login-password-caps' : '']
                    .filter(Boolean)
                    .join(' ') || undefined
                "
                class="h-12 w-full rounded-[10px] border bg-white px-4 pr-12 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:ring-4 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400"
                :class="
                  errors.password
                    ? 'border-red-300 focus:border-red-500 focus:ring-red-50'
                    : 'border-slate-200 focus:border-brand-600 focus:ring-red-50'
                "
                placeholder="Your password"
                @blur="validateField('password')"
                @keyup="onPasswordKey"
                @keydown="onPasswordKey"
              >

              <button
                type="button"
                :disabled="isSubmitting"
                class="absolute right-3 top-1/2 -translate-y-1/2 rounded-md p-1.5 text-slate-400 outline-none transition hover:text-slate-700 focus-visible:ring-4 focus-visible:ring-red-50 disabled:cursor-not-allowed disabled:opacity-50"
                :aria-label="showPassword ? 'Hide password' : 'Show password'"
                :aria-pressed="showPassword"
                @click="showPassword = !showPassword"
              >
                <HugeiconsIcon
                  :icon="showPassword ? ViewOffIcon : ViewIcon"
                  :size="20"
                  :stroke-width="1.8"
                  aria-hidden="true"
                />
              </button>
            </div>

            <!-- Caps Lock warning -->
            <p
              v-if="capsLockOn && !errors.password"
              id="login-password-caps"
              class="mt-1.5 flex items-center gap-1.5 text-xs text-amber-600"
            >
              <HugeiconsIcon :icon="AlertCircleIcon" :size="14" aria-hidden="true" />
              Caps Lock is on.
            </p>

            <p
              v-if="errors.password"
              id="login-password-error"
              class="mt-1.5 flex items-center gap-1.5 text-xs text-red-600"
            >
              <HugeiconsIcon :icon="AlertCircleIcon" :size="14" aria-hidden="true" />
              {{ errors.password }}
            </p>
          </div>

          <!-- API error -->
          <div
            v-if="formError"
            class="flex items-start gap-3 rounded-xl border border-red-100 bg-red-50 px-4 py-3"
            role="alert"
            aria-live="assertive"
          >
            <HugeiconsIcon
              :icon="AlertCircleIcon"
              :size="18"
              class="mt-0.5 shrink-0 text-red-500"
              aria-hidden="true"
            />
            <p class="min-w-0 text-sm leading-5 text-red-700">
              {{ formError }}
            </p>
          </div>

          <!-- Submit -->
          <button
            type="submit"
            :disabled="isSubmitting"
            class="flex h-12 w-full items-center justify-center gap-2 rounded-[10px] bg-brand-600 px-4 text-sm font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100 disabled:cursor-not-allowed disabled:opacity-60"
          >
            <HugeiconsIcon
              v-if="isSubmitting"
              :icon="Loading03Icon"
              :size="18"
              class="animate-spin motion-reduce:animate-none"
              aria-hidden="true"
            />
            <span>{{ isSubmitting ? 'Signing in…' : 'Sign in' }}</span>
          </button>
        </form>

        <!-- Register -->
        <p class="mt-6 text-center text-sm text-slate-500">
          New to VidFixa?
          <NuxtLink
            to="/register"
            class="font-semibold text-brand-600 outline-none transition hover:text-brand-700 focus-visible:rounded focus-visible:ring-4 focus-visible:ring-red-100"
          >
            Create an account
          </NuxtLink>
        </p>

        <!-- Back -->
        <NuxtLink
          to="/"
          class="mx-auto mt-8 flex w-fit items-center gap-1.5 rounded-md text-sm text-slate-400 outline-none transition hover:text-slate-700 focus-visible:ring-4 focus-visible:ring-red-100"
        >
          <HugeiconsIcon :icon="ArrowLeft01Icon" :size="16" aria-hidden="true" />
          Back to home
        </NuxtLink>
      </div>
    </section>

    <!-- ─── Desktop visual ─────────────────────────────────────────── -->
    <aside
      class="relative hidden overflow-hidden bg-slate-50 p-12 lg:flex lg:items-center"
      aria-hidden="true"
    >
      <div
        class="absolute -right-28 -top-24 size-[500px] rounded-full bg-red-100/60 blur-3xl"
      />

      <div class="relative mx-auto max-w-md">
        <div
          class="grid aspect-[1.15] place-items-center rounded-3xl border border-slate-200 bg-white p-8 shadow-card"
        >
          <div class="w-full rounded-2xl border border-slate-200 bg-white p-5">
            <div class="text-sm font-semibold text-slate-900">
              Your downloads, in one place
            </div>

            <div class="mt-5 space-y-3">
              <div
                v-for="n in 3"
                :key="n"
                class="flex items-center gap-3 rounded-xl bg-slate-50 p-3"
              >
                <span class="media-thumb h-10 w-14 shrink-0 rounded-md" />
                <span class="min-w-0 flex-1">
                  <i class="mb-2 block h-2 w-28 max-w-full rounded bg-slate-300" />
                  <i class="block h-1.5 w-16 rounded bg-slate-200" />
                </span>
                <span class="size-2 shrink-0 rounded-full bg-green-500" />
              </div>
            </div>
          </div>
        </div>

        <h2 class="mt-8 text-2xl font-bold tracking-tight text-slate-900">
          Simple, fast video downloads.
        </h2>
        <p class="mt-2 leading-6 text-slate-500">
          Pick up where you left off and keep your downloads organized.
        </p>
      </div>
    </aside>
  </main>
</template>

<style scoped>
.media-thumb {
  background: linear-gradient(
    150deg,
    #bae6fd,
    #38bdf8 45%,
    #0f766e 46%,
    #164e63
  );
}
</style>