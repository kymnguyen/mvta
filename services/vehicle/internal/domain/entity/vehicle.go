package entity

import (
	"fmt"
	"time"

	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/event"
	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/valueobject"
)

// Vehicle is the main aggregate root for the vehicle bounded context.
type Vehicle struct {
	id                valueobject.VehicleID
	vin               string
	vehicleName       string
	vehicleModel      string
	licenseNumber     valueobject.LicenseNumber
	status            valueobject.VehicleStatus
	currentLocation   valueobject.Location
	mileage           valueobject.Mileage
	fuelLevel         valueobject.FuelLevel
	version           valueobject.Version
	createdAt         time.Time
	updatedAt         time.Time
	uncommittedEvents []interface{}
}

// NewVehicle creates a new Vehicle entity.
func NewVehicle(
	id valueobject.VehicleID,
	vin string,
	vehicleName string,
	vehicleModel string,
	licenseNumber valueobject.LicenseNumber,
	status valueobject.VehicleStatus,
	location valueobject.Location,
	mileage valueobject.Mileage,
	fuelLevel valueobject.FuelLevel,
) (*Vehicle, error) {
	if vin == "" {
		return nil, fmt.Errorf("vin cannot be empty")
	}
	if vehicleName == "" {
		return nil, fmt.Errorf("vehicle name cannot be empty")
	}
	if vehicleModel == "" {
		return nil, fmt.Errorf("vehicle model cannot be empty")
	}

	now := time.Now().UTC()
	v := &Vehicle{
		id:              id,
		vin:             vin,
		vehicleName:     vehicleName,
		vehicleModel:    vehicleModel,
		licenseNumber:   licenseNumber,
		status:          status,
		currentLocation: location,
		mileage:         mileage,
		fuelLevel:       fuelLevel,
		version:         valueobject.Version{},
		createdAt:       now,
		updatedAt:       now,
	}

	// Emit domain event on creation
	v.uncommittedEvents = append(v.uncommittedEvents, &event.VehicleCreatedEvent{
		VehicleID:     id.String(),
		VIN:           vin,
		VehicleName:   vehicleName,
		VehicleModel:  vehicleModel,
		LicenseNumber: licenseNumber.String(),
		Status:        string(status),
		Latitude:      location.Latitude(),
		Longitude:     location.Longitude(),
		Mileage:       mileage.Kilometers(),
		FuelLevel:     fuelLevel.Percentage(),
		Timestamp:     now.Unix(),
	})

	return v, nil
}

// ID returns the vehicle ID.
func (v *Vehicle) ID() valueobject.VehicleID {
	return v.id
}

// VIN returns the vehicle identification number.
func (v *Vehicle) VIN() string {
	return v.vin
}

// VehicleName returns the vehicle name.
func (v *Vehicle) VehicleName() string {
	return v.vehicleName
}

// VehicleModel returns the vehicle model.
func (v *Vehicle) VehicleModel() string {
	return v.vehicleModel
}

// LicenseNumber returns the license number.
func (v *Vehicle) LicenseNumber() valueobject.LicenseNumber {
	return v.licenseNumber
}

// Status returns the current status.
func (v *Vehicle) Status() valueobject.VehicleStatus {
	return v.status
}

// CurrentLocation returns the current location.
func (v *Vehicle) CurrentLocation() valueobject.Location {
	return v.currentLocation
}

// Mileage returns the vehicle mileage.
func (v *Vehicle) Mileage() valueobject.Mileage {
	return v.mileage
}

// FuelLevel returns the fuel level.
func (v *Vehicle) FuelLevel() valueobject.FuelLevel {
	return v.fuelLevel
}

// Version returns the optimistic concurrency version.
func (v *Vehicle) Version() valueobject.Version {
	return v.version
}

// CreatedAt returns the creation timestamp.
func (v *Vehicle) CreatedAt() time.Time {
	return v.createdAt
}

// UpdatedAt returns the last update timestamp.
func (v *Vehicle) UpdatedAt() time.Time {
	return v.updatedAt
}

// UpdateLocation updates the vehicle's location and emits a domain event.
func (v *Vehicle) UpdateLocation(location valueobject.Location) error {
	if location.Equals(v.currentLocation) {
		return nil // No change
	}

	v.currentLocation = location
	v.updatedAt = time.Now().UTC()
	v.version = v.version.Next()

	v.uncommittedEvents = append(v.uncommittedEvents, &event.VehicleLocationUpdatedEvent{
		VehicleID: v.id.String(),
		Latitude:  location.Latitude(),
		Longitude: location.Longitude(),
		Altitude:  location.Altitude(),
		Timestamp: location.Timestamp(),
		UpdatedAt: v.updatedAt.Unix(),
		Version:   v.version.Value(),
	})

	return nil
}

// UpdateMileage updates the vehicle mileage and emits a domain event.
func (v *Vehicle) UpdateMileage(newMileage valueobject.Mileage) error {
	if newMileage.Equals(v.mileage) {
		return nil // No change
	}

	v.mileage = newMileage
	v.updatedAt = time.Now().UTC()
	v.version = v.version.Next()

	v.uncommittedEvents = append(v.uncommittedEvents, &event.VehicleMileageUpdatedEvent{
		VehicleID: v.id.String(),
		Mileage:   newMileage.Kilometers(),
		UpdatedAt: v.updatedAt.Unix(),
		Version:   v.version.Value(),
	})

	return nil
}

