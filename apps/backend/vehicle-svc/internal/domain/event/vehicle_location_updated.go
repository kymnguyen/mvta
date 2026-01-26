package event

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
