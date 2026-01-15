# Story 1.3: Session Validation & Sign-Out

Status: ready-for-dev

## Story

As an **authenticated user**,
I want to **check my session status and sign out**,
So that **I can verify I'm logged in and securely end my session**.

## Acceptance Criteria

1. **Given** a user with a valid session cookie **When** they request get-session endpoint **Then** the session is looked up in PostgreSQL by token **And** the response returns `{ user, session }` if valid and not expired **And** the session `expiresAt` is checked against current time

2. **Given** a user with an expired or invalid session cookie **When** they request get-session endpoint **Then** the response returns `null` (not an error)

3. **Given** a user with no session cookie **When** they request get-session endpoint **Then** the response returns `null`

4. **Given** an authenticated user **When** they request sign-out endpoint **Then** the session is deleted from PostgreSQL **And** the `better-auth.session_token` cookie is cleared **And** the response returns `{ "success": true }`

5. **Given** a user without a valid session **When** they request sign-out endpoint **Then** the response returns `{ "success": true }` (idempotent)

## Tasks / Subtasks

- [ ] Task 1: Implement GetSession handler (AC: 1, 2, 3)
  - [ ] 1.1: Add GetSession method to AuthHandler struct
  - [ ] 1.2: Extract session token from `better-auth.session_token` cookie
  - [ ] 1.3: If no cookie, return `null` response (not error)
  - [ ] 1.4: Look up session by token (GetSessionByToken)
  - [ ] 1.5: Check if session is expired (expiresAt < now)
  - [ ] 1.6: If expired or not found, return `null` response
  - [ ] 1.7: Look up user by session.userId (GetUserById)
  - [ ] 1.8: Return `{ user, session }` in better-auth format

- [ ] Task 2: Add GetUserById query (AC: 1)
  - [ ] 2.1: Add GetUserById query to db/queries/users.sql
  - [ ] 2.2: Run `make sqlc` to regenerate

- [ ] Task 3: Implement SignOut handler (AC: 4, 5)
  - [ ] 3.1: Add SignOut method to AuthHandler struct
  - [ ] 3.2: Extract session token from cookie
  - [ ] 3.3: If no cookie, return `{ "success": true }` (idempotent)
  - [ ] 3.4: Delete session from database (DeleteSessionByToken)
  - [ ] 3.5: Clear the `better-auth.session_token` cookie (set empty, expired)
  - [ ] 3.6: Return `{ "success": true }`

- [ ] Task 4: Add DeleteSessionByToken query (AC: 4)
  - [ ] 4.1: Add DeleteSessionByToken query to db/queries/sessions.sql
  - [ ] 4.2: Run `make sqlc` to regenerate

- [ ] Task 5: Wire routes in server (AC: 1, 4)
  - [ ] 5.1: Add GET /api/auth/get-session route
  - [ ] 5.2: Add POST /api/auth/sign-out route

- [ ] Task 6: Create auth middleware for future stories (Foundation)
  - [ ] 6.1: Create `internal/middleware/auth.go`
  - [ ] 6.2: Implement RequireAuth middleware that validates session
  - [ ] 6.3: Store user/session in request context for handlers
  - [ ] 6.4: This will be used by stories 2.2, 4.1, 5.1-5.3

- [ ] Task 7: Write tests (AC: 1, 2, 3, 4, 5)
  - [ ] 7.1: Test get-session with valid session returns user + session
  - [ ] 7.2: Test get-session with expired session returns null
  - [ ] 7.3: Test get-session with no cookie returns null
  - [ ] 7.4: Test get-session with invalid token returns null
  - [ ] 7.5: Test sign-out deletes session and clears cookie
  - [ ] 7.6: Test sign-out without session returns success (idempotent)
  - [ ] 7.7: Test auth middleware blocks unauthenticated requests

## Dev Notes

### Critical Architecture Requirements

**better-auth Compatibility (NON-NEGOTIABLE):**
- Cookie name: `better-auth.session_token`
- get-session returns `null` for invalid/missing sessions (NOT an error)
- Error format: `{ "error": { "code": "ERROR_CODE", "message": "Human readable" } }`

**Endpoints:**
- `GET /api/auth/get-session`
- `POST /api/auth/sign-out`

### GetSession Response Formats

**Valid Session Response (200):**
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

**Invalid/Missing Session Response (200):**
```json
null
```

Note: Returns HTTP 200 with body `null`, NOT an error response.

