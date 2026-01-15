---
stepsCompleted: ['step-01-validate-prerequisites', 'step-02-design-epics', 'step-03-create-stories', 'step-04-final-validation']
status: complete
completedAt: '2026-01-14'
inputDocuments:
  - docs/planning-artifacts/prd.md
  - docs/planning-artifacts/architecture.md
  - docs/index.md
  - docs/project-overview.md
  - docs/integration-architecture.md
  - docs/api-contracts.md
  - docs/data-models.md
---

# AI-editor - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for AI-editor, decomposing the requirements from the PRD and Architecture into implementable stories.

**Project Context:** Brownfield - extending existing Go auth system to achieve better-auth compatibility

## Requirements Inventory

### Functional Requirements

**Account Management (Rewrite - Currently In-Memory Hack)**
- FR1: Users can create an account with email, password, and display name
- FR2: System generates unique user ID upon account creation
- FR3: System stores user profile information (name, email, image URL)
- FR4: System tracks account creation and update timestamps

**Authentication (Rewrite - Currently In-Memory Hack)**
- FR5: Users can sign in with email and password
- FR6: System validates credentials against stored password hash
- FR7: System creates a new session upon successful authentication
- FR8: Users can sign out and terminate their current session
- FR9: System can validate whether a request has a valid session
- FR10: System captures IP address and user agent for each session

**Email Verification (New)**
- FR11: System can send verification emails to users
- FR12: Users can verify their email by clicking a verification link
- FR13: System validates verification tokens and marks email as verified
- FR14: System tracks email verification status on user accounts
- FR15: Verification tokens expire after a configurable time period

**Password Management (New)**
- FR16: Users can request a password reset via email
- FR17: System sends password reset emails with secure tokens
- FR18: Users can reset their password using a valid reset token
- FR19: System invalidates reset tokens after use or expiration
- FR20: Authenticated users can change their password by providing current password
- FR21: Password changes can optionally revoke all other sessions

**Session Management (New)**
- FR22: Users can view a list of all their active sessions
- FR23: Session list displays device info (IP address, user agent)
- FR24: Users can revoke a specific session by its identifier
- FR25: Users can revoke all sessions except the current one
- FR26: System automatically expires sessions after configurable duration
- FR27: System supports JWT/JWE encoded session tokens for caching

**Email Delivery (New)**
- FR28: System abstracts email sending behind a provider interface
- FR29: System supports SendGrid as an email provider
- FR30: System supports Resend as an email provider
- FR31: Email provider is selected via configuration at startup
- FR32: System validates email configuration on startup

**Configuration & Security (Partially Exists)**
- FR33: All tunables (token expiry, session duration) are configurable via environment
- FR34: System uses bcrypt for password hashing
- FR35: Session tokens are stored in HttpOnly cookies
- FR36: System supports CORS with configurable allowed origins

### NonFunctional Requirements

**Performance**
- NFR1: Session validation using cached JWT tokens completes in < 5ms
- NFR2: Auth endpoints (sign-in, sign-up) respond within 500ms under normal load
- NFR3: Database queries use indexes for user/session lookups (no table scans)
- NFR4: Email sending is asynchronous — does not block API response

**Security**
- NFR5: Passwords are hashed using bcrypt with minimum cost factor of 10
- NFR6: Session tokens are cryptographically signed (JWT) or encrypted (JWE)
- NFR7: All session cookies are HttpOnly, SameSite=Lax, Secure in production
- NFR8: Verification and reset tokens expire after 1 hour (configurable)
- NFR9: Failed authentication attempts do not reveal whether email exists
- NFR10: Password reset invalidates the token immediately after use
- NFR11: CORS allows only configured origins (no wildcard in production)

**Integration (better-auth Compatibility)**
- NFR12: All auth endpoints match better-auth API specification exactly
- NFR13: Response shapes are compatible with better-auth TypeScript client
- NFR14: Database schema matches better-auth PostgreSQL adapter
- NFR15: Email provider interface allows adding new providers without core changes

**Reliability**
- NFR16: Invalid or expired sessions return consistent error responses (not crashes)
- NFR17: Email provider failures are logged but do not fail the parent operation
- NFR18: Database connection failures are handled gracefully with appropriate errors
- NFR19: Configuration validation fails fast at startup (not runtime)

### Additional Requirements

**From Architecture - In-Memory to PostgreSQL Migration (Critical):**
- Current sign-up/sign-in/sign-out/get-session are placeholder hacks using in-memory storage
- Must rewrite to use repository layer with sqlc-generated queries
- Must use existing PostgreSQL schema (user, session, account tables)
- This is foundation work before any new features can be built

