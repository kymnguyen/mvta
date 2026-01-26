package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/services/vehicle/internal/application/dto"
	"github.com/kymnguyen/mvta/services/vehicle/internal/application/query"
	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/repository"
	"github.com/kymnguyen/mvta/services/vehicle/internal/domain/valueobject"
)

// GetVehicleQueryHandler handles queries for retrieving a single vehicle.
type GetVehicleQueryHandler struct {
	vehicleRepo repository.VehicleRepository
}

// NewGetVehicleQueryHandler creates a new query handler.
func NewGetVehicleQueryHandler(vehicleRepo repository.VehicleRepository) *GetVehicleQueryHandler {
	return &GetVehicleQueryHandler{vehicleRepo: vehicleRepo}
}

// Handle processes the get vehicle query.
func (h *GetVehicleQueryHandler) Handle(ctx context.Context, q query.Query) (query.QueryResult, error) {
	getQuery, ok := q.(*query.GetVehicleQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type for GetVehicleQueryHandler")
	}

	vehicleID, err := valueobject.NewVehicleID(getQuery.VehicleID)
	if err != nil {
		return nil, fmt.Errorf("invalid vehicle id: %w", err)
	}

	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vehicle: %w", err)
	}

	return &dto.VehicleResponse{
		ID:            vehicle.ID().String(),
		VIN:           vehicle.VIN(),
		VehicleName:   vehicle.VehicleName(),
		VehicleModel:  vehicle.VehicleModel(),
		LicenseNumber: vehicle.LicenseNumber().String(),
		Status:        string(vehicle.Status()),
		Latitude:      vehicle.CurrentLocation().Latitude(),
		Longitude:     vehicle.CurrentLocation().Longitude(),
		Altitude:      vehicle.CurrentLocation().Altitude(),
		Mileage:       vehicle.Mileage().Kilometers(),
		FuelLevel:     vehicle.FuelLevel().Percentage(),
		Version:       vehicle.Version().Value(),
		CreatedAt:     vehicle.CreatedAt(),
		UpdatedAt:     vehicle.UpdatedAt(),
	}, nil
}

// GetAllVehiclesQueryHandler handles queries for retrieving all vehicles.
type GetAllVehiclesQueryHandler struct {
	vehicleRepo repository.VehicleRepository
}

// NewGetAllVehiclesQueryHandler creates a new query handler.
func NewGetAllVehiclesQueryHandler(vehicleRepo repository.VehicleRepository) *GetAllVehiclesQueryHandler {
	return &GetAllVehiclesQueryHandler{vehicleRepo: vehicleRepo}
}

// Handle processes the get all vehicles query.
func (h *GetAllVehiclesQueryHandler) Handle(ctx context.Context, q query.Query) (query.QueryResult, error) {
	allQuery, ok := q.(*query.GetAllVehiclesQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type for GetAllVehiclesQueryHandler")
	}

	vehicles, err := h.vehicleRepo.FindAll(ctx, allQuery.Limit, allQuery.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find vehicles: %w", err)
	}

	var responses []*dto.VehicleResponse
	for _, vehicle := range vehicles {
		responses = append(responses, &dto.VehicleResponse{
			ID:            vehicle.ID().String(),
			VIN:           vehicle.VIN(),
			VehicleName:   vehicle.VehicleName(),
			VehicleModel:  vehicle.VehicleModel(),
			LicenseNumber: vehicle.LicenseNumber().String(),
			Status:        string(vehicle.Status()),
			Latitude:      vehicle.CurrentLocation().Latitude(),
			Longitude:     vehicle.CurrentLocation().Longitude(),
			Altitude:      vehicle.CurrentLocation().Altitude(),
			Mileage:       vehicle.Mileage().Kilometers(),
			FuelLevel:     vehicle.FuelLevel().Percentage(),
			Version:       vehicle.Version().Value(),
			CreatedAt:     vehicle.CreatedAt(),
			UpdatedAt:     vehicle.UpdatedAt(),
		})
	}

	return responses, nil
}
