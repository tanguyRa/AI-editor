---
stepsCompleted: ['step-01-init', 'step-02-discovery', 'step-03-success', 'step-04-journeys', 'step-07-project-type', 'step-08-scoping', 'step-09-functional', 'step-10-nonfunctional', 'step-11-complete']
inputDocuments:
  - docs/api-contracts.md
  - docs/data-models.md
workflowType: 'prd'
lastStep: 11
documentCounts:
  briefs: 0
  research: 0
  projectDocs: 2
---

# Product Requirements Document - AI-editor

**Author:** Tanguy
**Date:** 2026-01-14

## Executive Summary

AI-editor's auth system is a Go-native implementation of better-auth server, providing a drop-in replacement for the TypeScript server while maintaining full API and database compatibility with better-auth clients.

Phase 1 focuses on completing core authentication: email verification, password reset, password change, and multi-session management — all backed by PostgreSQL with JWT/JWE token support for session caching.

### What Makes This Special

Better-auth is the most comprehensive authentication framework in the TypeScript ecosystem, but its TypeScript-only approach limits adoption for teams using other backend languages. This project brings better-auth's well-designed API and database schema to Go, offering:

- **Drop-in compatibility** — Same REST endpoints, same database schema, same client libraries work unchanged
- **Go advantages** — Faster performance, stronger type safety, lower memory footprint, better security posture
- **Reusable library** — Designed to be extracted and dropped into any Go project needing auth
- **No vendor lock-in** — Use better-auth's excellent frontend client with a backend you control

## Project Classification

**Technical Type:** api_backend
**Domain:** general
**Complexity:** Medium
**Project Context:** Brownfield - extending existing boilerplate with production auth system

This is a backend library/service project focused on API compatibility and security. The existing codebase provides the foundation (Nuxt frontend, Go API structure, PostgreSQL schema), and this PRD defines the completion of the auth system as a reusable component.

## Success Criteria

### User Success

- **Drop-in experience:** A Go developer can replace the better-auth TypeScript server and have their existing better-auth client work unchanged
- **Middleware integration:** Authenticated requests can be validated directly in Go middleware — no server proxying, no JWT forwarding to external services
- **Quick integration:** From "I found this" to "it's running in my project" in under an hour

### Business Success

- **Personal utility:** Usable in the author's own Go projects as the primary auth solution
- **Learning outcome:** Deep understanding of better-auth protocols and auth system design
- **Foundation for growth:** Architecture supports gradual feature expansion toward full better-auth parity

### Technical Success

- **API compatibility:** Core better-auth endpoints return identical response shapes — existing clients work without modification
- **Database compatibility:** Uses better-auth's database schema — can share DB with TypeScript server if needed
- **Security best practices:** bcrypt password hashing, secure session tokens (JWT/JWE), HttpOnly cookies, proper CORS
- **Clean architecture:** Email provider abstraction, config-driven behavior, sqlc for type-safe queries

### Measurable Outcomes

- All existing better-auth client methods for core auth work unchanged
- Session validation adds < 5ms overhead when using cached JWT tokens
- Zero breaking changes to existing sign-up/sign-in/sign-out flows during migration

## Product Scope

### MVP - Minimum Viable Product (Phase 1)

**Core Auth with PostgreSQL:**
- Wire PostgreSQL persistence (replace in-memory)
- Email verification flow (send, verify)
- Password reset flow (request, reset)
- Change password (authenticated)
- Multi-session management (list, revoke, revoke others)
- JWT/JWE session tokens with cookie caching
- Email sender abstraction with SendGrid + Resend implementations
- Config validation for email provider selection

### Growth Features (Phase 2)

**OAuth Providers:**
- Social login support (Google, GitHub, etc.)
- OAuth token storage and refresh
- Account linking

### Vision (Phase 3+)

**Full Better-Auth Parity:**
- Passkeys / WebAuthn
- Organizations / multi-tenancy
- 2FA (TOTP, SMS)
- Additional OAuth providers
- Plugin architecture mirroring better-auth's extensibility

## User Journeys

### Journey 1: Alex Chen — First-Time Sign-Up

Alex discovers an app built with this auth system and decides to create an account. She enters her email, name, and a strong password. The system creates her account and immediately sends a verification email. Alex checks her inbox, clicks the verification link, and lands back in the app with a "Email verified!" confirmation. She's now fully onboarded and can access all features.

*What happens if she doesn't verify?* She can still sign in, but sees a persistent banner reminding her to verify. Some features may be gated until verification is complete.

### Journey 2: Marcus Webb — Forgot My Password

Marcus hasn't logged in for a few weeks and can't remember his password. On the sign-in page, he clicks "Forgot password?" and enters his email. Within seconds, he receives an email with a reset link. He clicks it, enters a new password twice, and submits. The system confirms the reset, invalidates his old sessions for security, and redirects him to sign in with his new credentials. He's back in.

