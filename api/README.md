# Backend Documentation

## Overview

Go backend built on PocketBase with **schema-driven code generation**. Define your data models in YAML, generate the entire stack automatically.

## Schema-Driven Development

Define your data models in YAML schemas, then generate Go models, PocketBase collections, and JavaScript services automatically.

**Workflow:**
1. Create schema in `cmd/codegen/schemas/`
2. Code generation creates models, collections, and services
3. Use generated code in handlers and components

**What gets generated:**
- `internal/models/*_gen.go` - Go models with full CRUD
- `internal/server/collections/*_gen.go` - PocketBase collection setup
- `web/static/services/db/*_db_gen.js` - Frontend DB services

See example schemas in `cmd/codegen/schemas/` for reference.

## Documentation

For detailed schema documentation and examples, see the schema files in `cmd/codegen/schemas/`.

## Project Structure
```tree
app/
├── README.md
├── main.go
└── internal/
    ├── config/
    │   ├── README.md
    │   └── config.go
    ├── handlers/
    │   ├── README.md
    │   ├── auth.go
    │   ├── handlers.go
    │   └── response.go
    ├── middleware/
    │   ├── README.md
    │   ├── chain.go
    │   ├── chain_test.go
    │   └── cors.go
    ├── repository/
    │   ├── README.md
    │   ├── accounts.sql.go
    │   ├── db.go
    │   ├── jwks.sql.go
    │   ├── models.go
    │   ├── sessions.sql.go
    │   └── users.sql.go
    ├── server/
    │   ├── README.md
    │   ├── middleware.go
    │   ├── routes.go
    │   ├── server.go
    │   └── utils.go
    ├── session/
    │   ├── README.md
    │   └── context.go
    └── utils/
        ├── README.md
        └── date_parser.go
```
