# Development Guide

## Prerequisites

### Required Software

| Software | Version | Purpose |
|----------|---------|---------|
| Docker | Latest | Container runtime |
| Docker Compose | v2+ | Service orchestration |
| Make | Any | Build automation |

### Optional (for local development without Docker)

| Software | Version | Purpose |
|----------|---------|---------|
| Bun | Latest | Frontend package manager & runtime |
| Go | 1.24+ | Backend language |
| PostgreSQL | 14+ | Database |

## Quick Start

### 1. Clone and Setup

```bash
# Clone repository
git clone <repository-url>
cd AI-editor

# Copy environment file
cp .env.example .env.dev
```

### 2. Configure Environment

Edit `.env.dev`:

```bash
# Project identification
PROJECT_NAME=budhapp

# Docker build target
DOCKER_TARGET=development

# Authentication
BETTER_AUTH_SECRET=<generate-with-openssl-rand-base64-32>
BETTER_AUTH_URL=http://localhost:3000

# Database
DATABASE_URL=postgresql://<user>:<password>@<host>:<port>/<database>
```

### 3. Start Services

```bash
# Build and start all services
make start

# Or separately
make build
docker compose --env-file .env.dev up
```

### 4. Run Migrations

```bash
# Apply database migrations
make migrate-up
```

### 5. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

## Makefile Commands

### Service Management

| Command | Description |
|---------|-------------|
| `make build` | Build Docker images |
| `make start` | Build and start all services |
| `make stop` | Stop all services |

### Database

| Command | Description |
|---------|-------------|
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-drop` | Drop all tables (dangerous!) |
| `make migration name="name"` | Create new migration file |
| `make seed` | Run database seed |

### Code Generation

| Command | Description |
|---------|-------------|
| `make sqlc` | Regenerate sqlc queries |
| `make generate` | Alias for sqlc |

### Utilities

| Command | Description |
|---------|-------------|
| `make front` | Shell into frontend container |
| `make back` | Shell into backend container |
| `make db-url` | Print database URL |

## Development Workflow

### Frontend Development

The frontend runs with hot module replacement:

```bash
# Start services
make start

# Frontend auto-reloads on file changes
# Access at http://localhost:3000
```

**Key files to modify:**
- Pages: `app/app/pages/*.vue`
- Components: `app/app/components/*.vue`
- Styles: `app/app/assets/css/main.css`
- Config: `app/nuxt.config.ts`

### Backend Development

The backend uses Air for hot reload:

```bash
# Start services
make start

# Backend auto-reloads on file changes
# API at http://localhost:8080
```

**Key files to modify:**
- Handlers: `api/internal/server/main.go`
- Config: `api/internal/config/config.go`
- SQL Queries: `api/db/queries/*.sql`

### Database Changes

1. Create a new migration:
   ```bash
   make migration name="add_new_table"
   ```

2. Edit the migration files in `api/db/migrations/`

3. Apply the migration:
   ```bash
   make migrate-up
   ```

4. Update SQL queries in `api/db/queries/`

5. Regenerate Go code:
   ```bash
   make sqlc
   ```

## Project Structure Best Practices

### Frontend (Nuxt/Vue)

- **Pages**: Use file-based routing in `app/pages/`
- **Components**: Create reusable components in `app/components/`
- **Utilities**: Add shared code to `app/lib/`
- **Middleware**: Add route guards in `app/middleware/`
- **Server Routes**: Add API routes in `server/api/`

### Backend (Go)

- **Public API**: Add to `pkg/` for reusable packages
- **Private Code**: Add to `internal/` for app-specific code
- **Entry Points**: Add executables in `cmd/`
- **SQL Queries**: Define in `db/queries/`, generate with sqlc
- **Migrations**: Add to `db/migrations/`

## Environment Variables

### Required

| Variable | Description |
|----------|-------------|
| `PROJECT_NAME` | Docker project name prefix |
| `DATABASE_URL` | PostgreSQL connection string |
| `BETTER_AUTH_SECRET` | Auth encryption secret |
| `BETTER_AUTH_URL` | Base URL for auth |

### Optional

| Variable | Default | Description |
|----------|---------|-------------|
| `DOCKER_TARGET` | `development` | Docker build stage |
| `DOCKER_REGISTRY` | `registry.budhapp.com` | Image registry |
| `TAG` | `latest` | Image tag |
| `ENVIRONMENT` | `dev` | Runtime environment |
| `ENCRYPTION_KEY` | - | 256-bit encryption key |

## Debugging

### Frontend

```bash
# Access container shell
make front

# View logs
docker logs budhapp-app -f
```

### Backend

```bash
# View logs
docker logs budhapp-api -f

# Database connection
make db-url
```

### Database

```bash
# Connect to database
docker exec -it budhapp-api psql "$DATABASE_URL"
```

## Testing

### Frontend

```bash
# In container or locally
bun run lint
bun run typecheck
```

### Backend

```bash
# In container
go test ./...
```

## Production Deployment

### Build Production Images

```bash
# Set production target
export DOCKER_TARGET=production

# Build images
make build
```

### Environment Configuration

Use `.env.prod` for production settings:

```bash
DOCKER_TARGET=production
BETTER_AUTH_URL=https://your-domain.com
DATABASE_URL=postgresql://...
```

### Deploy

```bash
# Use production compose
docker compose --env-file .env.prod up -d
```
