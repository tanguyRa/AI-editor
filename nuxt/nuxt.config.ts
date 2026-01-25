// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@nuxt/eslint',
    '@nuxt/ui',
    "@polar-sh/nuxt"
  ],

  devtools: {
    enabled: true
  },

  colorMode: {
    preference: 'system',
    fallback: 'dark'
  },

  runtimeConfig: {
    // Server-side only (not exposed to client)
    goApiUrl: process.env.GO_API_URL || 'http://api:8080'
  },

  css: ['~/assets/css/main.css'],

  routeRules: {
    '/': { prerender: true }
  },

  compatibilityDate: '2026-01-19',

  eslint: {
    config: {
      stylistic: {
        commaDangle: 'never',
        braceStyle: '1tbs'
      }
    }
  },

  vite: {
    server: {
      allowedHosts: ['app', 'app-budhapp']
    }
  },

  nitro: {
    compressPublicAssets: false
  }
})
