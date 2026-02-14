# Project Structure

## Root Organization

```
beetool.dev-go-starter/
├── cmd/                    # CLI entry points
├── internal/               # Application code
│   ├── common/            # Shared utilities
│   ├── config/            # Configuration
│   ├── modules/           # Feature modules
│   └── server/            # Bootstrap
├── pkgs/                  # Reusable components
├── database/              # Migrations
├── deployment/            # Docker configs
├── api-docs/              # Swagger output
└── scripts/               # Utility scripts
```

## Module Structure (Clean Architecture)

Each feature module follows this pattern:

```
internal/modules/{module_name}/
├── application/           # Commands, queries, DTOs
├── domain/                # Entities, repository interfaces
├── infrastructure/        # Repo implementations
├── presentation/          # HTTP handlers, routes
└── fx_module.go           # FX dependency registration
```

## Layer Mapping

| Layer | Path Pattern | Purpose |
|-------|--------------|---------|
| Domain | `internal/modules/*/domain/` | Entities, value objects, repository interfaces |
| Application | `internal/modules/*/application/` | Use cases, commands, queries, DTOs |
| Infrastructure | `internal/modules/*/infrastructure/` | Persistence, external services |
| Presentation | `internal/modules/*/presentation/http/` | HTTP handlers, routing |

## Key Directories

- `pkgs/components/` — Reusable framework components (gin_comp, gorm_comp, cache_comp, etc.)
- `pkgs/middlewares/gin/` — HTTP middleware (auth, cors, logging, tracing)
- `pkgs/logger/` — Logging utilities
- `database/migrations/` — SQL migration files

## Bootstrap

`internal/server/boostrap.go` builds the FX application with:
- Global config and logger
- GORM database component
- Gin HTTP component
- Swagger component
- Feature modules
- HTTP server lifecycle

## Layer Independence Rules

- **Domain layer** (`internal/modules/**/domain/`):
  - Define new models here instead of importing from other modules
  - Do NOT import from presentation or infrastructure layers
  - Contains business logic and models

- **Application layer** (`internal/modules/**/application/`):
  - Define new models here instead of importing from other modules
  - Do NOT import from presentation or infrastructure layers
  - Command pattern: `command_{actor}_{action}` (e.g., `command_admin_delete_blog`)

- Application and domain stay independent of presentation and infrastructure.

## API Structure

When creating new APIs, follow `ai-docs/API-STRUCTURE.md`:
- Prefixes: `/private`, `/admin`, `/public`, `/` (authenticated)
- Path pattern: `prefix/{rest}` (e.g., `/v1/notes`, `/admin/v1/users`)