**From Architecture - Server Refactor (Foundation Work):**
- Extract handlers from monolithic `internal/server/main.go` to `internal/handlers/`
- Extract CORS to `internal/middleware/cors.go`
- Create response helpers (`respondJSON`, `respondError`) in `internal/handlers/response.go`
- Wire new structure in `internal/server/server.go`

**From Architecture - New SQL Queries:**
- Add verification queries: CreateVerification, GetVerification, DeleteVerification, DeleteExpiredVerifications
- Extend session queries: ListUserSessions, DeleteSessionByID, DeleteOtherSessions

**From Architecture - New Packages:**
- `pkg/session/cache.go` - Cookie cache sign/verify (HMAC-SHA256)
- `pkg/token/verification.go` - JWT verification tokens
- `internal/email/provider.go` - EmailProvider interface
- `internal/email/sendgrid.go` - SendGrid implementation
- `internal/email/resend.go` - Resend implementation

**From Architecture - Locked Patterns (Non-Negotiable):**
- Cookie name: `better-auth.session_token`
- Error format: `{ "error": { "code": "...", "message": "..." } }`
- JSON fields: camelCase
- Database columns: snake_case
- All auth endpoints match better-auth paths exactly

**From Architecture - New Dependencies:**
- `golang-jwt/jwt/v5` for JWT/JWE tokens

### FR Coverage Map

| FR | Epic | Description |
|----|------|-------------|
| FR1-FR4 | Epic 1 | Account creation & profile storage |
| FR5-FR10 | Epic 1 | Sign-in/out, session validation |
| FR11-FR15 | Epic 2 | Email verification flow |
| FR16-FR19 | Epic 3 | Password reset flow |
| FR20-FR21 | Epic 4 | Password change + session revocation |
| FR22-FR27 | Epic 5 | Session listing & management |
| FR28-FR32 | Epic 2 | Email provider abstraction |
| FR33-FR36 | Epic 1 | Configuration & security baseline |

## Epic List

### Epic 1: Authentication Foundation
**User Outcome:** Users can register, sign in, sign out, and maintain authenticated sessions with real database persistence.

**What gets built:**
- Server refactor (handlers/middleware/router structure)
- PostgreSQL persistence via repository (replace in-memory hacks)
- Response helpers with better-auth error format
- Cookie-based session management

**FRs covered:** FR1-FR10, FR33-FR36
**NFRs addressed:** NFR2-NFR3, NFR5-NFR7, NFR11-NFR14, NFR16, NFR18-NFR19

---

### Epic 2: Email Verification
**User Outcome:** Users can verify their email address after registration. *(Alex Chen journey)*

**What gets built:**
- Email provider interface + SendGrid/Resend implementations
- Verification token generation (JWT)
- Send verification email endpoint
- Verify email endpoint

**FRs covered:** FR11-FR15, FR28-FR32
**NFRs addressed:** NFR4, NFR8, NFR15, NFR17

---

### Epic 3: Password Recovery
**User Outcome:** Users who forgot their password can recover account access. *(Marcus Webb journey)*

**What gets built:**
- Forget password endpoint (request reset)
- Reset password endpoint (with token)
- Reset tokens are single-use and expire

**FRs covered:** FR16-FR19
**NFRs addressed:** NFR8-NFR10

---

### Epic 4: Password Change
**User Outcome:** Authenticated users can proactively change their password for security. *(Priya Sharma journey)*

**What gets built:**
- Change password endpoint (requires current password)
- Optional "revoke other sessions" on password change

**FRs covered:** FR20-FR21
**NFRs addressed:** (covered by Epic 1 security patterns)

---

### Epic 5: Session Management
**User Outcome:** Users can view and control all their active sessions across devices. *(David Park journey)*

**What gets built:**
- List sessions endpoint (with device info)
- Revoke specific session endpoint
- Revoke all other sessions endpoint
- JWT/JWE session token caching for <5ms validation

**FRs covered:** FR22-FR27
**NFRs addressed:** NFR1, NFR6

---

## Epic Dependencies

```
Epic 1 (Foundation) ──┬──► Epic 2 (Email Verification)
                      │         │
                      │         └──► Epic 3 (Password Recovery) [uses email]
                      │
                      ├──► Epic 4 (Password Change)
                      │
                      └──► Epic 5 (Session Management)
```

