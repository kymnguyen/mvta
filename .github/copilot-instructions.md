# AI Copilot Instructions for MVTA

## Project Overview

**MVTA** is a Go-based microservices monorepo for vehicle tracking and authentication systems. The project uses Go 1.25.6 with Go workspaces to manage multiple independent services.

## Architecture & Structure

### Workspace Organization
- **`/services`** - Go workspace root containing all microservices
  - Uses `go.work` for local development workspace management (enables building/testing multiple modules without publishing)
  - Each service is an independent Go module with its own `go.mod`
- **`/services/auth`** - Authentication service (scaffolding in progress)
- **`/services/tracking`** - Tracking service (scaffolding in progress)
- **`/services/vehicle`** - Core vehicle service
  - `cmd/` - CLI/main entry points
  - `domain/` - Business logic and domain models
- **`/infra`** - Infrastructure configuration (currently empty; reserved for deployment scripts, IaC, etc.)

### Infrastructure
- **MongoDB 7.0** for primary data storage (configured in `docker-compose.yaml`)
- Local development uses Docker Compose to spin up MongoDB with authentication
- Network: `qaas-workout` for inter-service communication

## Critical Developer Workflows

### Initial Setup
```sh
cd services
go work init
```
This creates the workspace file that links all services for local development.

### Building & Testing Services
Each service should have its own build/test commands. Until scaffolding is complete:
- Build a service: `cd services/<service-name> && go build ./cmd/...`
- Test a service: `cd services/<service-name> && go test ./...`

### Local Development Database
```sh
docker-compose up -d mongo
```
Starts MongoDB with default credentials (user: `mongo`, password: `password`). Environment variables can override these via `.env` file.

## Project Conventions

### Code Organization Pattern
- **`cmd/`** - Application entry points and CLI handlers; minimal business logic
- **`domain/`** - Core business logic, domain models, and interfaces
- Keep services loosely coupled; use well-defined APIs for inter-service communication

### Go Module Management
- Each service must have its own `go.mod` in its root directory
- Use Go workspaces for local development; do NOT commit modules vendored into the workspace
- When adding dependencies, add them to the specific service's `go.mod`

### Naming Conventions
- Services follow lowercase naming with hyphens if needed: `auth-service`, `vehicle-service`
- Packages use lowercase single words when possible (Go convention)

## Integration Points

### Service Communication
- Services should communicate via well-defined boundaries (HTTP APIs, gRPC, or event systems)
- Document service endpoints/contracts explicitly when adding new services

### Data Layer
- MongoDB is the primary data store; coordinate schema design across services via version control
- Use appropriate database drivers for Go (e.g., `mongo-go-client` for MongoDB)

## Important Files Reference

- **`services/go.work`** - Go workspace definition; edit when adding/removing services
- **`docker-compose.yaml`** - Local dev environment; update when adding new infrastructure services
- **`README.md`** - High-level project documentation (minimal; being expanded)

## Architectural Layers (DDD + Clean Architecture)

### 1. Domain Layer (`domain/`)
- **Value Objects** (`valueobject/`): `VehicleID`, `VehicleStatus`, `Location`, `Version`
  - Immutable, self-validating, business logic encapsulated
- **Entities** (`entity/`): `Vehicle` aggregate root
  - Generates and collects uncomitted domain events
  - Implements optimistic concurrency with `Version` value object
- **Events** (`event/`): `VehicleCreatedEvent`, `VehicleLocationUpdatedEvent`, `VehicleStatusChangedEvent`
  - Domain events for asynchronous inter-service communication
- **Repositories** (`repository/`): Abstract interfaces (no implementation)
  - `VehicleRepository`, `OutboxRepository`, `UnitOfWork`

### 2. Application Layer (`application/`)
- **DTOs** (`dto/`): Request/response transfer objects for API contracts
- **Commands** (`command/`): Intent-based operations (`CreateVehicleCommand`, `UpdateVehicleLocationCommand`)
  - Command handlers implement use cases, dispatch through command bus
- **Queries** (`query/`): Read operations (`GetVehicleQuery`, `GetAllVehiclesQuery`)
  - Query handlers fetch data, return query results
- **Services** (`service/`): Use case orchestration
  - Command handlers: domain logic + event publishing to outbox
  - Query handlers: read-only data retrieval

### 3. Infrastructure Layer (`infrastructure/`)
- **Persistence** (`persistence/`): MongoDB implementations
  - `MongoVehicleRepository`: Aggregate persistence with optimistic concurrency
  - `MongoOutboxRepository`: Transactional Outbox pattern for reliable async events
