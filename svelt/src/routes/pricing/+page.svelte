<script lang="ts">
    import { authClient, useSession, checkout } from "$lib/auth-client";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";

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

    const session = useSession();

    let products = $state<Product[]>([]);
    let loading = $state(true);
    let checkoutLoading = $state<string | null>(null);
    let error = $state("");

    onMount(async () => {
        try {
            const response = await fetch("/api/polar/products");
            const data = await response.json();
            if (data.products) {
                products = data.products;
            } else {
                error = "Failed to load pricing plans";
            }
        } catch (e) {
            error = "Failed to load pricing plans";
        } finally {
            loading = false;
        }
    });

    function formatPrice(product: Product): string {
        const price = product.prices[0];
        if (!price || price.priceAmount === null || price.priceAmount === 0) {
            return "Free";
        }
        return `$${price.priceAmount / 100}`;
    }

    function getBillingCycle(product: Product): string {
        const price = product.prices[0];
        if (!price || price.priceAmount === null || price.priceAmount === 0) {
            return "";
        }
        if (price.recurringInterval === "month") return "/month";
        if (price.recurringInterval === "year") return "/year";
        return "";
    }

    function getBillingPeriod(product: Product): string | undefined {
        const price = product.prices[0];
        if (price?.recurringInterval === "year") {
            return "billed annually";
        }
        return undefined;
    }

    function getDescription(product: Product): string {
        if (product.description) return product.description;
        const descriptions: Record<string, string> = {
            "Free": "Get started with basic features",
            "Premium": "For professionals who need more power",
            "Premium-Annual": "Best value - save 2 months"
        };
        return descriptions[product.slug] || "Access all features";
    }

    function getButtonLabel(product: Product): string {
        const price = product.prices[0];
        if (!price || price.priceAmount === null || price.priceAmount === 0) {
            return "Get Started";
        }
        return "Subscribe";
    }

    async function handleCheckout(slug: string) {
        if (!$session.data?.user) {
            goto("/login");
            return;
        }

        checkoutLoading = slug;
        try {
            // Use the polar checkout function from auth-client
            // This redirects to Polar's checkout page
            await checkout({ slug });
        } catch (e) {
            console.error("Checkout error:", e);
            checkoutLoading = null;
        }
    }
</script>