Each epic is standalone after Epic 1. Epics 2-5 all depend on Epic 1 but not on each other.

---

## Epic 1: Authentication Foundation

**Goal:** Users can register, sign in, sign out, and maintain authenticated sessions with real database persistence.

**FRs covered:** FR1-FR10, FR33-FR36
**NFRs addressed:** NFR2-NFR3, NFR5-NFR7, NFR11-NFR14, NFR16, NFR18-NFR19

---

### Story 1.1: User Registration with Database Persistence

As a **new user**,
I want to **create an account with my email, password, and name**,
So that **I can access the application with persistent credentials**.

**Acceptance Criteria:**

**Given** a user visits the registration endpoint
**When** they submit valid email, password, and name
**Then** a new user record is created in PostgreSQL with UUID
**And** a new account record is created with hashed password (bcrypt cost 10+)
**And** a new session is created in the database
**And** the response returns `{ user, session }` matching better-auth format
**And** the `better-auth.session_token` cookie is set (HttpOnly, SameSite=Lax)

**Given** a user tries to register with an existing email
**When** they submit the registration request
**Then** the response returns 400 with `{ "error": { "code": "USER_ALREADY_EXISTS", "message": "..." } }`

**Given** a user submits invalid or missing fields
**When** the request is processed
**Then** the response returns 400 with `{ "error": { "code": "INVALID_BODY", "message": "..." } }`

**Technical Notes:**
- Refactor server structure: create `internal/handlers/auth.go`, `internal/handlers/response.go`
- Create `internal/middleware/cors.go` (extract from monolith)
- Wire handlers through `internal/server/server.go`
- Use existing sqlc-generated repository queries
- All JSON responses use camelCase field names

---

### Story 1.2: User Sign-In with Session Creation

As a **registered user**,
I want to **sign in with my email and password**,
So that **I can access my account and authenticated features**.

**Acceptance Criteria:**

**Given** a registered user with valid credentials
**When** they submit email and password to sign-in endpoint
**Then** credentials are validated against stored bcrypt hash
**And** a new session is created in PostgreSQL with IP address and user agent
**And** the response returns `{ user, session }` matching better-auth format
**And** the `better-auth.session_token` cookie is set

**Given** a user submits incorrect password
**When** the request is processed
**Then** the response returns 401 with `{ "error": { "code": "INVALID_CREDENTIALS", "message": "Invalid email or password" } }`
**And** the error message does NOT reveal whether email exists (NFR9)

**Given** a user submits non-existent email
**When** the request is processed
**Then** the response returns 401 with same generic `INVALID_CREDENTIALS` error
**And** response timing is consistent to prevent timing attacks

**Technical Notes:**
- Endpoint: `POST /api/auth/sign-in/email`
- Capture `X-Forwarded-For` or remote IP for session record
- Capture `User-Agent` header for session record

---

### Story 1.3: Session Validation & Sign-Out

As an **authenticated user**,
I want to **check my session status and sign out**,
So that **I can verify I'm logged in and securely end my session**.

**Acceptance Criteria:**

**Given** a user with a valid session cookie
**When** they request get-session endpoint
**Then** the session is looked up in PostgreSQL by token
**And** the response returns `{ user, session }` if valid and not expired
**And** the session `expiresAt` is checked against current time

**Given** a user with an expired or invalid session cookie
**When** they request get-session endpoint
**Then** the response returns `null` (not an error)

**Given** a user with no session cookie
**When** they request get-session endpoint
**Then** the response returns `null`

**Given** an authenticated user
**When** they request sign-out endpoint
**Then** the session is deleted from PostgreSQL
**And** the `better-auth.session_token` cookie is cleared
**And** the response returns `{ "success": true }`

**Given** a user without a valid session
**When** they request sign-out endpoint
**Then** the response returns `{ "success": true }` (idempotent)

**Technical Notes:**
- Endpoint: `GET /api/auth/get-session`
- Endpoint: `POST /api/auth/sign-out`
- Create `internal/middleware/auth.go` for session validation middleware (for use in later stories)
- Session expiry default: 7 days

---

## Epic 2: Email Verification

**Goal:** Users can verify their email address after registration. *(Alex Chen journey)*

**FRs covered:** FR11-FR15, FR28-FR32
**NFRs addressed:** NFR4, NFR8, NFR15, NFR17

---

### Story 2.1: Email Provider Infrastructure

As a **system administrator**,
I want to **configure email delivery through SendGrid or Resend**,
So that **the application can send transactional emails reliably**.

**Acceptance Criteria:**

