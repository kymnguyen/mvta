# Detailed Design: Core Workflow Engine for MVTA (v2.0)

**Document:** MVTA-WF-ENGINE-DD  
**Version:** 2.0  
**Date:** January 31, 2025  
**Status:** Detailed Design  
**Supersedes:** [WORKFLOW_ENGINE_DETAILED_DESIGN.md](./WORKFLOW_ENGINE_DETAILED_DESIGN.md) (v1.0) — retained for reference

---

## Document Control

| Version | Date       | Changes                                           |
|---------|------------|---------------------------------------------------|
| 1.0     | 2025-01-31 | Initial design                                    |
| 2.0     | 2025-01-31 | Format abstraction, canonical model, BPMN roadmap |

---

## 1. Executive Summary

This document describes the technical design for the **Core Workflow Engine** used by the MVTA (Management Vehicle Tracking Application) platform. The engine orchestrates multi-step business processes via a state machine driven by events and actions.

**Design goals:**
- **Format-agnostic execution**: Engine operates on an internal canonical model, not tied to YAML.
- **Extensibility**: Support multiple definition formats (YAML, BPMN, JSON) via pluggable parsers.
- **Clear separation of concerns**: Parsing, execution, persistence, and integration are decoupled.

---

## 2. Scope and Definitions

### 2.1 Scope

| In Scope                                             | Out of Scope (Phase 1)                    |
|------------------------------------------------------|-------------------------------------------|
| State-machine workflows                              | Full BPMN 2.0 (parallel flows, timers)    |
| YAML definition format                               | Visual workflow designer                  |
| Event-driven and action-driven transitions           | Workflow versioning and migration         |
| MongoDB persistence                                  | Built-in notification/email               |
| Kafka integration                                    | Human task assignment engine              |

### 2.2 Definitions

| Term                 | Definition                                                                 |
|----------------------|-----------------------------------------------------------------------------|
| **Workflow Definition** | Static description of a process: nodes, edges, triggers.                  |
| **Workflow Instance**   | One running execution of a definition with its own state and history.     |
| **Node**                | A state or activity in the workflow (Start, Task, End).                   |
| **Edge**                | A directed transition between nodes, triggered by an event or action.     |
| **Trigger**             | An event or action that causes an edge to be followed.                    |
| **Canonical Model**     | Internal, format-agnostic representation of a workflow definition.        |

---

## 3. Design Principles

### 3.1 Format Abstraction

The engine does not depend on YAML, BPMN, or any specific format. It operates on a **canonical execution model** produced by **definition parsers**.

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  YAML Parser    │     │  BPMN Parser    │     │  JSON Parser    │
│  (Phase 1)      │     │  (Future)       │     │  (Future)       │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                                 ▼
                    ┌────────────────────────┐
                    │   CANONICAL MODEL      │
                    │   (Execution Graph)    │
                    └────────────┬───────────┘
                                 │
                                 ▼
                    ┌────────────────────────┐
                    │   WORKFLOW ENGINE      │
                    │   (State Machine)      │
                    └────────────────────────┘
```

### 3.2 Extensibility Points

| Extension Point        | Purpose                                      | Example                       |
|------------------------|----------------------------------------------|-------------------------------|
| **Definition Parser**  | Add new definition formats                   | YAML, BPMN, JSON              |
| **Trigger Handler**    | Add new trigger types                        | Timer, signal, message        |
| **Action Handler**     | Execute side effects on transitions          | Kafka publish, HTTP call      |
| **Storage Adapter**    | Change persistence backend                   | MongoDB, PostgreSQL           |
| **Event Source**       | Connect to different messaging systems       | Kafka, NATS, RabbitMQ         |

---

## 4. Canonical Execution Model

### 4.1 Concept

The engine uses a **canonical model** that all parsers must produce. This keeps the engine logic independent of definition format and prepares for BPMN.

### 4.2 Node Types

| Type          | Symbol   | Description                               | BPMN Mapping           |
|---------------|----------|-------------------------------------------|------------------------|
| **Start**     | ○        | Single entry point                        | Start Event            |
| **Task**      | □        | Work unit / state                         | User Task, Service Task|
| **End**       | ◎        | Terminal state                            | End Event              |
| **Gateway**   | ◇        | Reserved for future use                   | Exclusive/Parallel     |

Phase 1 implements Start, Task, End only. Gateway support is deferred.

### 4.3 Edge Model

Each edge represents a possible transition: `from_node_id ──[trigger]──> to_node_id`

A **trigger** is either:
- **Event**: External or system event (e.g., `vehicle.created`)
- **Action**: User/system action (e.g., `approve`, `reject`)

---

## 5. System Context

### 5.1 Context Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                         EXTERNAL ACTORS                                  │
│  ┌──────────┐  ┌───────────┐  ┌──────────┐  ┌──────────────────────────┐ │
│  │ Admin UI │  │vehicle-svc│  │ Kafka    │  │ BPMN / YAML Definitions  │ │
│  │ Clients  │  │ tracking  │  │ Topics   │  │ (future: wf designer).   │ │
│  └────┬─────┘  └────┬──────┘  └────┬─────┘  └────────────┬─────────────┘ │
└───────┼─────────────┼─────────────┼─────────────────────┼───────────────┘
        │ REST        │ Events      │ Consume/Publish     │ Load definitions
        ▼             ▼             ▼                     ▼
┌──────────────────────────────────────────────────────────────────────────┐
│                      WORKFLOW ENGINE SERVICE                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │ API Layer   │  │ Engine Core │  │ Parsers     │  │ Integrations    │  │
│  │ (REST)      │  │ (Canonical) │  │ (YAML/BPMN) │  │ (Kafka, etc.)   │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────────┘  │
└────────────────────────────────────────────────┬─────────────────────────┘
                                                 │
                                                 ▼
┌──────────────────────────────────────────────────────────────────────────┐
│  MongoDB  │  Definition Store (files/DB)  │  Kafka                       │
└──────────────────────────────────────────────────────────────────────────┘
```