// UpdateFuelLevel updates the fuel level and emits a domain event.
func (v *Vehicle) UpdateFuelLevel(fuelLevel valueobject.FuelLevel) error {
	if fuelLevel.Equals(v.fuelLevel) {
		return nil // No change
	}

	v.fuelLevel = fuelLevel
	v.updatedAt = time.Now().UTC()
	v.version = v.version.Next()

	v.uncommittedEvents = append(v.uncommittedEvents, &event.VehicleFuelLevelUpdatedEvent{
		VehicleID: v.id.String(),
		FuelLevel: fuelLevel.Percentage(),
		IsLow:     fuelLevel.IsLow(),
		UpdatedAt: v.updatedAt.Unix(),
		Version:   v.version.Value(),
	})

	return nil
}

// ChangeStatus changes the vehicle status and emits a domain event.
func (v *Vehicle) ChangeStatus(newStatus valueobject.VehicleStatus) error {
	if newStatus == v.status {
		return nil // No change
	}

	oldStatus := v.status
	v.status = newStatus
	v.updatedAt = time.Now().UTC()
	v.version = v.version.Next()

	v.uncommittedEvents = append(v.uncommittedEvents, &event.VehicleStatusChangedEvent{
		VehicleID: v.id.String(),
		OldStatus: string(oldStatus),
		NewStatus: string(newStatus),
		ChangedAt: v.updatedAt.Unix(),
		Version:   v.version.Value(),
	})

	return nil
}

// UncommittedEvents returns and clears uncommitted domain events.
func (v *Vehicle) UncommittedEvents() []interface{} {
	events := v.uncommittedEvents
	v.uncommittedEvents = []interface{}{}
	return events
}

// LoadFromHistory reconstructs a Vehicle from persisted state.
func LoadFromHistory(
	id valueobject.VehicleID,
	vin string,
	vehicleName string,
	vehicleModel string,
	licenseNumber valueobject.LicenseNumber,
	status valueobject.VehicleStatus,
	location valueobject.Location,
	mileage valueobject.Mileage,
	fuelLevel valueobject.FuelLevel,
	version valueobject.Version,
	createdAt, updatedAt time.Time,
) *Vehicle {
	return &Vehicle{
		id:              id,
		vin:             vin,
		vehicleName:     vehicleName,
		vehicleModel:    vehicleModel,
		licenseNumber:   licenseNumber,
		status:          status,
		currentLocation: location,
		mileage:         mileage,
		fuelLevel:       fuelLevel,
		version:         version,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

// ID returns the vehicle ID.
func (v *Vehicle) ID() valueobject.VehicleID {
	return v.id
}

// VIN returns the vehicle identification number.
func (v *Vehicle) VIN() string {
	return v.vin
}

// Status returns the current status.
func (v *Vehicle) Status() valueobject.VehicleStatus {
	return v.status
}

// CurrentLocation returns the current location.
func (v *Vehicle) CurrentLocation() valueobject.Location {
	return v.currentLocation
}

// Version returns the optimistic concurrency version.
func (v *Vehicle) Version() valueobject.Version {
	return v.version
}

// CreatedAt returns the creation timestamp.
func (v *Vehicle) CreatedAt() time.Time {
	return v.createdAt
}

// UpdatedAt returns the last update timestamp.
func (v *Vehicle) UpdatedAt() time.Time {
	return v.updatedAt
}

// UpdateLocation updates the vehicle's location and emits a domain event.
func (v *Vehicle) UpdateLocation(location valueobject.Location) error {
	if location.Equals(v.currentLocation) {
		return nil // No change
	}

	v.currentLocation = location
	v.updatedAt = time.Now().UTC()
	v.version = v.version.Next()

	v.uncommittedEvents = append(v.uncommittedEvents, &event.VehicleLocationUpdatedEvent{
		VehicleID: v.id.String(),
		Latitude:  location.Latitude(),
		Longitude: location.Longitude(),
		Altitude:  location.Altitude(),
		Timestamp: location.Timestamp(),
		UpdatedAt: v.updatedAt.Unix(),
		Version:   v.version.Value(),
	})

	return nil
}

// ChangeStatus changes the vehicle status and emits a domain event.
func (v *Vehicle) ChangeStatus(newStatus valueobject.VehicleStatus) error {
	if newStatus == v.status {
		return nil // No change
	}

	oldStatus := v.status
	v.status = newStatus
	v.updatedAt = time.Now().UTC()
	v.version = v.version.Next()

	v.uncommittedEvents = append(v.uncommittedEvents, &event.VehicleStatusChangedEvent{
		VehicleID: v.id.String(),
		OldStatus: string(oldStatus),
		NewStatus: string(newStatus),
		ChangedAt: v.updatedAt.Unix(),
		Version:   v.version.Value(),
	})

	return nil
}

// UncommittedEvents returns and clears uncommitted domain events.
func (v *Vehicle) UncommittedEvents() []interface{} {
	events := v.uncommittedEvents
	v.uncommittedEvents = []interface{}{}
	return events
}

// LoadFromHistory reconstructs a Vehicle from persisted state.
func LoadFromHistory(
	id valueobject.VehicleID,
	vin string,
	status valueobject.VehicleStatus,
	location valueobject.Location,
	version valueobject.Version,
	createdAt, updatedAt time.Time,
) *Vehicle {
	return &Vehicle{
		id:              id,
		vin:             vin,
		status:          status,
		currentLocation: location,
		version:         version,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}
