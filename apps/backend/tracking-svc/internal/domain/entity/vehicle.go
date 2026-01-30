package entity

import (
	"fmt"
	"time"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/event"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/valueobject"
)

type Vehicle struct {
	id                valueobject.VehicleID
	refId             string
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

func NewVehicle(
	id valueobject.VehicleID,
	refId string,
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
		refId:           refId,
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

	v.uncommittedEvents = append(v.uncommittedEvents, &event.VehicleCreatedEvent{
		VehicleID:     id.String(),
		VIN:           vin,
		VehicleName:   vehicleName,
		VehicleModel:  vehicleModel,
		LicenseNumber: licenseNumber.String(),
		Status:        string(status),
		Latitude:      location.Latitude(),
		Longitude:     location.Longitude(),
		Mileage:       float64(mileage.Kilometers()),
		FuelLevel:     float64(fuelLevel.Percentage()),
		Timestamp:     now.Unix(),
	})

	return v, nil
}

func (v *Vehicle) ID() valueobject.VehicleID {
	return v.id
}

func (v *Vehicle) RefID() string {
	return v.refId
}

func (v *Vehicle) VIN() string {
	return v.vin
}

func (v *Vehicle) VehicleName() string {
	return v.vehicleName
}

func (v *Vehicle) VehicleModel() string {
	return v.vehicleModel
}

func (v *Vehicle) LicenseNumber() valueobject.LicenseNumber {
	return v.licenseNumber
}

func (v *Vehicle) Status() valueobject.VehicleStatus {
	return v.status
}

func (v *Vehicle) CurrentLocation() valueobject.Location {
	return v.currentLocation
}

func (v *Vehicle) Mileage() valueobject.Mileage {
	return v.mileage
}

func (v *Vehicle) FuelLevel() valueobject.FuelLevel {
	return v.fuelLevel
}

func (v *Vehicle) Version() valueobject.Version {
	return v.version
}

func (v *Vehicle) CreatedAt() time.Time {
	return v.createdAt
}

func (v *Vehicle) UpdatedAt() time.Time {
	return v.updatedAt
}

func (v *Vehicle) UpdateLocation(location valueobject.Location) error {
	if location.Equals(v.currentLocation) {
		return nil
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

func (v *Vehicle) UpdateMileage(newMileage valueobject.Mileage) error {
	if newMileage.Equals(v.mileage) {
		return nil
	}

	v.mileage = newMileage
	v.updatedAt = time.Now().UTC()
	v.version = v.version.Next()

	v.uncommittedEvents = append(v.uncommittedEvents, &event.VehicleMileageUpdatedEvent{
		VehicleID: v.id.String(),
		Mileage:   float64(newMileage.Kilometers()),
		UpdatedAt: v.updatedAt.Unix(),
		Version:   v.version.Value(),
	})

	return nil
}

func (v *Vehicle) UpdateFuelLevel(fuelLevel valueobject.FuelLevel) error {
	if fuelLevel.Equals(v.fuelLevel) {
		return nil // No change
	}

	v.fuelLevel = fuelLevel
	v.updatedAt = time.Now().UTC()
	v.version = v.version.Next()

	v.uncommittedEvents = append(v.uncommittedEvents, &event.VehicleFuelLevelUpdatedEvent{
		VehicleID: v.id.String(),
		FuelLevel: float64(fuelLevel.Percentage()),
		IsLow:     fuelLevel.IsLow(),
		UpdatedAt: v.updatedAt.Unix(),
		Version:   v.version.Value(),
	})

	return nil
}

func (v *Vehicle) ChangeStatus(newStatus valueobject.VehicleStatus) error {
	if newStatus == v.status {
		return nil
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

func (v *Vehicle) UncommittedEvents() []interface{} {
	events := v.uncommittedEvents
	v.uncommittedEvents = []interface{}{}
	return events
}

func LoadFromHistory(
	id valueobject.VehicleID,
	refId string,
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
		refId:           refId,
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
