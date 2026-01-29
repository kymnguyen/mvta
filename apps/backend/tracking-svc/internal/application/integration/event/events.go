package event

type VehicleCreatedEvent struct {
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

type VehicleLocationUpdatedEvent struct {
	VehicleID string  `json:"vehicleId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Timestamp int64   `json:"timestamp"`
	UpdatedAt int64   `json:"updatedAt"`
	Version   int64   `json:"version"`
}

type VehicleStatusChangedEvent struct {
	VehicleID string `json:"vehicleId"`
	OldStatus string `json:"oldStatus"`
	NewStatus string `json:"newStatus"`
	ChangedAt int64  `json:"changedAt"`
	Version   int64  `json:"version"`
}

type VehicleMileageUpdatedEvent struct {
	VehicleID string  `json:"vehicleId"`
	Mileage   float64 `json:"mileage"`
	UpdatedAt int64   `json:"updatedAt"`
	Version   int64   `json:"version"`
}

type VehicleFuelLevelUpdatedEvent struct {
	VehicleID string  `json:"vehicleId"`
	FuelLevel float64 `json:"fuelLevel"`
	IsLow     bool    `json:"isLow"`
	UpdatedAt int64   `json:"updatedAt"`
	Version   int64   `json:"version"`
}

type TrackingCorrectionAppliedEvent struct {
	VehicleID string `json:"vehicle_id"`
	Field     string `json:"field"` // e.g., "mileage", "fuel_level"
	OldValue  string `json:"old_value"`
	NewValue  string `json:"new_value"`
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}

type TrackingAlertEvent struct {
	VehicleID string `json:"vehicle_id"`
	AlertType string `json:"alert_type"` // e.g., "low_fuel", "high_mileage"
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
