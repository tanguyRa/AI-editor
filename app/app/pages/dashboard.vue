<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
})

import { authClient } from "~/lib/auth-client";
const session = authClient.useSession()

console.log("Current session:", session);
</script>


<template>
  <UPageSection
      id="session"
      title="Current session information"
      :features="[{
        icon: 'i-lucide-rocket',
        title: 'id',
        description: session.data?.user?.id
      }, {
        icon: 'i-lucide-palette',
        title: 'email',
        description: session.data?.user?.email
      }]"
    />
</template>