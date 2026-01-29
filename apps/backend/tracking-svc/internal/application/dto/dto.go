package dto

import "time"

type CreateVehicleRequest struct {
	VIN           string  `json:"vin" binding:"required"`
	VehicleName   string  `json:"vehicleName" binding:"required"`
	VehicleModel  string  `json:"vehicleModel" binding:"required"`
	LicenseNumber string  `json:"licenseNumber" binding:"required"`
	Status        string  `json:"status" binding:"required"`
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
	Altitude      float64 `json:"altitude"`
	Mileage       float64 `json:"mileage"`
	FuelLevel     float64 `json:"fuelLevel"`
}

type CreateVehicleResponse struct {
	ID            string    `json:"id"`
	VIN           string    `json:"vin"`
	VehicleName   string    `json:"vehicleName"`
	VehicleModel  string    `json:"vehicleModel"`
	LicenseNumber string    `json:"licenseNumber"`
	Status        string    `json:"status"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Altitude      float64   `json:"altitude"`
	Mileage       float64   `json:"mileage"`
	FuelLevel     float64   `json:"fuelLevel"`
	CreatedAt     time.Time `json:"createdAt"`
}

type UpdateVehicleLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Altitude  float64 `json:"altitude"`
	Timestamp int64   `json:"timestamp" binding:"required"`
}

type UpdateVehicleLocationResponse struct {
	ID        string    `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Version   int64     `json:"version"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ChangeVehicleStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type ChangeVehicleStatusResponse struct {
	ID        string    `json:"id"`
	OldStatus string    `json:"oldStatus"`
	NewStatus string    `json:"newStatus"`
	Version   int64     `json:"version"`
	ChangedAt time.Time `json:"changedAt"`
}

type VehicleResponse struct {
	ID            string    `json:"id"`
	VIN           string    `json:"vin"`
	VehicleName   string    `json:"vehicleName"`
	VehicleModel  string    `json:"vehicleModel"`
	LicenseNumber string    `json:"licenseNumber"`
	Status        string    `json:"status"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Altitude      float64   `json:"altitude"`
	Mileage       float64   `json:"mileage"`
	FuelLevel     float64   `json:"fuelLevel"`
	Version       int64     `json:"version"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type UpdateVehicleMileageRequest struct {
	Mileage float64 `json:"mileage" binding:"required"`
}

type UpdateVehicleMileageResponse struct {
	ID        string    `json:"id"`
	Mileage   float64   `json:"mileage"`
	Version   int64     `json:"version"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateVehicleFuelLevelRequest struct {
	FuelLevel float64 `json:"fuelLevel" binding:"required"`
}

type UpdateVehicleFuelLevelResponse struct {
	ID        string    `json:"id"`
	FuelLevel float64   `json:"fuelLevel"`
	IsLow     bool      `json:"isLow"`
	Version   int64     `json:"version"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type VehicleChangeHistoryResponse struct {
	VehicleID string                `json:"vehicleId"`
	Changes   []VehicleChangeRecord `json:"changes"`
	Total     int                   `json:"total"`
}

type VehicleChangeRecord struct {
	VehicleID  string                 `json:"vehicleId"`
	VIN        string                 `json:"vin"`
	ChangeType string                 `json:"changeType"`
	OldValue   map[string]interface{} `json:"oldValue"`
	NewValue   map[string]interface{} `json:"newValue"`
	ChangedAt  string                 `json:"changedAt"`
	Version    int64                  `json:"version"`
}
