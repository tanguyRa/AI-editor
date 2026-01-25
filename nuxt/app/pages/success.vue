<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
})

// Auto-redirect to dashboard after a few seconds
const countdown = ref(5);

onMounted(() => {
  const interval = setInterval(() => {
    countdown.value--;
    if (countdown.value <= 0) {
      clearInterval(interval);
      navigateTo('/dashboard');
    }
  }, 1000);
});
</script>

<template>
  <div class="min-h-[60vh] flex items-center justify-center px-4">
    <div class="text-center max-w-md">
      <div class="mb-6">
        <div class="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4">
          <UIcon name="i-lucide-check-circle" class="w-8 h-8 text-primary" />
        </div>
        <h1 class="text-3xl font-bold mb-2">Payment Successful!</h1>
        <p class="text-muted">
          Thank you for your purchase. Your subscription is now active.
        </p>
      </div>

      <UAlert
        color="primary"
        variant="subtle"
        icon="i-lucide-info"
        class="mb-6"
      >
        <template #description>
          Redirecting to dashboard in {{ countdown }} seconds...
        </template>
      </UAlert>

      <div class="flex gap-3 justify-center">
        <UButton
          to="/dashboard"
          color="primary"
        >
          Go to Dashboard
        </UButton>
        <UButton
          to="/pricing"
          color="neutral"
          variant="outline"
        >
          View Plans
        </UButton>
      </div>
    </div>
  </div>
</template>
