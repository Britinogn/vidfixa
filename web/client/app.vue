<script setup lang="ts">
const siteUrl = 'https://vidfixa.onrender.com'
const siteName = 'VidFixa'
const defaultTitle = 'VidFixa — Download videos. Simple. Fast.'

const route = useRoute()
const canonical = computed(() => `${siteUrl}${route.path === '/' ? '' : route.path}`)

const defaultDescription =
  'Save videos from Instagram, Facebook, X, and LinkedIn in seconds. High quality, no watermark, secure and private.'

useHead({
  htmlAttrs: { lang: 'en' },

  titleTemplate: (title) =>
    title ? `${title} · ${siteName}` : defaultTitle,

  link: [
    // SVG favicon retired — the .ico below is the live icon. Kept for reference.
    { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
    //?v=2 cache-busts browsers that pinned the old SVG icon.
    // { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico?v=2' },
    { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
    // No favicon-16x16.png exists in public/ — the multi-size .ico covers 16px.

    // Apple touch icons — ascending order so older devices pick the closest match.
    { rel: 'apple-touch-icon', sizes: '152x152', href: '/apple-touch-icon-152x152.png' },
    { rel: 'apple-touch-icon', sizes: '167x167', href: '/apple-touch-icon-167x167.png' },
    { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },

    { rel: 'manifest', href: '/site.webmanifest' },

    // { rel: 'canonical', href: siteUrl },
    { rel: 'canonical', href: canonical },

    { rel: 'mask-icon', href: '/safari-pinned-tab.svg', color: '#E11D48' },
    { rel: 'sitemap', type: 'application/xml', href: '/sitemap.xml' },

    { rel: 'preconnect', href: 'https://www.sabilytics.com', crossorigin: '' },
    { rel: 'dns-prefetch', href: 'https://www.sabilytics.com' },
  ],

  meta: [
    // { name: 'theme-color', content: '#E11D48' },
    { name: 'color-scheme', content: 'light' },
    { name: 'theme-color', content: '#E11D48', media: '(prefers-color-scheme: light)' },
    { name: 'theme-color', content: '#0F172A', media: '(prefers-color-scheme: dark)' },
    { name: 'msapplication-TileColor', content: '#E11D48' },
    { name: 'msapplication-TileImage', content: '/mstile-144x144.png' },
    { name: 'apple-mobile-web-app-title', content: 'VidFixa' },
    { name: 'apple-mobile-web-app-capable', content: 'yes' },        // iOS < 16.4
    { name: 'mobile-web-app-capable', content: 'yes' },              // standard, iOS 16.4+
    { name: 'apple-mobile-web-app-status-bar-style', content: 'default' },
    { name: 'referrer', content: 'strict-origin-when-cross-origin' },
    { name: 'format-detection', content: 'telephone=no' },
    {
      name: 'keywords',
      content:
        'video downloader, download instagram video, download facebook video, no watermark, VidFixa',
    },
  ],

  script: [
    // ── JSON-LD structured data ─────────────────────────────────────
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'WebSite',
        name: siteName,
        url: siteUrl,
        description: defaultDescription,
        publisher: {
          '@type': 'Organization',
          name: siteName,
          logo: {
            '@type': 'ImageObject',
            url: `${siteUrl}/logo.png`,
          },
        },
      }),
    },
    // ── Analytics ───────────────────────────────────────────────────
    {
      src: 'https://www.sabilytics.com/script.js',
      async: true,
      'data-site': 'w3hfo4ayyk5a',
      'data-domain': 'vidfixa.onrender.com',
    },
  ],
})

useSeoMeta({
  title: defaultTitle,
  description: defaultDescription,
  robots: 'index, follow, max-snippet:-1, max-image-preview:large, max-video-preview:-1',
  author: 'VidFixa',
  ogType: 'website',
  ogSiteName: siteName,
  ogTitle: defaultTitle,
  ogDescription: defaultDescription,
  ogUrl: canonical,
  ogImage: `${siteUrl}/og-image.png`,
  ogImageType: 'image/png',
  ogImageWidth: 1200,
  ogImageHeight: 630,
  ogImageAlt: 'VidFixa — Download videos. Simple. Fast.',
  ogLocale: 'en_US',

  twitterCard: 'summary_large_image',
  twitterTitle: defaultTitle,
  twitterDescription: defaultDescription,
  twitterImage: `${siteUrl}/og-image.png`,
  twitterImageAlt: 'VidFixa — Download videos. Simple. Fast.',
  twitterSite: '@Britinogn',
  twitterCreator: '@Britinogn',
})
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<!-- <link rel="icon" type="image/png" href="/favicon.ico/favicon-96x96.png?v=20260929" sizes="96x96" />
<link rel="icon" type="image/svg+xml" href="/favicon.ico/favicon.svg?v=20260929" />
<link rel="shortcut icon" href="/favicon.ico/favicon.ico?v=20260929" />
<link rel="apple-touch-icon" sizes="180x180" href="/favicon.ico/apple-touch-icon.png?v=20260929" />
<meta name="apple-mobile-web-app-title" content="VidFixa" />
<link rel="manifest" href="/favicon.ico/site.webmanifest?v=20260929" /> -->