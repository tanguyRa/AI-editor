---
stepsCompleted: ['step-01-init', 'step-02-context', 'step-03-starter', 'step-04-decisions', 'step-05-patterns', 'step-06-structure', 'step-07-validation', 'step-08-complete']
status: 'complete'
completedAt: '2026-01-14'
inputDocuments:
  - docs/planning-artifacts/prd.md
  - docs/index.md
  - docs/architecture-api.md
  - docs/api-contracts.md
  - docs/data-models.md
  - docs/integration-architecture.md
workflowType: 'architecture'
project_name: 'AI-editor'
user_name: 'Tanguy'
date: '2026-01-14'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements (36 total):**
- Account Management (FR1-FR4): User creation, profile storage - exists
- Authentication (FR5-FR10): Sign-in/out, session validation - exists
- Email Verification (FR11-FR15): Send/verify emails, token expiry - new
- Password Management (FR16-FR21): Reset flow, change password, session revocation - new
- Session Management (FR22-FR27): List sessions, revoke specific/all, JWT caching - new
- Email Delivery (FR28-FR32): Provider abstraction, SendGrid + Resend - new
- Configuration (FR33-FR36): Env-driven config, bcrypt, CORS - partially exists

**Non-Functional Requirements (19 total):**
- Performance (NFR1-NFR4): <5ms session validation via JWT caching, async email
- Security (NFR5-NFR11): bcrypt cost 10+, HttpOnly cookies, timing-safe token validation
- Integration (NFR12-NFR15): Exact better-auth API compatibility, shared DB schema
- Reliability (NFR16-NFR19): Graceful failures, email errors don't block operations

**Scale & Complexity:**
- Primary domain: API backend (Go, PostgreSQL, REST)
- Complexity level: Medium
- New endpoints: 8 (on top of 4 existing)
- New database queries: ~10 (verification CRUD, session listing/revocation)

### Technical Constraints & Dependencies

**Hard Constraints (better-auth compatibility):**
- API response shapes must match exactly
- Database schema must remain compatible
- Cookie name and format: `better-auth.session_token`
- Error format: `{ "error": { "code": "...", "message": "..." } }`

**Existing Infrastructure:**
- PostgreSQL with pgx driver, sqlc code generation
- `verification` table exists (needs queries)
- `session` table has ipAddress/userAgent fields (ready for listing)
- Database adapter pattern established

**Technical Debt to Address:**
- Monolithic `internal/server/main.go` needs restructuring
- No middleware chain - CORS is inline
- No organized handler packages
- No email infrastructure

### Cross-Cutting Concerns

| Concern | Affects | Architectural Impact |
|---------|---------|---------------------|
| Token Generation | Verification, Password Reset, Sessions | Shared secure token utilities |
| Email Delivery | Verification, Password Reset | Async provider abstraction |
| Session Security | All authenticated endpoints | Middleware, JWT validation |
| Error Handling | All endpoints | Consistent better-auth format |
| Timing Safety | Password reset, email verification | Constant-time comparisons |

### Server Architecture Refactor

**Current State:** Single-file server with inline handlers and CORS

**Target State (Go idioms):**
- `internal/handlers/` - Organized by domain (auth, session, health)
- `internal/middleware/` - CORS, auth validation, request logging
- `internal/router/` - Route registration, middleware chain
- Clean separation of concerns, testable handlers

## Technical Foundation (Brownfield)

### Existing Stack (Unchanged)

| Layer | Technology | Version | Notes |
|-------|------------|---------|-------|
| Language | Go | 1.24 | |
| HTTP Router | net/http ServeMux | stdlib | Pattern matching (Go 1.22+) |
| Database Driver | pgx | 5.8 | PostgreSQL |
| Code Generation | sqlc | latest | Type-safe queries |
| Migrations | golang-migrate | latest | |
| Password Hashing | bcrypt | golang.org/x/crypto | |

