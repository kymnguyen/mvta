package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TopicConfig struct {
	Name              string
	NumPartitions     int
	ReplicationFactor int
}

func InitializeTopics(brokers []string, logger *zap.Logger) error {
	topics := []TopicConfig{
		{Name: "vehicle.created", NumPartitions: 3, ReplicationFactor: 1},
		{Name: "vehicle.location.updated", NumPartitions: 3, ReplicationFactor: 1},
		{Name: "vehicle.status.changed", NumPartitions: 3, ReplicationFactor: 1},
		{Name: "vehicle.mileage.updated", NumPartitions: 3, ReplicationFactor: 1},
		{Name: "vehicle.fuel.updated", NumPartitions: 3, ReplicationFactor: 1},
	}

	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to controller: %w", err)
	}
	defer controllerConn.Close()

	controllerConn.SetDeadline(time.Now().Add(10 * time.Second))

	topicConfigs := make([]kafka.TopicConfig, len(topics))
	for i, topic := range topics {
		topicConfigs[i] = kafka.TopicConfig{
			Topic:             topic.Name,
			NumPartitions:     topic.NumPartitions,
			ReplicationFactor: topic.ReplicationFactor,
		}
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		logger.Warn("some topics may already exist", zap.Error(err))
	}

	logger.Info("kafka topics initialized", zap.Int("count", len(topics)))
	return nil
}

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
