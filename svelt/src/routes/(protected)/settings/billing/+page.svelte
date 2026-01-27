<script lang="ts">
    import { checkout, customer, useSession } from "$lib/auth-client";
    import { sessionMiddleware } from "better-auth/api";
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

    interface Subscription {
        id: string;
        status: string;
        productId: string;
        priceId: string;
        currentPeriodEnd: string | null;
        cancelAtPeriodEnd: boolean;
    }

    interface CustomerState {
        id: string;
        subscriptions: Subscription[];
    }

    let products = $state<Product[]>([]);
    let customerState = $state<CustomerState | null>(null);
    let loading = $state(true);
    let customerLoading = $state(true);
    let checkoutLoading = $state<string | null>(null);
    let portalLoading = $state(false);
    let error = $state("");

    const session = useSession();

    onMount(async () => {
        await Promise.all([fetchProducts(), fetchCustomerState()]);
    });

    async function fetchProducts() {
        try {
            const response = await fetch("/api/polar/products");
            const data = await response.json();
            if (data.products) {
                products = data.products;
            }
        } catch {
            error = "Failed to load plans";
        } finally {
            loading = false;
        }
    }

    async function fetchCustomerState() {
        try {
            customerState = {
                id: $session.data?.user.id || "",
                subscriptions: customer.subscriptions.list as Subscription[],
            };
            // if (customerError) {
            //     console.error("Customer state error:", customerError);
            // } else if (data) {
            //     customerState = data as CustomerState;
            // }
        } catch (e) {
            console.error("Failed to fetch customer state:", e);
        } finally {
            customerLoading = false;
        }
    }

    function getCurrentSubscription(): Subscription | null {
        if (!customerState?.subscriptions?.length) return null;
        return (
            customerState.subscriptions.find((s) => s.status === "active") ||
            null
        );
    }

    function getCurrentProduct(): Product | null {
        const subscription = getCurrentSubscription();
        if (!subscription) return null;
        return products.find((p) => p.id === subscription.productId) || null;
    }

    function isCurrentPlan(product: Product): boolean {
        const currentProduct = getCurrentProduct();
        return currentProduct?.id === product.id;
    }

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

    function formatDate(dateString: string | null): string {
        if (!dateString) return "";
        return new Date(dateString).toLocaleDateString("en-US", {
            year: "numeric",
            month: "long",
            day: "numeric",
        });
    }

    function getStatusBadgeClass(status: string): string {
        switch (status) {
            case "active":
                return "badge-success";
            case "past_due":
                return "badge-warning";
            case "canceled":
                return "badge-error";
            default:
                return "badge-neutral";
        }
    }

    async function handleCheckout(slug: string) {
        checkoutLoading = slug;
        try {
            await checkout({ slug });
        } catch (e) {
            console.error("Checkout error:", e);
        } finally {
            checkoutLoading = null;
        }
    }

    async function handleOpenPortal() {
        portalLoading = true;
        try {
            const response = await fetch("/api/auth/portal", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
            });
            const data = await response.json();
            if (data.url) {
                window.location.href = data.url;
            }
        } catch (e) {
            console.error("Portal error:", e);
        } finally {
            portalLoading = false;
        }
    }
</script>