### New Dependencies for Phase 1

| Purpose | Library | Rationale |
|---------|---------|-----------|
| JWT/JWE Tokens | `golang-jwt/jwt/v5` | Industry standard, well-maintained |
| Email (SendGrid) | HTTP client (stdlib) | Direct API calls, no SDK bloat |
| Email (Resend) | HTTP client (stdlib) | Direct API calls, no SDK bloat |

### Design Principle

**Prefer stdlib over external dependencies.** Go 1.22+ `http.ServeMux` supports method-based routing (`POST /api/auth/sign-in`) and path parameters - sufficient for this project without chi/gorilla overhead.

## Core Architectural Decisions

### Decision Summary

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Session caching | Compact (base64 + HMAC-SHA256) | better-auth default, fast validation |
| Verification tokens | JWT (signed with secret) | better-auth compatibility |
| Email delivery | Fire-and-forget goroutine | Simple, non-blocking |
| Server structure | Domain-split handlers | Go idioms, testable |

### Authentication & Session Strategy

**Session Token Flow (better-auth compatible):**

1. On sign-in: Create session in DB, return signed cookie cache
2. On request: Validate signature from cookie (no DB hit)
3. On cache expiry (5min default): Refresh from DB, update cookie
4. On sign-out/revoke: Delete from DB (cookie becomes stale on next refresh)

**Cookie Cache Format:**
- Strategy: `compact` (base64url + HMAC-SHA256)
- Cookie name: `better-auth.session_token`
- Cache TTL: 5 minutes (configurable)
- Contains: userId, sessionId, expiresAt

**Implementation:**
- `pkg/session/cache.go` - Sign/verify cookie payload
- Use `BETTER_AUTH_SECRET` for HMAC key

### Verification Token Strategy

**Token Format:** JWT signed with `BETTER_AUTH_SECRET`

**Token Types:**

| Type | Identifier Format | Default Expiry |
|------|------------------|----------------|
| Email verification | `email-verification:{userId}` | 1 hour |
| Password reset | `reset-password:{userId}` | 1 hour |

**Storage:** `verification` table (already exists)

**Implementation:**
- `pkg/token/verification.go` - Create/validate verification JWTs
- Store in DB for single-use enforcement
- Delete after use or expiration

### Email Provider Architecture

**Pattern:** Fire-and-forget with logging

```go
// Non-blocking email send
go func() {
    if err := emailService.Send(ctx, email); err != nil {
        log.Error("email send failed", "error", err, "to", email.To)
    }
}()
```

**Interface:**

```go
type EmailProvider interface {
    Send(ctx context.Context, email Email) error
}
```

**Implementations:** `internal/email/sendgrid.go`, `internal/email/resend.go`

### Server Package Structure

```
internal/
├── handlers/
│   ├── auth.go          # sign-up, sign-in, sign-out, get-session
│   ├── verification.go  # send-verification, verify-email
│   ├── password.go      # forget-password, reset-password, change-password
│   └── session.go       # list-sessions, revoke-session, revoke-others
├── middleware/
│   ├── cors.go
│   ├── auth.go          # session validation
│   └── chain.go         # middleware composition
├── server/
│   └── server.go        # route registration, Start()
└── email/
    ├── provider.go      # interface
    ├── sendgrid.go
    └── resend.go
```

## Implementation Patterns & Consistency Rules

### Locked Patterns (better-auth Compatibility)

| Area | Pattern | Example |
|------|---------|---------|
| Database tables | snake_case | `user`, `session`, `verification` |
| Database columns | snake_case | `user_id`, `expires_at`, `email_verified` |
| JSON fields | camelCase | `userId`, `expiresAt`, `emailVerified` |
| API endpoints | better-auth paths | `/api/auth/sign-in/email` |
| Error format | better-auth shape | `{ "error": { "code": "...", "message": "..." } }` |
| Cookie name | better-auth | `better-auth.session_token` |

