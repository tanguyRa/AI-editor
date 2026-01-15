# AI-Editor Documentation Index

## Project Overview

- **Type:** Multi-part (client/server)
- **Primary Language:** TypeScript (frontend), Go (backend)
- **Architecture:** Nuxt 4 + Go REST API

## Quick Reference

### Frontend (app/)

- **Type:** Web Application
- **Tech Stack:** Nuxt 4.2, Vue 3, TypeScript, Nuxt UI
- **Root:** `./app/`
- **Port:** 3000

### Backend (api/)

- **Type:** REST API
- **Tech Stack:** Go 1.24, net/http, PostgreSQL, sqlc
- **Root:** `./api/`
- **Port:** 8080

## Generated Documentation

### Architecture

- [Project Overview](./project-overview.md) - Executive summary and architecture diagram
- [Frontend Architecture](./architecture-app.md) - Nuxt/Vue application structure
- [Backend Architecture](./architecture-api.md) - Go API structure and patterns
- [Integration Architecture](./integration-architecture.md) - Service communication and data flow

### Technical Reference

- [API Contracts](./api-contracts.md) - REST API endpoints and schemas
- [Data Models](./data-models.md) - Database schema and entity relationships
- [Source Tree Analysis](./source-tree-analysis.md) - Complete directory structure

### Development

- [Development Guide](./development-guide.md) - Setup, commands, and workflow

## Existing Documentation

### Frontend

- [app/README.md](../app/README.md) - Nuxt Starter Template docs

### Backend

- [api/README.md](../api/README.md) - Backend overview and structure
- [api/internal/repository/README.md](../api/internal/repository/README.md) - Repository layer (sqlc)
- [api/internal/config/README.md](../api/internal/config/README.md) - Configuration
- [api/internal/server/README.md](../api/internal/server/README.md) - HTTP server
- [api/pkg/database/README.md](../api/pkg/database/README.md) - Database adapter

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Make

### Quick Start

```bash
# 1. Clone and configure
cp .env.example .env.dev
# Edit .env.dev with your settings

# 2. Start services
make start

# 3. Run migrations
make migrate-up

# 4. Access the app
open http://localhost:3000
```

### Common Commands

| Command | Description |
|---------|-------------|
| `make start` | Start all services |
| `make stop` | Stop all services |
| `make migrate-up` | Apply database migrations |
| `make sqlc` | Regenerate SQL code |
| `make front` | Shell into frontend container |

## Project Structure

```
AI-editor/
├── app/                 # Nuxt frontend (port 3000)
│   ├── app/             # Application source
│   │   ├── pages/       # Routes (index, login, register, dashboard)
│   │   ├── components/  # Vue components
│   │   ├── lib/         # Auth utilities
│   │   └── middleware/  # Route guards
│   └── server/api/      # Server routes (auth handler)
├── api/                 # Go backend (port 8080)
│   ├── cmd/server/      # Entry point
│   ├── internal/        # Private packages
│   │   ├── server/      # HTTP handlers
│   │   ├── repository/  # Database (sqlc)
│   │   └── config/      # Configuration
│   ├── pkg/             # Public packages
│   └── db/              # Migrations & queries
├── docs/                # This documentation
├── compose.yml          # Docker orchestration
└── Makefile             # Development commands
```

## Technology Stack

| Layer | Technology |
|-------|------------|
| Frontend | Nuxt 4.2, Vue 3, TypeScript, Nuxt UI |
| Backend | Go 1.24, net/http |
| Database | PostgreSQL, pgx, sqlc |
| Auth | better-auth (client + compatible API) |
| Infrastructure | Docker, Docker Compose |
| Dev Tools | Bun, Air, ESLint |

## Key Features

- **Authentication**: Email/password with session cookies
- **Type Safety**: TypeScript (frontend), sqlc (backend)
- **Hot Reload**: Bun dev server, Air for Go
- **Containerized**: Docker Compose for all services

---

*Documentation generated on 2026-01-14 by BMAD Document Project workflow*
