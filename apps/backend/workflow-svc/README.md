# Workflow Service

## Overview

The Workflow Service is a standalone microservice that manages state-driven workflows using YAML-based definitions. It supports event-driven transitions via Kafka, timeout handling, and provides a RESTful API for workflow operations.

## Architecture Decisions

- **Service Type**: Standalone microservice
- **Database**: MongoDB for workflow instances and event deduplication
- **HTTP Framework**: Go stdlib (net/http)
- **Kafka Integration**: Dedicated consumer with correlation_id lookup
- **Event Deduplication**: Separate MongoDB collection with event_id tracking
- **Error Handling**: DLQ for failed events + idempotency checks
- **Timeout Handling**: Background worker with time-based queries
- **Authentication**: JWT validation middleware
- **Observability**: OpenTelemetry (Prometheus metrics + Jaeger tracing)
- **Testing**: In-memory implementations for unit tests
- **Hot Reload**: Manual `/admin/reload` endpoint

## Features

### Core MVP Features
- Start workflow instances
- Process events (via Kafka consumer)
- Process actions (via REST API)
- Get instance details
- List instances with filtering
- Workflow definitions API
- Kafka transition publisher
- Timeout handling worker
- Event deduplication
- Dead Letter Queue (DLQ)

## Quick Start

### Prerequisites
- Go 1.21+
- MongoDB
- Kafka (with Zookeeper)
- Jaeger (optional, for tracing)

### Configuration

Copy `.env.example` to `.env` and update values:

```bash
cp .env.example .env
```

Key configuration:
- `MONGO_URI`: MongoDB connection string
- `KAFKA_BROKERS`: Comma-separated Kafka broker addresses
- `KAFKA_TOPIC`: Topic for incoming events (e.g., `vehicle.events`)
- `JWT_SECRET`: Secret key for JWT validation
- `WORKFLOW_DIR`: Directory containing YAML workflow definitions

### Run Locally

```bash
go run ./cmd/main.go
```

### Build

```bash
go build -o workflow-svc ./cmd/main.go
```

### Docker

```bash
docker build -t workflow-svc .
docker run -p 50003:50003 --env-file .env workflow-svc
```

## API Endpoints

### Workflow Definitions

#### List Workflows
```bash
GET /api/workflows
```

#### Get Workflow Definition
```bash
GET /api/workflows/:name
```

#### Reload Workflows (Admin)
```bash
POST /admin/reload
```

### Workflow Instances

#### Start Workflow
```bash
POST /api/workflows/:name/start
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "correlation_id": "vehicle-123",
  "context": {
    "vehicle_id": "vehicle-123",
    "user_id": "user-456"
  }
}
```

#### Get Instance
```bash
GET /api/instances/:id
Authorization: Bearer <jwt_token>
```

#### List Instances
```bash
GET /api/instances?workflow_name=vehicle_approval&state=pending_approval
Authorization: Bearer <jwt_token>
```

#### Process Action
```bash
POST /api/instances/:id/actions/:action
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "context": {
    "approved_by": "manager-789"
  }
}
```

## YAML Workflow Definition

Example workflow at `config/workflows/vehicle_approval.yaml`:

```yaml
name: vehicle_approval
version: "1.0"
description: Vehicle approval workflow

states:
  draft:
    name: draft
    type: initial
    timeout: 72h

  pending_approval:
    name: pending_approval
    type: intermediate
    timeout: 48h

  approved:
    name: approved
    type: intermediate

  rejected:
    name: rejected
    type: terminal

  active:
    name: active
    type: terminal

transitions:
  - from: draft
    to: pending_approval
    action: submit_for_approval

  - from: pending_approval
    to: approved
    action: approve

  - from: pending_approval
    to: rejected
    action: reject

  - from: approved
    to: active
    event: vehicle_activated
```

## Kafka Integration

### Incoming Events

The service consumes events from `KAFKA_TOPIC` with the following structure:

```json
{
  "event_id": "evt-123",
  "event_type": "vehicle_activated",
  "correlation_id": "vehicle-123",
  "payload": {
    "activated_at": "2026-01-31T10:00:00Z"
  },
  "timestamp": "2026-01-31T10:00:00Z"
}
```

- `event_id`: Unique event ID for deduplication
- `correlation_id`: Used to lookup workflow instance (e.g., vehicle_id)
- `event_type`: Triggers workflow transition

### Outgoing Events

Publishes state transitions to `workflow.transitions` topic:

```json
{
  "instance_id": "inst-456",
  "workflow_name": "vehicle_approval",
  "correlation_id": "vehicle-123",
  "from_state": "pending_approval",
  "to_state": "approved",
  "trigger_type": "action",
  "trigger_name": "approve",
  "context": {}
}
```

### Dead Letter Queue

Failed events are sent to `KAFKA_DLQ_TOPIC` with error headers for manual inspection.

## Timeout Handling

The timeout worker runs every `TIMEOUT_WORKER_INTERVAL` (default: 30s):
1. Queries instances with `timeout_at <= now()`
2. Processes timeout event for each instance
3. Transitions to configured timeout state

## Event Deduplication

Events with `event_id` are tracked in `processed_events` collection:
- Prevents duplicate processing
- Auto-expires after 7 days (TTL index)

## Authentication

All API endpoints (except `/health`, `/metrics`) require JWT:

```bash
Authorization: Bearer <jwt_token>
```

JWT claims expected:
```json
{
  "user_id": "user-123",
  "roles": ["admin", "manager"]
}
```

## Observability

### Metrics

Prometheus metrics available at `/metrics`:
- Workflow instance counts by state
- Transition event counts
- Consumer lag
- Repository operation durations

### Tracing

Distributed traces sent to Jaeger:
- HTTP request spans
- Kafka message processing spans
- MongoDB query spans

## Development

### Project Structure

```
workflow-svc/
├── cmd/
│   ├── main.go
│   └── config/
│       └── config.go
├── internal/
│   ├── domain/
│   │   ├── workflow/        # Domain models
│   │   └── repository/      # Repository interfaces
│   ├── application/
│   │   ├── loader/          # YAML loader
│   │   ├── registry/        # Workflow registry
│   │   └── service/         # Core service
│   ├── infrastructure/
│   │   ├── persistence/     # MongoDB implementations
│   │   ├── messaging/       # Kafka consumer/publisher
│   │   ├── worker/          # Timeout worker
│   │   ├── middleware/      # JWT middleware
│   │   └── observability/   # OpenTelemetry
│   └── api/
│       ├── handler/         # HTTP handlers
│       └── route/           # Router setup
├── config/
│   └── workflows/           # YAML workflow definitions
├── Dockerfile
├── .env.example
└── README.md
```

### Adding New Workflows

1. Create YAML file in `config/workflows/`
2. Restart service OR call `/admin/reload`
3. Workflow automatically registered

### Testing

Use in-memory implementations for unit tests:

```go
// In-memory repository for testing
type InMemoryInstanceRepository struct {
    instances map[string]*workflow.WorkflowInstance
}

// In-memory event deduplicator
type InMemoryDeduplicator struct {
    processedEvents map[string]bool
}
```

## Production Checklist

- [ ] Change `JWT_SECRET` to strong random value
- [ ] Configure MongoDB replica set for production
- [ ] Set up Kafka cluster with replication
- [ ] Configure proper Jaeger endpoints
- [ ] Add health check endpoint
- [ ] Configure resource limits in Kubernetes
- [ ] Set up monitoring alerts
- [ ] Configure backup for MongoDB
- [ ] Review timeout intervals for production load
- [ ] Enable TLS for MongoDB and Kafka connections

## License

MIT
