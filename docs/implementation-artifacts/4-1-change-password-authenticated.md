# Story 4.1: Change Password (Authenticated)

Status: ready-for-dev

## Story

As an **authenticated user**,
I want to **change my password by providing my current password**,
So that **I can improve my account security proactively**.

## Acceptance Criteria

1. **Given** an authenticated user with valid current password **When** they submit current password and new password to change-password endpoint **Then** the current password is validated against stored hash **And** the new password is hashed with bcrypt (cost 10+) **And** the password is updated in `account` table **And** the response returns `{ "success": true }`

2. **Given** an authenticated user submits incorrect current password **When** the request is processed **Then** the response returns 400 with `{ "error": { "code": "INVALID_PASSWORD", "message": "Current password is incorrect" } }`

3. **Given** an authenticated user wants to revoke other sessions **When** they include `revokeOtherSessions: true` in the request **Then** the password is changed **And** all sessions EXCEPT the current one are deleted **And** the response returns `{ "success": true }`

4. **Given** an authenticated user changes password without revoking sessions **When** they omit or set `revokeOtherSessions: false` **Then** other sessions remain active

5. **Given** an unauthenticated user **When** they request change-password endpoint **Then** the response returns 401 `UNAUTHORIZED`

## Tasks / Subtasks

- [ ] Task 1: Implement ChangePassword handler (AC: 1, 2, 3, 4)
  - [ ] 1.1: Add ChangePassword method to PasswordHandler
  - [ ] 1.2: Extract user from auth middleware context
  - [ ] 1.3: Validate request body (currentPassword, newPassword required)
  - [ ] 1.4: Get user's account credentials (GetAccountByUserIdAndProvider)
  - [ ] 1.5: Validate current password with bcrypt.CompareHashAndPassword
  - [ ] 1.6: Hash new password with bcrypt
  - [ ] 1.7: Update password in account table
  - [ ] 1.8: If revokeOtherSessions=true, delete other sessions (keep current)
  - [ ] 1.9: Return success

- [ ] Task 2: Wire route with auth middleware (AC: 5)
  - [ ] 2.1: Add POST /api/auth/change-password route
  - [ ] 2.2: Wrap with RequireAuth middleware

- [ ] Task 3: Write tests (AC: 1, 2, 3, 4, 5)
  - [ ] 3.1: Test valid current password changes password
  - [ ] 3.2: Test wrong current password returns INVALID_PASSWORD
  - [ ] 3.3: Test revokeOtherSessions=true deletes other sessions
  - [ ] 3.4: Test revokeOtherSessions=false keeps other sessions
  - [ ] 3.5: Test unauthenticated request returns 401
  - [ ] 3.6: Test empty passwords return INVALID_BODY

## Dev Notes

### Dependencies

**Requires from previous stories:**
- Auth middleware (Story 1.3)
- PasswordHandler (Story 3.1)
- UpdateAccountPassword query (Story 3.2)
- DeleteOtherSessions query (from architecture spec)

### Endpoint & API Contract

**Endpoint:** `POST /api/auth/change-password`

**Request Body:**
```json
{
  "currentPassword": "oldPassword123",
  "newPassword": "newSecurePassword456",
  "revokeOtherSessions": false
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
| `UNAUTHORIZED` | 401 | No valid session |
| `INVALID_BODY` | 400 | Missing currentPassword or newPassword |
| `INVALID_PASSWORD` | 400 | Current password is incorrect |
| `INTERNAL_ERROR` | 500 | Database or bcrypt failure |

### SQL Query to Add (if not already present)

```sql
-- db/queries/sessions.sql

-- name: DeleteOtherSessions :exec
DELETE FROM session WHERE user_id = $1 AND id != $2;
```

### Handler Implementation

```go
// internal/handlers/password.go (add to existing file)

type ChangePasswordRequest struct {
    CurrentPassword     string `json:"currentPassword"`
    NewPassword         string `json:"newPassword"`
    RevokeOtherSessions bool   `json:"revokeOtherSessions"`
}

