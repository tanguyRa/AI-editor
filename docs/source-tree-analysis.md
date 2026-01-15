# Source Tree Analysis

## Complete Project Structure

```
AI-editor/
├── .claude/                      # Claude Code configuration
│   └── commands/                 # BMAD workflow commands
├── .env.dev                      # Development environment
├── .env.example                  # Environment template
├── .env.prod                     # Production environment
├── .git/                         # Git repository
├── .github/                      # GitHub configuration
├── .gitignore                    # Git ignore rules
├── Makefile                      # Development commands
├── _bmad/                        # BMAD framework files
├── _bmad-output/                 # BMAD output artifacts
├── api/                          # 🔧 Go Backend
│   ├── .air.toml                 # Air hot reload config
│   ├── Dockerfile                # Container build
│   ├── README.md                 # Backend documentation
│   ├── cmd/                      # Entry points
│   │   ├── docs/                 # Documentation generator
│   │   │   ├── discovery.go      # File discovery logic
│   │   │   ├── generators.go     # Doc generators
│   │   │   ├── main.go           # Entry point
│   │   │   ├── parser.go         # Go code parser
│   │   │   ├── types.go          # Type definitions
│   │   │   └── utils.go          # Utilities
│   │   └── server/               # Server entry
│   │       └── main.go           # ⚡ Server bootstrap
│   ├── db/                       # Database layer
│   │   ├── migrations/           # Schema migrations
│   │   │   ├── 000001_init_auth.down.sql
│   │   │   └── 000001_init_auth.up.sql
│   │   └── queries/              # SQL queries (sqlc input)
│   │       ├── sessions.sql      # Session CRUD
│   │       └── users.sql         # User CRUD
│   ├── go.mod                    # Go module definition
│   ├── go.sum                    # Dependency checksums
│   ├── internal/                 # Private packages
│   │   ├── config/               # Configuration
│   │   │   ├── README.md
│   │   │   └── config.go         # Config loading
│   │   ├── handlers/             # HTTP handlers
│   │   │   ├── README.md
│   │   │   └── utils.go          # Handler utilities
│   │   ├── repository/           # Data access (sqlc generated)
│   │   │   ├── README.md
│   │   │   ├── db.go             # DBTX interface
│   │   │   ├── models.go         # Generated models
│   │   │   ├── sessions.sql.go   # Generated session queries
│   │   │   └── users.sql.go      # Generated user queries
│   │   ├── server/               # HTTP server
│   │   │   ├── README.md
│   │   │   └── main.go           # ⚡ Server + Auth handlers
│   │   └── utils/                # Utilities
│   │       ├── README.md
│   │       └── date_parser.go    # Date parsing
│   ├── main.go                   # Root entry (placeholder)
│   ├── pkg/                      # Public packages
│   │   ├── auth/                 # Auth package (WIP)
│   │   │   ├── README.md
│   │   │   └── auth.go
│   │   ├── database/             # Database abstraction
│   │   │   ├── README.md
│   │   │   ├── adapter.go        # Adapter interface
│   │   │   └── main.go
│   │   ├── oauth/                # OAuth providers (WIP)
│   │   │   ├── README.md
│   │   │   └── provider.go
│   │   └── session/              # Session management (WIP)
│   │       ├── README.md
│   │       └── session.go
│   ├── sqlc.yml                  # sqlc configuration
│   └── tmp/                      # Air build artifacts
├── app/                          # 🎨 Nuxt Frontend
│   ├── .editorconfig             # Editor settings
│   ├── .github/                  # GitHub workflows
│   ├── .npmrc                    # npm configuration
│   ├── .nuxt/                    # Nuxt build cache
│   ├── Dockerfile                # Container build
│   ├── LICENSE                   # License file
│   ├── README.md                 # Frontend documentation
│   ├── app/                      # Application source
│   │   ├── app.config.ts         # App configuration
│   │   ├── app.vue               # ⚡ Root component
│   │   ├── assets/               # Static assets
│   │   │   └── css/
│   │   │       └── main.css      # Global styles
│   │   ├── components/           # Vue components
│   │   │   ├── AppLogo.vue       # Logo component
│   │   │   ├── TemplateComponent.vue
│   │   │   └── TemplateMenu.vue  # Navigation menu
│   │   ├── lib/                  # Shared utilities
│   │   │   ├── auth-client.ts    # 🔐 Auth client exports
│   │   │   └── auth.ts           # 🔐 Server-side auth
│   │   ├── middleware/           # Route middleware
│   │   │   └── auth.ts           # 🔐 Auth middleware
│   │   └── pages/                # File-based routing
│   │       ├── (template)/       # Template route group
│   │       │   └── test.vue      # Test page
│   │       ├── dashboard.vue     # 🔐 Protected dashboard
│   │       ├── index.vue         # Landing page
│   │       ├── login.vue         # 🔐 Login page
│   │       └── register.vue      # 🔐 Registration page
│   ├── bun.lock                  # Bun lockfile
│   ├── eslint.config.mjs         # ESLint config
│   ├── node_modules/             # Dependencies
│   ├── nuxt.config.ts            # Nuxt configuration
│   ├── package.json              # Dependencies
│   ├── pnpm-lock.yaml            # pnpm lockfile (unused)
│   ├── pnpm-workspace.yaml       # pnpm workspace (unused)
│   ├── public/                   # Public static files
│   ├── renovate.json             # Renovate bot config
│   ├── server/                   # Nuxt server
│   │   └── api/                  # Server API routes
│   │       └── auth/             # Auth API
│   │           ├── [...all].ts   # ⚡ Auth catch-all handler
│   │           └── health.ts     # Health check
│   └── tsconfig.json             # TypeScript config
├── compose.yml                   # 🐳 Docker Compose
├── data/                         # Shared data directory
│   └── shared/                   # Mounted in containers
└── docs/                         # 📚 Documentation (this folder)
    ├── project-scan-report.json  # Scan state
    ├── project-overview.md       # Project overview
    ├── architecture-app.md       # Frontend architecture
    ├── architecture-api.md       # Backend architecture
    ├── api-contracts.md          # API documentation
    ├── data-models.md            # Database schema
    ├── source-tree-analysis.md   # This file
    ├── development-guide.md      # Dev guide
    ├── integration-architecture.md # Integration docs
    └── index.md                  # Master index
```

