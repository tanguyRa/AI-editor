<script setup>
useHead({
  meta: [
    { name: 'viewport', content: 'width=device-width, initial-scale=1' }
  ],
  link: [
    { rel: 'icon', href: '/favicon.ico' }
  ],
  htmlAttrs: {
    lang: 'en'
  }
})

const title = 'Flowmail - Email, reimagined'
const description = 'A next-generation email client that connects to Gmail, Outlook, and iCloud. Smart filters, multiple inboxes, and AI agents make email a pleasure again.'

useSeoMeta({
  title,
  description,
  ogTitle: title,
  ogDescription: description,
  twitterCard: 'summary_large_image'
})
</script>

<template>
  <UApp>
    <UHeader>
      <template #left>
        <NuxtLink to="/" class="flex items-center gap-2">
          <UIcon name="i-lucide-mail" class="w-6 h-6 text-primary-500" />
          <span class="font-semibold text-lg">Flowmail</span>
        </NuxtLink>
      </template>

      <template #center>
        <UButton
          to="/#features"
          color="neutral"
          variant="ghost"
          size="sm"
        >
          Features
        </UButton>
        <UButton
          to="/pricing"
          color="neutral"
          variant="ghost"
          size="sm"
        >
          Pricing
        </UButton>
      </template>

      <template #right>
        <UColorModeButton />

        <UserMenu />
      </template>
    </UHeader>

    <UMain>
      <NuxtPage />
    </UMain>

    <USeparator icon="i-lucide-mail" />

    <UFooter>
      <template #left>
        <p class="text-sm text-muted">
          Flowmail • © {{ new Date().getFullYear() }} • Email, reimagined.
        </p>
      </template>

      <template #right>
        <UButton
          to="https://twitter.com"
          target="_blank"
          icon="i-simple-icons-x"
          aria-label="Twitter"
          color="neutral"
          variant="ghost"
        />
      </template>
    </UFooter>
  </UApp>
</template>
