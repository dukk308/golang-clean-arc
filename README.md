# Golang Clean Architecture

A Go application built with Clean Architecture principles, featuring dependency injection using Uber FX, Gin web framework, and GORM for database operations.

## Overview

This project demonstrates a clean, maintainable architecture for Go applications, following Domain-Driven Design (DDD) principles and separation of concerns. The architecture is organized into distinct layers: Domain, Application, Infrastructure, and Presentation.

## Project Structure

![Project Structure](doc/project-structure.png)

```
golang-clean-arc/
├── cmd/                    # Application entry points
│   ├── root.go            # Root command
│   ├── server.go          # Server command
│   └── worker.go          # Worker command
├── deployment/             # Docker and deployment configurations
│   ├── development/
│   ├── local/
│   └── production/
├── doc/                    # Documentation and diagrams
│   ├── clean-architecture.png
│   └── project-structure.png
├── internal/               # Internal application code
│   ├── common/            # Shared utilities
│   ├── config/            # Configuration management
│   ├── modules/           # Feature modules
│   │   └── user/          # User module
│   │       ├── application/        # Application layer (use cases)
│   │       ├── domain/             # Domain layer (entities, repositories)
│   │       ├── infrastructure/     # Infrastructure layer (implementations)
│   │       └── presentation/       # Presentation layer (HTTP handlers)
│   └── server/            # Server bootstrap
├── pkgs/                   # Shared packages
│   ├── components/        # Reusable components
│   │   ├── gin_comp/      # Gin HTTP component
│   │   └── gorm_comp/     # GORM database component
│   ├── ddd/               # DDD utilities
│   └── logger/            # Logging utilities
├── go.mod
├── go.sum
├── main.go
└── Makefile
```

## Architecture

![Clean Architecture](doc/clean-architecture.png)

The project follows Clean Architecture principles with the following layers:

### Domain Layer

- **Entities**: Core business models (User, Viewer, Editor, Admin)
- **Value Objects**: Email, Role
- **Repository Interfaces**: Define data access contracts
- **Domain Errors**: Business logic errors

### Application Layer

- **Commands**: Use cases and business logic orchestration
- **DTOs**: Data transfer objects for use cases

### Infrastructure Layer

- **Persistence**: Database implementations (GORM)
- **External Services**: Third-party integrations

### Presentation Layer

- **HTTP Handlers**: REST API endpoints
- **Route Registration**: API routing setup

## Features

- ✅ Clean Architecture implementation
- ✅ Dependency Injection with Uber FX
- ✅ RESTful API with Gin
- ✅ Database operations with GORM
- ✅ User management with role-based access (Viewer, Editor, Admin)
- ✅ Domain-Driven Design patterns
- ✅ Modular structure for scalability

## Prerequisites

- Go 1.24.4 or higher
- MySQL 8.0 or higher
- Docker (optional, for containerized deployment)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/dukk308/golang-clean-arch-starter.git
cd golang-clean-arc
```

2. Install dependencies:

```bash
go mod download
```

3. Configure the database:

   - Update database settings in `internal/config/config.go` or use environment variables
   - Create the database schema

4. Run the application:

```bash
go run main.go serve
```

## Usage

### Start Server

```bash
go run main.go serve
```

The server will start on port `8080` by default.

## Development

### Project Structure Guidelines

1. **Domain Layer** (`internal/modules/*/domain/`)

   - Contains pure business logic
   - No dependencies on external frameworks
   - Defines repository interfaces

2. **Application Layer** (`internal/modules/*/application/`)

   - Implements use cases
   - Orchestrates domain logic
   - Depends only on domain layer

3. **Infrastructure Layer** (`internal/modules/*/infrastructure/`)

   - Implements repository interfaces
   - Handles external dependencies (database, APIs)
   - Depends on domain and application layers

4. **Presentation Layer** (`internal/modules/*/presentation/`)
   - HTTP handlers and routing
   - Request/response transformation
   - Depends on application layer

### Adding New Features

1. Create a new module in `internal/modules/`
2. Follow the layer structure (domain, application, infrastructure, presentation)
3. Register the module in `internal/modules/fx_features.go`
4. Add routes in the presentation layer

## Configuration

Configuration is managed in `internal/config/config.go`. You can customize:

- Server port
- Database connection settings
- Service name and environment

## Docker

### Development

```bash
cd deployment/development
docker-compose up
```

### Production

```bash
cd deployment/production
docker build -t golang-clean-arc .
docker run -p 8080:8080 golang-clean-arc
```

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [Uber FX](https://github.com/uber-go/fx) - Dependency injection
- [Cobra](https://github.com/spf13/cobra) - CLI framework

## License

This project is open source and available under the MIT License.
