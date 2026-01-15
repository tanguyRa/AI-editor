# Story 2.2: Send Verification Email

Status: ready-for-dev

## Story

As a **registered user**,
I want to **receive a verification email after registration**,
So that **I can prove ownership of my email address**.

## Acceptance Criteria

1. **Given** an authenticated user with unverified email **When** they request send-verification-email endpoint **Then** a JWT verification token is generated with 1-hour expiry **And** the token is stored in `verification` table with identifier `email-verification:{userId}` **And** an email is sent asynchronously (fire-and-forget) **And** the response returns immediately with `{ "success": true }` (NFR4)

2. **Given** an authenticated user requests verification email **When** a previous verification token exists for this user **Then** the old token is deleted before creating a new one **And** only one active verification token exists per user

3. **Given** the email provider fails to send **When** the error occurs **Then** the error is logged with context (NFR17) **And** the API response is NOT affected (already returned success)

4. **Given** an unauthenticated user **When** they request send-verification-email endpoint **Then** the response returns 401 `UNAUTHORIZED`

5. **Given** an already-verified user **When** they request send-verification-email endpoint **Then** the response returns 400 with appropriate message

## Tasks / Subtasks

- [ ] Task 1: Create verification token package (AC: 1)
  - [ ] 1.1: Create `pkg/token/verification.go`
  - [ ] 1.2: Implement `GenerateVerificationToken(secret, userId, identifier, expiry) (string, error)`
  - [ ] 1.3: Implement `ValidateVerificationToken(secret, token) (claims, error)`
  - [ ] 1.4: Use `golang-jwt/jwt/v5` for JWT operations
  - [ ] 1.5: Include userId, identifier, exp claims in token

- [ ] Task 2: Add verification SQL queries (AC: 1, 2)
  - [ ] 2.1: Create `db/queries/verifications.sql`
  - [ ] 2.2: Add CreateVerification query
  - [ ] 2.3: Add GetVerificationByIdentifier query
  - [ ] 2.4: Add DeleteVerificationByIdentifier query
  - [ ] 2.5: Add DeleteExpiredVerifications query
  - [ ] 2.6: Run `make sqlc` to generate repository code

- [ ] Task 3: Create verification handler (AC: 1, 2, 4, 5)
  - [ ] 3.1: Create `internal/handlers/verification.go`
  - [ ] 3.2: Create VerificationHandler struct with db, email, config dependencies
  - [ ] 3.3: Implement SendVerificationEmail handler
  - [ ] 3.4: Extract user from auth middleware context
  - [ ] 3.5: Check if user email is already verified (return 400 if yes)
  - [ ] 3.6: Delete any existing verification token for this user
  - [ ] 3.7: Generate new JWT verification token
  - [ ] 3.8: Store token in verification table
  - [ ] 3.9: Send email asynchronously (goroutine)
  - [ ] 3.10: Return success immediately

- [ ] Task 4: Create email template for verification (AC: 1)
  - [ ] 4.1: Create verification email HTML template
  - [ ] 4.2: Include verification link: `{FRONTEND_URL}/verify-email?token={token}`
  - [ ] 4.3: Include plain text alternative

- [ ] Task 5: Wire route with auth middleware (AC: 4)
  - [ ] 5.1: Add POST /api/auth/send-verification-email route
  - [ ] 5.2: Wrap with RequireAuth middleware from Story 1.3
  - [ ] 5.3: Inject VerificationHandler into server

- [ ] Task 6: Write tests (AC: 1, 2, 3, 4, 5)
  - [ ] 6.1: Test token generation creates valid JWT
  - [ ] 6.2: Test send-verification creates token in DB
  - [ ] 6.3: Test old tokens are deleted on resend
  - [ ] 6.4: Test unauthenticated request returns 401
  - [ ] 6.5: Test already-verified user returns 400
  - [ ] 6.6: Test email send failure is logged but doesn't fail request

## Dev Notes

### Dependencies

**Requires from previous stories:**
- Auth middleware (Story 1.3) - for protected endpoint
- Email provider (Story 2.1) - for sending emails

**New dependency to add:**
```bash
go get github.com/golang-jwt/jwt/v5
```

### Verification Token Format

**JWT Claims:**
```go
type VerificationClaims struct {
    jwt.RegisteredClaims
    UserID     string `json:"userId"`
    Identifier string `json:"identifier"` // e.g., "email-verification:{userId}"
}
```

**Token Generation:**
```go
// pkg/token/verification.go

package token

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type VerificationClaims struct {
    jwt.RegisteredClaims
    UserID     string `json:"userId"`
    Identifier string `json:"identifier"`
}

func GenerateVerificationToken(secret string, userID string, identifier string, expiry time.Duration) (string, error) {
    claims := VerificationClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
        UserID:     userID,
        Identifier: identifier,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ValidateVerificationToken(secret string, tokenString string) (*VerificationClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &VerificationClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*VerificationClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token claims")
}
```

### Database Schema

**Existing `verification` table:**
```sql
CREATE TABLE "verification" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier VARCHAR NOT NULL,    -- "email-verification:{userId}"
    value VARCHAR NOT NULL,         -- the JWT token
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### SQL Queries

```sql
-- db/queries/verifications.sql

