package messaging

import (
	"context"
	"encoding/json"
	"time"

	"workflow-svc/internal/application/service"
	"workflow-svc/internal/infrastructure/persistence"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type IncomingEvent struct {
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	CorrelationID string                 `json:"correlation_id"` // vehicle_id
	Payload       map[string]interface{} `json:"payload"`
	Timestamp     time.Time              `json:"timestamp"`
}

type KafkaEventConsumer struct {
	reader       *kafka.Reader
	dlqWriter    *kafka.Writer
	workflowSvc  *service.WorkflowService
	deduplicator *persistence.EventDeduplicator
	logger       *zap.Logger
}

func NewKafkaEventConsumer(
	brokers []string,
	topic string,
	groupID string,
	dlqTopic string,
	workflowSvc *service.WorkflowService,
	deduplicator *persistence.EventDeduplicator,
	logger *zap.Logger,
) *KafkaEventConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})

	dlqWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    dlqTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaEventConsumer{
		reader:       reader,
		dlqWriter:    dlqWriter,
		workflowSvc:  workflowSvc,
		deduplicator: deduplicator,
		logger:       logger,
	}
}

func (c *KafkaEventConsumer) Start(ctx context.Context) error {
	c.logger.Info("Starting Kafka event consumer")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Stopping Kafka event consumer")
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				c.logger.Error("Failed to fetch message", zap.Error(err))
				continue
			}

			if err := c.processMessage(ctx, msg); err != nil {
				c.logger.Error("Failed to process message",
					zap.Error(err),
					zap.String("topic", msg.Topic),
					zap.Int("partition", msg.Partition),
					zap.Int64("offset", msg.Offset))

				// Send to DLQ
				if err := c.sendToDLQ(ctx, msg, err); err != nil {
					c.logger.Error("Failed to send message to DLQ", zap.Error(err))
				}
			}

			// Commit message
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.logger.Error("Failed to commit message", zap.Error(err))
			}
		}
	}
}

func (c *KafkaEventConsumer) processMessage(ctx context.Context, msg kafka.Message) error {
	var event IncomingEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	// Check for duplicate
	if event.EventID != "" {
		processed, err := c.deduplicator.IsProcessed(ctx, event.EventID)
		if err != nil {
			return err
		}
		if processed {
			c.logger.Info("Event already processed, skipping",
				zap.String("event_id", event.EventID))
			return nil
		}
	}

	// Process event using correlation_id (vehicle_id) to find instance
	instance, err := c.workflowSvc.ProcessEvent(ctx, event.CorrelationID, event.EventType, event.Payload)
	if err != nil {
		return err
	}

	// Mark as processed
	if event.EventID != "" {
		if err := c.deduplicator.MarkProcessed(ctx, event.EventID, instance.ID); err != nil {
			c.logger.Warn("Failed to mark event as processed",
				zap.String("event_id", event.EventID),
				zap.Error(err))
		}
	}

	c.logger.Info("Successfully processed event",
		zap.String("event_id", event.EventID),
		zap.String("event_type", event.EventType),
		zap.String("correlation_id", event.CorrelationID),
		zap.String("instance_id", instance.ID))

	return nil
}

func (c *KafkaEventConsumer) sendToDLQ(ctx context.Context, msg kafka.Message, processingErr error) error {
	dlqMsg := kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
		Headers: append(msg.Headers,
			kafka.Header{Key: "error", Value: []byte(processingErr.Error())},
			kafka.Header{Key: "original_topic", Value: []byte(msg.Topic)},
			kafka.Header{Key: "original_partition", Value: []byte{byte(msg.Partition)}},
		),
	}

	return c.dlqWriter.WriteMessages(ctx, dlqMsg)
}

func (c *KafkaEventConsumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.dlqWriter.Close()
}
