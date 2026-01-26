package event

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
		BaseDomainEvent: InitBaseDomainEvent("vehicle.status.changed", vehicleID),
		VehicleID:       vehicleID,
		OldStatus:       oldStatus,
		NewStatus:       newStatus,
		ChangedAt:       changedAt,
		Version:         version,
	}
}