### 5.2 Deployment Options

| Option | Description | Use Case |
|--------|-------------|----------|
| **A. Dedicated workflow-svc** | Standalone service with its own DB and Kafka consumer | Recommended |
| **B. Embedded** | Engine as a library inside another service | Prototype or single-service scenarios |

---

## 6. Component Architecture

### 6.1 Layered View

```
┌───────────────────────────────────────────────────────────────────────────┐
│ API LAYER                                                                 │
│  WorkflowHandler │ InstanceHandler │ ActionHandler                        │
└────────────────────────────────────┬──────────────────────────────────────┘
                                     │
                                     v
┌───────────────────────────────────────────────────────────────────────────┐
│ APPLICATION LAYER                  │                                       │
│  WorkflowService (orchestrates)    │  DefinitionRegistry (CanonicalModel)  │
│  TransitionValidator               │  TriggerResolver                      │
└────────────────────────────────────┼───────────────────────────────────────┘
                                     │
                                     v
┌───────────────────────────────────────────────────────────────────────────┐
│ DOMAIN LAYER                       │                                       │
│  CanonicalModel │ Instance │ Edge  │ Node │ Trigger                         │
└────────────────────────────────────┼───────────────────────────────────────┘
                                     │
                                     v
┌─────────────────────────────────────────────────────────────────────────┐
│ INFRASTRUCTURE LAYER                                                     │
│  DefinitionParser (YAML) │ InstanceRepo │ KafkaConsumer │ ActionHandlers  │
└──────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Definition Registry

The registry stores **CanonicalModels**, not raw YAML/BPMN. Parsers are used only at load time.

```go
type DefinitionRegistry interface {
    Get(name string) (*CanonicalModel, error)
    Register(model *CanonicalModel) error
    List() ([]string, error)
    Reload() error
}
```

### 6.3 Engine Interface

The engine operates only on the canonical model:

```go
type Engine interface {
    Start(ctx context.Context, workflowName string, payload, context map[string]interface{}) (*Instance, error)
    ProcessEvent(ctx context.Context, instanceID, eventName string, payload map[string]interface{}) error
    ProcessAction(ctx context.Context, instanceID, actionName, actor string) error
    GetInstance(ctx context.Context, instanceID string) (*Instance, error)
    ListInstances(ctx context.Context, filter InstanceFilter) ([]*Instance, error)
}
```

---

## 7. Data Structures

### 7.1 Canonical Model (Definition)

```go
type CanonicalModel struct {
    ID          string
    Name        string
    Version     string
    Nodes       map[string]Node
    Edges       []Edge
    StartNodeID string
}

type Node struct {
    ID      string
    Type    NodeType // Start, Task, End, Gateway
    Meta    map[string]string
}

type Edge struct {
    ID       string
    From     string
    To       string
    Trigger  Trigger
}

