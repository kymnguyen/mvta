package valueobject

import (
	"fmt"

	"github.com/google/uuid"
)

type VehicleID struct {
	value string
}

func NewVehicleID(id string) (VehicleID, error) {
	if id == "" {
		return VehicleID{}, fmt.Errorf("vehicle id cannot be empty")
	}
	if _, err := uuid.Parse(id); err != nil {
		return VehicleID{}, fmt.Errorf("invalid vehicle id format: %w", err)
	}
	return VehicleID{value: id}, nil
}

func GenerateVehicleID() VehicleID {
	return VehicleID{value: uuid.New().String()}
}

func (v VehicleID) String() string {
	return v.value
}

func (v VehicleID) Equals(other VehicleID) bool {
	return v.value == other.value
}

type VehicleStatus string

const (
	StatusActive      VehicleStatus = "active"
	StatusInactive    VehicleStatus = "inactive"
	StatusMaintenance VehicleStatus = "maintenance"
	StatusRetired     VehicleStatus = "retired"
)

func NewVehicleStatus(status string) (VehicleStatus, error) {
	s := VehicleStatus(status)
	switch s {
	case StatusActive, StatusInactive, StatusMaintenance, StatusRetired:
		return s, nil
	default:
		return "", fmt.Errorf("invalid vehicle status: %s", status)
	}
}

type Location struct {
	latitude  float64
	longitude float64
	altitude  float64
	timestamp int64
}

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

func (l Location) Latitude() float64 {
	return l.latitude
}

func (l Location) Longitude() float64 {
	return l.longitude
}

func (l Location) Altitude() float64 {
	return l.altitude
}

func (l Location) Timestamp() int64 {
	return l.timestamp
}

func (l Location) Equals(other Location) bool {
	return l.latitude == other.latitude &&
		l.longitude == other.longitude &&
		l.altitude == other.altitude &&
		l.timestamp == other.timestamp
}

type Version struct {
	value int64
}

func NewVersion(value int64) (Version, error) {
	if value < 0 {
		return Version{}, fmt.Errorf("version cannot be negative: %d", value)
	}
	return Version{value: value}, nil
}

func (v Version) Value() int64 {
	return v.value
}

func (v Version) Next() Version {
	return Version{value: v.value + 1}
}

func (v Version) Equals(other Version) bool {
	return v.value == other.value
}

type Mileage struct {
	kilometers float64
}

func NewMileage(kilometers float64) (Mileage, error) {
	if kilometers < 0 {
		return Mileage{}, fmt.Errorf("mileage cannot be negative: %f", kilometers)
	}
	return Mileage{kilometers: kilometers}, nil
}

func (m Mileage) Kilometers() float64 {
	return m.kilometers
}

func (m Mileage) AddKilometers(km float64) (Mileage, error) {
	if km < 0 {
		return Mileage{}, fmt.Errorf("cannot add negative kilometers")
	}
	return NewMileage(m.kilometers + km)
}

func (m Mileage) Equals(other Mileage) bool {
	return m.kilometers == other.kilometers
}

type FuelLevel struct {
	percentage float64
}

func NewFuelLevel(percentage float64) (FuelLevel, error) {
	if percentage < 0 || percentage > 100 {
		return FuelLevel{}, fmt.Errorf("fuel level must be between 0 and 100: %f", percentage)
	}
	return FuelLevel{percentage: percentage}, nil
}

func (f FuelLevel) Percentage() float64 {
	return f.percentage
}

func (f FuelLevel) IsLow() bool {
	return f.percentage < 15
}

func (f FuelLevel) Equals(other FuelLevel) bool {
	return f.percentage == other.percentage
}

type LicenseNumber struct {
	value string
}

func NewLicenseNumber(value string) (LicenseNumber, error) {
	if value == "" {
		return LicenseNumber{}, fmt.Errorf("license number cannot be empty")
	}
	if len(value) < 3 || len(value) > 20 {
		return LicenseNumber{}, fmt.Errorf("license number must be between 3 and 20 characters")
	}
	return LicenseNumber{value: value}, nil
}

func (l LicenseNumber) String() string {
	return l.value
}

func (l LicenseNumber) Equals(other LicenseNumber) bool {
	return l.value == other.value
}
