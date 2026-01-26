import { json } from '@sveltejs/kit';
import { polarClient } from '$lib/auth';
import type { RequestHandler } from './$types';

// Product slugs configured in auth.ts for checkout
const productSlugs: Record<string, string> = {
    "e54c3dec-3fa6-4a6d-b359-35fafdfe4b30": "Premium-Annual",
    "a741f0a8-929d-4420-8329-2e880fa2ecf8": "Premium",
    "015ddd64-2330-4fc7-a59d-c8cfcd9751ed": "Free"
};

export const GET: RequestHandler = async () => {
    try {
        // Fetch products from Polar API
        const response = await polarClient.products.list({
            isArchived: false,
        });

        // Map products with checkout slugs and sort by price
        const products = response.result.items
            .filter(product => productSlugs[product.id])
            .map(product => ({
                id: product.id,
                slug: productSlugs[product.id],
                name: product.name,
                description: product.description,
                prices: product.prices.map(price => ({
                    id: price.id,
                    type: price.type,
                    amountType: price.amountType,
                    priceAmount: price.amountType === 'fixed' ? price.priceAmount : null,
                    priceCurrency: price.amountType === 'fixed' ? price.priceCurrency : null,
                    recurringInterval: price.type === 'recurring' ? price.recurringInterval : null,
                })),
                benefits: product.benefits.map(benefit => ({
                    id: benefit.id,
                    description: benefit.description,
                    type: benefit.type,
                })),
                isRecurring: product.isRecurring,
                isHighlighted: (product as { isHighlighted?: boolean }).isHighlighted ?? false,
            }))
            .sort((a, b) => {
                // Sort by price: free first, then monthly, then annual
                const getMinPrice = (p: typeof products[0]) => {
                    const fixedPrice = p.prices.find(pr => pr.amountType === 'fixed' && pr.priceAmount !== null);
                    return fixedPrice?.priceAmount ?? 0;
                };
                return getMinPrice(a) - getMinPrice(b);
            });

        return json({ products });
    } catch (error) {
        console.error("Failed to fetch Polar products:", error);
        return json({ error: "Failed to fetch products" }, { status: 500 });
    }
};