<div class="pricing-container">
    <div class="pricing-content">
        <div class="pricing-header">
            <h1>Choose Your Plan</h1>
            <p>Select the plan that best fits your needs</p>
        </div>

        {#if loading}
            <div class="loading">
                <div class="spinner"></div>
                <p>Loading plans...</p>
            </div>
        {:else if error}
            <div class="error-message">{error}</div>
        {:else if products.length === 0}
            <div class="empty-state">
                <p>No plans available at the moment.</p>
            </div>
        {:else}
            <div class="plans-grid">
                {#each products as product}
                    {@const isHighlighted = product.isHighlighted || product.slug === "Premium"}
                    {@const isFree = product.prices[0]?.priceAmount === 0 || product.prices[0]?.priceAmount === null}
                    <div class="plan-card" class:highlighted={isHighlighted}>
                        {#if isHighlighted}
                            <div class="badge">Most Popular</div>
                        {/if}

                        <div class="plan-header">
                            <h2>{product.name}</h2>
                            <p class="description">{getDescription(product)}</p>
                        </div>

                        <div class="plan-price">
                            <span class="amount">{formatPrice(product)}</span>
                            <span class="cycle">{getBillingCycle(product)}</span>
                            {#if getBillingPeriod(product)}
                                <p class="period">{getBillingPeriod(product)}</p>
                            {/if}
                        </div>

                        <ul class="features">
                            {#each product.benefits as benefit}
                                <li>
                                    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2">
                                        <polyline points="20 6 9 17 4 12"></polyline>
                                    </svg>
                                    {benefit.description}
                                </li>
                            {/each}
                            {#if product.benefits.length === 0}
                                <li>
                                    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2">
                                        <polyline points="20 6 9 17 4 12"></polyline>
                                    </svg>
                                    Access to basic features
                                </li>
                            {/if}
                        </ul>

                        <button
                            class="plan-button"
                            class:primary={isHighlighted}
                            onclick={() => handleCheckout(product.slug)}
                            disabled={checkoutLoading !== null}
                        >
                            {#if checkoutLoading === product.slug}
                                <span class="button-spinner"></span>
                                Processing...
                            {:else}
                                {getButtonLabel(product)}
                            {/if}
                        </button>
                    </div>
                {/each}
            </div>
        {/if}

        <div class="pricing-footer">
            <p>
                All plans include a 14-day money-back guarantee.
                {#if $session.data?.user}
                    <a href="/dashboard">View your current subscription</a>
                {/if}
            </p>
        </div>
    </div>
</div>

<style>
    .pricing-container {
        min-height: 100vh;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        padding: 2rem 1rem;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
    }

    .pricing-content {
        max-width: 1200px;
        margin: 0 auto;
    }

    .pricing-header {
        text-align: center;
        margin-bottom: 3rem;
        color: white;
    }

    .pricing-header h1 {
        margin: 0 0 0.75rem;
        font-size: 2.5rem;
        font-weight: 700;
    }

    .pricing-header p {
        margin: 0;
        font-size: 1.125rem;
        opacity: 0.9;
    }

    .loading {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 4rem;
        color: white;
    }

    .spinner {
        width: 40px;
        height: 40px;
        border: 3px solid rgba(255, 255, 255, 0.3);
        border-top-color: white;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin-bottom: 1rem;
    }

    @keyframes spin {
        to { transform: rotate(360deg); }
    }

    .error-message {
        background: #fef2f2;
        border: 1px solid #fecaca;
        color: #dc2626;
        padding: 1rem 1.5rem;
        border-radius: 12px;
        text-align: center;
        max-width: 400px;
        margin: 0 auto;
    }

    .empty-state {
        text-align: center;
        color: white;
        padding: 4rem;
        opacity: 0.9;
    }

    .plans-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 1.5rem;
        align-items: stretch;
    }

    .plan-card {
        background: white;
        border-radius: 16px;
        padding: 2rem;
        display: flex;
        flex-direction: column;
        position: relative;
        transition: transform 0.2s, box-shadow 0.2s;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
    }

    .plan-card:hover {
        transform: translateY(-4px);
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
    }

    .plan-card.highlighted {
        border: 2px solid #667eea;
        transform: scale(1.02);
    }

    .plan-card.highlighted:hover {
        transform: scale(1.02) translateY(-4px);
    }

    .badge {
        position: absolute;
        top: -12px;
        left: 50%;
        transform: translateX(-50%);
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        padding: 0.375rem 1rem;
        border-radius: 20px;
        font-size: 0.75rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .plan-header {
        margin-bottom: 1.5rem;
    }

    .plan-header h2 {
        margin: 0 0 0.5rem;
        font-size: 1.5rem;
        font-weight: 700;
        color: #1a1a2e;
    }

    .description {
        margin: 0;
        color: #6b7280;
        font-size: 0.95rem;
        line-height: 1.5;
    }

    .plan-price {
        margin-bottom: 1.5rem;
        padding-bottom: 1.5rem;
        border-bottom: 1px solid #e5e7eb;
    }

    .amount {
        font-size: 3rem;
        font-weight: 700;
        color: #1a1a2e;
        line-height: 1;
    }

    .cycle {
        font-size: 1.125rem;
        color: #6b7280;
        margin-left: 0.25rem;
    }

    .period {
        margin: 0.5rem 0 0;
        font-size: 0.875rem;
        color: #9ca3af;
    }

    .features {
        list-style: none;
        padding: 0;
        margin: 0 0 2rem;
        flex: 1;
    }

    .features li {
        display: flex;
        align-items: flex-start;
        gap: 0.75rem;
        padding: 0.625rem 0;
        color: #374151;
        font-size: 0.95rem;
        line-height: 1.4;
    }

    .features li svg {
        flex-shrink: 0;
        color: #10b981;
        margin-top: 2px;
    }

    .plan-button {
        width: 100%;
        padding: 1rem 1.5rem;
        border: 2px solid #e5e7eb;
        border-radius: 10px;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        background: white;
        color: #374151;
    }

    .plan-button:hover:not(:disabled) {
        border-color: #667eea;
        color: #667eea;
    }

    .plan-button.primary {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        border: none;
        color: white;
    }

    .plan-button.primary:hover:not(:disabled) {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
    }

    .plan-button:disabled {
        opacity: 0.7;
        cursor: not-allowed;
    }

    .button-spinner {
        width: 18px;
        height: 18px;
        border: 2px solid transparent;
        border-top-color: currentColor;
        border-radius: 50%;
        animation: spin 0.8s linear infinite;
    }

    .pricing-footer {
        text-align: center;
        margin-top: 3rem;
        color: white;
    }

    .pricing-footer p {
        font-size: 0.95rem;
        opacity: 0.9;
    }

    .pricing-footer a {
        color: white;
        font-weight: 600;
        text-decoration: underline;
        text-underline-offset: 2px;
    }

    .pricing-footer a:hover {
        opacity: 0.8;
    }

    @media (max-width: 768px) {
        .pricing-header h1 {
            font-size: 2rem;
        }

        .plans-grid {
            grid-template-columns: 1fr;
        }

        .plan-card.highlighted {
            transform: none;
        }

        .plan-card.highlighted:hover {
            transform: translateY(-4px);
        }

        .amount {
            font-size: 2.5rem;
        }
    }
</style>
