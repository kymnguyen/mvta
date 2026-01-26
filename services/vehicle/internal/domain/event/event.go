package event

import "time"

// DomainEvent is the base interface for all domain events.
type DomainEvent interface {
	EventName() string
	Timestamp() time.Time
	AggregateID() string
}

// BaseDomainEvent provides common functionality for all events.
type BaseDomainEvent struct {
	eventName   string
	timestamp   time.Time
	aggregateID string
}

// NewBaseDomainEvent creates a new base domain event.
func NewBaseDomainEvent(eventName, aggregateID string) BaseDomainEvent {
	return BaseDomainEvent{
		eventName:   eventName,
		timestamp:   time.Now().UTC(),
		aggregateID: aggregateID,
	}
}

// EventName returns the event name.
func (b BaseDomainEvent) EventName() string {
	return b.eventName
}

// Timestamp returns the event timestamp.
func (b BaseDomainEvent) Timestamp() time.Time {
	return b.timestamp
}

// AggregateID returns the aggregate ID.
func (b BaseDomainEvent) AggregateID() string {
	return b.aggregateID
}

// VehicleCreatedEvent is published when a new vehicle is created.
type VehicleCreatedEvent struct {
	BaseDomainEvent
	VehicleID     string  `json:"vehicleId"`
	VIN           string  `json:"vin"`
	VehicleName   string  `json:"vehicleName"`
	VehicleModel  string  `json:"vehicleModel"`
	LicenseNumber string  `json:"licenseNumber"`
	Status        string  `json:"status"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Mileage       float64 `json:"mileage"`
	FuelLevel     float64 `json:"fuelLevel"`
	Timestamp     int64   `json:"timestamp"`
}

// NewVehicleCreatedEvent creates a new VehicleCreatedEvent.
func NewVehicleCreatedEvent(vehicleID, vin, vehicleName, vehicleModel, licenseNumber, status string, latitude, longitude float64, mileage, fuelLevel float64, timestamp int64) *VehicleCreatedEvent {
	return &VehicleCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("vehicle.created", vehicleID),
		VehicleID:       vehicleID,
		VIN:             vin,
		VehicleName:     vehicleName,
		VehicleModel:    vehicleModel,
		LicenseNumber:   licenseNumber,
		Status:          status,
		Latitude:        latitude,
		Longitude:       longitude,
		Mileage:         mileage,
		FuelLevel:       fuelLevel,
		Timestamp:       timestamp,
	}
}

// VehicleLocationUpdatedEvent is published when vehicle location is updated.
type VehicleLocationUpdatedEvent struct {
	BaseDomainEvent
	VehicleID string  `json:"vehicleId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Timestamp int64   `json:"timestamp"`
	UpdatedAt int64   `json:"updatedAt"`
	Version   int64   `json:"version"`
}

// NewVehicleLocationUpdatedEvent creates a new VehicleLocationUpdatedEvent.
func NewVehicleLocationUpdatedEvent(vehicleID string, latitude, longitude, altitude float64, timestamp, updatedAt, version int64) *VehicleLocationUpdatedEvent {
	return &VehicleLocationUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("vehicle.location.updated", vehicleID),
		VehicleID:       vehicleID,
		Latitude:        latitude,
		Longitude:       longitude,
		Altitude:        altitude,
		Timestamp:       timestamp,
		UpdatedAt:       updatedAt,
		Version:         version,
	}
}

// VehicleStatusChangedEvent is published when vehicle status changes.
type VehicleStatusChangedEvent struct {
	BaseDomainEvent
	VehicleID string `json:"vehicleId"`
	OldStatus string `json:"oldStatus"`
	NewStatus string `json:"newStatus"`
	ChangedAt int64  `json:"changedAt"`
	Version   int64  `json:"version"`
}

// NewVehicleStatusChangedEvent creates a new VehicleStatusChangedEvent.
func NewVehicleStatusChangedEvent(vehicleID, oldStatus, newStatus string, changedAt, version int64) *VehicleStatusChangedEvent {
	return &VehicleStatusChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent("vehicle.status.changed", vehicleID),
		VehicleID:       vehicleID,
		OldStatus:       oldStatus,
		NewStatus:       newStatus,
		ChangedAt:       changedAt,
		Version:         version,
	}
}

// VehicleMileageUpdatedEvent is published when vehicle mileage is updated.
type VehicleMileageUpdatedEvent struct {
	BaseDomainEvent
	VehicleID string  `json:"vehicleId"`
	Mileage   float64 `json:"mileage"`
	UpdatedAt int64   `json:"updatedAt"`
	Version   int64   `json:"version"`
}

// NewVehicleMileageUpdatedEvent creates a new VehicleMileageUpdatedEvent.
func NewVehicleMileageUpdatedEvent(vehicleID string, mileage float64, updatedAt, version int64) *VehicleMileageUpdatedEvent {
	return &VehicleMileageUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("vehicle.mileage.updated", vehicleID),
		VehicleID:       vehicleID,
		Mileage:         mileage,
		UpdatedAt:       updatedAt,
		Version:         version,
	}
}

// VehicleFuelLevelUpdatedEvent is published when vehicle fuel level is updated.
type VehicleFuelLevelUpdatedEvent struct {
	BaseDomainEvent
	VehicleID string  `json:"vehicleId"`
	FuelLevel float64 `json:"fuelLevel"`
	IsLow     bool    `json:"isLow"`
	UpdatedAt int64   `json:"updatedAt"`
	Version   int64   `json:"version"`
}

// NewVehicleFuelLevelUpdatedEvent creates a new VehicleFuelLevelUpdatedEvent.
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
