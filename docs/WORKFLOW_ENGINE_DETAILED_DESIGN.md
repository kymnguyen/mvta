# Detailed Design: Core Workflow Engine for MVTA

**Document Version:** 1.0  
**Date:** January 31, 2025  
**Status:** Detailed Design

---

## 1. Introduction

This document provides the detailed technical design for the core workflow engine proposed for the MVTA vehicle management system. It specifies component design, data structures, sequence flows, error handling, concurrency, and operational concerns.

---

## 2. System Context

### 2.1 Context Diagram

```
                    ┌─────────────────────────────────────────────┐
                    │              EXTERNAL ACTORS                │
                    │  ┌─────────┐  ┌─────────┐  ┌─────────────┐  │
                    │  │ Admin   │  │ vehicle │  │ Kafka       │  │
                    │  │ UI /    │  │ -svc    │  │ Topics      │  │
                    │  │ Client  │  │         │  │             │  │
                    │  └────┬────┘  └────┬────┘  └──────┬──────┘  │
                    └───────┼────────────┼──────────────┼─────────┘
                            │ REST       │ Event        │ Event
                            │            │ (trigger)    │ (consume)
                            ▼            ▼              ▼
                    ┌─────────────────────────────────────────────┐
                    │           WORKFLOW ENGINE                   │
                    │  ┌──────────────────────────────────────┐   │
                    │  │  API  │  Engine  │  Repository       │   │
                    │  └──────────────────────────────────────┘   │
                    └─────────────────────┬───────────────────────┘
                                          │
                                          ▼
                    ┌─────────────────────────────────────────────┐
                    │  MongoDB  │  YAML Files  │  Kafka (publish) │
                    └─────────────────────────────────────────────┘
```

### 2.2 Deployment Options

| Option | Description |
|--------|-------------|
| **A. Embedded in vehicle-svc** | Engine as internal package; same process, DB, Kafka |
| **B. Dedicated workflow-svc** | Standalone service; owns workflow collections; subscribes to Kafka |

This design supports both; interfaces are implementation-agnostic.

---

## 3. Component Design

### 3.1 Component Diagram

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              API LAYER                                       │
│  ┌─────────────────┐  ┌─────────────────┐   ┌─────────────────────────────┐  │
│  │ WorkflowHandler │  │ InstanceHandler │   │ ActionHandler               │  │
│  │ - StartWorkflow │  │ - GetInstance   │   │ - ExecuteAction             │  │
│  │                 │  │ - ListInstances │   │                             │  │
│  └────────┬────────┘  └────────┬────────┘   └──────────────┬──────────────┘  │
└───────────┼────────────────────┼───────────────────────────┼─────────────────┘
            │                    │                           │
            ▼                    ▼                           ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         APPLICATION LAYER                                   │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    WorkflowService (implements Engine)                │  │
│  │  - Start()  - ProcessEvent()  - ProcessAction()  - GetInstance()      │  │
│  └───────┬───────────────────────────────────────────────────────────────┘  │
│          │ uses                                                             │
│          ▼                                                                  │
│  ┌───────────────────┐  ┌───────────────────┐   ┌─────────────────────────┐ │
│  │ DefinitionRegistry│  │TransitionValidator│   │ ActionHandler (optional)│ │
│  │ - Get(name)       │  │- CanTransition()  │   │ - OnTransition()        │ │
│  └─────────┬─────────┘  └───────────────────┘   └────────────┬────────────┘ │
└────────────┼─────────────────────────────────────────────────┼──────────────┘
             │                                                 │
             ▼                                                 │
┌─────────────────────────────────────────────────────────────────────────────┐
│                           DOMAIN LAYER                                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────────┐  │
│  │ WorkflowDef     │  │ WorkflowInstance│  │ Transition                  │  │
│  │ State           │  │ StateTransition │  │ Event, Action               │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
             │                                                 │
             ▼                                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      INFRASTRUCTURE LAYER                                   │
│  ┌───────────────────┐  ┌───────────────────┐  ┌─────────────────────────┐  │
│  │ InstanceRepository│  │ YAMLDefinition    │  │ KafkaPublisher /        │  │
│  │ (MongoDB)         │  │ Loader            │  │ NoOpPublisher           │  │
│  └───────────────────┘  └───────────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Definition Registry

- Central in-memory cache of workflow definitions.
- Populated at startup from YAML files or MongoDB.
- Optional: periodic reload or file watcher for hot-reload.

```go
// DefinitionRegistry provides thread-safe access to workflow definitions
type DefinitionRegistry interface {
    Get(name string) (*WorkflowDefinition, error)
    Register(def *WorkflowDefinition) error
    Reload() error
}
```

### 3.3 Transition Validator

