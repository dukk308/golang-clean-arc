# Product

Beetool.Dev is a Go application built with Clean Architecture principles. It provides a REST API with authentication, blog, and note management features. The project demonstrates domain-driven design with clear separation of concerns across domain, application, infrastructure, and presentation layers.

## Key Features
- Clean Architecture with DDD-style modules (auth, note, user)
- REST API with Gin and Swagger documentation
- Role-based access control (Viewer, Editor, Admin)
- Optional integrations: Redis/Valkey caching, RabbitMQ messaging, OpenTelemetry observability

## Architecture Principles

1. **Layer Independence**: Domain and application layers must not import from presentation or infrastructure layers
2. **Module Isolation**: Each feature module (auth, blog, note) has its own domain models
3. **Dependency Injection**: All dependencies wired through Uber FX
4. **API Organization**: Routes follow prefix pattern (`/private`, `/admin`, `/public`, `/`)