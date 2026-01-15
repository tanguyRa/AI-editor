# Story 3.2: Reset Password with Token

Status: ready-for-dev

## Story

As a **user with a reset link**,
I want to **set a new password using the link**,
So that **I can access my account with new credentials**.

## Acceptance Criteria

1. **Given** a user submits a valid reset token and new password **When** the reset-password endpoint is called **Then** the JWT token signature is validated **And** the token is looked up in `verification` table **And** the token expiry is checked (1 hour, NFR8) **And** the new password is hashed with bcrypt (cost 10+) **And** the user's password is updated in `account` table **And** the reset token is deleted immediately (NFR10 - single use) **And** ALL user sessions are invalidated (security best practice) **And** the response returns `{ "success": true }`

2. **Given** a user submits an expired reset token **When** the reset-password endpoint is called **Then** the response returns 400 with `INVALID_TOKEN`

3. **Given** a user submits an invalid/tampered token **When** the reset-password endpoint is called **Then** the response returns 400 with `INVALID_TOKEN`

4. **Given** a user submits a token that was already used **When** the reset-password endpoint is called **Then** the response returns 400 with `INVALID_TOKEN`

5. **Given** a user submits a weak or empty password **When** the reset-password endpoint is called **Then** the response returns 400 with `INVALID_BODY`

## Tasks / Subtasks

- [ ] Task 1: Add UpdateAccountPassword query (AC: 1)
  - [ ] 1.1: Add query to db/queries/accounts.sql
  - [ ] 1.2: Run `make sqlc` to regenerate

- [ ] Task 2: Add DeleteUserSessions query (AC: 1)
  - [ ] 2.1: Add query to db/queries/sessions.sql to delete all sessions for a user
  - [ ] 2.2: Run `make sqlc` to regenerate

- [ ] Task 3: Implement ResetPassword handler (AC: 1, 2, 3, 4, 5)
  - [ ] 3.1: Add ResetPassword method to PasswordHandler
  - [ ] 3.2: Validate request body (token, newPassword required)
  - [ ] 3.3: Validate JWT signature
  - [ ] 3.4: Extract userId and identifier from claims
  - [ ] 3.5: Look up token in verification table
  - [ ] 3.6: Check token expiry
  - [ ] 3.7: Hash new password with bcrypt (cost 10+)
  - [ ] 3.8: Update password in account table
  - [ ] 3.9: Delete reset token (single-use)
  - [ ] 3.10: Delete ALL user sessions (security)
  - [ ] 3.11: Return success

- [ ] Task 4: Wire route (AC: 1)
  - [ ] 4.1: Add POST /api/auth/reset-password route
  - [ ] 4.2: No auth middleware (token provides auth)

- [ ] Task 5: Write tests (AC: 1, 2, 3, 4, 5)
  - [ ] 5.1: Test valid token resets password
  - [ ] 5.2: Test all sessions are invalidated after reset
  - [ ] 5.3: Test token is deleted after use
  - [ ] 5.4: Test expired token returns INVALID_TOKEN
  - [ ] 5.5: Test invalid/tampered token returns INVALID_TOKEN
  - [ ] 5.6: Test already-used token returns INVALID_TOKEN
  - [ ] 5.7: Test empty password returns INVALID_BODY

## Dev Notes

### Dependencies

**Requires from previous stories:**
- Verification token package (Story 2.2)
- Verification SQL queries (Story 2.2)
- PasswordHandler (Story 3.1)

### Endpoint & API Contract

**Endpoint:** `POST /api/auth/reset-password`

**Request Body:**
```json
{
  "token": "jwt-reset-token-here",
  "newPassword": "newSecurePassword123"
}
```

**Success Response (200):**
```json
{
  "success": true
}
```

**Error Responses:**
| Code | Status | When |
|------|--------|------|
| `INVALID_BODY` | 400 | Missing token or newPassword |
| `INVALID_TOKEN` | 400 | Expired, invalid, tampered, or already used token |
| `INTERNAL_ERROR` | 500 | Database or bcrypt failure |

### SQL Queries to Add

```sql
-- db/queries/accounts.sql

-- name: UpdateAccountPassword :exec
UPDATE account
SET password = $2, "updatedAt" = NOW()
WHERE "userId" = $1 AND "providerId" = 'credential';
```

```sql
-- db/queries/sessions.sql

-- name: DeleteUserSessions :exec
DELETE FROM session WHERE user_id = $1;
```

### Handler Implementation

