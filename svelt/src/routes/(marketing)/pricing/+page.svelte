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
            Free: "Get started with basic features",
            Premium: "For professionals who need more power",
            "Premium-Annual": "Best value - save 2 months",
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
            await checkout({ slug });
        } catch (e) {
            console.error("Checkout error:", e);
            checkoutLoading = null;
        }
    }
</script>

<div class="pricing-container">
    <div class="pricing-content">
        <div class="section-header">
            <h1>Choose Your Plan</h1>
            <p>Select the plan that best fits your needs</p>
        </div>

        {#if loading}
            <div class="loading">
                <div class="spinner"></div>
                <p>Loading plans...</p>
            </div>
        {:else if error}
            <div class="error-message centered">{error}</div>
        {:else if products.length === 0}
            <div class="empty-state">
                <p>No plans available at the moment.</p>
            </div>
        {:else}
            <div class="plans-grid">
                {#each products as product}
                    {@const isHighlighted =
                        product.isHighlighted || product.slug === "Premium"}
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
                            <span class="cycle">{getBillingCycle(product)}</span
                            >
                            {#if getBillingPeriod(product)}
                                <p class="period">
                                    {getBillingPeriod(product)}
                                </p>
                            {/if}
                        </div>

                        <ul class="features">
                            {#each product.benefits as benefit}
                                <li>
                                    <svg
                                        viewBox="0 0 24 24"
                                        width="20"
                                        height="20"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="2"
                                    >
                                        <polyline points="20 6 9 17 4 12"
                                        ></polyline>
                                    </svg>
                                    {benefit.description}
                                </li>
                            {/each}
                            {#if product.benefits.length === 0}
                                <li>
                                    <svg
                                        viewBox="0 0 24 24"
                                        width="20"
                                        height="20"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="2"
                                    >
                                        <polyline points="20 6 9 17 4 12"
                                        ></polyline>
                                    </svg>
                                    Access to basic features
                                </li>
                            {/if}
                        </ul>

                        <button
                            class="btn plan-button"
                            class:btn-primary={isHighlighted}
                            onclick={() => handleCheckout(product.slug)}
                            disabled={checkoutLoading !== null}
                        >
                            {#if checkoutLoading === product.slug}
                                <span class="spinner spinner-sm"></span>
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
    /* Pricing Page Specific Styles */
    .pricing-container {
        min-height: 100vh;
        background: var(--gradient-primary);
        padding: var(--spacing-xl) var(--spacing-md);
    }

    .pricing-content {
        max-width: 1200px;
        margin: 0 auto;
    }

    .pricing-content .section-header {
        color: white;
    }

    .pricing-content .section-header h1 {
        font-size: var(--font-size-4xl);
        margin-bottom: 0.75rem;
    }

    .pricing-content .section-header p {
        color: rgba(255, 255, 255, 0.9);
    }

    .loading {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: var(--spacing-3xl);
        color: white;
    }

    .loading p {
        margin-top: var(--spacing-md);
    }

    .centered {
        max-width: 400px;
        margin: 0 auto;
    }

    .empty-state {
        text-align: center;
        color: white;
        padding: var(--spacing-3xl);
        opacity: 0.9;
    }

    .plans-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: var(--spacing-lg);
        align-items: stretch;
    }

    .plan-card {
        background: var(--color-bg);
        border-radius: var(--radius-xl);
        padding: var(--spacing-xl);
        display: flex;
        flex-direction: column;
        position: relative;
        transition:
            transform var(--transition-base),
            box-shadow var(--transition-base);
        box-shadow: var(--shadow-md);
    }

    .plan-card:hover {
        transform: translateY(-4px);
        box-shadow: var(--shadow-lg);
    }

    .plan-card.highlighted {
        border: 2px solid var(--color-primary);
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
        background: var(--gradient-primary);
        color: white;
        padding: 0.375rem 1rem;
        border-radius: var(--radius-full);
        font-size: var(--font-size-xs);
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .plan-header {
        margin-bottom: var(--spacing-lg);
    }

    .plan-header h2 {
        font-size: var(--font-size-2xl);
        margin-bottom: var(--spacing-sm);
        color: var(--color-text);
    }

    .description {
        color: var(--color-text-muted);
        font-size: 0.95rem;
        line-height: 1.5;
    }

    .plan-price {
        margin-bottom: var(--spacing-lg);
        padding-bottom: var(--spacing-lg);
        border-bottom: 1px solid var(--color-border);
    }

    .amount {
        font-size: 3rem;
        font-weight: 700;
        color: var(--color-text);
        line-height: 1;
    }

    .cycle {
        font-size: var(--font-size-lg);
        color: var(--color-text-muted);
        margin-left: var(--spacing-xs);
    }

    .period {
        margin-top: var(--spacing-sm);
        font-size: var(--font-size-sm);
        color: var(--color-text-light);
    }

    .features {
        list-style: none;
        padding: 0;
        margin: 0 0 var(--spacing-xl);
        flex: 1;
    }

    .features li {
        display: flex;
        align-items: flex-start;
        gap: 0.75rem;
        padding: 0.625rem 0;
        color: var(--color-text-secondary);
        font-size: 0.95rem;
        line-height: 1.4;
    }

    .features li svg {
        flex-shrink: 0;
        color: var(--color-success);
        margin-top: 2px;
    }

    .plan-button {
        width: 100%;
        border: 2px solid var(--color-border);
        background: var(--color-bg);
        color: var(--color-text-secondary);
    }

    .plan-button:hover:not(:disabled) {
        border-color: var(--color-primary);
        color: var(--color-primary);
    }

    .plan-button.btn-primary {
        border: none;
    }

    .pricing-footer {
        text-align: center;
        margin-top: var(--spacing-2xl);
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
