import { defineNuxtConfig } from 'nuxt/config'

export default defineNuxtConfig({
  compatibilityDate: '2026-09-25',
  srcDir: '.',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api',
    },
  },
  typescript: {
    strict: true,
    typeCheck: true,
  },
  app: {
    head: {
      title: 'VidFixa — Download videos. Simple. Fast.',
      meta: [
        { name: 'description', content: 'Save videos from supported platforms in a few seconds.' },
        { name: 'theme-color', content: '#ffffff' },
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap' },
      ],
    },
  },
})
