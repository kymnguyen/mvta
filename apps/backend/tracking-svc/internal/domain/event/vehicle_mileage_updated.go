package event

type VehicleMileageUpdatedEvent struct {
	BaseDomainEvent
	VehicleID string  `json:"vehicleId"`
	Mileage   float64 `json:"mileage"`
	UpdatedAt int64   `json:"updatedAt"`
	Version   int64   `json:"version"`
}

func NewVehicleMileageUpdatedEvent(vehicleID string, mileage float64, updatedAt, version int64) *VehicleMileageUpdatedEvent {
	return &VehicleMileageUpdatedEvent{
		BaseDomainEvent: InitBaseDomainEvent("vehicle.mileage.updated", vehicleID),
		VehicleID:       vehicleID,
		Mileage:         mileage,
		UpdatedAt:       updatedAt,
		Version:         version,
	}
}
