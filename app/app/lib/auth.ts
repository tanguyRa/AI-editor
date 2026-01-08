import { ClientBase, Pool } from "pg";
import { betterAuth } from "better-auth";

const pool = new Pool({
    connectionString: process.env.DATABASE_URL,
    // Pool configuration
    max: 20, // Maximum 20 concurrent connections
    idleTimeoutMillis: 30000, // Close idle connections after 30s
    maxLifetimeSeconds: 3600, // Max connection lifetime 1 hour
    // connectionTimeout: 10, // Connection timeout 10s
});


export const auth = betterAuth({
    baseUrl: "http://localhost:8080",
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
})