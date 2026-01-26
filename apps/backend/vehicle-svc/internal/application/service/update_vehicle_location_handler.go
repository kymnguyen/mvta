package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/valueobject"
)

type UpdateVehicleLocationCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

func NewUpdateVehicleLocationCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *UpdateVehicleLocationCommandHandler {
	return &UpdateVehicleLocationCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

func (h *UpdateVehicleLocationCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	updateCmd, ok := cmd.(*command.UpdateVehicleLocationCommand)
	if !ok {
		return fmt.Errorf("invalid command type for UpdateVehicleLocationCommandHandler")
	}

	vehicleID, err := valueobject.NewVehicleID(updateCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	location, err := valueobject.NewLocation(
		updateCmd.Latitude,
		updateCmd.Longitude,
		updateCmd.Altitude,
		updateCmd.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("invalid location: %w", err)
	}

	if err := vehicle.UpdateLocation(location); err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	if err := h.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to save vehicle: %w", err)
	}

	for _, event := range vehicle.UncommittedEvents() {
		if err := h.outboxRepo.SaveOutboxEvent(ctx, vehicleID.String(), event); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}
	}

	return nil
}
