# Core Engine Design Detail

**Document:** MVTA-WF-CORE-DD  
**Version:** 1.0  
**Date:** January 31, 2025  
**Status:** Detailed Design  
**Related:** [WORKFLOW_ENGINE_DETAILED_DESIGN.md](./WORKFLOW_ENGINE_DETAILED_DESIGN.md) (v1.0) | [WORKFLOW_ENGINE_DETAILED_DESIGN_V2.md](./WORKFLOW_ENGINE_DETAILED_DESIGN_V2.md) (v2.0)

---

## 1. Core Engine Scope

### 1.1 Responsibility Boundary

The core engine is the **pure execution logic** that computes the next state of a workflow instance given a trigger. It:

| Does | Does Not |
|------|----------|
| Executes a workflow instance from current node given a trigger | Parse YAML/BPMN |
| Validates that transitions are allowed before applying them | Persist instances |
| Produces deterministic next state from (definition + current node + trigger) | Handle Kafka/HTTP |
| Is stateless and has no I/O | Call external services |

### 1.2 Core vs. Surrounding Components

```
┌─────────────────────────────────────────────────────────────────────┐
│                        OUTSIDE CORE (orchestration)                  │
│  API │ Kafka Consumer │ Repository │ ActionHandler │ DefinitionLoader│
└─────────────────────────────────────────┬───────────────────────────┘
                                          │
                                          ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         CORE ENGINE                                  │
│                                                                      │
│   Input:  (Definition, InstanceSnapshot, Trigger)                    │
│   Output: (NextNode, TransitionRecord) | Error                       │
│                                                                      │
│   Pure function: no I/O, no side effects                             │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. Core Data Model

### 2.1 Definition (Read-Only, Immutable)

```go
// Definition is the execution graph. Immutable during execution.
type Definition struct {
    ID          string
    Name        string
    Version     string
    Nodes       map[string]Node
    Edges       []Edge
    StartNodeID string
}

type Node struct {
    ID      string
    Type    NodeType // Start | Task | End
    IsStart bool
    IsEnd   bool
}

type Edge struct {
    ID      string
    From    string
    To      string
    Trigger Trigger
}

type Trigger struct {
    Type TriggerType // Event | Action
    Name string      // e.g., "vehicle.created", "approve"
}

type NodeType int
const (
    NodeTypeStart NodeType = iota + 1
    NodeTypeTask
    NodeTypeEnd
)

type TriggerType int
const (
    TriggerEvent  TriggerType = 1
    TriggerAction TriggerType = 2
)
```

### 2.2 Instance Snapshot (Input to Core)

```go
// InstanceSnapshot is the minimal view the core needs to compute a transition.
// Supplied by the caller (e.g., repository).
type InstanceSnapshot struct {
    ID          string
    WorkflowID  string
    CurrentNode string
    Version     int64
}
```

### 2.3 Core Output

```go
// TransitionResult is the output of a successful transition.
type TransitionResult struct {
    FromNode   string
    ToNode     string
    Trigger    Trigger
    Record     TransitionRecord
    IsTerminal bool
}

type TransitionRecord struct {
    From      string
    To        string
    Trigger   string
    Actor     string
    Timestamp time.Time
}

// StartResult is the output of EngineCore.Start.
type StartResult struct {
    InitialNode string
    Record      *TransitionRecord // nil for start
}
```

---

## 3. Core Interface

### 3.1 EngineCore Contract

```go
// EngineCore is the pure execution logic. No I/O. Thread-safe.
type EngineCore interface {
    // Start produces the initial node for a new workflow instance.
    Start(ctx context.Context, def *Definition, payload, ctxMap map[string]interface{}) (*StartResult, error)

    // Transition computes the next node for a trigger. Pure function; no I/O.
    Transition(ctx context.Context, def *Definition, instance *InstanceSnapshot, trigger Trigger, actor string) (*TransitionResult, error)

    // CanTransition returns true if a transition exists (read-only check).
    CanTransition(def *Definition, instance *InstanceSnapshot, trigger Trigger) bool
}
```

### 3.2 Design Notes

- `EngineCore` receives `*Definition` and `*InstanceSnapshot`; it does not fetch them.
- `Start` and `Transition` are pure: same inputs → same outputs.
- `ctx` is passed for future extensibility (e.g., deadlines, tracing); core does not use it for I/O.

---

## 4. Trigger Resolution Algorithm

### 4.1 Matching Rules

For a given `(currentNodeID, trigger)`:

1. Find edges where `edge.From == currentNodeID`.
2. Match trigger: `edge.Trigger.Type == trigger.Type && edge.Trigger.Name == trigger.Name`.
3. Return the first matching edge (or none).

**Determinism:** Exactly one edge per `(from, trigger)`; duplicates are invalid in the definition.

### 4.2 Pseudocode

```
FUNCTION ResolveEdge(def, currentNodeID, trigger):
    FOR each edge IN def.Edges:
        IF edge.From != currentNodeID:
            CONTINUE
        IF edge.Trigger.Type != trigger.Type:
            CONTINUE
        IF edge.Trigger.Name != trigger.Name:
            CONTINUE
        RETURN edge
    RETURN nil