-- name: CreateVerification :one
INSERT INTO verification (identifier, value, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVerificationByIdentifier :one
SELECT * FROM verification WHERE identifier = $1;

-- name: GetVerificationByIdentifierAndValue :one
SELECT * FROM verification WHERE identifier = $1 AND value = $2;

-- name: DeleteVerificationByIdentifier :exec
DELETE FROM verification WHERE identifier = $1;

-- name: DeleteVerification :exec
DELETE FROM verification WHERE id = $1;

-- name: DeleteExpiredVerifications :exec
DELETE FROM verification WHERE expires_at < NOW();
```

### Handler Implementation

```go
// internal/handlers/verification.go

package handlers

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "time"

    "your-module/internal/config"
    "your-module/internal/email"
    "your-module/internal/middleware"
    "your-module/internal/repository"
    "your-module/pkg/token"
)

type VerificationHandler struct {
    queries *repository.Queries
    email   email.EmailProvider
    config  *config.Config
}

func NewVerificationHandler(q *repository.Queries, e email.EmailProvider, c *config.Config) *VerificationHandler {
    return &VerificationHandler{queries: q, email: e, config: c}
}

func (h *VerificationHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        respondError(w, "METHOD_NOT_ALLOWED", "Method not allowed", 405)
        return
    }

    // Get user from auth middleware context
    user, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        respondError(w, "UNAUTHORIZED", "Authentication required", 401)
        return
    }

    // Check if already verified
    if user.EmailVerified {
        respondError(w, "ALREADY_VERIFIED", "Email is already verified", 400)
        return
    }

    ctx := r.Context()
    identifier := fmt.Sprintf("email-verification:%s", user.ID)

    // Delete any existing verification token
    h.queries.DeleteVerificationByIdentifier(ctx, identifier)

    // Generate new token (1 hour expiry)
    expiry := time.Hour
    tokenString, err := token.GenerateVerificationToken(
        h.config.BetterAuthSecret,
        user.ID.String(),
        identifier,
        expiry,
    )
    if err != nil {
        slog.Error("failed to generate verification token", "error", err, "userId", user.ID)
        respondError(w, "INTERNAL_ERROR", "Failed to generate verification token", 500)
        return
    }

    // Store in database
    _, err = h.queries.CreateVerification(ctx, repository.CreateVerificationParams{
        Identifier: identifier,
        Value:      tokenString,
        ExpiresAt:  time.Now().Add(expiry),
    })
    if err != nil {
        slog.Error("failed to store verification token", "error", err, "userId", user.ID)
        respondError(w, "INTERNAL_ERROR", "Failed to create verification", 500)
        return
    }

    // Send email asynchronously (fire-and-forget)
    go func() {
        verifyURL := fmt.Sprintf("%s/verify-email?token=%s", h.config.FrontendURL, tokenString)

        emailMsg := email.Email{
            To:       user.Email,
            From:     h.config.EmailFrom,
            Subject:  "Verify your email address",
            HTMLBody: fmt.Sprintf(`<p>Click the link below to verify your email:</p><p><a href="%s">Verify Email</a></p><p>This link expires in 1 hour.</p>`, verifyURL),
            TextBody: fmt.Sprintf("Verify your email by visiting: %s\n\nThis link expires in 1 hour.", verifyURL),
        }

        if err := h.email.Send(context.Background(), emailMsg); err != nil {
            slog.Error("failed to send verification email", "error", err, "to", user.Email, "userId", user.ID)
        } else {
            slog.Info("verification email sent", "to", user.Email, "userId", user.ID)
        }
    }()

    // Return immediately - don't wait for email
    respondJSON(w, map[string]bool{"success": true}, 200)
}
```

### Endpoint & API Contract

**Endpoint:** `POST /api/auth/send-verification-email`

**Request:** No body required (user from session)

**Success Response (200):**
```json
{
  "success": true
}
```

**Error Responses:**
| Code | Status | When |
|------|--------|------|
| `UNAUTHORIZED` | 401 | No valid session |
| `ALREADY_VERIFIED` | 400 | Email already verified |
| `INTERNAL_ERROR` | 500 | Token generation or DB failure |

### Configuration Required

```bash
# Required env vars (from Story 2.1)
BETTER_AUTH_SECRET=your-secret-key  # for JWT signing
FRONTEND_URL=http://localhost:3000   # for verification link
EMAIL_FROM=noreply@example.com
```

### Wiring in Server

```go
// internal/server/server.go

func (s *Server) setupRoutes() {
    // ... existing routes ...

    // Protected routes (require auth)
    verificationHandler := handlers.NewVerificationHandler(s.queries, s.email, s.config)

    // Wrap with auth middleware
    s.mux.Handle("POST /api/auth/send-verification-email",
        middleware.RequireAuth(s.queries, http.HandlerFunc(verificationHandler.SendVerificationEmail)))
}
```

### Testing

```go
func TestVerificationHandler_SendVerificationEmail(t *testing.T) {
    // Test: authenticated user with unverified email
    // - Should create verification in DB
    // - Should return success immediately
    // - Should trigger async email send

    // Test: already verified user
    // - Should return 400 ALREADY_VERIFIED

    // Test: unauthenticated request
    // - Should return 401 UNAUTHORIZED

    // Test: resend deletes old token
    // - Create user, send verification
    // - Send again
    // - Should only have 1 token in DB
}
```

### References

- [Source: docs/planning-artifacts/architecture.md#Verification Token Strategy]
- [Source: docs/planning-artifacts/epics.md#Story 2.2]
- [JWT RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519)

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
- api/pkg/token/verification.go
- api/pkg/token/verification_test.go
- api/db/queries/verifications.sql
- api/internal/handlers/verification.go
- api/internal/handlers/verification_test.go

**Expected Generated Files (via sqlc):**
- api/internal/repository/verifications.sql.go

**Expected Modified Files:**
- api/internal/server/server.go (add route + wire handler)
- api/go.mod (add golang-jwt/jwt/v5)
