<script setup lang="ts">
import { authClient, useSession, signOut } from "~/lib/auth-client";
import { polarClient } from "~/lib/auth";

const session = useSession();
const subscription = ref<any>(null);
const isLoadingSubscription = ref(false);

// Fetch customer subscription info when user is logged in
watch(() => session.value.data?.user, async (user) => {
  if (user) {
    await fetchSubscription();
  } else {
    subscription.value = null;
  }
}, { immediate: true });

async function fetchSubscription() {
  isLoadingSubscription.value = true;
  try {
    const response = await authClient.customer.subscriptions.list();
    if (response.data) {
      subscription.value = response.data;
    }
  } catch (error) {
    console.error("Failed to fetch subscription:", error);
  } finally {
    isLoadingSubscription.value = false;
  }
}

const currentPlan = computed(() => {
  if (!subscription.value?.subscription) {
    return { name: "Free", color: "neutral" as const };
  }

  const sub = subscription.value.subscription;
  const productName = sub.product?.name || sub.productId;

  if (productName?.toLowerCase().includes("annual")) {
    return { name: "Premium Annual", color: "primary" as const };
  } else if (productName?.toLowerCase().includes("premium")) {
    return { name: "Premium", color: "primary" as const };
  }

  return { name: "Free", color: "neutral" as const };
});

async function handleSignOut() {
  await signOut();
  navigateTo('/');
}

async function openPortal() {
  try {
    const result = await polarClient.customerSessions.create({
      customerId: session.value.data!.user!.id,
    });

    await navigateTo(result.customerPortalUrl, { external: true });
  } catch (error) {
    console.error("Portal error:", error);
  }
}

const menuItems = computed(() => [
  [{
    label: session.value.data?.user?.email || 'Account',
    slot: 'account',
    disabled: true
  }],
  [{
    label: 'Dashboard',
    icon: 'i-lucide-layout-dashboard',
    click: () => navigateTo('/dashboard')
  }, {
    label: 'Pricing',
    icon: 'i-lucide-credit-card',
    click: () => navigateTo('/pricing')
  }, {
    label: 'Manage Subscription',
    icon: 'i-lucide-settings',
    click: openPortal
  }],
  [{
    label: 'Sign Out',
    icon: 'i-lucide-log-out',
    click: handleSignOut
  }]
]);
</script>

<template>
  <div class="flex items-center gap-2">
    <template v-if="session.data?.user">
      <UBadge
        :color="currentPlan.color"
        variant="subtle"
        size="sm"
      >
        {{ currentPlan.name }}
      </UBadge>

      <UDropdownMenu :items="menuItems">
        <UButton
          color="neutral"
          variant="ghost"
          icon="i-lucide-user"
          aria-label="User menu"
        />

        <template #account>
          <div class="text-sm truncate max-w-[200px]">
            {{ session.data?.user?.email }}
          </div>
        </template>
      </UDropdownMenu>
    </template>

    <template v-else>
      <UButton
        to="/login"
        color="neutral"
        variant="ghost"
        size="sm"
      >
        Sign In
      </UButton>
      <UButton
        to="/register"
        color="primary"
        size="sm"
      >
        Get Started
      </UButton>
    </template>
  </div>
</template>