func (h *PasswordHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        respondError(w, "METHOD_NOT_ALLOWED", "Method not allowed", 405)
        return
    }

    // Get user and session from auth middleware context
    user, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        respondError(w, "UNAUTHORIZED", "Authentication required", 401)
        return
    }

    session, ok := middleware.GetSessionFromContext(r.Context())
    if !ok {
        respondError(w, "UNAUTHORIZED", "Authentication required", 401)
        return
    }

    var req ChangePasswordRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, "INVALID_BODY", "Invalid request body", 400)
        return
    }

    if req.CurrentPassword == "" {
        respondError(w, "INVALID_BODY", "Current password is required", 400)
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

    // Get current credentials
    account, err := h.queries.GetAccountByUserIdAndProvider(ctx, repository.GetAccountByUserIdAndProviderParams{
        UserID:     user.ID,
        ProviderID: "credential",
    })
    if err != nil {
        slog.Error("failed to get account", "error", err, "userId", user.ID)
        respondError(w, "INTERNAL_ERROR", "Failed to process request", 500)
        return
    }

    // Validate current password
    if account.Password == nil {
        respondError(w, "INVALID_PASSWORD", "Current password is incorrect", 400)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(req.CurrentPassword)); err != nil {
        respondError(w, "INVALID_PASSWORD", "Current password is incorrect", 400)
        return
    }

    // Hash new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        slog.Error("failed to hash password", "error", err)
        respondError(w, "INTERNAL_ERROR", "Failed to process request", 500)
        return
    }

    // Update password
    err = h.queries.UpdateAccountPassword(ctx, repository.UpdateAccountPasswordParams{
        UserID:   user.ID,
        Password: string(hashedPassword),
    })
    if err != nil {
        slog.Error("failed to update password", "error", err, "userId", user.ID)
        respondError(w, "INTERNAL_ERROR", "Failed to change password", 500)
        return
    }

    // Optionally revoke other sessions
    if req.RevokeOtherSessions {
        err = h.queries.DeleteOtherSessions(ctx, repository.DeleteOtherSessionsParams{
            UserID: user.ID,
            ID:     session.ID, // Keep current session
        })
        if err != nil {
            slog.Warn("failed to revoke other sessions", "error", err, "userId", user.ID)
            // Continue - password is already changed
        } else {
            slog.Info("revoked other sessions", "userId", user.ID)
        }
    }

    slog.Info("password changed successfully", "userId", user.ID, "revokedOthers", req.RevokeOtherSessions)
    respondJSON(w, map[string]bool{"success": true}, 200)
}
```

### Key Differences from Reset Password

| Aspect | Change Password | Reset Password |
|--------|-----------------|----------------|
| **Auth** | Requires session (auth middleware) | Token-based (no session needed) |
| **Verification** | Current password must match | Reset token must be valid |
| **Session Handling** | Optional revoke (user's choice) | Always revokes all sessions |
| **Use Case** | Proactive security improvement | Account recovery |

### Wiring in Server

```go
// internal/server/server.go

func (s *Server) setupRoutes() {
    // ... existing routes ...

    passwordHandler := handlers.NewPasswordHandler(s.queries, s.email, s.config)

    // Public routes
    s.mux.HandleFunc("POST /api/auth/forget-password", passwordHandler.ForgetPassword)
    s.mux.HandleFunc("POST /api/auth/reset-password", passwordHandler.ResetPassword)

    // Protected route
    s.mux.Handle("POST /api/auth/change-password",
        middleware.RequireAuth(s.queries, http.HandlerFunc(passwordHandler.ChangePassword)))
}
```

### Testing

```go
func TestPasswordHandler_ChangePassword(t *testing.T) {
    t.Run("valid password change", func(t *testing.T) {
        // Setup: create user with known password, sign in
        // Call: POST /api/auth/change-password with correct current password
        // Assert: can sign in with new password
        // Assert: response = { "success": true }
    })

    t.Run("wrong current password", func(t *testing.T) {
        // Setup: authenticated user
        // Call with wrong currentPassword
        // Assert: 400 INVALID_PASSWORD
    })

    t.Run("revoke other sessions", func(t *testing.T) {
        // Setup: user with 3 sessions
        // Call with revokeOtherSessions: true
        // Assert: only current session remains
    })

    t.Run("keep other sessions", func(t *testing.T) {
        // Setup: user with 3 sessions
        // Call with revokeOtherSessions: false
        // Assert: all 3 sessions remain
    })

    t.Run("unauthenticated request", func(t *testing.T) {
        // Call without session cookie
        // Assert: 401 UNAUTHORIZED
    })
}
```

### Frontend Integration Note

The frontend should have a "Change Password" form in account settings:
1. Fields: Current Password, New Password, Confirm New Password
2. Checkbox: "Sign out of all other devices"
3. On submit: call `POST /api/auth/change-password`
4. On success: show success message
5. On INVALID_PASSWORD: show "Current password is incorrect"

### References

- [Source: docs/planning-artifacts/epics.md#Story 4.1]
- [Source: docs/planning-artifacts/prd.md#FR20-FR21]

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
- api/internal/handlers/password.go (add ChangePassword method)
- api/internal/handlers/password_test.go (add tests)
- api/internal/server/server.go (add route)
- api/db/queries/sessions.sql (add DeleteOtherSessions if not present)

**Expected Generated Files (via sqlc):**
- api/internal/repository/sessions.sql.go (regenerated if query added)
