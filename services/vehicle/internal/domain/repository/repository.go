package repository

import (
	"context"

	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/entity"
	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/valueobject"
)

type VehicleRepository interface {
	Save(ctx context.Context, vehicle *entity.Vehicle) error

	FindByID(ctx context.Context, id valueobject.VehicleID) (*entity.Vehicle, error)

	FindAll(ctx context.Context, limit int, offset int) ([]*entity.Vehicle, error)

	Delete(ctx context.Context, id valueobject.VehicleID) error

	ExistsByVIN(ctx context.Context, vin string) (bool, error)
}

type EventStore interface {
	AppendEvent(ctx context.Context, aggregateID string, event interface{}, version int64) error

	GetEventsByAggregateID(ctx context.Context, aggregateID string) ([]interface{}, error)

	GetEventsSince(ctx context.Context, afterVersion int64) ([]interface{}, error)
}

type OutboxRepository interface {
	SaveOutboxEvent(ctx context.Context, aggregateID string, event interface{}) error

	GetPendingEvents(ctx context.Context, limit int) ([]OutboxEvent, error)

	MarkEventAsPublished(ctx context.Context, eventID string) error
}

type OutboxEvent struct {
	ID          string
	AggregateID string
	EventType   string
	EventData   []byte
	CreatedAt   int64
	PublishedAt *int64
}

type UnitOfWork interface {
	BeginTx(ctx context.Context) (Transaction, error)

	Commit(ctx context.Context) error

	Rollback(ctx context.Context) error
}

type Transaction interface {
	Commit(ctx context.Context) error

	Rollback(ctx context.Context) error
}
