# API Contracts

## Overview

The backend exposes a REST API compatible with [better-auth](https://better-auth.com/). All endpoints are prefixed with `/api/auth/`.

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: Configured via `BETTER_AUTH_URL`

## Authentication Endpoints

### POST /api/auth/sign-up/email

Create a new user account with email and password.

**Request**

```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}
```

**Response (200 OK)**

```json
{
  "user": {
    "id": "a1b2c3d4e5f6...",
    "email": "user@example.com",
    "name": "John Doe",
    "image": null,
    "emailVerified": false,
    "createdAt": "2026-01-14T12:00:00Z",
    "updatedAt": "2026-01-14T12:00:00Z"
  },
  "session": {
    "id": "s1e2s3s4i5o6...",
    "userId": "a1b2c3d4e5f6...",
    "token": "abc123...",
    "expiresAt": "2026-01-21T12:00:00Z",
    "createdAt": "2026-01-14T12:00:00Z",
    "updatedAt": "2026-01-14T12:00:00Z"
  }
}
```

**Errors**

| Status | Code | Message |
|--------|------|---------|
| 400 | `INVALID_BODY` | Invalid request body |
| 400 | `USER_ALREADY_EXISTS` | User with this email already exists |
| 500 | `INTERNAL_ERROR` | Failed to hash password |

**Side Effects**

- Sets `better-auth.session_token` cookie (HttpOnly, 7 days)

---

### POST /api/auth/sign-in/email

Authenticate with email and password.

**Request**

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response (200 OK)**

```json
{
  "user": {
    "id": "a1b2c3d4e5f6...",
    "email": "user@example.com",
    "name": "John Doe",
    "image": null,
    "emailVerified": false,
    "createdAt": "2026-01-14T12:00:00Z",
    "updatedAt": "2026-01-14T12:00:00Z"
  },
  "session": {
    "id": "s1e2s3s4i5o6...",
    "userId": "a1b2c3d4e5f6...",
    "token": "abc123...",
    "expiresAt": "2026-01-21T12:00:00Z",
    "createdAt": "2026-01-14T12:00:00Z",
    "updatedAt": "2026-01-14T12:00:00Z"
  }
}
```

**Errors**

| Status | Code | Message |
|--------|------|---------|
| 400 | `INVALID_BODY` | Invalid request body |
| 401 | `INVALID_CREDENTIALS` | Invalid email or password |

**Side Effects**

- Sets `better-auth.session_token` cookie (HttpOnly, 7 days)

---

### GET /api/auth/get-session

Get the current user session.

**Request**

Requires `better-auth.session_token` cookie.

**Response (200 OK) - Authenticated**

```json
{
  "user": {
    "id": "a1b2c3d4e5f6...",
    "email": "user@example.com",
    "name": "John Doe",
    "image": null,
    "emailVerified": false,
    "createdAt": "2026-01-14T12:00:00Z",
    "updatedAt": "2026-01-14T12:00:00Z"
  },
  "session": {
    "id": "s1e2s3s4i5o6...",
    "userId": "a1b2c3d4e5f6...",
    "token": "abc123...",
    "expiresAt": "2026-01-21T12:00:00Z",
    "createdAt": "2026-01-14T12:00:00Z",
    "updatedAt": "2026-01-14T12:00:00Z"
  }
}
```

**Response (200 OK) - Not Authenticated**

```json
null
```

---

### POST /api/auth/sign-out

Sign out and invalidate the current session.

**Request**

Requires `better-auth.session_token` cookie.

**Response (200 OK)**

```json
{
  "success": true
}
```

**Side Effects**

- Deletes session from storage
- Clears `better-auth.session_token` cookie

---

## Health Check

### GET /health

Check API health status.

**Response (200 OK)**

```json
{
  "status": "ok"
}
```

---

## CORS Configuration

The API supports CORS with the following settings:

| Header | Value |
|--------|-------|
| `Access-Control-Allow-Origin` | Request origin |
| `Access-Control-Allow-Credentials` | `true` |
| `Access-Control-Allow-Methods` | `GET, POST, PUT, DELETE, OPTIONS` |
| `Access-Control-Allow-Headers` | `Content-Type, Authorization` |

Preflight `OPTIONS` requests return `200 OK`.

---

## Error Response Format

All errors follow this format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

---

## Session Cookie

| Property | Value |
|----------|-------|
| Name | `better-auth.session_token` |
| Path | `/` |
| MaxAge | 7 days (604800 seconds) |
| HttpOnly | `true` |
| SameSite | `Lax` |
| Secure | `false` (enable in production) |
