import { Pool } from "pg";
import { betterAuth } from "better-auth";
import { jwt } from "better-auth/plugins";
import { polar, checkout, portal, usage, webhooks } from "@polar-sh/better-auth";
import { Polar } from "@polar-sh/sdk";


const pool = new Pool({
    connectionString: process.env.DATABASE_URL,
    // Pool configuration
    max: 20, // Maximum 20 concurrent connections
    idleTimeoutMillis: 30000, // Close idle connections after 30s
    maxLifetimeSeconds: 3600, // Max connection lifetime 1 hour
    // connectionTimeout: 10, // Connection timeout 10s
});

export const polarClient = new Polar({
    accessToken: process.env.POLAR_ACCESS_TOKEN,
    server: process.env.POLAR_SERVER === 'production' ? 'production' : 'sandbox'
});

export const auth = betterAuth({
    database: pool,
    advanced: {
        database: {
            generateId: false, // "serial" for auto-incrementing numeric IDs
        },
    },
    emailAndPassword: {
        enabled: true,
        async sendResetPassword(url, user) {
            console.log("Reset password url:", url);
        },
    },
    plugins: [
        jwt(),
        polar({
            client: polarClient,
            createCustomerOnSignUp: true,
            use: [
                checkout({
                    products: [
                        {
                            productId: "e54c3dec-3fa6-4a6d-b359-35fafdfe4b30",
                            slug: "Premium-Annual" // Custom slug for easy reference in Checkout URL, e.g. /checkout/Premium-Annual
                        },
                        {
                            productId: "a741f0a8-929d-4420-8329-2e880fa2ecf8",
                            slug: "Premium" // Custom slug for easy reference in Checkout URL, e.g. /checkout/Premium
                        },
                        {
                            productId: "015ddd64-2330-4fc7-a59d-c8cfcd9751ed",
                            slug: "Free" // Custom slug for easy reference in Checkout URL, e.g. /checkout/Free
                        }
                    ],
                    successUrl: process.env.POLAR_SUCCESS_URL,
                    authenticatedUsersOnly: true
                }),
                portal(),
                usage(),
                // webhooks({
                //     secret: process.env.POLAR_WEBHOOK_SECRET,
                //     //                     onCustomerStateChanged: (payload) => // Triggered when anything regarding a customer changes
                //     //                         onOrderPaid: (payload) => // Triggered when an order was paid (purchase, subscription renewal, etc.)
                //     //                     ...  // Over 25 granular webhook handlers
                //     // onPayload: (payload) => // Catch-all for all events
                // })
            ],
        })
    ]
})