package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/valueobject"
)

type UpdateVehicleFuelLevelCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

func NewUpdateVehicleFuelLevelCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *UpdateVehicleFuelLevelCommandHandler {
	return &UpdateVehicleFuelLevelCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

func (h *UpdateVehicleFuelLevelCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	fuelCmd, ok := cmd.(*command.UpdateVehicleFuelLevelCommand)
	if !ok {
		return fmt.Errorf("invalid command type for UpdateVehicleFuelLevelCommandHandler")
	}

	vehicleID, err := valueobject.NewVehicleID(fuelCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	newFuelLevel, err := valueobject.NewFuelLevel(fuelCmd.FuelLevel)
	if err != nil {
		return fmt.Errorf("invalid fuel level: %w", err)
	}

	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	if err := vehicle.UpdateFuelLevel(newFuelLevel); err != nil {
		return fmt.Errorf("failed to update fuel level: %w", err)
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