### SignOut Response Format

**Success Response (200):**
```json
{
  "success": true
}
```

### Cookie Clearing

```go
func clearSessionCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "better-auth.session_token",
        Value:    "",
        Path:     "/",
        MaxAge:   -1, // Delete immediately
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
        Secure:   false, // Set true in production
    })
}
```

### Auth Middleware Pattern

```go
// internal/middleware/auth.go

type contextKey string

const (
    UserContextKey    contextKey = "user"
    SessionContextKey contextKey = "session"
)

func RequireAuth(queries *repository.Queries, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("better-auth.session_token")
        if err != nil {
            respondError(w, "UNAUTHORIZED", "Authentication required", 401)
            return
        }

        session, err := queries.GetSessionByToken(r.Context(), cookie.Value)
        if err != nil || session.ExpiresAt.Before(time.Now()) {
            respondError(w, "UNAUTHORIZED", "Invalid or expired session", 401)
            return
        }

        user, err := queries.GetUserById(r.Context(), session.UserID)
        if err != nil {
            respondError(w, "UNAUTHORIZED", "User not found", 401)
            return
        }

        // Add to context for handler access
        ctx := context.WithValue(r.Context(), UserContextKey, user)
        ctx = context.WithValue(ctx, SessionContextKey, session)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Helper to get user from context
func GetUserFromContext(ctx context.Context) (*repository.User, bool) {
    user, ok := ctx.Value(UserContextKey).(*repository.User)
    return user, ok
}

func GetSessionFromContext(ctx context.Context) (*repository.Session, bool) {
    session, ok := ctx.Value(SessionContextKey).(*repository.Session)
    return session, ok
}
```

### New SQL Queries Required

**users.sql:**
```sql
-- name: GetUserById :one
SELECT * FROM "user" WHERE id = $1;
```

**sessions.sql:**
```sql
-- name: DeleteSessionByToken :exec
DELETE FROM session WHERE token = $1;
```

### Handler Implementations

**GetSession:**
```go
func (h *AuthHandler) GetSession(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        respondError(w, "METHOD_NOT_ALLOWED", "Method not allowed", 405)
        return
    }

    cookie, err := r.Cookie("better-auth.session_token")
    if err != nil {
        // No cookie = return null (not error)
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("null"))
        return
    }

    ctx := r.Context()

    session, err := h.queries.GetSessionByToken(ctx, cookie.Value)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("null"))
        return
    }

    // Check expiry
    if session.ExpiresAt.Before(time.Now()) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("null"))
        return
    }

    user, err := h.queries.GetUserById(ctx, session.UserID)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("null"))
        return
    }

    respondJSON(w, map[string]any{
        "user":    formatUser(user),
        "session": formatSession(session),
    }, 200)
}
```

**SignOut:**
```go
func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        respondError(w, "METHOD_NOT_ALLOWED", "Method not allowed", 405)
        return
    }

    cookie, err := r.Cookie("better-auth.session_token")
    if err == nil && cookie.Value != "" {
        // Delete session from DB (ignore errors - idempotent)
        h.queries.DeleteSessionByToken(r.Context(), cookie.Value)
    }

    // Always clear cookie
    clearSessionCookie(w)

    respondJSON(w, map[string]bool{"success": true}, 200)
}
```

### Session Expiry

- Default session duration: 7 days
- Check `expiresAt` field against current time
- Expired sessions should be treated as invalid (return null)

### References

- [Source: docs/planning-artifacts/architecture.md#Handler Pattern]
- [Source: docs/planning-artifacts/epics.md#Story 1.3]
- [Source: docs/planning-artifacts/architecture.md#Authentication & Session Strategy]

## Dev Agent Record

### Agent Model Used

(To be filled by dev agent)

### Debug Log References

(To be filled during implementation)

### Completion Notes List

(To be filled during implementation)

### Change Log

(To be filled during implementation)

### File List

**Expected New Files:**
- api/internal/middleware/auth.go

**Expected Modified Files:**
- api/internal/handlers/auth.go (add GetSession, SignOut methods)
- api/internal/handlers/auth_test.go (add tests)
- api/internal/server/server.go (add routes)
- api/db/queries/users.sql (add GetUserById)
- api/db/queries/sessions.sql (add DeleteSessionByToken)

**Expected Generated Files (via sqlc):**
- api/internal/repository/users.sql.go (regenerated)
- api/internal/repository/sessions.sql.go (regenerated)
