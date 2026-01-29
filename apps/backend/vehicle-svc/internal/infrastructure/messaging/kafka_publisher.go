package messaging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type EventPublisher interface {
	Publish(ctx context.Context, topic string, event interface{}) error
	Close() error
}

type KafkaPublisher struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func NewKafkaPublisher(brokers []string, logger *zap.Logger) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
		},
		logger: logger,
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("failed to marshal event", zap.Error(err))
		return err
	}

	var key []byte
	if m, ok := event.(map[string]interface{}); ok {
		if v, ok := m["id"].(string); ok && v != "" {
			key = []byte(v)
		}
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: payload,
		Time:  time.Now(),
	}); err != nil {
		p.logger.Error("failed to publish event",
			zap.String("topic", topic),
			zap.Error(err),
		)
		return err
	}

	p.logger.Debug("event published",
		zap.String("topic", topic),
		zap.String("key", string(key)),
	)
	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