**Given** the application starts with `EMAIL_PROVIDER=sendgrid` and valid SendGrid API key
**When** the configuration is loaded
**Then** the SendGrid email provider is initialized
**And** the provider is available for dependency injection into handlers

**Given** the application starts with `EMAIL_PROVIDER=resend` and valid Resend API key
**When** the configuration is loaded
**Then** the Resend email provider is initialized

**Given** the application starts with missing or invalid email configuration
**When** the configuration is validated
**Then** the application fails fast with a clear error message (NFR19)
**And** the error indicates which configuration is missing

**Given** an email provider is configured
**When** an email send is requested
**Then** the provider makes an HTTP POST to the provider's API
**And** the request includes proper authentication headers

**Technical Notes:**
- Create `internal/email/provider.go` with `EmailProvider` interface
- Create `internal/email/sendgrid.go` - direct HTTP calls to SendGrid API
- Create `internal/email/resend.go` - direct HTTP calls to Resend API
- Add to config: `EMAIL_PROVIDER`, `SENDGRID_API_KEY`, `RESEND_API_KEY`
- No external SDKs - use stdlib `net/http`

---

### Story 2.2: Send Verification Email

As a **registered user**,
I want to **receive a verification email after registration**,
So that **I can prove ownership of my email address**.

**Acceptance Criteria:**

**Given** an authenticated user with unverified email
**When** they request send-verification-email endpoint
**Then** a JWT verification token is generated with 1-hour expiry
**And** the token is stored in `verification` table with identifier `email-verification:{userId}`
**And** an email is sent asynchronously (fire-and-forget)
**And** the response returns immediately with `{ "success": true }` (NFR4)

**Given** an authenticated user requests verification email
**When** a previous verification token exists for this user
**Then** the old token is deleted before creating a new one
**And** only one active verification token exists per user

**Given** the email provider fails to send
**When** the error occurs
**Then** the error is logged with context (NFR17)
**And** the API response is NOT affected (already returned success)

**Given** an unauthenticated user
**When** they request send-verification-email endpoint
**Then** the response returns 401 `UNAUTHORIZED`

**Given** an already-verified user
**When** they request send-verification-email endpoint
**Then** the response returns 400 with appropriate message

**Technical Notes:**
- Endpoint: `POST /api/auth/send-verification-email`
- Requires auth middleware (from Story 1.3)
- Create `pkg/token/verification.go` for JWT token generation
- Use `BETTER_AUTH_SECRET` for JWT signing
- Add sqlc queries: `CreateVerification`, `DeleteVerificationByIdentifier`

---

### Story 2.3: Verify Email Address

As a **user with a verification link**,
I want to **click the link and verify my email**,
So that **my account is marked as verified and I can access all features**.

**Acceptance Criteria:**

**Given** a user clicks a valid verification link with token
**When** the verify-email endpoint is called
**Then** the JWT token signature is validated
**And** the token is looked up in `verification` table
**And** the token expiry is checked (1 hour default, NFR8)
**And** the user's `emailVerified` flag is set to `true`
**And** the verification record is deleted (single-use)
**And** the response returns `{ "success": true }` or redirects to success page

**Given** a user clicks an expired verification link
**When** the verify-email endpoint is called
**Then** the response returns 400 with `{ "error": { "code": "INVALID_TOKEN", "message": "Token has expired" } }`

**Given** a user clicks a verification link with invalid/tampered token
**When** the verify-email endpoint is called
**Then** the response returns 400 with `{ "error": { "code": "INVALID_TOKEN", "message": "Invalid token" } }`

**Given** a user clicks a verification link that was already used
**When** the verify-email endpoint is called
**Then** the response returns 400 with `INVALID_TOKEN` (token not found in DB)

**Technical Notes:**
- Endpoint: `GET /api/auth/verify-email?token={token}`
- Token validation: check signature, then check DB, then check expiry
- Add sqlc queries: `GetVerificationByIdentifierAndValue`, `DeleteVerification`
- Add sqlc query: `UpdateUserEmailVerified`

---

## Epic 3: Password Recovery

**Goal:** Users who forgot their password can recover account access. *(Marcus Webb journey)*

**FRs covered:** FR16-FR19
**NFRs addressed:** NFR8-NFR10

---

### Story 3.1: Request Password Reset

As a **user who forgot their password**,
I want to **request a password reset email**,
So that **I can regain access to my account**.

**Acceptance Criteria:**