- Ensures a transition exists for current state + event/action.
- Enforces:
  - Valid `from` state
  - Matching transition for event or action
  - Non-terminal state for non-terminal transitions
  - Idempotency (no transition if already in target state)

---

## 4. Data Structures (Detailed)

### 4.1 WorkflowDefinition

```go
type WorkflowDefinition struct {
    Name        string       `yaml:"name" json:"name"`
    Description string       `yaml:"description" json:"description"`
    Version     string       `yaml:"version" json:"version"`
    States      []StateDef   `yaml:"states" json:"states"`
    Transitions []Transition `yaml:"transitions" json:"transitions"`
}

type StateDef struct {
    ID      string `yaml:"id" json:"id"`
    Initial bool   `yaml:"initial" json:"initial"`
    Terminal bool  `yaml:"terminal" json:"terminal"`
}

type Transition struct {
    From   string `yaml:"from" json:"from"`
    To     string `yaml:"to" json:"to"`
    Event  string `yaml:"event,omitempty" json:"event,omitempty"`
    Action string `yaml:"action,omitempty" json:"action,omitempty"`
    // Exactly one of Event or Action must be set
}
```

Validation rules:

- Exactly one initial state.
- At least one terminal state.
- All `from`/`to` reference valid state IDs.
- Each transition has either `event` or `action`, not both.
- No duplicate transitions for same `(from, event)` or `(from, action)`.

### 4.2 WorkflowInstance

```go
type WorkflowInstance struct {
    ID           string                 `bson:"_id" json:"id"`
    WorkflowName string                 `bson:"workflow_name" json:"workflow_name"`
    State        string                 `bson:"state" json:"state"`
    Payload      map[string]interface{} `bson:"payload" json:"payload"`
    Context      map[string]string      `bson:"context" json:"context"`
    History      []StateTransition      `bson:"history" json:"history"`
    Version      int64                  `bson:"version" json:"-"` // optimistic lock
    CreatedAt    time.Time              `bson:"created_at" json:"created_at"`
    UpdatedAt    time.Time              `bson:"updated_at" json:"updated_at"`
}

type StateTransition struct {
    From      string    `bson:"from" json:"from"`
    To        string    `bson:"to" json:"to"`
    Trigger   string    `bson:"trigger" json:"trigger"`   // "event:vehicle.created" | "action:approve"
    Actor     string    `bson:"actor,omitempty" json:"actor,omitempty"`
    Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}
```

### 4.3 MongoDB Schema

**Collection: `workflow_instances`**

```json
{
  "_id": "uuid",
  "workflow_name": "vehicle_approval",
  "state": "pending_approval",
  "payload": { "vehicle_id": "v1", "vin": "..." },
  "context": { "correlation_id": "v1", "user_id": "u1" },
  "history": [
    { "from": "draft", "to": "pending_approval", "trigger": "event:vehicle.created", "timestamp": "..." }
  ],
  "version": 2,
  "created_at": "ISODate",
  "updated_at": "ISODate"
}
```

Indexes:

- `{ workflow_name: 1, state: 1 }` – list by workflow and state
- `{ "context.correlation_id": 1 }` – lookup by correlation
- `{ created_at: -1 }` – recent instances

---

## 5. Sequence Flows

### 5.1 Start Workflow

```
Client                API Handler         WorkflowService      Registry    Repository
   │                       │                      │                │            │
   │  POST /workflows/     │                      │                │            │
   │  vehicle_approval/    │                      │                │            │
   │  start                │                      │                │            │
   │──────────────────────>│                      │                │            │
   │                       │  Start(name, payload)│                │            │
   │                       │─────────────────────>│                │            │
   │                       │                      │  Get(name)     │            │
   │                       │                      │───────────────>│            │
   │                       │                      │<───────────────│            │
   │                       │                      │  validate def  │            │
   │                       │                      │  create instance (initial)  │
   │                       │                      │  Save(instance)│            │
   │                       │                      │────────────────────────────>│
   │                       │                      │<────────────────────────────│
   │                       │<─────────────────────│                │            │
   │<──────────────────────│  201 + instance      │                │            │
```

### 5.2 Process Event (Kafka → Transition)

```
Kafka Consumer       WorkflowService      Registry    Repository    ActionHandler
       │                     │                │            │              │
       │  event received     │                │            │              │
       │  (vehicle.created)  │                │            │              │
       │                     │                │            │              │
       │  ProcessEvent(      │                │            │              │
       │    instanceID,      │                │            │              │
       │    "vehicle.created")                │            │              │
       │────────────────────>│                │            │              │
       │                     │  FindByID()    │            │              │
       │                     │────────────────────────────>│              │
       │                     │<────────────────────────────│              │
       │                     │  Get(workflowName)          │              │
       │                     │───────────────>│            │              │
       │                     │<───────────────│            │              │
       │                     │  find transition(from, event)              │
       │                     │  validate      │            │              │
       │                     │  update state, append history              │
       │                     │  Update() (with version check)             │
       │                     │────────────────────────────>│              │
       │                     │<────────────────────────────│              │
       │                     │  OnTransition()             │              │
       │                     │───────────────────────────────────────────>│
       │                     │<───────────────────────────────────────────│
       │<────────────────────│  success       │            │              │
```

