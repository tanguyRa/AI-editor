# Story 3.1: Request Password Reset

Status: ready-for-dev

## Story

As a **user who forgot their password**,
I want to **request a password reset email**,
So that **I can regain access to my account**.

## Acceptance Criteria

1. **Given** a user submits their email to forget-password endpoint **When** the email exists in the system **Then** a JWT reset token is generated with 1-hour expiry **And** the token is stored in `verification` table with identifier `reset-password:{userId}` **And** a password reset email is sent asynchronously **And** the response returns `{ "success": true }`

2. **Given** a user submits an email that doesn't exist **When** the request is processed **Then** the response STILL returns `{ "success": true }` (NFR9 - no email enumeration) **And** no email is sent **And** response timing is consistent with successful case

3. **Given** a user requests reset when a previous token exists **When** the request is processed **Then** the old token is deleted before creating a new one

4. **Given** the email provider fails **When** the error occurs **Then** the error is logged but response is unaffected

## Tasks / Subtasks

- [ ] Task 1: Create password handler (AC: 1, 2, 3, 4)
  - [ ] 1.1: Create `internal/handlers/password.go`
  - [ ] 1.2: Create PasswordHandler struct with db, email, config dependencies
  - [ ] 1.3: Implement ForgetPassword handler
  - [ ] 1.4: Look up user by email (GetUserByEmail)
  - [ ] 1.5: If user not found, do timing-safe delay then return success
  - [ ] 1.6: Delete any existing reset token for this user
  - [ ] 1.7: Generate JWT reset token with identifier `reset-password:{userId}`
  - [ ] 1.8: Store token in verification table
  - [ ] 1.9: Send reset email asynchronously
  - [ ] 1.10: Return success immediately

- [ ] Task 2: Implement timing-safe response (AC: 2)
  - [ ] 2.1: Measure average time for successful flow
  - [ ] 2.2: Add delay for non-existent email to match timing
  - [ ] 2.3: Use constant-time operations where possible

- [ ] Task 3: Create reset email template (AC: 1)
  - [ ] 3.1: Create HTML template with reset link
  - [ ] 3.2: Link format: `{FRONTEND_URL}/reset-password?token={token}`
  - [ ] 3.3: Include plain text alternative

- [ ] Task 4: Wire route (AC: 1)
  - [ ] 4.1: Add POST /api/auth/forget-password route
  - [ ] 4.2: No auth middleware (public endpoint)

- [ ] Task 5: Write tests (AC: 1, 2, 3, 4)
  - [ ] 5.1: Test existing email creates token and sends email
  - [ ] 5.2: Test non-existent email returns success (no email sent)
  - [ ] 5.3: Test response times are similar for both cases
  - [ ] 5.4: Test old tokens are deleted on new request
  - [ ] 5.5: Test email failure doesn't affect response

## Dev Notes

### Dependencies

**Requires from previous stories:**
- Email provider (Story 2.1)
- Verification token package (Story 2.2)
- Verification SQL queries (Story 2.2)

### Endpoint & API Contract

**Endpoint:** `POST /api/auth/forget-password`

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Success Response (200):** Always returns success regardless of email existence
```json
{
  "success": true
}
```

**Error Responses:**
| Code | Status | When |
|------|--------|------|
| `INVALID_BODY` | 400 | Missing or invalid email |
| `INTERNAL_ERROR` | 500 | Database or token generation failure |

### Security: Timing Attack Prevention (NFR9)

The response must take the same amount of time whether the email exists or not. This prevents attackers from enumerating valid emails.

**Implementation Strategy:**
```go
func (h *PasswordHandler) ForgetPassword(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    // ... process request ...

    user, err := h.queries.GetUserByEmail(ctx, req.Email)
    if err != nil {
        // Email not found - but don't reveal this!
        // Wait to match timing of successful case
        elapsed := time.Since(start)
        minDuration := 200 * time.Millisecond // adjust based on testing
        if elapsed < minDuration {
            time.Sleep(minDuration - elapsed)
        }
        respondJSON(w, map[string]bool{"success": true}, 200)
        return
    }

    // ... create token, send email ...

    respondJSON(w, map[string]bool{"success": true}, 200)
}
```

### Handler Implementation

