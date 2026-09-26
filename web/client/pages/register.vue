```vue
<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'
import { HugeiconsIcon } from '@hugeicons/vue'
import {
  ViewIcon,
  ViewOffIcon,
} from '@hugeicons/core-free-icons'
import { getApiErrorMessage } from '~/utils/api-error'

definePageMeta({ layout: false })

const auth = useAuthStore()

const formError = ref('')
const showPassword = ref(false)

const schema = toTypedSchema(
  z.object({
    fullName: z
      .string()
      .trim()
      .min(5, 'Use at least 5 characters.')
      .max(25, 'Use no more than 25 characters.'),
    email: z.string().email('Enter a valid email.'),
    password: z
      .string()
      .min(8, 'Password must be at least 8 characters.'),
  }),
)

const {
  defineField,
  handleSubmit,
  errors,
  isSubmitting,
} = useForm({
  validationSchema: schema,
  initialValues: {
    fullName: '',
    email: '',
    password: '',
  },
})

const [fullName, fullNameAttrs] = defineField('fullName')
const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')

const passwordLength = computed(() => password.value?.length ?? 0)

const submit = handleSubmit(async (values) => {
  formError.value = ''

  try {
    const user = await auth.register(
      values.fullName,
      values.email,
      values.password,
    )

    await navigateTo(
      user.role === 'admin'
        ? '/admin'
        : '/dashboard',
    )
  } catch (error) {
    formError.value = getApiErrorMessage(
      error,
      'Could not create your account. Please try again.',
    )
  }
})
</script>

<template>
  <main class="grid min-h-screen lg:grid-cols-[1fr_1fr]">
    <!-- Register form -->
    <section
      class="flex items-center justify-center px-5 py-10 sm:px-8 sm:py-12"
    >
      <div class="w-full max-w-[460px]">
        <!-- Brand -->
        <NuxtLink
          to="/"
          class="inline-flex rounded-lg outline-none focus-visible:ring-4 focus-visible:ring-red-100"
        >
          <BrandMark />
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

        <form
          class="mt-8 space-y-5"
          novalidate
          @submit.prevent="submit"
        >
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
              v-model="fullName"
              v-bind="fullNameAttrs"
              type="text"
              autocomplete="name"
              :disabled="isSubmitting"
              :aria-invalid="Boolean(errors.fullName)"
              :aria-describedby="
                errors.fullName
                  ? 'register-full-name-error'
                  : undefined
              "
              class="mt-2 h-12 w-full rounded-[10px] border border-slate-200 bg-white px-4 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-brand-600 focus:ring-4 focus:ring-red-50 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400"
              placeholder="Your full name"
            >

            <p
              v-if="errors.fullName"
              id="register-full-name-error"
              class="mt-1.5 text-xs text-red-600"
            >
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
              :aria-describedby="
                errors.email
                  ? 'register-email-error'
                  : undefined
              "
              class="mt-2 h-12 w-full rounded-[10px] border border-slate-200 bg-white px-4 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-brand-600 focus:ring-4 focus:ring-red-50 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400"
              placeholder="you@example.com"
            >

            <p
              v-if="errors.email"
              id="register-email-error"
              class="mt-1.5 text-xs text-red-600"
            >
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
                  errors.password
                    ? 'register-password-error'
                    : 'register-password-requirement'
                "
                class="h-12 w-full rounded-[10px] border border-slate-200 bg-white px-4 pr-12 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-brand-600 focus:ring-4 focus:ring-red-50 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400"
                placeholder="Create a password"
              >

              <button
                type="button"
                :disabled="isSubmitting"
                class="absolute right-3 top-1/2 -translate-y-1/2 rounded-md p-1.5 text-slate-400 outline-none transition hover:text-slate-700 focus-visible:ring-4 focus-visible:ring-red-50 disabled:cursor-not-allowed disabled:opacity-50"
                :aria-label="
                  showPassword
                    ? 'Hide password'
                    : 'Show password'
                "
                @click="showPassword = !showPassword"
              >
                <HugeiconsIcon
                  :icon="
                    showPassword
                      ? ViewOffIcon
                      : ViewIcon
                  "
                  :size="20"
                  :stroke-width="1.8"
                />
              </button>
            </div>

            <p
              id="register-password-requirement"
              class="mt-2 text-xs"
              :class="
                passwordLength >= 8
                  ? 'text-green-600'
                  : 'text-slate-500'
              "
            >
              {{ passwordLength >= 8 ? '✓' : '•' }}
              At least 8 characters
            </p>

            <p
              v-if="errors.password"
              id="register-password-error"
              class="mt-1.5 text-xs text-red-600"
            >
              {{ errors.password }}
            </p>
          </div>

          <!-- API error -->
          <div
            v-if="formError"
            class="flex items-start gap-3 rounded-xl border border-red-100 bg-red-50 px-4 py-3"
            role="alert"
            aria-live="polite"
          >
            <div class="min-w-0">
              <p class="text-sm leading-5 text-red-700">
                {{ formError }}
              </p>
            </div>
          </div>

          <!-- Submit -->
          <button
            type="submit"
            :disabled="isSubmitting"
            class="flex h-12 w-full items-center justify-center rounded-[10px] bg-brand-600 px-4 text-sm font-semibold text-white outline-none transition hover:bg-brand-700 focus-visible:ring-4 focus-visible:ring-red-100 disabled:cursor-not-allowed disabled:opacity-60"
          >
            <span v-if="isSubmitting">
              Creating account…
            </span>

            <span v-else>
              Create account
            </span>
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
          class="mx-auto mt-8 flex w-fit items-center rounded-md text-sm text-slate-400 outline-none transition hover:text-slate-700 focus-visible:ring-4 focus-visible:ring-red-100"
        >
          <span aria-hidden="true">←</span>
          <span class="ml-1">Back to home</span>
        </NuxtLink>
      </div>
    </section>

    <!-- Desktop visual -->
    <aside
      class="relative hidden overflow-hidden bg-slate-50 p-12 lg:flex lg:items-center"
      aria-hidden="true"
    >
      <!-- Background accent -->
      <div
        class="absolute -right-28 -top-24 size-[500px] rounded-full bg-red-100/60 blur-3xl"
      />

      <div class="relative mx-auto max-w-md">
        <!-- Product preview -->
        <div
          class="grid aspect-[1.15] place-items-center rounded-3xl border border-slate-200 bg-white p-8 shadow-card"
        >
          <div class="w-full rounded-2xl border border-slate-200 bg-white p-5">
            <div class="flex items-center justify-between">
              <div>
                <div class="text-sm font-semibold text-slate-900">
                  Your VidFixa library
                </div>

                <div class="mt-1 text-xs text-slate-400">
                  Saved videos
                </div>
              </div>

              <div
                class="rounded-full bg-green-50 px-2.5 py-1 text-[11px] font-medium text-green-600"
              >
                Ready
              </div>
            </div>

            <div class="mt-5 space-y-3">
              <div
                v-for="n in 3"
                :key="n"
                class="flex items-center gap-3 rounded-xl bg-slate-50 p-3"
              >
                <span
                  class="media-thumb h-10 w-14 shrink-0 rounded-md"
                />

                <span class="min-w-0 flex-1">
                  <i
                    class="mb-2 block h-2 w-28 max-w-full rounded bg-slate-300"
                  />

                  <i
                    class="block h-1.5 w-16 rounded bg-slate-200"
                  />
                </span>

                <span
                  class="size-2 shrink-0 rounded-full bg-green-500"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Register-specific copy -->
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
```

This now gives **login and register the same visual system**, while the right-side content is specific to registration:

* Login: **“Your downloads, in one place”**
* Register: **“Keep your downloads organized.”**
* Same responsive two-column structure
* Same spacing, borders, shadows, focus states, and Hugeicons treatment
* Register keeps its password requirement and validation UX
* Mobile collapses cleanly to the form only
