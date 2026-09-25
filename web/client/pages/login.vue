<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'

definePageMeta({ layout: false })
const auth = useAuthStore()
const route = useRoute()
const formError = ref('')
const schema = toTypedSchema(z.object({ email: z.string().email('Enter a valid email.'), password: z.string().min(8, 'Password must be at least 8 characters.') }))
const { defineField, handleSubmit, errors, isSubmitting } = useForm({ validationSchema: schema, initialValues: { email: '', password: '' } })
const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')

const submit = handleSubmit(async (values) => {
  formError.value = ''
  try {
    const user = await auth.login(values.email, values.password)
    const target = typeof route.query.redirect === 'string' ? route.query.redirect : ''
    if (target.startsWith('/') && !target.startsWith('//')) await navigateTo(target)
    else await navigateTo(user.role === 'admin' ? '/admin' : '/dashboard')
  } catch (error) {
    const response = error as { data?: string; message?: string }
    formError.value = response.data || response.message || 'Could not sign in. Check your details and try again.'
  }
})
</script>

<template>
  <main class="grid min-h-screen lg:grid-cols-[1fr_1fr]">
    <section class="flex items-center justify-center px-5 py-12"><div class="w-full max-w-[420px]"><NuxtLink to="/"><BrandMark /></NuxtLink><div class="mt-12"><h1 class="text-3xl font-bold tracking-tight">Welcome back</h1><p class="mt-2 text-slate-500">Sign in to continue to VidFixa.</p></div>
      <form class="mt-8 space-y-5" @submit.prevent="submit"><label class="block text-sm font-medium text-slate-700">Email<input v-model="email" v-bind="emailAttrs" type="email" autocomplete="email" class="mt-2 h-12 w-full rounded-[10px] border border-slate-200 px-4 outline-none focus:border-brand-600 focus:ring-4 focus:ring-red-50" placeholder="you@example.com"><span v-if="errors.email" class="mt-1 block text-xs text-red-600">{{ errors.email }}</span></label><label class="block text-sm font-medium text-slate-700">Password<input v-model="password" v-bind="passwordAttrs" type="password" autocomplete="current-password" class="mt-2 h-12 w-full rounded-[10px] border border-slate-200 px-4 outline-none focus:border-brand-600 focus:ring-4 focus:ring-red-50" placeholder="Your password"><span v-if="errors.password" class="mt-1 block text-xs text-red-600">{{ errors.password }}</span></label><p v-if="formError" class="rounded-xl bg-red-50 p-3 text-sm text-red-700" role="alert">{{ formError }}</p><button :disabled="isSubmitting" class="w-full rounded-[10px] bg-brand-600 py-3.5 text-sm font-semibold text-white hover:bg-brand-700 disabled:opacity-60">{{ isSubmitting ? 'Signing in…' : 'Sign in' }}</button></form>
      <p class="mt-6 text-center text-sm text-slate-500">New to VidFixa? <NuxtLink to="/register" class="font-semibold text-brand-600 hover:text-brand-700">Create an account</NuxtLink></p><NuxtLink to="/" class="mt-8 block text-center text-sm text-slate-400 hover:text-slate-700">← Back to home</NuxtLink>
    </div></section>
    <aside class="relative hidden overflow-hidden bg-slate-50 p-12 lg:flex lg:items-center"><div class="absolute -right-28 -top-24 size-[500px] rounded-full bg-red-100/60 blur-3xl"/><div class="relative mx-auto max-w-md"><div class="grid aspect-[1.15] place-items-center rounded-3xl border border-slate-200 bg-white p-8 shadow-card"><div class="w-full rounded-2xl border border-slate-200 p-5"><div class="text-sm font-semibold">Your downloads, in one place</div><div class="mt-5 space-y-3"><div v-for="n in 3" :key="n" class="flex items-center gap-3 rounded-xl bg-slate-50 p-3"><span class="media-thumb h-10 w-14 rounded-md"/><span class="flex-1"><i class="mb-2 block h-2 w-28 rounded bg-slate-300"/><i class="block h-1.5 w-16 rounded bg-slate-200"/></span><span class="size-2 rounded-full bg-green-500"/></div></div></div></div><h2 class="mt-8 text-2xl font-bold">Simple, fast video downloads.</h2><p class="mt-2 leading-6 text-slate-500">Pick up where you left off and keep your downloads organized.</p></div></aside>
  </main>
</template>

<style scoped>.media-thumb{background:linear-gradient(150deg,#bae6fd,#38bdf8 45%,#0f766e 46%,#164e63)}</style>
