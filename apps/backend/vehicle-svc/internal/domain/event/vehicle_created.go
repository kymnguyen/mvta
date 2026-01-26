package event

type VehicleCreatedEvent struct {
	BaseDomainEvent
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

func NewVehicleCreatedEvent(vehicleID, vin, vehicleName, vehicleModel, licenseNumber, status string, latitude, longitude float64, mileage, fuelLevel float64, timestamp int64) *VehicleCreatedEvent {
	return &VehicleCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("vehicle.created", vehicleID),
		VehicleID:       vehicleID,
		VIN:             vin,
		VehicleName:     vehicleName,
		VehicleModel:    vehicleModel,
		LicenseNumber:   licenseNumber,
		Status:          status,
		Latitude:        latitude,
		Longitude:       longitude,
		Mileage:         mileage,
		FuelLevel:       fuelLevel,
		Timestamp:       timestamp,
	}
}
