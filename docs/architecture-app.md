# Frontend Architecture (app/)

## Overview

The frontend is a Nuxt 4 application using Vue 3 with TypeScript. It provides a modern, SSR-capable web interface with authentication capabilities.

## Technology Stack

| Category | Technology | Version |
|----------|------------|---------|
| Framework | Nuxt | 4.2.2 |
| UI Framework | Vue | 3.x |
| Language | TypeScript | 5.9.3 |
| UI Library | Nuxt UI | 4.3.0 |
| Auth Client | better-auth | 1.4.10 |
| Validation | Zod | 4.3.5 |
| Package Manager | Bun | latest |
| Icons | Lucide, Simple Icons | - |

## Directory Structure

```
app/
├── app/                     # Application source
│   ├── app.vue              # Root component
│   ├── app.config.ts        # App configuration
│   ├── assets/
│   │   └── css/
│   │       └── main.css     # Global styles
│   ├── components/
│   │   ├── AppLogo.vue      # Logo component
│   │   ├── TemplateMenu.vue # Navigation menu
│   │   └── TemplateComponent.vue
│   ├── lib/
│   │   ├── auth.ts          # Server-side better-auth setup
│   │   └── auth-client.ts   # Client-side auth utilities
│   ├── middleware/
│   │   └── auth.ts          # Auth route middleware
│   └── pages/
│       ├── index.vue        # Landing page
│       ├── login.vue        # Login page
│       ├── register.vue     # Registration page
│       ├── dashboard.vue    # Protected dashboard
│       └── (template)/
│           └── test.vue     # Template test page
├── server/
│   └── api/
│       └── auth/
│           ├── [...all].ts  # Catch-all auth handler
│           └── health.ts    # Health check endpoint
├── public/                  # Static assets
├── nuxt.config.ts           # Nuxt configuration
├── package.json             # Dependencies
├── tsconfig.json            # TypeScript config
├── eslint.config.mjs        # ESLint config
└── Dockerfile               # Container build
```

## Key Components

### Root Component (app/app.vue)

The root component sets up the app shell with:
- `<UApp>` - Nuxt UI app wrapper
- `<UHeader>` - Navigation header with logo and menu
- `<UMain>` - Main content area with `<NuxtPage>`
- `<UFooter>` - Footer section

### Authentication

#### Server-side (lib/auth.ts)

```typescript
// PostgreSQL connection pool for better-auth
const pool = new Pool({
    connectionString: process.env.DATABASE_URL,
    max: 20,
    idleTimeoutMillis: 30000,
    maxLifetimeSeconds: 3600,
});

export const auth = betterAuth({
    baseUrl: "http://localhost:8080",
    database: pool,
    emailAndPassword: { enabled: true },
});
```

#### Client-side (lib/auth-client.ts)

```typescript
export const authClient = createAuthClient({});

export const {
    signIn,
    signOut,
    signUp,
    useSession,
    requestPasswordReset,
    resetPassword,
} = authClient;
```

### Route Middleware (middleware/auth.ts)

Protects routes by checking session:

```typescript
export default defineNuxtRouteMiddleware(async (to, from) => {
    const { data: session } = await authClient.useSession(useFetch);
    if (!session.value) {
        return navigateTo("/");
    }
});
```

### Pages

| Page | Route | Description |
|------|-------|-------------|
| index.vue | `/` | Landing page with hero and features |
| login.vue | `/login` | Login form with email/password |
| register.vue | `/register` | Registration form |
| dashboard.vue | `/dashboard` | Protected user dashboard |

## Build & Deployment

### Development

```bash
bun install
bun run dev
```

### Production Build

Multi-stage Docker build:

1. **Staging**: Install dependencies with `bun install`
2. **Build**: Run `bun --bun run build`
3. **Production**: Copy `.output` folder, run with `bun --bun run /app/server/index.mjs`

### Docker

```dockerfile
FROM oven/bun:1 AS production
COPY --from=build /app/.output /app
ENTRYPOINT [ "bun", "--bun", "run", "/app/server/index.mjs" ]
```

## Configuration

### nuxt.config.ts

```typescript
export default defineNuxtConfig({
  modules: ['@nuxt/eslint', '@nuxt/ui'],
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],
  routeRules: { '/': { prerender: true } },
});
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `PORT` | Server port (default: 3000) |

## UI Components Used

From Nuxt UI:
- `UApp`, `UHeader`, `UMain`, `UFooter`
- `UPageHero`, `UPageSection`, `UPageCTA`, `UPageCard`
- `UAuthForm`, `ULink`, `UButton`
- `UColorModeButton`, `USeparator`