```go
// internal/handlers/password.go (add to existing file)

type ResetPasswordRequest struct {
    Token       string `json:"token"`
    NewPassword string `json:"newPassword"`
}

func (h *PasswordHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        respondError(w, "METHOD_NOT_ALLOWED", "Method not allowed", 405)
        return
    }

    var req ResetPasswordRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, "INVALID_BODY", "Invalid request body", 400)
        return
    }

    if req.Token == "" {
        respondError(w, "INVALID_BODY", "Token is required", 400)
        return
    }

    if req.NewPassword == "" {
        respondError(w, "INVALID_BODY", "New password is required", 400)
        return
    }

    // Optional: Add password strength validation
    if len(req.NewPassword) < 8 {
        respondError(w, "INVALID_BODY", "Password must be at least 8 characters", 400)
        return
    }

    ctx := r.Context()

    // Validate JWT signature and extract claims
    claims, err := token.ValidateVerificationToken(h.config.BetterAuthSecret, req.Token)
    if err != nil {
        slog.Warn("invalid reset token", "error", err)
        respondError(w, "INVALID_TOKEN", "Invalid token", 400)
        return
    }

    // Check JWT expiry
    if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
        respondError(w, "INVALID_TOKEN", "Token has expired", 400)
        return
    }

    // Look up token in database (ensures single-use)
    verification, err := h.queries.GetVerificationByIdentifierAndValue(ctx, repository.GetVerificationByIdentifierAndValueParams{
        Identifier: claims.Identifier,
        Value:      req.Token,
    })
    if err != nil {
        slog.Warn("reset token not found in DB", "identifier", claims.Identifier)
        respondError(w, "INVALID_TOKEN", "Invalid or already used token", 400)
        return
    }

    // Double-check DB expiry
    if verification.ExpiresAt.Before(time.Now()) {
        h.queries.DeleteVerification(ctx, verification.ID)
        respondError(w, "INVALID_TOKEN", "Token has expired", 400)
        return
    }

    // Parse userId
    userID, err := uuid.Parse(claims.UserID)
    if err != nil {
        slog.Error("invalid userId in token", "userId", claims.UserID)
        respondError(w, "INVALID_TOKEN", "Invalid token", 400)
        return
    }

    // Hash new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        slog.Error("failed to hash password", "error", err)
        respondError(w, "INTERNAL_ERROR", "Failed to process request", 500)
        return
    }

    // Update password in account table
    err = h.queries.UpdateAccountPassword(ctx, repository.UpdateAccountPasswordParams{
        UserID:   userID,
        Password: string(hashedPassword),
    })
    if err != nil {
        slog.Error("failed to update password", "error", err, "userId", userID)
        respondError(w, "INTERNAL_ERROR", "Failed to reset password", 500)
        return
    }

    // Delete reset token (single-use) - NFR10
    err = h.queries.DeleteVerification(ctx, verification.ID)
    if err != nil {
        slog.Warn("failed to delete reset token", "error", err, "id", verification.ID)
        // Continue - password is already reset
    }

    // Invalidate ALL sessions for security
    err = h.queries.DeleteUserSessions(ctx, userID)
    if err != nil {
        slog.Warn("failed to delete user sessions", "error", err, "userId", userID)
        // Continue - password is reset, sessions will expire naturally
    }

    slog.Info("password reset successful", "userId", userID)
    respondJSON(w, map[string]bool{"success": true}, 200)
}
```

### Security: Session Invalidation

When a password is reset, ALL existing sessions must be deleted. This ensures:
1. If the account was compromised, attacker sessions are terminated
2. User must sign in with new credentials
3. Follows security best practices

### Token Validation Flow

```
1. Extract token and newPassword from body
2. Validate JWT signature
3. Extract claims: userId, identifier (reset-password:{userId})
4. Look up in DB: verification WHERE identifier = ? AND value = ?
5. If not found → already used → INVALID_TOKEN
6. Check expiry
7. Hash new password with bcrypt
8. Update account.password
9. Delete verification record (single-use)
10. Delete ALL user sessions
11. Return success
```

### Wiring in Server

```go
// internal/server/server.go

func (s *Server) setupRoutes() {
    // ... existing routes ...

    passwordHandler := handlers.NewPasswordHandler(s.queries, s.email, s.config)

    s.mux.HandleFunc("POST /api/auth/forget-password", passwordHandler.ForgetPassword)
    s.mux.HandleFunc("POST /api/auth/reset-password", passwordHandler.ResetPassword)
}
```

### Testing

```go
func TestPasswordHandler_ResetPassword(t *testing.T) {
    t.Run("valid token resets password", func(t *testing.T) {
        // Setup: create user, create reset token
        // Call: POST /api/auth/reset-password
        // Assert: password is changed (can sign in with new password)
        // Assert: token is deleted from DB
        // Assert: all sessions are deleted
    })

    t.Run("expired token fails", func(t *testing.T) {
        // Setup: create token with past expiry
        // Assert: 400 INVALID_TOKEN
    })

    t.Run("tampered token fails", func(t *testing.T) {
        // Setup: create valid token, modify it
        // Assert: 400 INVALID_TOKEN
    })

    t.Run("already used token fails", func(t *testing.T) {
        // Setup: use token once
        // Call again with same token
        // Assert: 400 INVALID_TOKEN
    })

    t.Run("empty password fails", func(t *testing.T) {
        // Assert: 400 INVALID_BODY
    })

    t.Run("sessions are invalidated", func(t *testing.T) {
        // Setup: create user with multiple sessions
        // Reset password
        // Assert: all sessions deleted from DB
    })
}
```

### Frontend Integration Note

The frontend should handle the `/reset-password?token=xxx` route:
1. Extract token from URL
2. Show "Enter new password" form
3. On submit: call `POST /api/auth/reset-password` with token + newPassword
4. On success: show "Password reset!" message, redirect to sign-in
5. On error: show error message with option to request new reset

### References

- [Source: docs/planning-artifacts/architecture.md#Verification Token Strategy]
- [Source: docs/planning-artifacts/epics.md#Story 3.2]
- [Source: docs/planning-artifacts/prd.md#NFR10 - single-use tokens]

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
- api/internal/handlers/password.go (add ResetPassword method)
- api/internal/handlers/password_test.go (add tests)
- api/internal/server/server.go (add route)
- api/db/queries/accounts.sql (add UpdateAccountPassword)
- api/db/queries/sessions.sql (add DeleteUserSessions)

**Expected Generated Files (via sqlc):**
- api/internal/repository/accounts.sql.go (regenerated)
- api/internal/repository/sessions.sql.go (regenerated)
