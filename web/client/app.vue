<script setup lang="ts">
const siteUrl = 'https://vidfixa.onrender.com'
const siteName = 'VidFixa'
const defaultTitle = 'VidFixa — Download videos. Simple. Fast.'

const defaultDescription =
  'Save videos from Instagram, Facebook, X, and LinkedIn in seconds. High quality, no watermark, secure and private.'

useHead({
  htmlAttrs: { lang: 'en' },

  titleTemplate: (title) =>
    title ? `${title} · ${siteName}` : defaultTitle,

  link: [
    { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
    { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
    { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png' },
    { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
    { rel: 'manifest', href: '/site.webmanifest' },
    { rel: 'canonical', href: siteUrl },
  ],

  meta: [
    { name: 'theme-color', content: '#E11D48' },
    { name: 'msapplication-TileColor', content: '#E11D48' },
  ],
})

useSeoMeta({
  title: defaultTitle,
  description: defaultDescription,
  robots: 'index, follow',
  author: 'VidFixa',
  ogType: 'website',
  ogSiteName: siteName,
  ogTitle: defaultTitle,
  ogDescription: defaultDescription,
  ogUrl: siteUrl,
  ogImage: `${siteUrl}/og-image.png`,
  ogImageWidth: 1200,
  ogImageHeight: 630,
  ogImageAlt: 'VidFixa — Download videos. Simple. Fast.',
  ogLocale: 'en_US',

  twitterCard: 'summary_large_image',
  twitterTitle: defaultTitle,
  twitterDescription: defaultDescription,
  twitterImage: `${siteUrl}/og-image.png`,
  twitterImageAlt: 'VidFixa — Download videos. Simple. Fast.',
})

useHead({
  meta: [
    {
      name: 'keywords',
      content:
        'video downloader, download instagram video, download facebook video, no watermark, VidFixa',
    },
  ],
})
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>