### 5.3 Process Action (REST → Transition)

```
Client          API Handler       WorkflowService      Repository
   │                 │                    │                  │
   │  POST           │                    │                  │
   │  /instances/    │                    │                  │
   │  {id}/actions/  │                    │                  │
   │  approve        │                    │                  │
   │────────────────>│                    │                  │
   │                 │  ProcessAction(    │                  │
   │                 │    id, "approve",  │                  │
   │                 │    actor)          │                  │
   │                 │───────────────────>│                  │
   │                 │                    │  FindByID()      │
   │                 │                    │─────────────────>│
   │                 │                    │<─────────────────│
   │                 │                    │  find transition │
   │                 │                    │  (from, action)  │
   │                 │                    │  update + Save   │
   │                 │                    │─────────────────>│
   │                 │                    │<─────────────────│
   │<────────────────│  200 OK            │                  │
```

---

## 6. Event Correlation

### 6.1 Resolving Instance from Kafka Event

Option A – `instance_id` in event:

```json
{
  "event": "vehicle.created",
  "vehicle_id": "v1",
  "workflow_instance_id": "wf-123"
}
```

Option B – correlation via `context`:

1. On `Start`, set `context["correlation_id"] = vehicle_id`.
2. Create index on `context.correlation_id`.
3. On `vehicle.created`, query: `FindByCorrelationID(vehicle_id)`.
4. If multiple instances per correlation, add workflow name and/or state filter.

Recommended: Option A when workflow-svc starts the instance and knows the ID; Option B when the vehicle is created first and workflow starts asynchronously.

---

## 7. Error Handling

### 7.1 Error Types

| Error | HTTP | Handling |
|-------|------|----------|
| `ErrDefinitionNotFound` | 404 | No workflow definition for name |
| `ErrInstanceNotFound` | 404 | Instance ID not found |
| `ErrInvalidTransition` | 409 | No transition for current state + event/action |
| `ErrInstanceTerminal` | 409 | Instance already in terminal state |
| `ErrConcurrentModification` | 409 | Optimistic lock failed (version conflict) |
| `ErrValidation` | 400 | Invalid payload, missing required fields |
| Internal/DB error | 500 | Log, return generic error |

### 7.2 Idempotency

- `ProcessEvent(instanceID, eventName, eventPayload)` with optional `event_id`.
- Store processed `event_id`s per instance (or in a separate collection).
- If `event_id` seen, return success without changing state.

### 7.3 Retry for Kafka Consumer

- Use consumer group offset commit only after successful `ProcessEvent`.
- On error: log, optionally send to DLQ; do not commit offset.
- Consumer retries same message on restart.

---

## 8. Concurrency

### 8.1 Optimistic Locking

```go
// Repository Update
filter := bson.M{
    "_id": instanceID,
    "version": instance.Version,
}
instance.Version++
instance.UpdatedAt = time.Now().UTC()
result, err := col.UpdateOne(ctx, filter, bson.M{
    "$set": bson.M{
        "state":      instance.State,
        "history":    instance.History,
        "payload":    instance.Payload,
        "version":    instance.Version,
        "updated_at": instance.UpdatedAt,
    },
})
if result.MatchedCount == 0 {
    return ErrConcurrentModification
}
```

### 8.2 Thread Safety

- `DefinitionRegistry`: read-heavy; use `sync.RWMutex` or immutable snapshots.
- `WorkflowService`: stateless; safe for concurrent use.
- Repository: MongoDB driver handles connection pooling.

---

## 9. YAML Loader

### 9.1 Loading Strategy

```
config/
  workflows/
    vehicle_approval.yaml
    maintenance_request.yaml
```

- Load all `.yaml` files from a configured directory at startup.
- Validate each definition; fail fast if any invalid.
- Register in `DefinitionRegistry`.

### 9.2 Validation on Load

1. Parse YAML.
2. Check required fields: `name`, `states`, `transitions`.
3. Validate states: exactly one initial, at least one terminal.
4. Validate transitions: valid `from`/`to`, unique `(from, event)` and `(from, action)`.
5. If valid, register; otherwise return error with file/line info.

---

## 10. Action Handler Extension

### 10.1 Interface