```go
// internal/handlers/password.go

package handlers

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "net/http"
    "time"

    "your-module/internal/config"
    "your-module/internal/email"
    "your-module/internal/repository"
    "your-module/pkg/token"
)

type PasswordHandler struct {
    queries *repository.Queries
    email   email.EmailProvider
    config  *config.Config
}

func NewPasswordHandler(q *repository.Queries, e email.EmailProvider, c *config.Config) *PasswordHandler {
    return &PasswordHandler{queries: q, email: e, config: c}
}

type ForgetPasswordRequest struct {
    Email string `json:"email"`
}

func (h *PasswordHandler) ForgetPassword(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        respondError(w, "METHOD_NOT_ALLOWED", "Method not allowed", 405)
        return
    }

    start := time.Now()
    minDuration := 200 * time.Millisecond // Ensure consistent timing

    var req ForgetPasswordRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, "INVALID_BODY", "Invalid request body", 400)
        return
    }

    if req.Email == "" {
        respondError(w, "INVALID_BODY", "Email is required", 400)
        return
    }

    ctx := r.Context()

    // Look up user
    user, err := h.queries.GetUserByEmail(ctx, req.Email)
    if err != nil {
        // User not found - wait for timing consistency, then return success
        slog.Debug("password reset requested for non-existent email", "email", req.Email)
        elapsed := time.Since(start)
        if elapsed < minDuration {
            time.Sleep(minDuration - elapsed)
        }
        respondJSON(w, map[string]bool{"success": true}, 200)
        return
    }

    identifier := fmt.Sprintf("reset-password:%s", user.ID)

    // Delete any existing reset token
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
        slog.Error("failed to generate reset token", "error", err, "userId", user.ID)
        respondError(w, "INTERNAL_ERROR", "Failed to process request", 500)
        return
    }

    // Store in database
    _, err = h.queries.CreateVerification(ctx, repository.CreateVerificationParams{
        Identifier: identifier,
        Value:      tokenString,
        ExpiresAt:  time.Now().Add(expiry),
    })
    if err != nil {
        slog.Error("failed to store reset token", "error", err, "userId", user.ID)
        respondError(w, "INTERNAL_ERROR", "Failed to process request", 500)
        return
    }

    // Send email asynchronously
    go func() {
        resetURL := fmt.Sprintf("%s/reset-password?token=%s", h.config.FrontendURL, tokenString)

        emailMsg := email.Email{
            To:       user.Email,
            From:     h.config.EmailFrom,
            Subject:  "Reset your password",
            HTMLBody: fmt.Sprintf(`
                <p>You requested a password reset.</p>
                <p><a href="%s">Click here to reset your password</a></p>
                <p>This link expires in 1 hour.</p>
                <p>If you didn't request this, you can safely ignore this email.</p>
            `, resetURL),
            TextBody: fmt.Sprintf(
                "Reset your password by visiting: %s\n\nThis link expires in 1 hour.\n\nIf you didn't request this, you can safely ignore this email.",
                resetURL,
            ),
        }

        if err := h.email.Send(context.Background(), emailMsg); err != nil {
            slog.Error("failed to send reset email", "error", err, "to", user.Email, "userId", user.ID)
        } else {
            slog.Info("password reset email sent", "to", user.Email, "userId", user.ID)
        }
    }()

    // Ensure minimum duration for timing consistency
    elapsed := time.Since(start)
    if elapsed < minDuration {
        time.Sleep(minDuration - elapsed)
    }

    respondJSON(w, map[string]bool{"success": true}, 200)
}
```

### Token Identifier Format

Uses same verification table as email verification, but with different identifier prefix:
- Email verification: `email-verification:{userId}`
- Password reset: `reset-password:{userId}`

### Wiring in Server

```go
// internal/server/server.go

func (s *Server) setupRoutes() {
    // ... existing routes ...

    passwordHandler := handlers.NewPasswordHandler(s.queries, s.email, s.config)

    // Public route (no auth required)
    s.mux.HandleFunc("POST /api/auth/forget-password", passwordHandler.ForgetPassword)
}
```

### Testing Timing Attack Prevention

```go
func TestForgetPassword_TimingConsistency(t *testing.T) {
    // Test with existing email
    var existingTimes []time.Duration
    for i := 0; i < 10; i++ {
        start := time.Now()
        // POST /api/auth/forget-password with existing email
        existingTimes = append(existingTimes, time.Since(start))
    }
    avgExisting := average(existingTimes)

    // Test with non-existent email
    var nonExistentTimes []time.Duration
    for i := 0; i < 10; i++ {
        start := time.Now()
        // POST /api/auth/forget-password with non-existent email
        nonExistentTimes = append(nonExistentTimes, time.Since(start))
    }
    avgNonExistent := average(nonExistentTimes)

    // Times should be within 50ms of each other
    diff := avgExisting - avgNonExistent
    if diff < 0 {
        diff = -diff
    }
    assert.Less(t, diff, 50*time.Millisecond, "Response times should be similar")
}
```

### References

- [Source: docs/planning-artifacts/architecture.md#Verification Token Strategy]
- [Source: docs/planning-artifacts/epics.md#Story 3.1]
- [Source: docs/planning-artifacts/prd.md#NFR9 - email enumeration prevention]

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
- api/internal/handlers/password.go
- api/internal/handlers/password_test.go

**Expected Modified Files:**
- api/internal/server/server.go (add route + wire handler)
