<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  AlertCircleIcon,
  ArrowLeft01Icon,
  CheckmarkCircle01Icon,
  Loading03Icon,
  ViewIcon,
  ViewOffIcon,
} from '@hugeicons/core-free-icons'
import { getApiErrorMessage } from '~/utils/api-error'

useSeoMeta({ robots: 'noindex, nofollow' })

definePageMeta({ layout: false })

const auth = useAuthStore()

const formError = ref('')
const showPassword = ref(false)
const capsLockOn = ref(false)
const fullNameInput = ref<HTMLInputElement | null>(null)

const schema = toTypedSchema(
  z.object({
    fullName: z
      .string()
      .trim()
      .min(5, 'Use at least 5 characters.')
      .max(25, 'Use no more than 25 characters.'),
    email: z
      .string()
      .min(1, 'Email is required.')
      .email('Enter a valid email.'),
    password: z
      .string()
      .min(8, 'Password must be at least 8 characters.'),
    terms: z
      .boolean()
      .refine((v) => v === true, {
        message: 'You must accept the terms to continue.',
      }),
  }),
)

const { defineField, handleSubmit, errors, isSubmitting, validateField, setFieldError } =
  useForm({
    validationSchema: schema,
    initialValues: {
      fullName: '',
      email: '',
      password: '',
      terms: false,
    },
    validateOnMount: false,
  })

const [fullName, fullNameAttrs] = defineField('fullName', {
  validateOnModelUpdate: false,
})
const [email, emailAttrs] = defineField('email', {
  validateOnModelUpdate: false,
})
const [password, passwordAttrs] = defineField('password', {
  validateOnModelUpdate: false,
})
const [terms, termsAttrs] = defineField('terms')

/* Password strength: 4 checks → score 0..4 */
const strengthChecks = computed(() => [
  { id: 'len', label: 'At least 8 characters', ok: (password.value?.length ?? 0) >= 8 },
  { id: 'lower', label: 'Lowercase letter', ok: /[a-z]/.test(password.value ?? '') },
  { id: 'upper', label: 'Uppercase letter', ok: /[A-Z]/.test(password.value ?? '') },
  { id: 'num', label: 'Number or symbol', ok: /[\d\W]/.test(password.value ?? '') },
])

const strengthScore = computed(() => strengthChecks.value.filter((c) => c.ok).length)

const strength = computed(() => {
  const score = strengthScore.value
  if (score <= 1) return { label: 'Weak', color: 'bg-red-500', text: 'text-red-600', width: 'w-1/4' }
  if (score === 2) return { label: 'Fair', color: 'bg-amber-500', text: 'text-amber-600', width: 'w-2/4' }
  if (score === 3) return { label: 'Good', color: 'bg-lime-500', text: 'text-lime-600', width: 'w-3/4' }
  return { label: 'Strong', color: 'bg-green-500', text: 'text-green-600', width: 'w-full' }
})

/* Clear API error on edit */
watch([fullName, email, password], () => {
  if (formError.value) formError.value = ''
})

/* Caps Lock detection */
function onPasswordKey(e: KeyboardEvent) {
  if (typeof e.getModifierState === 'function') {
    capsLockOn.value = e.getModifierState('CapsLock')
  }
}

const submit = handleSubmit(async (values) => {
  formError.value = ''

  try {
    const user = await auth.register(values.fullName, values.email, values.password)
    await navigateTo(user.role === 'admin' ? '/admin' : '/dashboard')
  } catch (error) {
    const message = getApiErrorMessage(
      error,
      'Could not create your account. Please try again.',
    )
    formError.value = message

    // Pin known error types to their fields for a tighter loop.
    if (/email.*(in use|exists|taken)/i.test(message)) {
      setFieldError('email', 'An account with this email already exists.')
    }
  }
})

onMounted(() => {
  fullNameInput.value?.focus()
})
</script>

