# Story 1.1: User Registration with Database Persistence

Status: review

## Story

As a **new user**,
I want to **create an account with my email, password, and name**,
So that **I can access the application with persistent credentials**.

## Acceptance Criteria

1. **Given** a user visits the registration endpoint **When** they submit valid email, password, and name **Then** a new user record is created in PostgreSQL with UUID **And** a new account record is created with hashed password (bcrypt cost 10+) **And** a new session is created in the database **And** the response returns `{ user, session }` matching better-auth format **And** the `better-auth.session_token` cookie is set (HttpOnly, SameSite=Lax)

2. **Given** a user tries to register with an existing email **When** they submit the registration request **Then** the response returns 400 with `{ "error": { "code": "USER_ALREADY_EXISTS", "message": "..." } }`

3. **Given** a user submits invalid or missing fields **When** the request is processed **Then** the response returns 400 with `{ "error": { "code": "INVALID_BODY", "message": "..." } }`

## Tasks / Subtasks

- [x] Task 1: Server structure refactor - extract handlers from monolith (AC: Foundation)
  - [x] 1.1: Create `internal/handlers/auth.go` with AuthHandler struct
  - [x] 1.2: Move response helpers to `internal/handlers/response.go` (respondJSON, respondError)
  - [x] 1.3: Create `internal/middleware/cors.go` (extract from monolith)
  - [x] 1.4: Update `internal/server/server.go` to wire new handler structure

- [x] Task 2: Add account SQL queries for credential storage (AC: 1)
  - [x] 2.1: Create `db/queries/accounts.sql` with CreateAccount, GetAccountByUserIdAndProvider queries
  - [x] 2.2: Run `make sqlc` to generate repository code
  - [x] 2.3: Verify generated `internal/repository/accounts.sql.go`

- [x] Task 3: Add GetSessionByToken query (AC: 1)
  - [x] 3.1: Add `GetSessionByToken` query to `db/queries/sessions.sql`
  - [x] 3.2: Run `make sqlc` to regenerate

- [x] Task 4: Implement PostgreSQL-backed sign-up handler (AC: 1, 2, 3)
  - [x] 4.1: Create AuthHandler struct with db dependency injection
  - [x] 4.2: Implement SignUp handler using repository queries
  - [x] 4.3: Hash password with bcrypt (cost 10+)
  - [x] 4.4: Create user, account, and session in transaction
  - [x] 4.5: Capture IP address from X-Forwarded-For or RemoteAddr
  - [x] 4.6: Capture User-Agent header
  - [x] 4.7: Set better-auth.session_token cookie

- [x] Task 5: Write tests for sign-up endpoint (AC: 1, 2, 3)
  - [x] 5.1: Test successful registration returns user + session
  - [x] 5.2: Test duplicate email returns USER_ALREADY_EXISTS
  - [x] 5.3: Test missing fields returns INVALID_BODY
  - [x] 5.4: Test password is hashed (not stored plaintext)

## Dev Notes

### Critical Architecture Requirements

**better-auth Compatibility (NON-NEGOTIABLE):**
- Cookie name: `better-auth.session_token`
- Error format: `{ "error": { "code": "ERROR_CODE", "message": "Human readable" } }`
- JSON response fields: camelCase (`userId`, `emailVerified`, `expiresAt`)
- Database columns: snake_case (already correct in schema)

**Endpoint:** `POST /api/auth/sign-up/email`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "User Name"
}
```

**Success Response (200):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "User Name",
    "image": null,
    "emailVerified": false,
    "createdAt": "2026-01-14T...",
    "updatedAt": "2026-01-14T..."
  },
  "session": {
    "id": "uuid",
    "userId": "uuid",
    "token": "session-token",
    "expiresAt": "2026-01-21T...",
    "ipAddress": "127.0.0.1",
    "userAgent": "Mozilla/5.0...",
    "createdAt": "2026-01-14T...",
    "updatedAt": "2026-01-14T..."
  }
}
```

### Existing Code to Reuse

**Response helpers already exist:** `internal/handlers/utils.go` has `writeJSON` and `writeError` - use these.

**Repository queries exist:**
- `CreateUser` - `internal/repository/users.sql.go`
- `GetUserByEmail` - `internal/repository/users.sql.go`
- `CreateSession` - `internal/repository/sessions.sql.go`

**Missing queries (need to add):**
- `CreateAccount` - for storing hashed password
- `GetAccountByUserIdAndProvider` - for sign-in (future story)
- `GetSessionByToken` - for session validation

### Database Schema (already exists)

```sql
-- user table
CREATE TABLE "user" (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    "emailVerified" BOOLEAN NOT NULL DEFAULT false,
    image VARCHAR,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- account table (for credentials)
CREATE TABLE "account" (
    id UUID PRIMARY KEY,
    "userId" UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    "accountId" VARCHAR NOT NULL,
    "providerId" VARCHAR NOT NULL,
    password VARCHAR, -- bcrypt hashed
    ...
);

-- session table
CREATE TABLE "session" (
    id UUID PRIMARY KEY,
    token VARCHAR(255) UNIQUE NOT NULL,
    "userId" UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    "expiresAt" TIMESTAMPTZ NOT NULL,
    "ipAddress" VARCHAR(45),
    "userAgent" VARCHAR,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Handler Pattern to Follow

```go
// internal/handlers/auth.go
type AuthHandler struct {
    queries *repository.Queries
}

