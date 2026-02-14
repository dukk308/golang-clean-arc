# Tech Stack

## Core
- **Go** 1.25.1+
- **Gin** — HTTP framework
- **GORM** — ORM (Postgres, MySQL, SQLite, SQL Server)
- **Uber FX** — Dependency injection
- **Cobra** — CLI

## Optional Components
- **Redis/Valkey** — Caching via `cache_comp`
- **RabbitMQ** — Messaging via `rabbitmq_comp`
- **OpenTelemetry** — Tracing via `otel_comp`

## Tools
- **atlas** — Schema diff and migration validation
- **goose** — Run SQL migrations
- **swag** — Swagger documentation generation
- **air** — Hot reload (optional)
- **dlv** — Debugger

## Common Commands

```bash
# Start server
go run main.go serve

# Start worker
go run main.go worker

# Print environment help
go run main.go outenv

# Run migrations
make migration-up

# Rollback one migration
make migrate-down-1

# Generate Swagger docs
make swag

# Create migration
make migration-create-<name>

# Generate migration from GORM
make migration-gen-<name>

# Check migration status
make migration-status

# Validate migrations
make migration-validate
```

## Configuration

Environment variables and flags are defined in `.env.example`. Key flags:
- `-gin-port` — HTTP server port (default 8080)
- `-db-dsn` — Database connection string
- `-redis-addrs` — Redis/Valkey addresses
- `-rabbitmq-url` — RabbitMQ connection URL

## Development Tools

- **Debugger**: Connect to port 2345 (Air Debugger)
- **VS Code**: Use `vscode-go-air-reconnect` extension for debugger re-attach after Air rebuilds