# Story 2.3: Verify Email Address

Status: ready-for-dev

## Story

As a **user with a verification link**,
I want to **click the link and verify my email**,
So that **my account is marked as verified and I can access all features**.

## Acceptance Criteria

1. **Given** a user clicks a valid verification link with token **When** the verify-email endpoint is called **Then** the JWT token signature is validated **And** the token is looked up in `verification` table **And** the token expiry is checked (1 hour default, NFR8) **And** the user's `emailVerified` flag is set to `true` **And** the verification record is deleted (single-use) **And** the response returns `{ "success": true }` or redirects to success page

2. **Given** a user clicks an expired verification link **When** the verify-email endpoint is called **Then** the response returns 400 with `{ "error": { "code": "INVALID_TOKEN", "message": "Token has expired" } }`

3. **Given** a user clicks a verification link with invalid/tampered token **When** the verify-email endpoint is called **Then** the response returns 400 with `{ "error": { "code": "INVALID_TOKEN", "message": "Invalid token" } }`

4. **Given** a user clicks a verification link that was already used **When** the verify-email endpoint is called **Then** the response returns 400 with `INVALID_TOKEN` (token not found in DB)

## Tasks / Subtasks

- [ ] Task 1: Add UpdateUserEmailVerified query (AC: 1)
  - [ ] 1.1: Add query to db/queries/users.sql
  - [ ] 1.2: Run `make sqlc` to regenerate

- [ ] Task 2: Implement VerifyEmail handler (AC: 1, 2, 3, 4)
  - [ ] 2.1: Add VerifyEmail method to VerificationHandler
  - [ ] 2.2: Extract token from query parameter
  - [ ] 2.3: Validate JWT signature using ValidateVerificationToken
  - [ ] 2.4: Extract userId and identifier from token claims
  - [ ] 2.5: Look up token in verification table by identifier AND value
  - [ ] 2.6: Check if token has expired (expiresAt < now)
  - [ ] 2.7: Update user's emailVerified to true
  - [ ] 2.8: Delete the verification record (single-use)
  - [ ] 2.9: Return success response

- [ ] Task 3: Wire route (AC: 1)
  - [ ] 3.1: Add GET /api/auth/verify-email route
  - [ ] 3.2: No auth middleware needed (public endpoint with token)

- [ ] Task 4: Write tests (AC: 1, 2, 3, 4)
  - [ ] 4.1: Test valid token verifies user and deletes token
  - [ ] 4.2: Test expired token returns INVALID_TOKEN
  - [ ] 4.3: Test tampered/invalid JWT returns INVALID_TOKEN
  - [ ] 4.4: Test already-used token returns INVALID_TOKEN
  - [ ] 4.5: Test missing token parameter returns error

## Dev Notes

### Dependencies

**Requires from previous stories:**
- Verification token package (Story 2.2)
- Verification SQL queries (Story 2.2)
- VerificationHandler (Story 2.2)

### Endpoint & API Contract

**Endpoint:** `GET /api/auth/verify-email?token={jwt_token}`

**Success Response (200):**
```json
{
  "success": true
}
```

**Error Responses:**
| Code | Status | When |
|------|--------|------|
| `INVALID_TOKEN` | 400 | Token expired, invalid signature, tampered, or already used |
| `INVALID_BODY` | 400 | Missing token parameter |
| `INTERNAL_ERROR` | 500 | Database failure |

### SQL Query to Add

```sql
-- db/queries/users.sql

-- name: UpdateUserEmailVerified :exec
UPDATE "user"
SET "emailVerified" = $2, "updatedAt" = NOW()
WHERE id = $1;
```

### Handler Implementation