func NewAuthHandler(queries *repository.Queries) *AuthHandler {
    return &AuthHandler{queries: queries}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
    // 1. Decode request body
    // 2. Validate required fields
    // 3. Check if email exists (GetUserByEmail)
    // 4. Hash password with bcrypt
    // 5. Create user (CreateUser)
    // 6. Create account (CreateAccount)
    // 7. Create session (CreateSession)
    // 8. Set cookie
    // 9. Return response
}
```

### Security Requirements

- bcrypt cost: minimum 10 (use `bcrypt.DefaultCost` which is 10)
- Cookies: `HttpOnly: true`, `SameSite: Lax`, `Secure: true` in production
- Session expiry: 7 days default
- UUID generation: use `github.com/google/uuid`

### Error Codes

| Code | HTTP Status | When |
|------|-------------|------|
| `INVALID_BODY` | 400 | Malformed JSON or missing email/password/name |
| `USER_ALREADY_EXISTS` | 400 | Email already registered |
| `INTERNAL_ERROR` | 500 | Database or bcrypt failure |

### Project Structure Notes

**Current monolith location:** `internal/server/main.go` (372 lines)

**Target structure after refactor:**
```
api/
├── internal/
│   ├── handlers/
│   │   ├── auth.go          # AuthHandler with SignUp, SignIn, SignOut, GetSession
│   │   └── response.go      # writeJSON, writeError (move from utils.go)
│   ├── middleware/
│   │   └── cors.go          # Extract corsMiddleware from server/main.go
│   └── server/
│       └── server.go        # Route registration only
```

### Testing Approach

Use table-driven tests:
```go
func TestAuthHandler_SignUp(t *testing.T) {
    tests := []struct {
        name       string
        body       string
        wantStatus int
        wantCode   string
    }{
        {"valid registration", `{"email":"new@test.com","password":"pass123","name":"Test"}`, 200, ""},
        {"missing email", `{"password":"pass","name":"Test"}`, 400, "INVALID_BODY"},
        {"duplicate email", `{"email":"existing@test.com","password":"pass","name":"Test"}`, 400, "USER_ALREADY_EXISTS"},
    }
    // ...
}
```

### References

- [Source: docs/planning-artifacts/architecture.md#Handler Pattern]
- [Source: docs/planning-artifacts/architecture.md#Server Package Structure]
- [Source: docs/planning-artifacts/epics.md#Story 1.1]
- [Source: docs/project-context.md#Critical Rules]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

None - implementation completed without blocking issues.

### Completion Notes List

- **Task 1 Complete:** Refactored server structure by extracting handlers from monolith. Created AuthHandler struct with dependency injection pattern, moved response helpers to dedicated file, extracted CORS middleware, and updated server wiring to use new handler structure.

- **Task 2 Complete:** Added account SQL queries (CreateAccount, GetAccountByUserIdAndProvider, GetAccountById, DeleteAccount) in `db/queries/accounts.sql`. Generated repository code with sqlc.

- **Task 3 Complete:** Added GetSessionByToken query to `db/queries/sessions.sql` for session validation by cookie token.

- **Task 4 Complete:** Implemented PostgreSQL-backed SignUp handler with:
  - Request validation (email, password, name required)
  - Duplicate email check via GetUserByEmail
  - Password hashing with bcrypt (DefaultCost = 10)
  - User, account, and session creation with proper UUIDs
  - IP address capture from X-Forwarded-For or RemoteAddr
  - User-Agent header capture
  - better-auth.session_token cookie set (HttpOnly, SameSite=Lax)
  - Response in better-auth compatible format (camelCase JSON)

- **Task 5 Complete:** Wrote comprehensive tests covering:
  - Successful registration (verifies user + session response, cookie set)
  - Duplicate email returns USER_ALREADY_EXISTS
  - Missing/empty fields return INVALID_BODY (table-driven tests)
  - Invalid JSON returns INVALID_BODY
  - Password is hashed and verifiable with bcrypt
  - Wrong HTTP method returns METHOD_NOT_ALLOWED

### Change Log

- 2026-01-15: Implemented user registration with database persistence (Story 1.1)

### File List

**New Files:**
- api/internal/handlers/auth.go
- api/internal/handlers/auth_test.go
- api/internal/middleware/cors.go
- api/internal/server/server.go
- api/db/queries/accounts.sql

**Modified Files:**
- api/internal/config/config.go (added DatabaseURL field)
- api/internal/handlers/response.go (renamed from utils.go)
- api/db/queries/sessions.sql (added GetSessionByToken)

**Generated Files (via sqlc):**
- api/internal/repository/accounts.sql.go
- api/internal/repository/sessions.sql.go (regenerated)

**Deleted Files:**
- api/internal/server/main.go (replaced by server.go)