**Given** a user submits their email to forget-password endpoint
**When** the email exists in the system
**Then** a JWT reset token is generated with 1-hour expiry
**And** the token is stored in `verification` table with identifier `reset-password:{userId}`
**And** a password reset email is sent asynchronously
**And** the response returns `{ "success": true }`

**Given** a user submits an email that doesn't exist
**When** the request is processed
**Then** the response STILL returns `{ "success": true }` (NFR9 - no email enumeration)
**And** no email is sent
**And** response timing is consistent with successful case

**Given** a user requests reset when a previous token exists
**When** the request is processed
**Then** the old token is deleted before creating a new one

**Given** the email provider fails
**When** the error occurs
**Then** the error is logged but response is unaffected

**Technical Notes:**
- Endpoint: `POST /api/auth/forget-password`
- Request body: `{ "email": "user@example.com" }`
- Reuse `pkg/token/verification.go` with different identifier prefix
- Use constant-time operations to prevent timing attacks

---

### Story 3.2: Reset Password with Token

As a **user with a reset link**,
I want to **set a new password using the link**,
So that **I can access my account with new credentials**.

**Acceptance Criteria:**

**Given** a user submits a valid reset token and new password
**When** the reset-password endpoint is called
**Then** the JWT token signature is validated
**And** the token is looked up in `verification` table
**And** the token expiry is checked (1 hour, NFR8)
**And** the new password is hashed with bcrypt (cost 10+)
**And** the user's password is updated in `account` table
**And** the reset token is deleted immediately (NFR10 - single use)
**And** ALL user sessions are invalidated (security best practice)
**And** the response returns `{ "success": true }`

**Given** a user submits an expired reset token
**When** the reset-password endpoint is called
**Then** the response returns 400 with `INVALID_TOKEN`

**Given** a user submits an invalid/tampered token
**When** the reset-password endpoint is called
**Then** the response returns 400 with `INVALID_TOKEN`

**Given** a user submits a token that was already used
**When** the reset-password endpoint is called
**Then** the response returns 400 with `INVALID_TOKEN`

**Given** a user submits a weak or empty password
**When** the reset-password endpoint is called
**Then** the response returns 400 with `INVALID_BODY`

**Technical Notes:**
- Endpoint: `POST /api/auth/reset-password`
- Request body: `{ "token": "...", "newPassword": "..." }`
- Add sqlc query: `UpdateAccountPassword`
- Add sqlc query: `DeleteUserSessions` (invalidate all sessions)

---

## Epic 4: Password Change

**Goal:** Authenticated users can proactively change their password for security. *(Priya Sharma journey)*

**FRs covered:** FR20-FR21
**NFRs addressed:** (covered by Epic 1 security patterns)

---

### Story 4.1: Change Password (Authenticated)

As an **authenticated user**,
I want to **change my password by providing my current password**,
So that **I can improve my account security proactively**.

**Acceptance Criteria:**

**Given** an authenticated user with valid current password
**When** they submit current password and new password to change-password endpoint
**Then** the current password is validated against stored hash
**And** the new password is hashed with bcrypt (cost 10+)
**And** the password is updated in `account` table
**And** the response returns `{ "success": true }`

**Given** an authenticated user submits incorrect current password
**When** the request is processed
**Then** the response returns 400 with `{ "error": { "code": "INVALID_PASSWORD", "message": "Current password is incorrect" } }`

**Given** an authenticated user wants to revoke other sessions
**When** they include `revokeOtherSessions: true` in the request
**Then** the password is changed
**And** all sessions EXCEPT the current one are deleted
**And** the response returns `{ "success": true }`

**Given** an authenticated user changes password without revoking sessions
**When** they omit or set `revokeOtherSessions: false`
**Then** other sessions remain active

**Given** an unauthenticated user
**When** they request change-password endpoint
**Then** the response returns 401 `UNAUTHORIZED`

**Technical Notes:**
- Endpoint: `POST /api/auth/change-password`
- Request body: `{ "currentPassword": "...", "newPassword": "...", "revokeOtherSessions": false }`
- Requires auth middleware
- Reuse `DeleteOtherSessions` query from architecture spec

---

## Epic 5: Session Management

**Goal:** Users can view and control all their active sessions across devices. *(David Park journey)*

**FRs covered:** FR22-FR27
**NFRs addressed:** NFR1, NFR6

---

### Story 5.1: List Active Sessions

As an **authenticated user**,
I want to **see a list of all my active sessions**,
So that **I can monitor which devices have access to my account**.

**Acceptance Criteria:**

