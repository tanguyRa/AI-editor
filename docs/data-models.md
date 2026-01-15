# Data Models

## Overview

The database uses PostgreSQL with UUID primary keys. All tables support the better-auth schema requirements.

## Entity Relationship Diagram

```
┌─────────────────┐
│      user       │
├─────────────────┤
│ id (PK)         │
│ name            │
│ email (UNIQUE)  │
│ emailVerified   │
│ image           │
│ createdAt       │
│ updatedAt       │
└────────┬────────┘
         │
         │ 1:N
         ▼
┌─────────────────┐     ┌─────────────────┐
│    session      │     │    account      │
├─────────────────┤     ├─────────────────┤
│ id (PK)         │     │ id (PK)         │
│ userId (FK)     │     │ userId (FK)     │
│ token (UNIQUE)  │     │ accountId       │
│ expiresAt       │     │ providerId      │
│ ipAddress       │     │ accessToken     │
│ userAgent       │     │ refreshToken    │
│ createdAt       │     │ password        │
│ updatedAt       │     │ ...             │
└─────────────────┘     └─────────────────┘

┌─────────────────┐
│  verification   │
├─────────────────┤
│ id (PK)         │
│ identifier      │
│ value           │
│ expiresAt       │
│ createdAt       │
│ updatedAt       │
└─────────────────┘
```

## Tables

### user

Stores user account information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| name | VARCHAR(255) | NOT NULL | User's display name |
| email | VARCHAR(255) | UNIQUE, NOT NULL | Email address |
| emailVerified | BOOLEAN | NOT NULL, DEFAULT FALSE | Email verification status |
| image | TEXT | NULLABLE | Profile image URL |
| createdAt | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updatedAt | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

### session

Stores active user sessions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| userId | UUID | FK → user(id) ON DELETE CASCADE | Associated user |
| token | VARCHAR(255) | UNIQUE, NOT NULL | Session token |
| expiresAt | TIMESTAMPTZ | NOT NULL | Expiration timestamp |
| ipAddress | VARCHAR(45) | NULLABLE | Client IP address |
| userAgent | VARCHAR | NULLABLE | Client user agent |
| createdAt | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updatedAt | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Indexes:**
- `idx_session_user_id` ON (userId)
- `idx_session_token` ON (token)

### account

Stores OAuth provider accounts and credential accounts.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| userId | UUID | FK → user(id) ON DELETE CASCADE | Associated user |
| accountId | VARCHAR(255) | NOT NULL | Provider account ID |
| providerId | VARCHAR(50) | NOT NULL | Provider identifier |
| accessToken | TEXT | NULLABLE | OAuth access token |
| refreshToken | TEXT | NULLABLE | OAuth refresh token |
| accessTokenExpiresAt | TIMESTAMPTZ | NULLABLE | Access token expiry |
| refreshTokenExpiresAt | TIMESTAMPTZ | NULLABLE | Refresh token expiry |
| scope | TEXT | NULLABLE | OAuth scopes |
| idToken | TEXT | NULLABLE | OAuth ID token |
| password | TEXT | NULLABLE | Hashed password (credential provider) |
| createdAt | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updatedAt | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Constraints:**
- UNIQUE (providerId, accountId)

**Indexes:**
- `idx_accounts_user_id` ON (userId)

### verification

Stores email verification and password reset tokens.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() | Unique identifier |
| identifier | VARCHAR(255) | NOT NULL | Email or identifier |
| value | VARCHAR(255) | NOT NULL | Verification code/token |
| expiresAt | TIMESTAMPTZ | NOT NULL | Expiration timestamp |
| createdAt | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Creation timestamp |
| updatedAt | TIMESTAMPTZ | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Last update timestamp |

**Indexes:**
- `idx_verifications_identifier` ON (identifier)

## Go Models (sqlc generated)

### repository.User

```go
type User struct {
    ID            uuid.UUID  `json:"id"`
    Name          string     `json:"name"`
    Email         string     `json:"email"`
    EmailVerified bool       `json:"emailVerified"`
    Image         *string    `json:"image"`
    CreatedAt     time.Time  `json:"createdAt"`
    UpdatedAt     time.Time  `json:"updatedAt"`
}
```

### repository.Session

```go
type Session struct {
    ID        uuid.UUID  `json:"id"`
    UserId    uuid.UUID  `json:"userId"`
    Token     string     `json:"token"`
    ExpiresAt time.Time  `json:"expiresAt"`
    IpAddress *string    `json:"ipAddress"`
    UserAgent *string    `json:"userAgent"`
    CreatedAt time.Time  `json:"createdAt"`
    UpdatedAt time.Time  `json:"updatedAt"`
}
```

### repository.Account

```go
type Account struct {
    ID                    uuid.UUID   `json:"id"`
    UserId                uuid.UUID   `json:"userId"`
    AccountId             string      `json:"accountId"`
    ProviderId            string      `json:"providerId"`
    AccessToken           *string     `json:"accessToken"`
    RefreshToken          *string     `json:"refreshToken"`
    AccessTokenExpiresAt  *time.Time  `json:"accessTokenExpiresAt"`
    RefreshTokenExpiresAt *time.Time  `json:"refreshTokenExpiresAt"`
    Scope                 *string     `json:"scope"`
    IdToken               *string     `json:"idToken"`
    Password              *string     `json:"password"`
    CreatedAt             time.Time   `json:"createdAt"`
    UpdatedAt             time.Time   `json:"updatedAt"`
}
```

### repository.Verification

```go
type Verification struct {
    ID         uuid.UUID `json:"id"`
    Identifier string    `json:"identifier"`
    Value      string    `json:"value"`
    ExpiresAt  time.Time `json:"expiresAt"`
    CreatedAt  time.Time `json:"createdAt"`
    UpdatedAt  time.Time `json:"updatedAt"`
}
```

## Migrations

### Migration: 000001_init_auth

**Up Migration** (`000001_init_auth.up.sql`):
- Creates `user`, `session`, `account`, `verification` tables
- Creates indexes for performance

**Down Migration** (`000001_init_auth.down.sql`):
- Drops all tables and indexes

### Running Migrations

```bash
# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down

# Drop all (dangerous)
make migrate-drop
```

## SQL Queries

### User Queries (db/queries/users.sql)

| Query | Description |
|-------|-------------|
| `ListUsers` | Get all users |
| `CreateUser` | Insert new user |
| `GetUserByID` | Get user by UUID |
| `GetUserByEmail` | Get user by email |
| `UpdateUser` | Update user fields |
| `DeleteUser` | Delete user by ID |

### Session Queries (db/queries/sessions.sql)

| Query | Description |
|-------|-------------|
| `CreateSession` | Insert new session |
| `GetSession` | Get session by ID |
| `UpdateSession` | Update session |
| `DeleteSession` | Delete session by ID |
| `DeleteUserSessions` | Delete all sessions for user |