*What if the link expires?* The reset link is valid for 1 hour. If Marcus waits too long, he clicks the link and sees "This link has expired. Request a new one." — one click to restart the flow.

### Journey 3: Priya Sharma — Changing Her Password

Priya is already signed in but wants to update her password after reading about a data breach at another service. She navigates to account settings, enters her current password for verification, then enters her new password twice. The system confirms the change and asks if she wants to sign out of other devices. She chooses yes — all other sessions are revoked, but her current session stays active. Peace of mind restored.

### Journey 4: David Park — "What Devices Am I Signed In On?"

David gets a security alert email from another service and becomes paranoid. He signs into the app and navigates to "Active Sessions" in his security settings. He sees a list: his MacBook (current session), his iPhone, and... a Windows PC he doesn't recognize from an IP in another country. He immediately clicks "Revoke" on the suspicious session. The system confirms it's been terminated. He then clicks "Sign out all other devices" just to be safe, leaving only his current session active. He also changes his password for good measure.

### Journey Requirements Summary

| Journey | Capabilities Required |
|---------|----------------------|
| Sign-up & Verification | Create user, create session, send verification email, verify token, update emailVerified flag |
| Password Recovery | Request reset (create verification token), send reset email, validate token, update password, invalidate sessions |
| Password Change | Validate current password, update password, optionally revoke other sessions |
| Session Management | List user sessions, revoke specific session, revoke all other sessions |

## API Backend Specific Requirements

### Project-Type Overview

This is a REST API backend implementing better-auth server protocol in Go. The auth endpoints (`/api/auth/*`) strictly follow better-auth's API specification for client compatibility, while remaining independent of conventions used by other application endpoints.

### Endpoint Specification

**Auth Endpoints (Phase 1 MVP):**

| Endpoint | Method | Description | Status |
|----------|--------|-------------|--------|
| `/api/auth/sign-up/email` | POST | Create account with email/password | Exists |
| `/api/auth/sign-in/email` | POST | Authenticate with email/password | Exists |
| `/api/auth/sign-out` | POST | End current session | Exists |
| `/api/auth/get-session` | GET | Get current session | Exists |
| `/api/auth/send-verification-email` | POST | Send email verification | New |
| `/api/auth/verify-email` | GET | Verify email with token | New |
| `/api/auth/forget-password` | POST | Request password reset | New |
| `/api/auth/reset-password` | POST | Reset password with token | New |
| `/api/auth/change-password` | POST | Change password (authenticated) | New |
| `/api/auth/list-sessions` | GET | List all user sessions | New |
| `/api/auth/revoke-session` | POST | Revoke specific session | New |
| `/api/auth/revoke-other-sessions` | POST | Revoke all except current | New |

### Authentication Model

- **Session tokens:** JWT/JWE format with HMAC-SHA256 signing (compact) or full encryption (JWE)
- **Cookie handling:** `better-auth.session_token`, HttpOnly, SameSite=Lax, 7-day default expiry
- **Session caching:** Signed cookie payload reduces DB lookups; configurable cache TTL
- **Password hashing:** bcrypt with configurable cost factor

### Data Schemas

All request/response schemas match better-auth exactly:
- User object: `{ id, email, name, image, emailVerified, createdAt, updatedAt }`
- Session object: `{ id, userId, token, expiresAt, ipAddress, userAgent, createdAt, updatedAt }`
- Error format: `{ "error": { "code": "ERROR_CODE", "message": "Human readable" } }`

### Rate Limiting

**Defaults (configurable via environment variables):**

| Endpoint Pattern | Limit | Window | Env Variable |
|-----------------|-------|--------|--------------|
| `/api/auth/sign-in/*` | 5 requests | 1 minute | `RATE_LIMIT_SIGNIN` |
| `/api/auth/forget-password` | 3 requests | 1 minute | `RATE_LIMIT_RESET` |
| `/api/auth/sign-up/*` | 10 requests | 1 minute | `RATE_LIMIT_SIGNUP` |
| `/api/auth/*` (default) | 60 requests | 1 minute | `RATE_LIMIT_AUTH_DEFAULT` |

Rate limiting keyed by IP address. Returns `429 Too Many Requests` with `Retry-After` header.

### API Versioning

**Auth endpoints:** No versioning — `/api/auth/*` matches better-auth exactly
**Other endpoints:** May use `/api/v1/*` convention independently

### Implementation Considerations

- **Compatibility principle:** When in doubt, match better-auth behavior exactly
- **Error isolation:** Auth endpoints use better-auth error format; other endpoints can differ
- **Config-driven:** All tunables (rate limits, token expiry, email provider) configurable via env
- **Database:** PostgreSQL with sqlc-generated type-safe queries

## Project Scoping & Phased Development

### MVP Strategy & Philosophy

**MVP Approach:** Platform MVP — build the core authentication foundation that all future features layer onto

**Resource Model:** Solo developer, personal use first, potential extraction to library later

**Core Principle:** Ship working auth that handles the 4 essential user journeys reliably before adding bells and whistles