## Legend

- ⚡ = Entry point / Main file
- 🔐 = Authentication related
- 🔧 = Backend (Go)
- 🎨 = Frontend (Nuxt/Vue)
- 🐳 = Docker/Infrastructure
- 📚 = Documentation

## Critical Directories

### Frontend (app/)

| Directory | Purpose |
|-----------|---------|
| `app/app/` | Main application source |
| `app/app/pages/` | File-based routing (Nuxt pages) |
| `app/app/components/` | Reusable Vue components |
| `app/app/lib/` | Shared utilities (auth, etc.) |
| `app/app/middleware/` | Route middleware |
| `app/server/api/` | Server-side API routes |

### Backend (api/)

| Directory | Purpose |
|-----------|---------|
| `api/cmd/server/` | Main server entry point |
| `api/internal/server/` | HTTP server and handlers |
| `api/internal/repository/` | Database queries (sqlc generated) |
| `api/internal/config/` | Configuration management |
| `api/pkg/database/` | Database adapter interface |
| `api/db/migrations/` | Database migrations |
| `api/db/queries/` | SQL query definitions |

## Entry Points

### Frontend
- **Development**: `bun run dev` → Nuxt dev server
- **Production**: `bun --bun run /app/server/index.mjs`
- **Root Component**: `app/app/app.vue`

### Backend
- **Development**: `air` → Hot reload via `cmd/server/main.go`
- **Production**: `/server serve --http=0.0.0.0:8080`
- **Main Handler**: `api/internal/server/main.go`

## File Statistics

| Part | Go Files | Vue/TS Files | SQL Files | Config Files |
|------|----------|--------------|-----------|--------------|
| api/ | 21 | 0 | 4 | 3 |
| app/ | 0 | 16 | 0 | 6 |
| root | 0 | 0 | 0 | 4 |