<div class="settings-page">
    <header class="settings-header">
        <h1>Billing & Subscription</h1>
        <p>Manage your subscription and billing information</p>
    </header>

    {#if loading || customerLoading}
        <div class="loading-state">
            <div class="spinner spinner-dark"></div>
        </div>
    {:else}
        {@const subscription = getCurrentSubscription()}
        {@const currentProduct = getCurrentProduct()}

        <div class="settings-sections">
            <!-- Current Plan Section -->
            <section class="settings-section">
                <div class="section-header">
                    <h2>Current Plan</h2>
                    <p>Your active subscription</p>
                </div>

                {#if subscription && currentProduct}
                    <div class="current-plan">
                        <div class="plan-info">
                            <div class="plan-name-row">
                                <h3>{currentProduct.name}</h3>
                                <span
                                    class="badge {getStatusBadgeClass(
                                        subscription.status,
                                    )}"
                                >
                                    {subscription.status === "active"
                                        ? "Active"
                                        : subscription.status}
                                </span>
                            </div>
                            <p class="plan-price-info">
                                {formatPrice(currentProduct)}{getBillingCycle(
                                    currentProduct,
                                )}
                            </p>
                            {#if subscription.currentPeriodEnd}
                                <p class="plan-renewal">
                                    {#if subscription.cancelAtPeriodEnd}
                                        <span class="cancellation-notice">
                                            Cancels on {formatDate(
                                                subscription.currentPeriodEnd,
                                            )}
                                        </span>
                                    {:else}
                                        Renews on {formatDate(
                                            subscription.currentPeriodEnd,
                                        )}
                                    {/if}
                                </p>
                            {/if}
                        </div>
                        <div class="plan-actions">
                            <button
                                class="btn btn-secondary"
                                onclick={handleOpenPortal}
                                disabled={portalLoading}
                            >
                                {#if portalLoading}
                                    <span
                                        class="spinner spinner-sm spinner-dark"
                                    ></span>
                                {/if}
                                Manage Subscription
                            </button>
                        </div>
                    </div>
                {:else}
                    <div class="current-plan free-plan">
                        <div class="plan-info">
                            <div class="plan-name-row">
                                <h3>Free Plan</h3>
                                <span class="badge badge-neutral">Active</span>
                            </div>
                            <p class="plan-price-info">$0/month</p>
                            <p class="plan-description">
                                Basic features with limited access
                            </p>
                        </div>
                    </div>
                {/if}
            </section>

            <!-- Available Plans Section -->
            <section class="settings-section">
                <div class="section-header">
                    <h2>Available Plans</h2>
                    <p>Upgrade or change your subscription</p>
                </div>

                {#if error}
                    <div class="error-message">{error}</div>
                {:else if products.length === 0}
                    <p class="empty-state">No plans available.</p>
                {:else}
                    <div class="plans-list">
                        {#each products as product}
                            {@const isCurrent = isCurrentPlan(product)}
                            {@const isFree =
                                product.prices[0]?.priceAmount === 0 ||
                                product.prices[0]?.priceAmount === null}
                            <div class="plan-item" class:current={isCurrent}>
                                <div class="plan-item-info">
                                    <div class="plan-item-header">
                                        <h4>{product.name}</h4>
                                        {#if isCurrent}
                                            <span class="badge badge-primary"
                                                >Current</span
                                            >
                                        {/if}
                                    </div>
                                    <p class="plan-item-price">
                                        {formatPrice(product)}{getBillingCycle(
                                            product,
                                        )}
                                    </p>
                                    {#if product.benefits.length > 0}
                                        <ul class="plan-item-features">
                                            {#each product.benefits.slice(0, 3) as benefit}
                                                <li>{benefit.description}</li>
                                            {/each}
                                        </ul>
                                    {/if}
                                </div>
                                <div class="plan-item-action">
                                    {#if isCurrent}
                                        <button
                                            class="btn btn-secondary"
                                            disabled
                                        >
                                            Current Plan
                                        </button>
                                    {:else if isFree}
                                        <button
                                            class="btn btn-secondary"
                                            onclick={handleOpenPortal}
                                            disabled={portalLoading}
                                        >
                                            Downgrade
                                        </button>
                                    {:else}
                                        <button
                                            class="btn btn-primary"
                                            onclick={() =>
                                                handleCheckout(product.slug)}
                                            disabled={checkoutLoading !== null}
                                        >
                                            {#if checkoutLoading === product.slug}
                                                <span class="spinner spinner-sm"
                                                ></span>
                                            {/if}
                                            Upgrade
                                        </button>
                                    {/if}
                                </div>
                            </div>
                        {/each}
                    </div>
                {/if}
            </section>

            <!-- Payment Method Section -->
            <section class="settings-section">
                <div class="section-header">
                    <h2>Payment Method</h2>
                    <p>Manage your payment information</p>
                </div>

                <div class="payment-info">
                    <p class="payment-description">
                        Update your payment method, view invoices, or cancel
                        your subscription through the customer portal.
                    </p>
                    <button
                        class="btn btn-secondary"
                        onclick={handleOpenPortal}
                        disabled={portalLoading}
                    >
                        {#if portalLoading}
                            <span class="spinner spinner-sm spinner-dark"
                            ></span>
                        {/if}
                        Open Customer Portal
                    </button>
                </div>
            </section>
        </div>
    {/if}
</div>

<style>
    .settings-page {
        padding: var(--spacing-xl);
        max-width: 800px;
    }

    .settings-header {
        margin-bottom: var(--spacing-2xl);
    }

    .settings-header h1 {
        font-size: var(--font-size-3xl);
        color: var(--color-text);
        margin-bottom: var(--spacing-sm);
    }

    .settings-header p {
        color: var(--color-text-muted);
        font-size: var(--font-size-lg);
    }

    .settings-sections {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-xl);
    }

    .settings-section {
        background: var(--color-bg);
        border-radius: var(--radius-lg);
        padding: var(--spacing-xl);
        box-shadow: var(--shadow-sm);
    }

    .section-header {
        margin-bottom: var(--spacing-lg);
        padding-bottom: var(--spacing-md);
        border-bottom: 1px solid var(--color-border);
    }

    .section-header h2 {
        font-size: var(--font-size-xl);
        color: var(--color-text);
        margin-bottom: var(--spacing-xs);
    }

    .section-header p {
        color: var(--color-text-muted);
        font-size: var(--font-size-sm);
    }

    .loading-state {
        display: flex;
        justify-content: center;
        padding: var(--spacing-3xl);
    }

    /* Current Plan */
    .current-plan {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: var(--spacing-lg);
    }

    .plan-info {
        flex: 1;
    }

    .plan-name-row {
        display: flex;
        align-items: center;
        gap: var(--spacing-md);
        margin-bottom: var(--spacing-sm);
    }

    .plan-name-row h3 {
        font-size: var(--font-size-xl);
        color: var(--color-text);
        margin: 0;
    }

    .plan-price-info {
        font-size: var(--font-size-lg);
        color: var(--color-text-secondary);
        margin-bottom: var(--spacing-xs);
    }

    .plan-renewal {
        font-size: var(--font-size-sm);
        color: var(--color-text-muted);
    }

    .plan-description {
        font-size: var(--font-size-sm);
        color: var(--color-text-muted);
    }

    .cancellation-notice {
        color: var(--color-warning);
    }

    .plan-actions {
        flex-shrink: 0;
    }

    /* Badges */
    .badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        border-radius: var(--radius-full);
        font-size: var(--font-size-xs);
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .badge-success {
        background: #ecfdf5;
        color: var(--color-success);
    }

    .badge-warning {
        background: #fffbeb;
        color: var(--color-warning);
    }

    .badge-error {
        background: var(--color-error-bg);
        color: var(--color-error);
    }

    .badge-neutral {
        background: var(--color-bg-tertiary);
        color: var(--color-text-muted);
    }

    .badge-primary {
        background: rgba(102, 126, 234, 0.1);
        color: var(--color-primary);
    }

    /* Plans List */
    .plans-list {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);
    }

    .plan-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: var(--spacing-lg);
        border: 1px solid var(--color-border);
        border-radius: var(--radius-md);
        transition: border-color var(--transition-fast);
    }

    .plan-item:hover {
        border-color: var(--color-primary);
    }

    .plan-item.current {
        border-color: var(--color-primary);
        background: rgba(102, 126, 234, 0.02);
    }

    .plan-item-info {
        flex: 1;
    }

    .plan-item-header {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        margin-bottom: var(--spacing-xs);
    }

    .plan-item-header h4 {
        font-size: var(--font-size-lg);
        color: var(--color-text);
        margin: 0;
    }

    .plan-item-price {
        font-size: var(--font-size-base);
        color: var(--color-text-secondary);
        margin-bottom: var(--spacing-sm);
    }

    .plan-item-features {
        list-style: none;
        padding: 0;
        margin: 0;
        display: flex;
        flex-wrap: wrap;
        gap: var(--spacing-sm);
    }

    .plan-item-features li {
        font-size: var(--font-size-xs);
        color: var(--color-text-muted);
        background: var(--color-bg-tertiary);
        padding: 0.25rem 0.5rem;
        border-radius: var(--radius-sm);
    }

    .plan-item-action {
        flex-shrink: 0;
        margin-left: var(--spacing-lg);
    }

    /* Payment Info */
    .payment-info {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);
    }

    .payment-description {
        color: var(--color-text-muted);
        font-size: var(--font-size-sm);
    }

    .payment-info .btn {
        align-self: flex-start;
    }

    .empty-state {
        color: var(--color-text-muted);
        text-align: center;
        padding: var(--spacing-lg);
    }

    @media (max-width: 768px) {
        .settings-page {
            padding: var(--spacing-md);
        }

        .settings-section {
            padding: var(--spacing-lg);
        }

        .current-plan {
            flex-direction: column;
        }

        .plan-actions {
            width: 100%;
        }

        .plan-actions .btn {
            width: 100%;
        }

        .plan-item {
            flex-direction: column;
            align-items: flex-start;
            gap: var(--spacing-md);
        }

        .plan-item-action {
            margin-left: 0;
            width: 100%;
        }

        .plan-item-action .btn {
            width: 100%;
        }
    }
</style>
