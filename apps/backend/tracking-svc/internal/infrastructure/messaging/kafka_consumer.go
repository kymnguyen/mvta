package messaging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type EventHandler func(ctx context.Context, payload []byte) error

type KafkaConsumer struct {
	reader   *kafka.Reader
	logger   *zap.Logger
	handlers map[string]EventHandler
}

func NewKafkaConsumer(brokers []string, groupID string, topics []string, logger *zap.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:          brokers,
			GroupID:          groupID,
			GroupTopics:      topics,
			MinBytes:         1e3,
			MaxBytes:         10e6,
			CommitInterval:   1 * time.Second,
			StartOffset:      kafka.LastOffset,
			SessionTimeout:   20 * time.Second,
			RebalanceTimeout: 60 * time.Second,
		}),
		logger:   logger,
		handlers: make(map[string]EventHandler),
	}
}

func (c *KafkaConsumer) RegisterHandler(topic string, handler EventHandler) {
	c.handlers[topic] = handler
	c.logger.Info("event handler registered", zap.String("topic", topic))
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	c.logger.Info("kafka consumer started")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("kafka consumer context cancelled")
			return ctx.Err()
		default:
		}

		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			c.logger.Error("failed to fetch message", zap.Error(err))
			continue
		}

		handler, ok := c.handlers[msg.Topic]
		if !ok {
			c.logger.Warn("no handler registered", zap.String("topic", msg.Topic))
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.logger.Error("failed to commit message", zap.Error(err))
			}
			continue
		}

		if err := handler(ctx, msg.Value); err != nil {
			c.logger.Error("handler failed, will retry",
				zap.String("topic", msg.Topic),
				zap.Int64("offset", msg.Offset),
				zap.Error(err),
			)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error("failed to commit message", zap.Error(err))
		}

		c.logger.Debug("message processed",
			zap.String("topic", msg.Topic),
			zap.Int64("partition", int64(msg.Partition)),
			zap.Int64("offset", msg.Offset),
		)
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

func Decode[T any](data []byte, out *T) error {
	return json.Unmarshal(data, out)
}
