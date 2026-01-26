package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/dto"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/query"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/repository"
)

type GetAllVehiclesQueryHandler struct {
	vehicleRepo repository.VehicleRepository
}

func NewGetAllVehiclesQueryHandler(vehicleRepo repository.VehicleRepository) *GetAllVehiclesQueryHandler {
	return &GetAllVehiclesQueryHandler{vehicleRepo: vehicleRepo}
}

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
