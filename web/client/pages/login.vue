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
const route = useRoute()

const formError = ref('')
const showPassword = ref(false)

const schema = toTypedSchema(
  z.object({
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
    email: '',
    password: '',
  },
})

const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')

const submit = handleSubmit(async (values) => {
  formError.value = ''

  try {
    const user = await auth.login(
      values.email,
      values.password,
    )

    const target =
      typeof route.query.redirect === 'string'
        ? route.query.redirect
        : ''

    if (target.startsWith('/') && !target.startsWith('//')) {
      await navigateTo(target)
      return
    }

    await navigateTo(
      user.role === 'admin'
        ? '/admin'
        : '/dashboard',
    )
  } catch (error) {
    formError.value = getApiErrorMessage(
      error,
      'Could not sign in. Check your details and try again.',
    )
  }
})
</script>

<template>
  <main class="grid min-h-screen lg:grid-cols-[1fr_1fr]">
    <!-- Login form -->
    <section class="flex items-center justify-center px-5 py-10 sm:px-8 sm:py-12">
      <div class="w-full max-w-[420px]">
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
            Welcome back
          </h1>

          <p class="mt-2 text-sm leading-6 text-slate-500 sm:text-base">
            Sign in to continue to VidFixa.
          </p>
        </div>

        <form
          class="mt-8 space-y-5"
          novalidate
          @submit.prevent="submit"
        >
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
              class="mt-2 h-12 w-full rounded-[10px] border border-slate-200 bg-white px-4 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-brand-600 focus:ring-4 focus:ring-red-50 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400"
              placeholder="you@example.com"
            >

            <p
              v-if="errors.email"
              id="login-email-error"
              class="mt-1.5 text-xs text-red-600"
            >
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
                :aria-describedby="errors.password ? 'login-password-error' : undefined"
                class="h-12 w-full rounded-[10px] border border-slate-200 bg-white px-4 pr-12 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-brand-600 focus:ring-4 focus:ring-red-50 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400"
                placeholder="Your password"
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
              v-if="errors.password"
              id="login-password-error"
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
              Signing in…
            </span>

            <span v-else>
              Sign in
            </span>
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
      <div
        class="absolute -right-28 -top-24 size-[500px] rounded-full bg-red-100/60 blur-3xl"
      />

      <div class="relative mx-auto max-w-md">
        <!-- Download preview -->
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

        <!-- Supporting copy -->
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

<!-- ### What changed

* **Eye icon** using your Hugeicons setup.
* **Password show/hide** without changing the field layout.
* **Inputs lock during login** so users can't modify the form while the request is processing.
* **Better API error box** with a subtle border instead of just a red background.
* **Accessible error associations** with `aria-invalid` and `aria-describedby`.
* **Better keyboard focus states** on links and the eye button.
* **Better mobile spacing** without changing the overall design.
* **`novalidate`** so your VeeValidate/Zod messages are the ones users see instead of inconsistent browser-native messages.
* **Email input improvements**: `inputmode`, `autocapitalize`, `spellcheck`.
* **Submit button has a fixed height**, preventing layout movement when `Signing in…` appears.
* **Desktop marketing panel remains clean** rather than adding unnecessary UI.

I would keep the login page this simple. The next UX improvement worth doing is the **forgot-password flow**, once that backend endpoint exists. -->
