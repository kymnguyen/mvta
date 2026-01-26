package repository

import (
	"context"

	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/entity"
	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/valueobject"
)

// VehicleRepository defines repository operations for Vehicle aggregate.
type VehicleRepository interface {
	// Save persists a vehicle aggregate with optimistic concurrency control.
	Save(ctx context.Context, vehicle *entity.Vehicle) error

	// FindByID retrieves a vehicle by its ID.
	FindByID(ctx context.Context, id valueobject.VehicleID) (*entity.Vehicle, error)

	// FindAll retrieves all vehicles with pagination.
	FindAll(ctx context.Context, limit int, offset int) ([]*entity.Vehicle, error)

	// Delete removes a vehicle from storage.
	Delete(ctx context.Context, id valueobject.VehicleID) error

	// ExistsByVIN checks if a vehicle with the given VIN exists.
	ExistsByVIN(ctx context.Context, vin string) (bool, error)
}

// EventStore defines event sourcing operations.
type EventStore interface {
	// AppendEvent appends a domain event to the event log.
	AppendEvent(ctx context.Context, aggregateID string, event interface{}, version int64) error

	// GetEventsByAggregateID retrieves all events for an aggregate.
	GetEventsByAggregateID(ctx context.Context, aggregateID string) ([]interface{}, error)

	// GetEventsSince retrieves events since a specific version.
	GetEventsSince(ctx context.Context, afterVersion int64) ([]interface{}, error)
}

// OutboxRepository defines operations for the transactional outbox pattern.
type OutboxRepository interface {
	// SaveOutboxEvent saves a domain event to the outbox for asynchronous publication.
	SaveOutboxEvent(ctx context.Context, aggregateID string, event interface{}) error

	// GetPendingEvents retrieves unpublished events from the outbox.
	GetPendingEvents(ctx context.Context, limit int) ([]OutboxEvent, error)

	// MarkEventAsPublished marks an outbox event as successfully published.
	MarkEventAsPublished(ctx context.Context, eventID string) error
}

// OutboxEvent represents an event in the outbox.
type OutboxEvent struct {
	ID          string
	AggregateID string
	EventType   string
	EventData   []byte
	CreatedAt   int64
	PublishedAt *int64
}

// UnitOfWork coordinates transactional operations across repositories.
type UnitOfWork interface {
	// BeginTx starts a new transaction.
	BeginTx(ctx context.Context) (Transaction, error)

	// Commit commits the current transaction.
	Commit(ctx context.Context) error

	// Rollback rolls back the current transaction.
	Rollback(ctx context.Context) error
}

// Transaction represents a database transaction.
type Transaction interface {
	// Commit commits the transaction.
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction.
	Rollback(ctx context.Context) error
}
