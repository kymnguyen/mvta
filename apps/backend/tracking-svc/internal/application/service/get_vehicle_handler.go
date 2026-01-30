package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/dto"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/query"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/valueobject"
)

type GetVehicleQueryHandler struct {
	vehicleRepo repository.VehicleRepository
}

func NewGetVehicleQueryHandler(vehicleRepo repository.VehicleRepository) *GetVehicleQueryHandler {
	return &GetVehicleQueryHandler{vehicleRepo: vehicleRepo}
}

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
		RefID:         vehicle.RefID(),
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