### Go Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Exported functions | PascalCase | `HandleSignIn`, `ValidateSession` |
| Unexported functions | camelCase | `validatePassword`, `createToken` |
| Files | snake_case | `auth_handler.go`, `session_cache.go` |
| Test files | `*_test.go` | `auth_handler_test.go` |
| Packages | lowercase | `handlers`, `middleware`, `email` |

### Handler Pattern

**Struct receiver with dependency injection:**

```go
// internal/handlers/auth.go
type AuthHandler struct {
    db     *database.Adapter
    email  email.Provider
    config *config.Config
}

func NewAuthHandler(db *database.Adapter, email email.Provider, cfg *config.Config) *AuthHandler {
    return &AuthHandler{db: db, email: email, config: cfg}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
    // handler implementation
}
```

### Response Helpers

**Location:** `internal/handlers/response.go`

```go
// Success response
func respondJSON(w http.ResponseWriter, data any, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// Error response (better-auth format)
func respondError(w http.ResponseWriter, code string, message string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]any{
        "error": map[string]string{
            "code":    code,
            "message": message,
        },
    })
}
```

### Validation Pattern

**Inline validation in handlers (no external library):**

```go
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
    var req SignInRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, "INVALID_BODY", "Invalid request body", 400)
        return
    }
    if req.Email == "" || req.Password == "" {
        respondError(w, "INVALID_BODY", "Email and password required", 400)
        return
    }
    // business logic...
}
```

### Logging Pattern

**Use `slog` (stdlib, structured):**

```go
import "log/slog"

// Info level for normal operations
slog.Info("user signed in", "userId", user.ID, "email", user.Email)

// Error level with context
slog.Error("sign-in failed", "error", err, "email", req.Email)

// Warn for recoverable issues
slog.Warn("email send failed, continuing", "error", err, "to", email.To)
```

### Error Codes

**Standard error codes for better-auth compatibility:**

| Code | HTTP Status | Usage |
|------|-------------|-------|
| `INVALID_BODY` | 400 | Malformed JSON or missing fields |
| `USER_ALREADY_EXISTS` | 400 | Email already registered |
| `INVALID_CREDENTIALS` | 401 | Wrong email/password |
| `UNAUTHORIZED` | 401 | No valid session |
| `INVALID_TOKEN` | 400 | Expired or invalid verification token |
| `USER_NOT_FOUND` | 404 | Email not found (only for non-sensitive ops) |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

### Test Pattern

**Co-located tests with table-driven style:**

