package messaging

import (
	"context"
	"encoding/json"

	"workflow-svc/internal/domain/workflow"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TransitionEvent struct {
	InstanceID    string                 `json:"instance_id"`
	WorkflowName  string                 `json:"workflow_name"`
	CorrelationID string                 `json:"correlation_id"`
	FromState     string                 `json:"from_state"`
	ToState       string                 `json:"to_state"`
	TriggerType   string                 `json:"trigger_type"`
	TriggerName   string                 `json:"trigger_name"`
	Context       map[string]interface{} `json:"context"`
}

type KafkaTransitionPublisher struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func NewKafkaTransitionPublisher(brokers []string, topic string, logger *zap.Logger) *KafkaTransitionPublisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}

	return &KafkaTransitionPublisher{
		writer: writer,
		logger: logger,
	}
}

func (p *KafkaTransitionPublisher) Publish(ctx context.Context, event TransitionEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(event.CorrelationID),
		Value: data,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		p.logger.Error("Failed to publish transition event",
			zap.String("instance_id", event.InstanceID),
			zap.Error(err))
		return err
	}

	p.logger.Info("Published transition event",
		zap.String("instance_id", event.InstanceID),
		zap.String("workflow", event.WorkflowName),
		zap.String("from", event.FromState),
		zap.String("to", event.ToState))

	return nil
}

func (p *KafkaTransitionPublisher) OnTransition(ctx context.Context, instance *workflow.WorkflowInstance, transition workflow.StateTransition) error {
	event := TransitionEvent{
		InstanceID:    instance.ID,
		WorkflowName:  instance.WorkflowName,
		CorrelationID: instance.CorrelationID,
		FromState:     transition.FromState,
		ToState:       transition.ToState,
		TriggerType:   transition.Trigger.Type,
		TriggerName:   transition.Trigger.Name,
		Context:       transition.Context,
	}
	return p.Publish(ctx, event)
}

func (p *KafkaTransitionPublisher) Close() error {
	return p.writer.Close()
}
