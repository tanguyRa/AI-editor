<script setup lang="ts">
import { useSession } from "~/lib/auth-client";

interface ProductPrice {
  id: string;
  type: string;
  amountType: string;
  priceAmount: number | null;
  priceCurrency: string | null;
  recurringInterval: string | null;
}

interface ProductBenefit {
  id: string;
  description: string;
  type: string;
}

interface Product {
  id: string;
  slug: string;
  name: string;
  description: string | null;
  prices: ProductPrice[];
  benefits: ProductBenefit[];
  isRecurring: boolean;
  isHighlighted: boolean;
}

interface PricingPlan {
  title: string;
  description: string;
  price: string;
  discount?: string;
  billingCycle?: string;
  billingPeriod?: string;
  features: string[];
  button: {
    label: string;
    onClick: () => void;
  };
  highlight?: boolean;
  scale?: boolean;
  badge?: string;
  slug: string;
}

const session = useSession();
const isLoading = ref<string | null>(null);

const { data: productsData, pending: loadingProducts } = await useFetch<{ products: Product[] }>('/api/polar/products');

const plans = computed<PricingPlan[]>(() => {
  if (!productsData.value?.products) return [];

  return productsData.value.products.map((product: Product) => {
    const price = product.prices[0];
    const priceAmount = price?.priceAmount ? price.priceAmount / 100 : 0;
    const isAnnual = price?.recurringInterval === 'year';
    const isMonthly = price?.recurringInterval === 'month';
    const isFree = priceAmount === 0;

    const billingCycle = isFree ? '' : isMonthly ? '/month' : isAnnual ? '/year' : '';
    const billingPeriod = isAnnual ? 'billed annually' : undefined;

    return {
      slug: product.slug,
      title: product.name,
      description: product.description || getDefaultDescription(product.slug),
      price: isFree ? 'Free' : `$${priceAmount}`,
      billingCycle,
      billingPeriod,
      features: product.benefits.map((b: ProductBenefit) => b.description).filter(Boolean),
      highlight: product.isHighlighted || product.slug === 'Premium',
      scale: product.isHighlighted || product.slug === 'Premium',
      badge: product.isHighlighted ? 'Most popular' : undefined,
      button: {
        label: isFree ? 'Get Started' : 'Subscribe',
        onClick: () => handleCheckout(product.slug)
      }
    };
  });
});

function getDefaultDescription(slug: string): string {
  const descriptions: Record<string, string> = {
    'Free': 'Get started with basic features',
    'Premium': 'For professionals who need more power',
    'Premium-Annual': 'Best value - save 2 months'
  };
  return descriptions[slug] || 'Access all features';
}

async function handleCheckout(slug: string) {
  if (!session.value.data?.user) {
    navigateTo('/login');
    return;
  }

  isLoading.value = slug;
  try {
    await navigateTo(`/api/checkout/${slug}`, { external: true });
  } catch (error) {
    console.error("Checkout error:", error);
  } finally {
    isLoading.value = null;
  }
}
</script>

<template>
  <div class="py-16 px-4">
    <div class="max-w-6xl mx-auto">
      <div class="text-center mb-12">
        <h1 class="text-4xl font-bold mb-4">Choose Your Plan</h1>
        <p class="text-lg text-muted">
          Select the plan that best fits your needs
        </p>
      </div>

      <div v-if="loadingProducts" class="flex justify-center py-20">
        <UIcon name="i-lucide-loader-2" class="w-8 h-8 animate-spin text-primary" />
      </div>

      <div v-else-if="plans.length === 0" class="text-center py-20">
        <p class="text-muted">No plans available at the moment.</p>
      </div>

      <UPricingPlans v-else>
        <UPricingPlan
          v-for="plan in plans"
          :key="plan.slug"
          :title="plan.title"
          :description="plan.description"
          :price="plan.price"
          :billing-cycle="plan.billingCycle"
          :billing-period="plan.billingPeriod"
          :features="plan.features"
          :button="plan.button"
          :highlight="plan.highlight"
          :scale="plan.scale"
          :badge="plan.badge"
        />
      </UPricingPlans>

      <div class="text-center mt-12">
        <p class="text-sm text-muted">
          All plans include a 14-day money-back guarantee.
          <NuxtLink to="/dashboard" class="text-primary hover:underline">
            View your current subscription
          </NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>