<template>
  <main class="grid min-h-screen lg:grid-cols-[1fr_1fr]">
    <!-- ─── Register form ──────────────────────────────────────────── -->
    <section class="flex items-center justify-center px-5 py-10 sm:px-8 sm:py-12">
      <div class="w-full max-w-[460px]">
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
            Create your account
          </h1>
          <p class="mt-2 text-sm leading-6 text-slate-500 sm:text-base">
            Start saving your favorite videos with VidFixa.
          </p>
        </div>

        <form class="mt-8 space-y-5" novalidate @submit.prevent="submit">
          <!-- Full name -->
          <div>
            <label
              for="register-full-name"
              class="block text-sm font-medium text-slate-700"
            >
              Full name
            </label>

            <input
              id="register-full-name"
              ref="fullNameInput"
              v-model="fullName"
              v-bind="fullNameAttrs"
              type="text"
              autocomplete="name"
              :disabled="isSubmitting"
              :aria-invalid="Boolean(errors.fullName)"
              :aria-describedby="errors.fullName ? 'register-full-name-error' : undefined"
              class="mt-2 h-12 w-full rounded-[10px] border bg-white px-4 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:ring-4 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400"
              :class="
                errors.fullName
                  ? 'border-red-300 focus:border-red-500 focus:ring-red-50'
                  : 'border-slate-200 focus:border-brand-600 focus:ring-red-50'
              "
              placeholder="Your full name"
              @blur="validateField('fullName')"
            >

            <p
              v-if="errors.fullName"
              id="register-full-name-error"
              class="mt-1.5 flex items-center gap-1.5 text-xs text-red-600"
            >
              <HugeiconsIcon :icon="AlertCircleIcon" :size="14" aria-hidden="true" />
              {{ errors.fullName }}
            </p>
          </div>

          <!-- Email -->
          <div>
            <label
              for="register-email"
              class="block text-sm font-medium text-slate-700"
            >
              Email
            </label>

            <input
              id="register-email"
              v-model="email"
              v-bind="emailAttrs"
              type="email"
              inputmode="email"
              autocomplete="email"
              autocapitalize="none"
              spellcheck="false"
              :disabled="isSubmitting"
              :aria-invalid="Boolean(errors.email)"
              :aria-describedby="errors.email ? 'register-email-error' : undefined"
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
              id="register-email-error"
              class="mt-1.5 flex items-center gap-1.5 text-xs text-red-600"
            >
              <HugeiconsIcon :icon="AlertCircleIcon" :size="14" aria-hidden="true" />
              {{ errors.email }}
            </p>
          </div>

          <!-- Password -->
          <div>
            <label
              for="register-password"
              class="block text-sm font-medium text-slate-700"
            >
              Password
            </label>

            <div class="relative mt-2">
              <input
                id="register-password"
                v-model="password"
                v-bind="passwordAttrs"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="new-password"
                :disabled="isSubmitting"
                :aria-invalid="Boolean(errors.password)"
                :aria-describedby="
                  [errors.password ? 'register-password-error' : '', 'register-password-requirements', capsLockOn ? 'register-password-caps' : '']
                    .filter(Boolean)
                    .join(' ')
                "
                class="h-12 w-full rounded-[10px] border bg-white px-4 pr-12 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:ring-4 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400"
                :class="
                  errors.password
                    ? 'border-red-300 focus:border-red-500 focus:ring-red-50'
                    : 'border-slate-200 focus:border-brand-600 focus:ring-red-50'
                "
                placeholder="Create a password"
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

            <!-- Strength meter (only once the user starts typing) -->
            <div v-if="password" class="mt-2.5">
              <div class="flex items-center gap-2">
                <div class="h-1 flex-1 overflow-hidden rounded-full bg-slate-100">
                  <div
                    class="h-full rounded-full transition-all duration-300 motion-reduce:transition-none"
                    :class="[strength.color, strength.width]"
                  />
                </div>
                <span
                  class="w-14 text-right text-[11px] font-medium"
                  :class="strength.text"
                >
                  {{ strength.label }}
                </span>
              </div>

              <!-- Requirement checklist -->
              <ul
                id="register-password-requirements"
                class="mt-2 grid grid-cols-2 gap-x-3 gap-y-1 text-[11px]"
              >
                <li
                  v-for="check in strengthChecks"
                  :key="check.id"
                  class="flex items-center gap-1.5 transition-colors"
                  :class="check.ok ? 'text-green-600' : 'text-slate-400'"
                >
                  <HugeiconsIcon
                    :icon="check.ok ? CheckmarkCircle01Icon : AlertCircleIcon"
                    :size="12"
                    aria-hidden="true"
                  />
                  {{ check.label }}
                </li>
              </ul>
            </div>

            <!-- Caps Lock warning -->
            <p
              v-if="capsLockOn && !errors.password"
              id="register-password-caps"
              class="mt-1.5 flex items-center gap-1.5 text-xs text-amber-600"
            >
              <HugeiconsIcon :icon="AlertCircleIcon" :size="14" aria-hidden="true" />
              Caps Lock is on.
            </p>

            <p
              v-if="errors.password"
              id="register-password-error"
              class="mt-1.5 flex items-center gap-1.5 text-xs text-red-600"
            >
              <HugeiconsIcon :icon="AlertCircleIcon" :size="14" aria-hidden="true" />
              {{ errors.password }}
            </p>
          </div>

          <!-- Terms -->
          <div>
            <label class="flex cursor-pointer items-start gap-2.5 text-sm text-slate-600">
              <input
                v-model="terms"
                v-bind="termsAttrs"
                type="checkbox"
                :disabled="isSubmitting"
                :aria-invalid="Boolean(errors.terms)"
                :aria-describedby="errors.terms ? 'register-terms-error' : undefined"
                class="mt-0.5 size-4 shrink-0 cursor-pointer rounded border-slate-300 text-brand-600 outline-none transition focus-visible:ring-4 focus-visible:ring-red-50 disabled:cursor-not-allowed"
              >
              <span class="leading-5">
                I agree to the
                <NuxtLink
                  to="/terms"
                  class="font-medium text-brand-600 outline-none hover:text-brand-700 focus-visible:rounded focus-visible:ring-4 focus-visible:ring-red-100"
                  @click.stop
                >
                  Terms of Service
                </NuxtLink>
                and
                <NuxtLink
                  to="/privacy"
                  class="font-medium text-brand-600 outline-none hover:text-brand-700 focus-visible:rounded focus-visible:ring-4 focus-visible:ring-red-100"
                  @click.stop
                >
                  Privacy Policy
                </NuxtLink>.
              </span>
            </label>

            <p
              v-if="errors.terms"
              id="register-terms-error"
              class="mt-1.5 flex items-center gap-1.5 text-xs text-red-600"
            >
              <HugeiconsIcon :icon="AlertCircleIcon" :size="14" aria-hidden="true" />
              {{ errors.terms }}
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
            <span>{{ isSubmitting ? 'Creating account…' : 'Create account' }}</span>
          </button>
        </form>

        <!-- Login -->
        <p class="mt-6 text-center text-sm text-slate-500">
          Already have an account?
          <NuxtLink
            to="/login"
            class="font-semibold text-brand-600 outline-none transition hover:text-brand-700 focus-visible:rounded focus-visible:ring-4 focus-visible:ring-red-100"
          >
            Sign in
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
            <div class="flex items-center justify-between">
              <div>
                <div class="text-sm font-semibold text-slate-900">
                  Your VidFixa library
                </div>
                <div class="mt-1 text-xs text-slate-400">Saved videos</div>
              </div>

              <div class="rounded-full bg-green-50 px-2.5 py-1 text-[11px] font-medium text-green-600">
                Ready
              </div>
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
          Keep your downloads organized.
        </h2>
        <p class="mt-2 leading-6 text-slate-500">
          Create your account to keep track of your downloads and pick up
          where you left off.
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