package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/repository"
)

// OutboxWorker asynchronously publishes domain events from the outbox.
type OutboxWorker struct {
	outboxRepo     repository.OutboxRepository
	eventPublisher EventPublisher
	logger         *zap.Logger
	pollInterval   time.Duration
	batchSize      int
	done           chan struct{}
}

// EventPublisher publishes domain events to external systems.
type EventPublisher interface {
	// Publish publishes an event to a topic.
	Publish(ctx context.Context, topic string, event interface{}) error
}

// NewOutboxWorker creates a new outbox worker for event publishing.
func NewOutboxWorker(
	outboxRepo repository.OutboxRepository,
	eventPublisher EventPublisher,
	logger *zap.Logger,
	pollInterval time.Duration,
	batchSize int,
) *OutboxWorker {
	return &OutboxWorker{
		outboxRepo:     outboxRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
		pollInterval:   pollInterval,
		batchSize:      batchSize,
		done:           make(chan struct{}),
	}
}

// Start begins the outbox worker polling loop.
func (w *OutboxWorker) Start(ctx context.Context) {
	go w.pollLoop(ctx)
}

// Stop gracefully stops the outbox worker.
func (w *OutboxWorker) Stop() {
	close(w.done)
}

func (w *OutboxWorker) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.done:
			w.logger.Info("outbox worker stopped")
			return
		case <-ctx.Done():
			w.logger.Info("outbox worker context cancelled")
			return
		case <-ticker.C:
			if err := w.processPendingEvents(ctx); err != nil {
				w.logger.Error("failed to process outbox events", zap.Error(err))
			}
		}
	}
}

func (w *OutboxWorker) processPendingEvents(ctx context.Context) error {
	events, err := w.outboxRepo.GetPendingEvents(ctx, w.batchSize)
	if err != nil {
		return fmt.Errorf("failed to get pending events: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	w.logger.Debug("processing pending outbox events", zap.Int("count", len(events)))

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

func (w *OutboxWorker) publishEvent(ctx context.Context, event repository.OutboxEvent) error {
	var data map[string]interface{}
	if err := json.Unmarshal(event.EventData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	// Determine topic from event type
	topic := w.determineTopicFromEventType(event.EventType)

	// Publish event to topic
	if err := w.eventPublisher.Publish(ctx, topic, data); err != nil {
		return fmt.Errorf("failed to publish event to topic %s: %w", topic, err)
	}

	return nil
}

func (w *OutboxWorker) determineTopicFromEventType(eventType string) string {
	// Map event types to topics for multi-service communication
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
