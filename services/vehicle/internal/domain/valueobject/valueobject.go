package valueobject

import (
	"fmt"

	"github.com/google/uuid"
)

// VehicleID represents the unique identifier for a vehicle.
type VehicleID struct {
	value string
}

// NewVehicleID creates a new VehicleID from a string.
func NewVehicleID(id string) (VehicleID, error) {
	if id == "" {
		return VehicleID{}, fmt.Errorf("vehicle id cannot be empty")
	}
	if _, err := uuid.Parse(id); err != nil {
		return VehicleID{}, fmt.Errorf("invalid vehicle id format: %w", err)
	}
	return VehicleID{value: id}, nil
}

// GenerateVehicleID generates a new random VehicleID.
func GenerateVehicleID() VehicleID {
	return VehicleID{value: uuid.New().String()}
}

// String returns the string representation of VehicleID.
func (v VehicleID) String() string {
	return v.value
}

// Equals compares two VehicleIDs.
func (v VehicleID) Equals(other VehicleID) bool {
	return v.value == other.value
}

// VehicleStatus represents the operational status of a vehicle.
type VehicleStatus string

const (
	StatusActive      VehicleStatus = "active"
	StatusInactive    VehicleStatus = "inactive"
	StatusMaintenance VehicleStatus = "maintenance"
	StatusRetired     VehicleStatus = "retired"
)

// NewVehicleStatus creates a new VehicleStatus from a string.
func NewVehicleStatus(status string) (VehicleStatus, error) {
	s := VehicleStatus(status)
	switch s {
	case StatusActive, StatusInactive, StatusMaintenance, StatusRetired:
		return s, nil
	default:
		return "", fmt.Errorf("invalid vehicle status: %s", status)
	}
}

// Location represents a geographic coordinate.
type Location struct {
	latitude  float64
	longitude float64
	altitude  float64
	timestamp int64
}

// NewLocation creates a new Location value object.
func NewLocation(latitude, longitude, altitude float64, timestamp int64) (Location, error) {
	if latitude < -90 || latitude > 90 {
		return Location{}, fmt.Errorf("invalid latitude: %f", latitude)
	}
	if longitude < -180 || longitude > 180 {
		return Location{}, fmt.Errorf("invalid longitude: %f", longitude)
	}
	if timestamp < 0 {
		return Location{}, fmt.Errorf("invalid timestamp: %d", timestamp)
	}
	return Location{
		latitude:  latitude,
		longitude: longitude,
		altitude:  altitude,
		timestamp: timestamp,
	}, nil
}

// Latitude returns the latitude coordinate.
func (l Location) Latitude() float64 {
	return l.latitude
}

// Longitude returns the longitude coordinate.
func (l Location) Longitude() float64 {
	return l.longitude
}

// Altitude returns the altitude value.
func (l Location) Altitude() float64 {
	return l.altitude
}

// Timestamp returns the timestamp of the location.
func (l Location) Timestamp() int64 {
	return l.timestamp
}

// Equals compares two Location value objects.
func (l Location) Equals(other Location) bool {
	return l.latitude == other.latitude &&
		l.longitude == other.longitude &&
		l.altitude == other.altitude &&
		l.timestamp == other.timestamp
}

// Version represents the optimistic concurrency version.
type Version struct {
	value int64
}

// NewVersion creates a new Version.
func NewVersion(value int64) (Version, error) {
	if value < 0 {
		return Version{}, fmt.Errorf("version cannot be negative: %d", value)
	}
	return Version{value: value}, nil
}

// Value returns the version number.
func (v Version) Value() int64 {
	return v.value
}

// Next returns the next version.
func (v Version) Next() Version {
	return Version{value: v.value + 1}
}

// Equals compares two versions.
func (v Version) Equals(other Version) bool {
	return v.value == other.value
}

// Mileage represents the vehicle mileage in kilometers.
type Mileage struct {
	kilometers float64
}

// NewMileage creates a new Mileage value object.
func NewMileage(kilometers float64) (Mileage, error) {
	if kilometers < 0 {
		return Mileage{}, fmt.Errorf("mileage cannot be negative: %f", kilometers)
	}
	return Mileage{kilometers: kilometers}, nil
}

// Kilometers returns the mileage in kilometers.
func (m Mileage) Kilometers() float64 {
	return m.kilometers
}

// AddKilometers adds kilometers to current mileage.
func (m Mileage) AddKilometers(km float64) (Mileage, error) {
	if km < 0 {
		return Mileage{}, fmt.Errorf("cannot add negative kilometers")
	}
	return NewMileage(m.kilometers + km)
}

// Equals compares two mileage values.
func (m Mileage) Equals(other Mileage) bool {
	return m.kilometers == other.kilometers
}

// FuelLevel represents the fuel level percentage (0-100).
type FuelLevel struct {
	percentage float64
}

// NewFuelLevel creates a new FuelLevel value object.
func NewFuelLevel(percentage float64) (FuelLevel, error) {
	if percentage < 0 || percentage > 100 {
		return FuelLevel{}, fmt.Errorf("fuel level must be between 0 and 100: %f", percentage)
	}
	return FuelLevel{percentage: percentage}, nil
}

// Percentage returns the fuel level percentage.
func (f FuelLevel) Percentage() float64 {
	return f.percentage
}

// IsLow returns true if fuel level is below 15%.
func (f FuelLevel) IsLow() bool {
	return f.percentage < 15
}

// Equals compares two fuel levels.
func (f FuelLevel) Equals(other FuelLevel) bool {
	return f.percentage == other.percentage
}

// LicenseNumber represents a vehicle license plate number.
type LicenseNumber struct {
	value string
}

// NewLicenseNumber creates a new LicenseNumber value object.
func NewLicenseNumber(value string) (LicenseNumber, error) {
	if value == "" {
		return LicenseNumber{}, fmt.Errorf("license number cannot be empty")
	}
	if len(value) < 3 || len(value) > 20 {
		return LicenseNumber{}, fmt.Errorf("license number must be between 3 and 20 characters")
	}
	return LicenseNumber{value: value}, nil
}

// String returns the string representation.
func (l LicenseNumber) String() string {
	return l.value
}

// Equals compares two license numbers.
func (l LicenseNumber) Equals(other LicenseNumber) bool {
	return l.value == other.value
}
