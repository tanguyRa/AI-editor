<script setup lang="ts">
import { authClient } from "~/lib/auth-client";

interface UserResponse {
  id: string
  email: string
  name: string
}

const session = authClient.useSession()
const userData = ref<UserResponse | null>(null)
const apiError = ref("")
const apiLoading = ref(false)

async function testGoAPI() {
  apiLoading.value = true
  userData.value = null
  apiError.value = ""

  try {
    const response = await $fetch<UserResponse>("/api/secured/ping")
    console.log("Go API response:", response)
    userData.value = response
  } catch (e: any) {
    apiError.value = e.message || String(e)
  }

  apiLoading.value = false
}
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-4 p-4">
    <UPageCard class="w-full max-w-md">
      <template #title>Go API Auth</template>
      <template #description>
        Test authentication against a Go API server
      </template>

      <div class="space-y-4">
        <UButton
          block
          color="primary"
          :disabled="!session.data || apiLoading"
          :loading="apiLoading"
          @click="testGoAPI"
        >
          {{ apiLoading ? "Testing..." : "Test Go API" }}
        </UButton>

        <p v-if="!session.data" class="text-sm text-muted">
          Sign in to test the Go API
        </p>

        <UPageCard v-if="userData" variant="soft" class="mt-4">
          <template #title>
            <span class="text-green-500">User Info</span>
          </template>

          <dl class="space-y-3">
            <div class="flex justify-between">
              <dt class="text-muted font-medium">ID</dt>
              <dd class="font-mono text-sm">{{ userData.id }}</dd>
            </div>
            <div class="flex justify-between">
              <dt class="text-muted font-medium">Email</dt>
              <dd class="font-mono text-sm">{{ userData.email }}</dd>
            </div>
            <div class="flex justify-between">
              <dt class="text-muted font-medium">Name</dt>
              <dd class="font-mono text-sm">{{ userData.name }}</dd>
            </div>
          </dl>
        </UPageCard>

        <UAlert
          v-if="apiError"
          color="error"
          icon="i-lucide-alert-circle"
          title="API Error"
          :description="apiError"
        />
      </div>
    </UPageCard>
  </div>
</template>
