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

type VehicleStatusChangedEvent struct {
	BaseDomainEvent
	VehicleID string `json:"vehicleId"`
	OldStatus string `json:"oldStatus"`
	NewStatus string `json:"newStatus"`
	ChangedAt int64  `json:"changedAt"`
	Version   int64  `json:"version"`
}

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
