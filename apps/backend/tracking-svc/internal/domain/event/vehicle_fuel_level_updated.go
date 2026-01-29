package event

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
		BaseDomainEvent: InitBaseDomainEvent("vehicle.fuel.updated", vehicleID),
		VehicleID:       vehicleID,
		FuelLevel:       fuelLevel,
		IsLow:           isLow,
		UpdatedAt:       updatedAt,
		Version:         version,
	}
}