```go
// internal/handlers/verification.go (add to existing file)

func (h *VerificationHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        respondError(w, "METHOD_NOT_ALLOWED", "Method not allowed", 405)
        return
    }

    // Extract token from query parameter
    tokenString := r.URL.Query().Get("token")
    if tokenString == "" {
        respondError(w, "INVALID_BODY", "Token parameter is required", 400)
        return
    }

    ctx := r.Context()

    // Validate JWT signature and extract claims
    claims, err := token.ValidateVerificationToken(h.config.BetterAuthSecret, tokenString)
    if err != nil {
        slog.Warn("invalid verification token", "error", err)
        respondError(w, "INVALID_TOKEN", "Invalid token", 400)
        return
    }

    // Check if token is expired (JWT library checks this, but double-check)
    if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
        respondError(w, "INVALID_TOKEN", "Token has expired", 400)
        return
    }

    // Look up token in database (ensures single-use)
    verification, err := h.queries.GetVerificationByIdentifierAndValue(ctx, repository.GetVerificationByIdentifierAndValueParams{
        Identifier: claims.Identifier,
        Value:      tokenString,
    })
    if err != nil {
        // Token not found = already used or never existed
        slog.Warn("verification token not found in DB", "identifier", claims.Identifier)
        respondError(w, "INVALID_TOKEN", "Invalid or already used token", 400)
        return
    }

    // Double-check DB expiry
    if verification.ExpiresAt.Before(time.Now()) {
        // Clean up expired token
        h.queries.DeleteVerification(ctx, verification.ID)
        respondError(w, "INVALID_TOKEN", "Token has expired", 400)
        return
    }

    // Parse userId from claims
    userID, err := uuid.Parse(claims.UserID)
    if err != nil {
        slog.Error("invalid userId in token", "userId", claims.UserID)
        respondError(w, "INVALID_TOKEN", "Invalid token", 400)
        return
    }

    // Update user's emailVerified status
    err = h.queries.UpdateUserEmailVerified(ctx, repository.UpdateUserEmailVerifiedParams{
        ID:            userID,
        EmailVerified: true,
    })
    if err != nil {
        slog.Error("failed to update email verified status", "error", err, "userId", userID)
        respondError(w, "INTERNAL_ERROR", "Failed to verify email", 500)
        return
    }

    // Delete the verification record (single-use)
    err = h.queries.DeleteVerification(ctx, verification.ID)
    if err != nil {
        // Log but don't fail - user is already verified
        slog.Warn("failed to delete verification record", "error", err, "id", verification.ID)
    }

    slog.Info("email verified successfully", "userId", userID)
    respondJSON(w, map[string]bool{"success": true}, 200)
}
```

### Token Validation Flow

```
1. Extract token from ?token= query param
2. Validate JWT signature (using BETTER_AUTH_SECRET)
3. Extract claims: userId, identifier, exp
4. Look up in DB: verification WHERE identifier = ? AND value = ?
5. If not found → already used or invalid → return INVALID_TOKEN
6. Check expiry (both JWT and DB)
7. Update user.emailVerified = true
8. Delete verification record
9. Return success
```

### Security Considerations

**Why validate both JWT and DB:**
- JWT signature proves token wasn't tampered with
- DB lookup ensures single-use (token is deleted after use)
- DB expiry is authoritative (in case of clock skew)

**Why no auth required:**
- User may not be logged in when clicking email link
- The token itself proves ownership (sent to their email)
- Token is bound to specific userId via claims

### Testing

```go
func TestVerificationHandler_VerifyEmail(t *testing.T) {
    t.Run("valid token verifies user", func(t *testing.T) {
        // Setup: create user, create verification token
        // Call: GET /api/auth/verify-email?token={token}
        // Assert: user.emailVerified = true
        // Assert: verification record deleted
        // Assert: response = { "success": true }
    })

    t.Run("expired token returns error", func(t *testing.T) {
        // Setup: create verification with past expiry
        // Call: GET /api/auth/verify-email?token={token}
        // Assert: 400 INVALID_TOKEN "Token has expired"
    })

    t.Run("tampered token returns error", func(t *testing.T) {
        // Setup: create valid token, modify payload
        // Call: GET /api/auth/verify-email?token={tampered}
        // Assert: 400 INVALID_TOKEN
    })

    t.Run("already used token returns error", func(t *testing.T) {
        // Setup: create and verify token once
        // Call: GET /api/auth/verify-email?token={same token}
        // Assert: 400 INVALID_TOKEN (not found in DB)
    })

    t.Run("missing token returns error", func(t *testing.T) {
        // Call: GET /api/auth/verify-email (no token param)
        // Assert: 400 INVALID_BODY
    })
}
```

### Wiring in Server

```go
// internal/server/server.go

func (s *Server) setupRoutes() {
    // ... existing routes ...

    verificationHandler := handlers.NewVerificationHandler(s.queries, s.email, s.config)

    // Protected route (from Story 2.2)
    s.mux.Handle("POST /api/auth/send-verification-email",
        middleware.RequireAuth(s.queries, http.HandlerFunc(verificationHandler.SendVerificationEmail)))

    // Public route (token provides auth)
    s.mux.HandleFunc("GET /api/auth/verify-email", verificationHandler.VerifyEmail)
}
```

### Frontend Integration Note

The frontend should handle the `/verify-email?token=xxx` route:
1. Extract token from URL
2. Call `GET /api/auth/verify-email?token={token}`
3. On success: show "Email verified!" message, redirect to dashboard
4. On error: show error message with option to resend

### References

- [Source: docs/planning-artifacts/architecture.md#Verification Token Strategy]
- [Source: docs/planning-artifacts/epics.md#Story 2.3]
- [Source: docs/planning-artifacts/prd.md#NFR8 - token expiry]

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

**Expected Modified Files:**
- api/internal/handlers/verification.go (add VerifyEmail method)
- api/internal/handlers/verification_test.go (add tests)
- api/internal/server/server.go (add route)
- api/db/queries/users.sql (add UpdateUserEmailVerified)

**Expected Generated Files (via sqlc):**
- api/internal/repository/users.sql.go (regenerated)