```go
// internal/handlers/auth_handler_test.go
func TestAuthHandler_SignIn(t *testing.T) {
    tests := []struct {
        name       string
        body       string
        wantStatus int
        wantCode   string
    }{
        {"valid credentials", `{"email":"test@example.com","password":"pass"}`, 200, ""},
        {"missing email", `{"password":"pass"}`, 400, "INVALID_BODY"},
        {"wrong password", `{"email":"test@example.com","password":"wrong"}`, 401, "INVALID_CREDENTIALS"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Project Structure & Boundaries

### Target API Directory Structure

```
api/
├── cmd/
│   └── server/
│       └── main.go                    # Entry point, wire dependencies
├── internal/
│   ├── config/
│   │   └── config.go                  # Config loading (exists)
│   ├── handlers/
│   │   ├── auth.go                    # FR1-FR10: sign-up, sign-in, sign-out, get-session
│   │   ├── verification.go            # FR11-FR15: send-verification, verify-email
│   │   ├── password.go                # FR16-FR21: forget-password, reset-password, change-password
│   │   ├── session.go                 # FR22-FR27: list-sessions, revoke-session, revoke-others
│   │   ├── health.go                  # Health check
│   │   └── response.go                # respondJSON, respondError helpers
│   ├── middleware/
│   │   ├── cors.go                    # CORS handling
│   │   ├── auth.go                    # Session validation middleware
│   │   └── chain.go                   # Middleware composition
│   ├── repository/
│   │   ├── db.go                      # DBTX interface (exists)
│   │   ├── models.go                  # sqlc models (exists)
│   │   ├── users.sql.go               # User queries (exists)
│   │   ├── sessions.sql.go            # Session queries (exists, needs extension)
│   │   └── verifications.sql.go       # NEW: Verification queries
│   ├── email/
│   │   ├── provider.go                # EmailProvider interface
│   │   ├── sendgrid.go                # SendGrid implementation
│   │   └── resend.go                  # Resend implementation
│   └── server/
│       └── server.go                  # Route registration, middleware stack
├── pkg/
│   ├── database/
│   │   ├── adapter.go                 # Database adapter (exists, needs extension)
│   │   └── main.go                    # DB init (exists)
│   ├── session/
│   │   └── cache.go                   # NEW: Cookie cache sign/verify
│   └── token/
│       └── verification.go            # NEW: JWT verification tokens
├── db/
│   ├── migrations/
│   │   └── 000001_init_auth.up.sql    # Existing schema
│   └── queries/
│       ├── users.sql                  # User queries (exists)
│       ├── sessions.sql               # Session queries (extend)
│       └── verifications.sql          # NEW: Verification queries
└── sqlc.yml
```

### PRD to File Mapping

| PRD Requirement | Files |
|----------------|-------|
| FR1-FR4 (Account) | `handlers/auth.go`, `repository/users.sql.go` |
| FR5-FR10 (Auth) | `handlers/auth.go`, `middleware/auth.go`, `pkg/session/cache.go` |
| FR11-FR15 (Verification) | `handlers/verification.go`, `pkg/token/verification.go`, `repository/verifications.sql.go` |
| FR16-FR21 (Password) | `handlers/password.go`, `pkg/token/verification.go` |
| FR22-FR27 (Sessions) | `handlers/session.go`, `repository/sessions.sql.go` |
| FR28-FR32 (Email) | `internal/email/*.go` |
| NFR1 (<5ms validation) | `pkg/session/cache.go`, `middleware/auth.go` |

### New SQL Queries Required

**verifications.sql:**

```sql
-- name: CreateVerification :one
INSERT INTO verification (identifier, value, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVerification :one
SELECT * FROM verification WHERE identifier = $1 AND value = $2;

-- name: DeleteVerification :exec
DELETE FROM verification WHERE id = $1;

-- name: DeleteExpiredVerifications :exec
DELETE FROM verification WHERE expires_at < NOW();
```

**sessions.sql (additions):**

```sql
-- name: ListUserSessions :many
SELECT * FROM session WHERE user_id = $1 ORDER BY created_at DESC;

-- name: DeleteSessionByID :exec
DELETE FROM session WHERE id = $1 AND user_id = $2;

-- name: DeleteOtherSessions :exec
DELETE FROM session WHERE user_id = $1 AND id != $2;
```

### Architectural Boundaries

| Boundary | Rule |
|----------|------|
| Handlers → Repository | Via database adapter interface only |
| Handlers → Email | Via `EmailProvider` interface only |
| External → Handlers | Through middleware chain (CORS → Auth → Handler) |
| Session validation | Middleware decodes cookie cache, passes session to handler context |
| Token operations | `pkg/token` and `pkg/session` are standalone, no DB dependency |

### Data Flow

```
Request → CORS Middleware → Auth Middleware → Handler → Repository → PostgreSQL
                              ↓
                        Cookie Cache (no DB)
                              or
                        DB Lookup (cache miss)
```

## Architecture Validation

### Coherence Validation ✅

**Decision Compatibility:**
- Go 1.24 + stdlib `http.ServeMux` + pgx 5.8 + sqlc + `golang-jwt/jwt/v5` - all compatible
- No version conflicts or incompatibilities identified

**Pattern Consistency:**
- Handler patterns (struct receivers, DI) align with Go idioms
- Naming conventions (snake_case DB, camelCase JSON, PascalCase exports) are consistent
- Error handling follows better-auth format throughout

**Structure Alignment:**
- Standard Go project layout (`cmd/`, `internal/`, `pkg/`, `db/`)
- Clear separation: handlers → adapter → repository → database
- Boundaries enforced via interfaces

### Requirements Coverage ✅

| Requirement Group | Status | Implementation Location |
|------------------|--------|------------------------|
| FR1-FR10 (Account/Auth) | ✅ Covered | `handlers/auth.go` |
| FR11-FR15 (Email Verification) | ✅ Covered | `handlers/verification.go`, `pkg/token/` |
| FR16-FR21 (Password Management) | ✅ Covered | `handlers/password.go`, `pkg/token/` |
| FR22-FR27 (Session Management) | ✅ Covered | `handlers/session.go`, `pkg/session/` |
| FR28-FR32 (Email Delivery) | ✅ Covered | `internal/email/` |
| NFR1-NFR4 (Performance) | ✅ Covered | Cookie cache, async email |
| NFR5-NFR11 (Security) | ✅ Covered | bcrypt, HMAC, HttpOnly, timing-safe |
| NFR12-NFR15 (Integration) | ✅ Covered | better-auth compatibility locked |
| NFR16-NFR19 (Reliability) | ✅ Covered | Graceful error handling patterns |

### Implementation Readiness ✅

**Completeness Checklist:**

- [x] All critical technology decisions documented with versions
- [x] Implementation patterns have concrete code examples
- [x] Every functional requirement mapped to specific file
- [x] Architectural boundaries defined via interfaces
- [x] better-auth compatibility requirements locked
- [x] SQL queries specified for new functionality
- [x] Error codes standardized

### Deferred Items (Post-MVP)

| Item | Rationale |
|------|-----------|
| Rate limiting | Explicitly Phase 2 per PRD |
| OAuth providers | Phase 2 scope |
| Email templates | Implementation detail, not architectural |

### Architecture Readiness

**Status:** READY FOR IMPLEMENTATION

**Confidence Level:** High - well-scoped brownfield project with clear boundaries

**Key Strengths:**
- Clear better-auth compatibility constraints reduce ambiguity
- Existing codebase provides proven patterns to follow
- Focused MVP scope (4 user journeys, 8 new endpoints)
- Strong separation of concerns via interfaces

## Implementation Handoff

### First Implementation Priority

**Server Refactor (Foundation):**
1. Extract handlers from `internal/server/main.go` → `internal/handlers/`
2. Extract CORS to `internal/middleware/cors.go`
3. Create `internal/handlers/response.go` with `respondJSON`, `respondError`
4. Wire new structure in `internal/server/server.go`

**Then New Functionality:**
1. Add verification queries → `db/queries/verifications.sql` → `make sqlc`
2. Extend session queries → `db/queries/sessions.sql` → `make sqlc`
3. Implement `pkg/session/cache.go` (cookie signing)
4. Implement `pkg/token/verification.go` (JWT tokens)
5. Implement `internal/email/` (provider interface + implementations)
6. Add new handlers one by one, following patterns

### AI Agent Guidelines

When implementing stories against this architecture:

1. **Read this document first** - all decisions are binding
2. **Follow better-auth API shapes exactly** - compatibility is non-negotiable
3. **Use the patterns defined** - handler structs, response helpers, slog logging
4. **Respect boundaries** - handlers → adapter → repository, never skip layers
5. **Match file locations** - put code where the structure says it goes

### Document Maintenance

Update this architecture document when:
- New dependencies are added
- Patterns need clarification based on implementation experience
- Scope changes require architectural adjustments

---

**Architecture Status:** COMPLETE ✅
**Ready for:** Implementation Phase