type Trigger struct {
    Type   TriggerType // Event | Action
    Name   string
}
```

### 7.2 Instance (Runtime)

```go
type Instance struct {
    ID           string
    WorkflowID   string
    CurrentNode  string
    Payload      map[string]interface{}
    Context      map[string]string
    History      []TransitionRecord
    Version      int64
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type TransitionRecord struct {
    From      string
    To        string
    Trigger   string
    Actor     string
    Timestamp time.Time
}
```

### 7.3 MongoDB Schema

**Collection: `workflow_instances`**

```json
{
  "_id": "uuid",
  "workflow_id": "vehicle_approval",
  "current_node": "pending_approval",
  "payload": {},
  "context": { "correlation_id": "v1" },
  "history": [],
  "version": 1,
  "created_at": "ISODate",
  "updated_at": "ISODate"
}
```

Indexes: `{ workflow_id: 1, current_node: 1 }`, `{ "context.correlation_id": 1 }`, `{ created_at: -1 }`

---

## 8. Definition Format: Parser Interface

### 8.1 Parser Contract

```go
type DefinitionParser interface {
    Format() string
    Parse(raw []byte) (*CanonicalModel, error)
    Validate(raw []byte) error
}
```

### 8.2 YAML Format (Phase 1)

```yaml
name: vehicle_approval
version: "1.0"
states:
  - id: draft
    type: start
  - id: pending_approval
    type: task
  - id: approved
    type: end
  - id: rejected
    type: end
transitions:
  - from: draft
    to: pending_approval
    event: vehicle.created
  - from: pending_approval
    to: approved
    action: approve
  - from: pending_approval
    to: rejected
    action: reject
```

### 8.3 BPMN Extension Path (Future)

1. Implement `DefinitionParser` for BPMN XML.
2. Map BPMN elements to `CanonicalModel`.
3. Extend `NodeType` for Gateways when engine supports branching.

---

## 9. Execution Flow

### 9.1 Start Workflow

`Client → POST /workflows/{name}/start → Registry.Get → Resolve StartNodeID → Create Instance → Save → Return`

### 9.2 Process Event / Action

`Trigger → FindByID → Registry.Get → Resolve edge → Update instance → Append history → Update (optimistic lock) → OnTransition (optional)`

---

## 10. Error Handling and Concurrency

### 10.1 Errors

| Error | HTTP | Description |
|-------|------|-------------|
| `ErrDefinitionNotFound` | 404 | Unknown workflow name |
| `ErrInstanceNotFound` | 404 | Unknown instance ID |
| `ErrInvalidTransition` | 409 | No edge for current node + trigger |
| `ErrInstanceTerminal` | 409 | Instance already at End node |
| `ErrConcurrentModification` | 409 | Version conflict on update |
| `ErrParse` | 400 | Invalid definition (YAML/BPMN) |

### 10.2 Concurrency

- Optimistic locking via `version` on instances.
- Thread-safe `DefinitionRegistry`.
- Stateless engine; safe for concurrent use.

---

## 11. Configuration and Observability

### 11.1 Configuration

| Variable | Description |
|----------|-------------|
| `WORKFLOW_DEFINITIONS_PATH` | Path for YAML definitions |
| `WORKFLOW_DEFINITION_FORMAT` | Default format: `yaml` |
| `MONGO_URI`, `MONGO_DATABASE` | MongoDB connection |
| `KAFKA_BROKERS` | Kafka brokers |

### 11.2 Observability

- Structured logging (zap) with `workflow_id`, `instance_id`, `node`, `trigger`.
- Future metrics: instance count by workflow/state, transition rate.

---

## 12. File Layout

```
workflow-svc/
├── cmd/main.go
├── config/workflows/
│   └── vehicle_approval.yaml
├── internal/
│   ├── api/
│   ├── application/
│   │   ├── service/
│   │   ├── registry/
│   │   └── validator/
│   ├── domain/
│   │   ├── engine/         # Core engine (see WORKFLOW_ENGINE_CORE_DESIGN.md)
│   │   ├── model/
│   │   └── repository/
│   └── infrastructure/
│       ├── parser/         # YAMLParser, (future) BPMNParser
│       ├── persistence/
│       ├── messaging/
│       └── handler/
├── go.mod
├── Dockerfile
└── .env.example
```

---

## 13. Implementation Phases

| Phase | Scope | Outcome |
|-------|-------|---------|
| **1** | Canonical model, YAML parser, engine core | Working state-machine engine |
| **2** | REST API, Kafka integration, MongoDB | End-to-end with vehicle-svc |
| **3** | Action handlers, observability | Production-ready workflow-svc |
| **4** | BPMN parser (subset) | BPMN support for selected flows |
| **5** | Gateways, timers, subprocesses | Advanced BPMN support |

---

## 14. Related Documents

- [WORKFLOW_ENGINE_DETAILED_DESIGN.md](./WORKFLOW_ENGINE_DETAILED_DESIGN.md) — Original design (v1.0).
- [WORKFLOW_ENGINE_CORE_DESIGN.md](./WORKFLOW_ENGINE_CORE_DESIGN.md) — Design detail for the core engine (pure execution logic).

---

**End of Detailed Design (v2.0)**
