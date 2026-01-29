// This file was renamed from outbox_worker.go to domain_event_worker.go
// See the original file for full implementation.
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/infrastructure/resilience"
)

var (
	outboxCircuitBreaker = resilience.NewCircuitBreaker(3, 2, 10*time.Second)
	outboxRetryPolicy    = resilience.NewRetryPolicy(3, 100*time.Millisecond, 1*time.Second)
)

type DomainEventWorker struct {
	outboxRepo     repository.OutboxRepository
	eventPublisher EventPublisher
	logger         *zap.Logger
	pollInterval   time.Duration
	batchSize      int
	done           chan struct{}
}

type EventPublisher interface {
	Publish(ctx context.Context, topic string, event interface{}) error
	Close() error
}

func NewDomainEventWorker(
	outboxRepo repository.OutboxRepository,
	eventPublisher EventPublisher,
	logger *zap.Logger,
	pollInterval time.Duration,
	batchSize int,
) *DomainEventWorker {
	return &DomainEventWorker{
		outboxRepo:     outboxRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
		pollInterval:   pollInterval,
		batchSize:      batchSize,
		done:           make(chan struct{}),
	}
}

func (w *DomainEventWorker) Start(ctx context.Context) {
	go w.pollLoop(ctx)
}

func (w *DomainEventWorker) Stop() {
	close(w.done)
}

func (w *DomainEventWorker) Close() error {
	w.Stop()
	return w.eventPublisher.Close()
}

func (w *DomainEventWorker) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.done:
			w.logger.Info("domain event worker stopped")
			return
		case <-ctx.Done():
			w.logger.Info("domain event worker context cancelled")
			return
		case <-ticker.C:
			if err := w.processPendingEvents(ctx); err != nil {
				w.logger.Error("failed to process domain events", zap.Error(err))
			}
		}
	}
}

func (w *DomainEventWorker) processPendingEvents(ctx context.Context) error {
	events, err := w.outboxRepo.GetPendingEvents(ctx, w.batchSize)
	if err != nil {
		return fmt.Errorf("failed to get pending events: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	w.logger.Debug("processing pending domain events", zap.Int("count", len(events)))

	for _, event := range events {
		if err := w.publishEvent(ctx, event); err != nil {
			w.logger.Error("failed to publish event",
				zap.String("eventId", event.ID),
				zap.Error(err))
			continue
		}

		if err := w.outboxRepo.MarkEventAsPublished(ctx, event.ID); err != nil {
			w.logger.Error("failed to mark event as published",
				zap.String("eventId", event.ID),
				zap.Error(err))
		}
	}

	return nil
}

func (w *DomainEventWorker) publishEvent(ctx context.Context, event repository.OutboxEvent) error {
	var data map[string]interface{}
	if err := json.Unmarshal(event.EventData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	topic := w.determineTopicFromEventType(event.EventType)

	err := outboxCircuitBreaker.Execute(func() error {
		return outboxRetryPolicy.Execute(ctx, func() error {
			return w.eventPublisher.Publish(ctx, topic, data)
		})
	})
	if err != nil {
		return fmt.Errorf("failed to publish event to topic %s: %w", topic, err)
	}
	return nil
}

func (w *DomainEventWorker) determineTopicFromEventType(eventType string) string {
	topicMap := map[string]string{
		"*event.VehicleCreatedEvent":          "vehicle.created",
		"*event.VehicleLocationUpdatedEvent":  "vehicle.location.updated",
		"*event.VehicleStatusChangedEvent":    "vehicle.status.changed",
		"*event.VehicleMileageUpdatedEvent":   "vehicle.mileage.updated",
		"*event.VehicleFuelLevelUpdatedEvent": "vehicle.fuel.updated",
	}

	if topic, exists := topicMap[eventType]; exists {
		return topic
	}

	return "vehicle.events"
}
