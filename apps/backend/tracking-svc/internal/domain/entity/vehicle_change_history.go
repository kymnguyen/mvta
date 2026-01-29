package entity

import (
	"time"
)

// VehicleChangeHistory records a single change to a vehicle
type VehicleChangeHistory struct {
	ID         string                 `bson:"_id,omitempty"`
	VehicleID  string                 `bson:"vehicleId"`
	VIN        string                 `bson:"vin"`
	ChangeType string                 `bson:"changeType"` // created, location_updated, status_changed, mileage_updated, fuel_updated
	OldValue   map[string]interface{} `bson:"oldValue"`   // Previous state
	NewValue   map[string]interface{} `bson:"newValue"`   // New state
	ChangedAt  time.Time              `bson:"changedAt"`
	Version    int64                  `bson:"version"` // Vehicle version at time of change
}

// NewVehicleChangeHistory creates a new change history record
func NewVehicleChangeHistory(
	vehicleID string,
	vin string,
	changeType string,
	oldValue map[string]interface{},
	newValue map[string]interface{},
	version int64,
) *VehicleChangeHistory {
	return &VehicleChangeHistory{
		VehicleID:  vehicleID,
		VIN:        vin,
		ChangeType: changeType,
		OldValue:   oldValue,
		NewValue:   newValue,
		ChangedAt:  time.Now().UTC(),
		Version:    version,
	}
}
