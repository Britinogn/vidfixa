import { defineNuxtConfig } from 'nuxt/config'

export default defineNuxtConfig({
  compatibilityDate: '2026-09-25',
  srcDir: '.',

  devtools: { enabled: true },

  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],

  css: ['~/assets/css/main.css'],

  runtimeConfig: {
    public: {
      apiBaseUrl:
        process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api',
    },
  },

  typescript: {
    strict: true,
    // Skip vue-tsc on every dev-server change — editor already surfaces errors live.
    // Type checking still runs on `nuxt build` so CI/prod catches regressions.
    typeCheck: process.env.NODE_ENV === 'production',
  },

  experimental: {
    // Makes <NuxtLink to="/dasboard"> a TypeScript error instead of a runtime 404.
    typedPages: true,
    // Smaller initial payload on client-side navigation.
    payloadExtraction: true,
  },

  app: {
    head: {
      meta: [
        // Nuxt does NOT inject viewport by default — without this, mobile Safari
        // renders at desktop width and scales down, breaking responsive layout.
        // viewport-fit=cover enables env(safe-area-inset-*) for notch/island spacing.
        {
          name: 'viewport',
          content: 'width=device-width, initial-scale=1, viewport-fit=cover',
        },
      ],
      link: [
        {
          rel: 'preconnect',
          href: 'https://fonts.googleapis.com',
        },
        {
          rel: 'preconnect',
          href: 'https://fonts.gstatic.com',
          crossorigin: '',
        },
        {
          rel: 'stylesheet',
          href:
            'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap',
        },
      ],
    },
  },

  // Render.com web service preset — wires up port binding and health checks.
  // If you're deploying as a Render *static site*, delete this block.
  nitro: {
    preset: 'render_com',
    compressPublicAssets: true,
  },

  routeRules: {
    // Marketing pages: prerender at build time for instant TTFB on Render free tier.
    '/': { prerender: true },
    '/about': { prerender: true },

    // Auth pages: SSR so <meta name="robots" content="noindex"> lands in initial HTML.
    '/login': { ssr: true },
    '/register': { ssr: true },

    // Auth-gated: SPA-rendered, no SSR needed.
    '/dashboard/**': { ssr: false },
    '/admin/**': { ssr: false },

    // Baseline security headers on every route.
    '/**': {
      headers: {
        'X-Content-Type-Options': 'nosniff',
        'X-Frame-Options': 'SAMEORIGIN',
        'Referrer-Policy': 'strict-origin-when-cross-origin',
        'Permissions-Policy': 'camera=(), microphone=(), geolocation=()',
      },
    },
  },
})