```

### 4.3 Edge Lookup Optimization

Build an index at load time for O(1) lookup:

```go
// EdgeIndex: (fromNodeID, triggerType, triggerName) -> Edge
type EdgeIndex map[string]map[TriggerType]map[string]Edge

func BuildEdgeIndex(edges []Edge) EdgeIndex {
    idx := make(EdgeIndex)
    for _, e := range edges {
        if idx[e.From] == nil {
            idx[e.From] = make(map[TriggerType]map[string]Edge)
        }
        if idx[e.From][e.Trigger.Type] == nil {
            idx[e.From][e.Trigger.Type] = make(map[string]Edge)
        }
        idx[e.From][e.Trigger.Type][e.Trigger.Name] = e
    }
    return idx
}

func (idx EdgeIndex) Lookup(from string, t Trigger) (Edge, bool) {
    e, ok := idx[from][t.Type][t.Name]
    return e, ok
}
```

---

## 5. Transition Algorithm (Core Logic)

### 5.1 Preconditions (Checked in Order)

| Check | Error | When |
|-------|-------|------|
| Definition exists | ErrDefinitionNotFound | def == nil |
| Instance exists | ErrInstanceNotFound | instance == nil |
| Current node in definition | ErrInvalidState | def.Nodes[instance.CurrentNode] missing |
| Node is not End | ErrInstanceTerminal | node.IsEnd |
| Edge exists | ErrInvalidTransition | ResolveEdge returns nil |

### 5.2 Transition Steps

```
FUNCTION Transition(def, instance, trigger, actor):
    1. node := def.Nodes[instance.CurrentNode]
       IF node == nil THEN RETURN ErrInvalidState

    2. IF node.IsEnd THEN RETURN ErrInstanceTerminal

    3. edge := ResolveEdge(def, instance.CurrentNode, trigger)
       IF edge == nil THEN RETURN ErrInvalidTransition

    4. toNode := def.Nodes[edge.To]
       IF toNode == nil THEN RETURN ErrInvalidDefinition

    5. record := TransitionRecord{
           From:      instance.CurrentNode,
           To:        edge.To,
           Trigger:   formatTrigger(trigger),
           Actor:     actor,
           Timestamp: now(),
       }

    6. RETURN TransitionResult{
           FromNode:   instance.CurrentNode,
           ToNode:     edge.To,
           Trigger:    trigger,
           Record:     record,
           IsTerminal: toNode.IsEnd,
       }
```

### 5.3 Start Algorithm

```
FUNCTION Start(def, payload, ctx):
    1. IF def.StartNodeID == "" THEN RETURN ErrInvalidDefinition
    2. node := def.Nodes[def.StartNodeID]
       IF node == nil THEN RETURN ErrInvalidDefinition
    3. RETURN StartResult{ InitialNode: def.StartNodeID, Record: nil }
