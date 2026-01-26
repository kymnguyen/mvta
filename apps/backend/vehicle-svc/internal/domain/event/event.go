package event

import "time"

type DomainEvent interface {
	EventName() string
	Timestamp() time.Time
	AggregateID() string
}

type BaseDomainEvent struct {
	eventName   string
	timestamp   time.Time
	aggregateID string
}

func NewBaseDomainEvent(eventName, aggregateID string) BaseDomainEvent {
	return BaseDomainEvent{
		eventName:   eventName,
		timestamp:   time.Now().UTC(),
		aggregateID: aggregateID,
	}
}

func (b BaseDomainEvent) EventName() string {
	return b.eventName
}

func (b BaseDomainEvent) Timestamp() time.Time {
	return b.timestamp
}

func (b BaseDomainEvent) AggregateID() string {
	return b.aggregateID
}