**Given** an authenticated user
**When** they request list-sessions endpoint
**Then** all sessions for the user are retrieved from PostgreSQL
**And** the response returns an array of session objects
**And** each session includes: id, createdAt, expiresAt, ipAddress, userAgent
**And** the current session is identifiable (e.g., `isCurrent: true` flag)

**Given** a user with multiple sessions
**When** they request the list
**Then** sessions are ordered by createdAt descending (most recent first)

**Given** an unauthenticated user
**When** they request list-sessions endpoint
**Then** the response returns 401 `UNAUTHORIZED`

**Technical Notes:**
- Endpoint: `GET /api/auth/list-sessions`
- Requires auth middleware
- Add sqlc query: `ListUserSessions`
- Compare session token to mark current session

---

### Story 5.2: Revoke Specific Session

As an **authenticated user**,
I want to **revoke a specific session by its ID**,
So that **I can terminate access on a device I no longer trust**.

**Acceptance Criteria:**

**Given** an authenticated user with a valid session ID
**When** they request revoke-session endpoint with the session ID
**Then** the session is deleted from PostgreSQL
**And** the response returns `{ "success": true }`

**Given** a user tries to revoke another user's session
**When** the request is processed
**Then** the session is NOT deleted (query filters by userId)
**And** the response returns 404 or success (don't leak info)

**Given** a user tries to revoke a non-existent session
**When** the request is processed
**Then** the response returns `{ "success": true }` (idempotent)

**Given** a user tries to revoke their current session
**When** the request is processed
**Then** the session is deleted
**And** they are effectively signed out

**Given** an unauthenticated user
**When** they request revoke-session endpoint
**Then** the response returns 401 `UNAUTHORIZED`

**Technical Notes:**
- Endpoint: `POST /api/auth/revoke-session`
- Request body: `{ "sessionId": "uuid-here" }`
- Requires auth middleware
- Add sqlc query: `DeleteSessionByID` (with userId filter for security)

---

### Story 5.3: Revoke All Other Sessions

As an **authenticated user**,
I want to **revoke all sessions except my current one**,
So that **I can secure my account while staying logged in**.

**Acceptance Criteria:**

**Given** an authenticated user with multiple sessions
**When** they request revoke-other-sessions endpoint
**Then** all sessions EXCEPT the current one are deleted
**And** the current session remains active
**And** the response returns `{ "success": true }`

**Given** a user with only one session (current)
**When** they request revoke-other-sessions endpoint
**Then** no sessions are deleted
**And** the response returns `{ "success": true }`

**Given** an unauthenticated user
**When** they request revoke-other-sessions endpoint
**Then** the response returns 401 `UNAUTHORIZED`

**Technical Notes:**
- Endpoint: `POST /api/auth/revoke-other-sessions`
- Requires auth middleware
- Add sqlc query: `DeleteOtherSessions` (userId + exclude current sessionId)

---

### Story 5.4: JWT Session Token Caching

As a **system**,
I want to **cache session data in signed cookies**,
So that **session validation completes in <5ms without database hits (NFR1)**.

**Acceptance Criteria:**

**Given** a user signs in or session is validated
**When** the session is confirmed valid in database
**Then** a signed cookie cache is created with: userId, sessionId, expiresAt
**And** the cache uses HMAC-SHA256 with `BETTER_AUTH_SECRET`
**And** the cache is stored in the existing `better-auth.session_token` cookie

**Given** a request with a valid cached session cookie
**When** the auth middleware processes the request
**Then** the signature is validated (no DB hit)
**And** the cache expiry is checked (5 min default)
**And** if cache is valid, request proceeds without DB lookup
**And** validation completes in <5ms

**Given** a cached session cookie that has expired (>5 min)
**When** the auth middleware processes the request
**Then** the session is validated against the database
**And** if still valid, a fresh cache is generated
**And** the new cache is set in the response cookie

**Given** a cached session cookie with invalid signature
**When** the auth middleware processes the request
**Then** the cache is rejected
**And** the session is validated against database as fallback

**Given** a session that was revoked
**When** the cached cookie is still valid but session deleted from DB
**Then** on next cache refresh (after 5 min), the session is invalidated
**And** user is effectively logged out

**Technical Notes:**
- Create `pkg/session/cache.go`
- Cookie format: base64url(payload) + "." + base64url(signature)
- Payload: JSON with userId, sessionId, expiresAt, cacheExpiresAt
- Config: `SESSION_CACHE_TTL` (default 5 minutes)
- This is the "compact" strategy from better-auth
