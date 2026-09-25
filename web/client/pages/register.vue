<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'

definePageMeta({ layout: false })
const auth = useAuthStore()
const formError = ref('')
const schema = toTypedSchema(z.object({ fullName: z.string().trim().min(5, 'Use at least 5 characters.').max(25, 'Use no more than 25 characters.'), email: z.string().email('Enter a valid email.'), password: z.string().min(8, 'Password must be at least 8 characters.') }))
const { defineField, handleSubmit, errors, isSubmitting } = useForm({ validationSchema: schema, initialValues: { fullName: '', email: '', password: '' } })
const [fullName, fullNameAttrs] = defineField('fullName')
const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')

const submit = handleSubmit(async (values) => {
  formError.value = ''
  try {
    const user = await auth.register(values.fullName, values.email, values.password)
    await navigateTo(user.role === 'admin' ? '/admin' : '/dashboard')
  } catch (error) {
    const response = error as { data?: string; message?: string }
    formError.value = response.data || response.message || 'Could not create your account. Please try again.'
  }
})
</script>

<template>
  <main class="flex min-h-screen items-center justify-center bg-slate-50 px-5 py-12"><div class="w-full max-w-[460px] rounded-2xl border border-slate-200 bg-white p-7 shadow-card sm:p-10"><NuxtLink to="/"><BrandMark /></NuxtLink><h1 class="mt-10 text-3xl font-bold tracking-tight">Create your account</h1><p class="mt-2 text-slate-500">Start saving your favorite videos.</p>
    <form class="mt-7 space-y-4" @submit.prevent="submit"><label class="block text-sm font-medium">Full name<input v-model="fullName" v-bind="fullNameAttrs" autocomplete="name" class="mt-2 h-12 w-full rounded-[10px] border border-slate-200 px-4 outline-none focus:border-brand-600 focus:ring-4 focus:ring-red-50"><span v-if="errors.fullName" class="mt-1 block text-xs text-red-600">{{ errors.fullName }}</span></label><label class="block text-sm font-medium">Email<input v-model="email" v-bind="emailAttrs" type="email" autocomplete="email" class="mt-2 h-12 w-full rounded-[10px] border border-slate-200 px-4 outline-none focus:border-brand-600 focus:ring-4 focus:ring-red-50"><span v-if="errors.email" class="mt-1 block text-xs text-red-600">{{ errors.email }}</span></label><label class="block text-sm font-medium">Password<input v-model="password" v-bind="passwordAttrs" type="password" autocomplete="new-password" class="mt-2 h-12 w-full rounded-[10px] border border-slate-200 px-4 outline-none focus:border-brand-600 focus:ring-4 focus:ring-red-50"><span v-if="errors.password" class="mt-1 block text-xs text-red-600">{{ errors.password }}</span></label><p v-if="formError" class="rounded-xl bg-red-50 p-3 text-sm text-red-700" role="alert">{{ formError }}</p><button :disabled="isSubmitting" class="mt-2 w-full rounded-[10px] bg-brand-600 py-3.5 text-sm font-semibold text-white hover:bg-brand-700 disabled:opacity-60">{{ isSubmitting ? 'Creating account…' : 'Create account' }}</button></form>
    <p class="mt-6 text-center text-sm text-slate-500">Already have an account? <NuxtLink to="/login" class="font-semibold text-brand-600">Sign in</NuxtLink></p>
  </div></main>
</template>
