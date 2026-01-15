# AI-Editor Project Overview

## Executive Summary

AI-Editor is a full-stack web application built with a modern TypeScript/Vue frontend and a Go backend. The project follows a multi-part architecture with separate frontend and backend services, orchestrated via Docker Compose.

## Project Information

| Property | Value |
|----------|-------|
| **Project Name** | AI-Editor (budhapp) |
| **Repository Type** | Multi-part (client/server) |
| **Primary Languages** | TypeScript, Go |
| **Package Manager** | Bun (frontend), Go Modules (backend) |

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        User Browser                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (app/)                          │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Nuxt 4.2 + Vue 3 + TypeScript                      │    │
│  │  - Nuxt UI Components                               │    │
│  │  - better-auth Client                               │    │
│  │  - SSR/SSG capable                                  │    │
│  └─────────────────────────────────────────────────────┘    │
│                         Port 3000                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend (api/)                           │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Go 1.24 + net/http                                 │    │
│  │  - better-auth compatible REST API                  │    │
│  │  - Session-based authentication                     │    │
│  │  - sqlc generated type-safe queries                 │    │
│  └─────────────────────────────────────────────────────┘    │
│                         Port 8080                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      PostgreSQL                             │
│  - Users, Sessions, Accounts, Verifications                 │
└─────────────────────────────────────────────────────────────┘
```

## Technology Stack Summary

### Frontend (app/)
- **Framework**: Nuxt 4.2 with Vue 3
- **UI Library**: Nuxt UI 4.3
- **Language**: TypeScript 5.9
- **Auth**: better-auth client
- **Validation**: Zod
- **Package Manager**: Bun

### Backend (api/)
- **Language**: Go 1.24
- **HTTP**: net/http (stdlib)
- **Database**: PostgreSQL via pgx/v5
- **Code Generation**: sqlc
- **Migrations**: golang-migrate
- **Dev Server**: Air (hot reload)

## Key Features

1. **User Authentication**
   - Email/password registration and login
   - Session-based authentication with secure cookies
   - Password hashing with bcrypt

2. **Database**
   - PostgreSQL with UUID primary keys
   - Type-safe SQL queries via sqlc
   - Database migrations with golang-migrate

3. **Development Experience**
   - Docker Compose for local development
   - Hot reload on both frontend and backend
   - Makefile for common operations

## Repository Structure

```
AI-editor/
├── app/                 # Nuxt frontend
├── api/                 # Go backend
├── data/                # Shared data directory
├── docs/                # Documentation (this folder)
├── compose.yml          # Docker Compose orchestration
├── Makefile             # Development commands
└── .env.*               # Environment configurations
```

## Quick Start

```bash
# Start all services
make start

# Stop all services
make stop

# Run database migrations
make migrate-up
```

## Links

- [Frontend Architecture](./architecture-app.md)
- [Backend Architecture](./architecture-api.md)
- [API Contracts](./api-contracts.md)
- [Data Models](./data-models.md)
- [Development Guide](./development-guide.md)
