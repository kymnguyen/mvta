package dto

import "time"

// CreateVehicleRequest is the DTO for creating a new vehicle.
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

// CreateVehicleResponse is the response DTO after vehicle creation.
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

// UpdateVehicleLocationRequest is the DTO for updating vehicle location.
type UpdateVehicleLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Altitude  float64 `json:"altitude"`
	Timestamp int64   `json:"timestamp" binding:"required"`
}

// UpdateVehicleLocationResponse is the response DTO after location update.
type UpdateVehicleLocationResponse struct {
	ID        string    `json:"id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Version   int64     `json:"version"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ChangeVehicleStatusRequest is the DTO for changing vehicle status.
type ChangeVehicleStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// ChangeVehicleStatusResponse is the response DTO after status change.
type ChangeVehicleStatusResponse struct {
	ID        string    `json:"id"`
	OldStatus string    `json:"oldStatus"`
	NewStatus string    `json:"newStatus"`
	Version   int64     `json:"version"`
	ChangedAt time.Time `json:"changedAt"`
}

// VehicleResponse is the standard vehicle data transfer object.
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

// UpdateVehicleMileageRequest is the DTO for updating vehicle mileage.
type UpdateVehicleMileageRequest struct {
	Mileage float64 `json:"mileage" binding:"required"`
}

// UpdateVehicleMileageResponse is the response DTO after mileage update.
type UpdateVehicleMileageResponse struct {
	ID        string    `json:"id"`
	Mileage   float64   `json:"mileage"`
	Version   int64     `json:"version"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UpdateVehicleFuelLevelRequest is the DTO for updating vehicle fuel level.
type UpdateVehicleFuelLevelRequest struct {
	FuelLevel float64 `json:"fuelLevel" binding:"required"`
}

// UpdateVehicleFuelLevelResponse is the response DTO after fuel level update.
type UpdateVehicleFuelLevelResponse struct {
	ID        string    `json:"id"`
	FuelLevel float64   `json:"fuelLevel"`
	IsLow     bool      `json:"isLow"`
	Version   int64     `json:"version"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ErrorResponse is the standard error response format.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
