# Integration Architecture

## Overview

AI-Editor follows a multi-part architecture where the frontend (app/) and backend (api/) are separate services that communicate via HTTP REST APIs.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Docker Network                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌──────────────────────┐         ┌──────────────────────┐         │
│   │    app (frontend)    │         │    api (backend)     │         │
│   │    oven/bun:1        │         │    golang:latest     │         │
│   │    Port: 3000        │◄───────►│    Port: 8080        │         │
│   │                      │  HTTP   │                      │         │
│   │  ┌────────────────┐  │   API   │  ┌────────────────┐  │         │
│   │  │  Nuxt Server   │  │         │  │  Go HTTP Server│  │         │
│   │  │  (SSR/SSG)     │  │         │  │  (net/http)    │  │         │
│   │  └────────────────┘  │         │  └────────────────┘  │         │
│   │         │            │         │         │            │         │
│   │         ▼            │         │         ▼            │         │
│   │  ┌────────────────┐  │         │  ┌────────────────┐  │         │
│   │  │  better-auth   │  │         │  │  PostgreSQL    │  │         │
│   │  │  (client)      │  │         │  │  (via pgx)     │  │         │
│   │  └────────────────┘  │         │  └────────────────┘  │         │
│   └──────────────────────┘         └──────────────────────┘         │
│              │                               │                      │
│              │                               │                      │
│              └───────────────┬───────────────┘                      │
│                              │                                      │
│                              ▼                                      │
│                    ┌──────────────────┐                             │
│                    │    PostgreSQL    │                             │
│                    │   (external DB)  │                             │
│                    └──────────────────┘                             │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Integration Points

### 1. Authentication Flow

The frontend uses `better-auth` client library to communicate with the Go backend, which implements better-auth compatible endpoints.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Browser   │     │   app/      │     │   api/      │
│   (User)    │     │   (Nuxt)    │     │   (Go)      │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       │  1. Login Form    │                   │
       │──────────────────►│                   │
       │                   │                   │
       │                   │  2. POST /api/auth/sign-in/email
       │                   │──────────────────►│
       │                   │                   │
       │                   │                   │ 3. Validate credentials
       │                   │                   │    Create session
       │                   │                   │
       │                   │  4. { user, session } + Set-Cookie
       │                   │◄──────────────────│
       │                   │                   │
       │  5. Redirect      │                   │
       │◄──────────────────│                   │
       │                   │                   │
```

### 2. Session Management

Sessions are managed via HTTP-only cookies:

| Component | Role |
|-----------|------|
| Frontend | Sends cookies with each request via `credentials: 'include'` |
| Backend | Creates/validates sessions, sets `better-auth.session_token` cookie |
| Database | Stores session records with expiration |

### 3. API Routes

#### Frontend Server Routes (Nuxt)

The frontend has server-side API routes that proxy to better-auth:

```typescript
// app/server/api/auth/[...all].ts
export default defineEventHandler((event) => {
    return auth.handler(toWebRequest(event));
});
```

#### Backend REST API

The Go backend exposes:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/sign-up/email` | User registration |
| POST | `/api/auth/sign-in/email` | User login |
| GET | `/api/auth/get-session` | Get current session |
| POST | `/api/auth/sign-out` | Logout |
| GET | `/health` | Health check |

## Docker Compose Integration

### Service Dependencies

```yaml
services:
  app:
    depends_on:
      - api
    ports:
      - 3000:3000

  api:
    ports:
      - 8080:8080
```

### Shared Resources

| Resource | Mount Path | Purpose |
|----------|------------|---------|
| `./data/shared` | `/data/shared` (api) | Shared file storage |

### Environment Sharing

Both services share database configuration:

```yaml
environment:
  - DATABASE_URL=${DATABASE_URL}
```

## Communication Protocols

### HTTP/REST

All inter-service communication uses HTTP REST:

- **Content-Type**: `application/json`
- **Authentication**: Cookie-based sessions
- **CORS**: Enabled with credentials

### CORS Configuration

The backend allows cross-origin requests:

```go
w.Header().Set("Access-Control-Allow-Origin", origin)
w.Header().Set("Access-Control-Allow-Credentials", "true")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
```

## Data Flow

### User Registration

```
1. User fills form → app/pages/register.vue
2. Form submits → authClient.signUp.email()
3. Request sent → POST /api/auth/sign-up/email
4. Backend creates user → PostgreSQL
5. Backend creates session → PostgreSQL
6. Response with cookie → Set-Cookie header
7. Redirect to dashboard → app/pages/dashboard.vue
```

### Session Validation

```
1. User visits protected page → app/pages/dashboard.vue
2. Middleware checks session → app/middleware/auth.ts
3. authClient.useSession() → GET /api/auth/get-session
4. Cookie sent automatically
5. Backend validates token → Check expiry, lookup user
6. Response: { user, session } or null
7. Middleware allows/redirects
```

## Error Handling

### Frontend

- Toast notifications for user feedback
- Redirect to login on 401 responses

### Backend

- Structured error responses: `{ error: { code, message } }`
- HTTP status codes for error types

## Scalability Considerations

### Current Architecture

- Stateful sessions stored in PostgreSQL
- Single instance of each service

### Future Improvements

1. **Session Store**: Move to Redis for faster session lookups
2. **Load Balancing**: Add nginx/traefik in front of services
3. **Database**: Connection pooling, read replicas
4. **Caching**: Add caching layer for frequently accessed data