```go
type TransitionHandler interface {
    OnTransition(ctx context.Context, instance *WorkflowInstance, from, to string, trigger string) error
}
```

### 10.2 Usage

- Register handlers per workflow or globally.
- After a successful transition, call `OnTransition`.
- Use for: publishing Kafka events, calling external APIs, sending notifications.
- On handler error: log; optionally persist for retry; do not roll back the transition (eventual consistency).

---

## 11. Configuration

### 11.1 Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `WORKFLOW_DEFINITIONS_PATH` | Path to YAML files | `./config/workflows` |
| `MONGO_URI` | MongoDB connection | `mongodb://localhost:27017` |
| `MONGO_DATABASE` | Database name | `workflow_db` |
| `WORKFLOW_INSTANCE_COLLECTION` | Collection for instances | `workflow_instances` |
| `KAFKA_BROKERS` | Kafka brokers | `localhost:9092` |
| `WORKFLOW_TOPIC_PREFIX` | Topic prefix for workflow events | `workflow.` |

### 11.2 Viper (Optional)

```yaml
workflow:
  definitions_path: ./config/workflows
  mongo:
    uri: ${MONGO_URI}
    database: workflow_db
  kafka:
    brokers: ${KAFKA_BROKERS}
```

---

## 12. Observability

### 12.1 Logging (zap)

- Log at INFO: workflow start, transitions, actions.
- Log at WARN: invalid transition attempts, retries.
- Log at ERROR: repository errors, handler failures.
- Structured fields: `workflow_name`, `instance_id`, `state`, `trigger`.

### 12.2 Metrics (Future)

- `workflow_instances_total{workflow, state}`
- `workflow_transitions_total{workflow, from, to}`
- `workflow_process_duration_seconds`

### 12.3 Tracing (Future)

- Trace ID in context from HTTP/gRPC middleware.
- Propagate to Kafka consumer and repository.

---

## 13. Testing Strategy

| Type | Scope | Tools |
|------|-------|-------|
| Unit | Engine, validator, loader | `testing`, mocks |
| Integration | Repository, Kafka consumer | testcontainers (MongoDB, Kafka) |
| E2E | Start → Event → Action → Complete | docker-compose, HTTP client |

---

## 14. File Layout (Final)

```
workflow-svc/   (or internal/workflow in vehicle-svc)
├── cmd/
│   └── main.go
├── config/
│   └── workflows/
│       └── vehicle_approval.yaml
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   ├── workflow.go
│   │   │   ├── instance.go
│   │   │   └── action.go
│   │   ├── middleware/
│   │   │   └── auth.go
│   │   └── route/
│   │       └── route.go
│   ├── application/
│   │   ├── service/
│   │   │   └── workflow_service.go
│   │   ├── loader/
│   │   │   └── yaml_loader.go
│   │   └── registry/
│   │       └── definition_registry.go
│   ├── domain/
│   │   ├── workflow/
│   │   │   ├── definition.go
│   │   │   ├── instance.go
│   │   │   ├── engine.go
│   │   │   ├── errors.go
│   │   │   └── validator.go
│   │   └── repository/
│   │       └── instance_repository.go
│   └── infrastructure/
│       ├── persistence/
│       │   └── mongo_instance_repository.go
│       ├── loader/
│       │   └── file_loader.go
│       └── handler/
│           └── kafka_transition_handler.go
├── go.mod
├── go.sum
├── Dockerfile
├── .env.example
└── README.md
```

---

## 15. Appendix: State Machine Rules (Pseudocode)

```
FUNCTION ProcessEvent(instanceID, eventName, payload):
    instance = repository.FindByID(instanceID)
    IF instance == nil THEN RETURN ErrInstanceNotFound
    IF instance.State IN terminal_states THEN RETURN ErrInstanceTerminal

    def = registry.Get(instance.WorkflowName)
    transition = def.FindTransition(instance.State, event: eventName)
    IF transition == nil THEN RETURN ErrInvalidTransition

    instance.State = transition.To
    instance.History.Append(from=transition.From, to=transition.To, trigger="event:"+eventName)
    instance.Version++

    IF optimistic_update_fails THEN RETURN ErrConcurrentModification

    handler.OnTransition(instance, transition.From, transition.To)
    RETURN nil
```

---

## 16. Related Documents

- [WORKFLOW_ENGINE_DETAILED_DESIGN_V2.md](./WORKFLOW_ENGINE_DETAILED_DESIGN_V2.md) — Updated design with format abstraction and BPMN extensibility (v2.0).
- [WORKFLOW_ENGINE_CORE_DESIGN.md](./WORKFLOW_ENGINE_CORE_DESIGN.md) — Design detail for the core engine (pure execution logic).

---

**End of Detailed Design**
