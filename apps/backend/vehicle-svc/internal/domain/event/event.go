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

type VehicleMileageUpdatedEvent struct {
	BaseDomainEvent
	VehicleID string  `json:"vehicleId"`
	Mileage   float64 `json:"mileage"`
	UpdatedAt int64   `json:"updatedAt"`
	Version   int64   `json:"version"`
}

func NewVehicleMileageUpdatedEvent(vehicleID string, mileage float64, updatedAt, version int64) *VehicleMileageUpdatedEvent {
	return &VehicleMileageUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("vehicle.mileage.updated", vehicleID),
		VehicleID:       vehicleID,
		Mileage:         mileage,
		UpdatedAt:       updatedAt,
		Version:         version,
	}
}

type VehicleFuelLevelUpdatedEvent struct {
	BaseDomainEvent
	VehicleID string  `json:"vehicleId"`
	FuelLevel float64 `json:"fuelLevel"`
	IsLow     bool    `json:"isLow"`
	UpdatedAt int64   `json:"updatedAt"`
	Version   int64   `json:"version"`
}

func NewVehicleFuelLevelUpdatedEvent(vehicleID string, fuelLevel float64, isLow bool, updatedAt, version int64) *VehicleFuelLevelUpdatedEvent {
	return &VehicleFuelLevelUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("vehicle.fuel.updated", vehicleID),
		VehicleID:       vehicleID,
		FuelLevel:       fuelLevel,
		IsLow:           isLow,
		UpdatedAt:       updatedAt,
		Version:         version,
	}
}