```

---

## 6. Invariants and Guarantees

### 6.1 Invariants (Always True)

| ID | Invariant |
|----|-----------|
| I1 | `instance.CurrentNode` is always a valid node ID in the definition |
| I2 | From a Start node, only event/action edges apply |
| I3 | At an End node, no transitions are allowed |
| I4 | For a given `(from, trigger)`, at most one edge exists |
| I5 | `Transition` is pure: same inputs → same output |
| I6 | No transitions modify the definition |

### 6.2 Guarantees

| ID | Guarantee |
|----|-----------|
| G1 | **Determinism**: Same (def, instance, trigger) → same result |
| G2 | **No I/O**: Core does not call DB, network, or file system |
| G3 | **Thread-safety**: Core is stateless; safe for concurrent use |
| G4 | **Idempotency**: Same trigger applied twice yields same result; caller enforces no double-apply via persistence |

---

## 7. State Diagram (Core Perspective)

```
                    ┌─────────────┐
                    │   START     │
                    │ (external)  │
                    └──────┬──────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  instance.CurrentNode  │
              │  = def.StartNodeID     │
              └────────────┬───────────┘
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
         ▼                 ▼                 ▼
    ┌─────────┐      ┌─────────┐      ┌─────────┐
    │ Trigger │      │ Trigger │      │ Trigger │
    │  A      │      │  B      │      │  C      │
    └────┬────┘      └────┬────┘      └────┬────┘
         │                 │                 │
         ▼                 ▼                 ▼
    ResolveEdge      ResolveEdge       ResolveEdge
         │                 │                 │
         ▼                 ▼                 ▼
    ┌─────────┐      ┌─────────┐      ┌─────────┐
    │ Next    │      │ Next    │      │ Next    │
    │ Node    │      │ Node    │      │ Node    │
    └────┬────┘      └────┬────┘      └────┬────┘
         │                 │                 │
         └─────────────────┼─────────────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │  TransitionResult      │
              │  (ToNode, Record)      │
              └────────────┬───────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
              ▼                         ▼
        IsTerminal=false          IsTerminal=true
              │                         │
              │                         ▼
              │                 ┌─────────────┐
              │                 │   END       │
              │                 │ (no more    │
              │                 │  triggers)  │
              │                 └─────────────┘
              │
              └─────► Loop: wait for next trigger
```

---

## 8. Error Handling (Core)

### 8.1 Error Types

```go
var (
    ErrDefinitionNotFound  = errors.New("workflow definition not found")
    ErrInstanceNotFound    = errors.New("workflow instance not found")
    ErrInvalidState        = errors.New("current node not found in definition")
    ErrInstanceTerminal    = errors.New("instance already in terminal state")
    ErrInvalidTransition   = errors.New("no transition for current state and trigger")
    ErrInvalidDefinition   = errors.New("definition has invalid structure")
)
```

### 8.2 Error Semantics

| Error | Meaning |
|-------|---------|
| ErrInstanceTerminal | Instance is complete; no further transitions allowed |
| ErrInvalidTransition | Trigger not valid from current node |
| ErrInvalidState | Definition/instance inconsistency (bug or corruption) |

---

## 9. Extensibility Within Core

### 9.1 Adding Trigger Types

Extend `TriggerType` and matching logic in `ResolveEdge`:

```go
const (
    TriggerEvent   TriggerType = 1
    TriggerAction  TriggerType = 2
    TriggerTimer   TriggerType = 3  // future
    TriggerSignal  TriggerType = 4  // future
)
```

### 9.2 Conditional Transitions (Future)

Add optional condition to edges:

```go
type Edge struct {
    ID        string
    From      string
    To        string
    Trigger   Trigger
    Condition Condition // optional; eval at runtime
}

type Condition interface {
    Evaluate(payload map[string]interface{}) bool
}
```

`ResolveEdge` would evaluate `Condition` when present.

### 9.3 Gateway Support (Future)

Introduce `NodeType` Gateway. Resolution splits into two phases:
1. Enter gateway: resolve incoming edge (same as today).
2. Exit gateway: gateway-specific logic (exclusive choice, parallel fork) selects outgoing edge(s).

---

## 10. Core Package Layout

```
internal/domain/engine/
├── core.go           # EngineCore interface + DefaultEngineCore implementation
├── definition.go     # Definition, Node, Edge, Trigger
├── instance.go       # InstanceSnapshot
├── result.go         # TransitionResult, StartResult
├── resolver.go       # Edge resolution, EdgeIndex
├── errors.go         # sentinel errors
└── core_test.go      # unit tests (no mocks, no I/O)
```

---

## 11. Testing Strategy for Core

### 11.1 Unit Tests (Pure)

- Fixed `Definition` + `InstanceSnapshot` + `Trigger` → assert `TransitionResult` or error.
- Cover: valid transition, invalid trigger, terminal state, missing node, duplicate edges.
- No database, Kafka, or file system.

### 11.2 Property-Based Tests

- For any valid definition and instance, applying a valid trigger yields a valid result.
- For any terminal instance, any trigger yields `ErrInstanceTerminal`.
- `Transition` is deterministic for same inputs.

---

## 12. Summary

| Aspect | Design Choice |
|--------|---------------|
| **Responsibility** | Pure execution: (Definition, Instance, Trigger) → Result |
| **I/O** | None; all persistence and messaging outside core |
| **State** | Stateless; safe for concurrent use |
| **Determinism** | Same inputs → same output |
| **Extensibility** | New trigger types, conditions, gateways via interfaces |
| **Testability** | Pure functions; no mocks needed for core logic |

The core engine is a small, focused, testable component that can be embedded in any service or workflow runtime.

---

**End of Core Engine Design Detail**