### MVP Feature Set (Phase 1)

**Core User Journeys Supported:**
- Sign-up & Email Verification (Alex Chen journey)
- Password Recovery (Marcus Webb journey)
- Password Change (Priya Sharma journey)
- Session Management (David Park journey)

**Must-Have Capabilities:**
- PostgreSQL persistence (replace in-memory storage)
- Email verification flow (send verification, verify token)
- Password reset flow (request reset, reset with token)
- Change password (authenticated, with optional session revocation)
- Multi-session management (list, revoke one, revoke others)
- JWT/JWE session tokens with cookie caching
- Email sender abstraction with SendGrid + Resend implementations
- Config-driven email provider selection

**Explicitly Deferred from MVP:**
- Rate limiting (Phase 2)
- OAuth providers (Phase 2)

### Post-MVP Features

**Phase 2 (Growth):**
- Rate limiting with configurable defaults
- OAuth providers (Google, GitHub)
- OAuth token storage and refresh
- Account linking (connect social to existing account)

**Phase 3 (Expansion):**
- Passkeys / WebAuthn
- 2FA (TOTP, SMS)
- Organizations / multi-tenancy
- Additional OAuth providers
- Plugin architecture

### Risk Mitigation Strategy

**Technical Risks:**
- JWT/JWE implementation complexity → Use well-tested Go libraries (golang-jwt)
- Email deliverability → Abstract provider, easy to swap if issues arise

**Scope Risks:**
- Feature creep → Strict better-auth API compatibility keeps scope bounded
- Perfectionism → Ship when 4 journeys work, iterate from there

## Functional Requirements

### Account Management

- FR1: Users can create an account with email, password, and display name
- FR2: System generates unique user ID upon account creation
- FR3: System stores user profile information (name, email, image URL)
- FR4: System tracks account creation and update timestamps

### Authentication

- FR5: Users can sign in with email and password
- FR6: System validates credentials against stored password hash
- FR7: System creates a new session upon successful authentication
- FR8: Users can sign out and terminate their current session
- FR9: System can validate whether a request has a valid session
- FR10: System captures IP address and user agent for each session

### Email Verification

- FR11: System can send verification emails to users
- FR12: Users can verify their email by clicking a verification link
- FR13: System validates verification tokens and marks email as verified
- FR14: System tracks email verification status on user accounts
- FR15: Verification tokens expire after a configurable time period

### Password Management

- FR16: Users can request a password reset via email
- FR17: System sends password reset emails with secure tokens
- FR18: Users can reset their password using a valid reset token
- FR19: System invalidates reset tokens after use or expiration
- FR20: Authenticated users can change their password by providing current password
- FR21: Password changes can optionally revoke all other sessions

### Session Management

- FR22: Users can view a list of all their active sessions
- FR23: Session list displays device info (IP address, user agent)
- FR24: Users can revoke a specific session by its identifier
- FR25: Users can revoke all sessions except the current one
- FR26: System automatically expires sessions after configurable duration
- FR27: System supports JWT/JWE encoded session tokens for caching

### Email Delivery

- FR28: System abstracts email sending behind a provider interface
- FR29: System supports SendGrid as an email provider
- FR30: System supports Resend as an email provider
- FR31: Email provider is selected via configuration at startup
- FR32: System validates email configuration on startup

### Configuration & Security

- FR33: All tunables (token expiry, session duration) are configurable via environment
- FR34: System uses bcrypt for password hashing
- FR35: Session tokens are stored in HttpOnly cookies
- FR36: System supports CORS with configurable allowed origins

## Non-Functional Requirements

### Performance

- NFR1: Session validation using cached JWT tokens completes in < 5ms
- NFR2: Auth endpoints (sign-in, sign-up) respond within 500ms under normal load
- NFR3: Database queries use indexes for user/session lookups (no table scans)
- NFR4: Email sending is asynchronous — does not block API response

### Security

- NFR5: Passwords are hashed using bcrypt with minimum cost factor of 10
- NFR6: Session tokens are cryptographically signed (JWT) or encrypted (JWE)
- NFR7: All session cookies are HttpOnly, SameSite=Lax, Secure in production
- NFR8: Verification and reset tokens expire after 1 hour (configurable)
- NFR9: Failed authentication attempts do not reveal whether email exists
- NFR10: Password reset invalidates the token immediately after use
- NFR11: CORS allows only configured origins (no wildcard in production)

### Integration

- NFR12: All auth endpoints match better-auth API specification exactly
- NFR13: Response shapes are compatible with better-auth TypeScript client
- NFR14: Database schema matches better-auth PostgreSQL adapter
- NFR15: Email provider interface allows adding new providers without core changes

### Reliability

- NFR16: Invalid or expired sessions return consistent error responses (not crashes)
- NFR17: Email provider failures are logged but do not fail the parent operation
- NFR18: Database connection failures are handled gracefully with appropriate errors
- NFR19: Configuration validation fails fast at startup (not runtime)