- **Messaging** (`messaging/`): In-memory command/query buses
  - `InMemoryCommandBus`, `InMemoryQueryBus` - can be swapped for distributed systems
- **Resilience** (`resilience/`): 
  - `CircuitBreaker`: Failure tolerance (Closed → Open → Half-Open states)
  - `RetryPolicy`: Exponential backoff with configurable attempts
  - `IdempotencyStore`: Idempotent operation tracking
- **Worker** (`worker/`): Background processes
  - `OutboxWorker`: Asynchronously publishes outbox events (polling-based)

### 4. API Layer (`api/`)
- **Handlers** (`handler/`): HTTP request/response mapping
  - Delegates to command/query bus, returns standardized responses
- **Routes** (`route/`): RESTful endpoint registration
  - `POST /api/v1/vehicles` - Create
  - `GET /api/v1/vehicles[/{id}]` - Retrieve
  - `PATCH /api/v1/vehicles/{id}/location` - Update location
  - `PATCH /api/v1/vehicles/{id}/status` - Change status

### 5. Entry Point (`cmd/main.go`)
- Bootstraps dependencies via DI container
- Starts MongoDB connection, command/query buses, outbox worker
- Registers HTTP routes and server
- Handles graceful shutdown

## Key Architectural Patterns Implemented

### Transactional Outbox Pattern
- Domain events saved atomically with aggregate updates
- `OutboxWorker` polls and publishes events asynchronously
- Ensures exactly-once event delivery even if publisher fails
- Placeholder: `noOpEventPublisher` - replace with Kafka/RabbitMQ/EventGrid

### Optimistic Concurrency Control
- Vehicle aggregate versioned with `Version` value object
- Repository save checks `version - 1` before update (CAS semantics)
- Prevents lost updates in concurrent scenarios

### Circuit Breaker & Retry Policies
- Available in `resilience/` for downstream service calls
- Fail-fast + backoff mechanisms prevent cascading failures

### Asynchronous Communication (Event-Driven)
- Domain events emitted on aggregate state changes
- Published via outbox worker to external systems
- Service decoupling: Vehicle service doesn't wait for downstream handlers

## Development Workflow

### Adding a New Command
1. Define command struct in `application/command/`
2. Implement command handler in `application/service/`
3. Register handler in DI container (`infrastructure/di/container.go`)
4. Add HTTP endpoint in `api/handler/` and `api/route/`

### Adding a New Domain Event
1. Define event struct extending `BaseDomainEvent` in `domain/event/`
2. Emit from aggregate root (`entity/`) on state change
3. Update outbox worker topic mapping in `infrastructure/worker/outbox_worker.go`
4. Implement external event handler (in tracking/auth services)

### Adding a New Value Object
1. Define immutable struct in `domain/valueobject/`
2. Implement validation in constructor
3. Use in aggregates and commands

## DDD / Clean Architecture Principles

- **Domain-Driven**: Business logic in aggregates/entities, not services
- **Dependency Inversion**: Repositories are interfaces, implementations in infrastructure
- **Separation of Concerns**: Each layer has single responsibility
- **No Framework Lock-in**: Core domain is Go-only, no external dependencies
- **Testability**: All logic can be unit tested without mocks (pure domain logic)

## Production Considerations

- Replace `noOpEventPublisher` with real message broker integration (Kafka, RabbitMQ, Azure EventGrid)
- Implement idempotency store (Redis, MongoDB) to prevent duplicate processing
- Add gRPC for inter-service communication (proto files in `services/vehicle/pkg/proto/`)
- Implement comprehensive logging/metrics (already using `zap`, add `prometheus`)
- Add database migrations (`goose` configured in `example.env`)

## Scaffolding Status

✅ **Vehicle Service**: Fully implemented DDD skeleton with all layers, ready for business logic expansion
- Remaining: Auth & Tracking services follow same pattern
- Use `go work use ./services/auth ./services/tracking ./services/vehicle` to update workspace

## Quick References

- **Language**: Go 1.25.6
- **Repository**: https://github.com/kymnguyen/mvta
- **Branch**: main (default)
- **Database**: MongoDB 7.0 on localhost:27017
- **API Port**: 8080 (configurable via PORT env var)
- **Outbox Poll Interval**: 5 seconds (configurable in `cmd/main.go`)
