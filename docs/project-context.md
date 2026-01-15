# Project Context - AI-editor

> Critical rules for AI agents implementing code in this project.

## Technology Stack

| Component | Technology | Version |
|-----------|------------|---------|
| Language | Go | 1.24 |
| HTTP | net/http ServeMux | stdlib |
| Database | PostgreSQL via pgx | 5.8 |
| Code Gen | sqlc | latest |
| JWT | golang-jwt/jwt/v5 | latest |
| Password | bcrypt | golang.org/x/crypto |

## Critical Rules

### better-auth Compatibility (NON-NEGOTIABLE)

- API response shapes MUST match better-auth exactly
- Cookie name: `better-auth.session_token`
- Error format: `{ "error": { "code": "ERROR_CODE", "message": "Human readable" } }`
- JSON fields: camelCase (`userId`, `emailVerified`, `expiresAt`)
- Database columns: snake_case (`user_id`, `email_verified`, `expires_at`)

### Go Patterns

**Handlers:**
```go
// Use struct receivers with dependency injection
type AuthHandler struct {
    db    *database.Adapter
    email email.Provider
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) { ... }
```

**Response Helpers (use these, don't create new ones):**
```go
respondJSON(w, data, 200)
respondError(w, "ERROR_CODE", "Message", 400)
```

**Logging:**
```go
slog.Info("action", "key", value)
slog.Error("failed", "error", err)
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Exported | PascalCase | `HandleSignIn`, `EmailProvider` |
| Unexported | camelCase | `validateToken`, `createSession` |
| Files | snake_case | `auth_handler.go`, `session_cache.go` |
| Tests | `*_test.go` | `auth_handler_test.go` |

### Security Rules

- bcrypt cost: minimum 10 (use `bcrypt.DefaultCost`)
- Cookies: `HttpOnly: true`, `SameSite: Lax`, `Secure: true` in production
- Password errors: always return same message (`"Invalid email or password"`)
- Token comparison: use constant-time comparison for verification tokens

### Testing

- Co-located tests: `foo.go` → `foo_test.go`
- Table-driven style required
- Test both success and error paths

### Anti-Patterns (DO NOT)

- Do NOT skip the database adapter layer
- Do NOT add external routing libraries (use stdlib ServeMux)
- Do NOT create new response helper functions
- Do NOT change error code strings (they're part of the API contract)
- Do NOT use `log` package - use `slog` only
- Do NOT store passwords without bcrypt hashing
