package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/valueobject"
)

type UpdateVehicleMileageCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

func NewUpdateVehicleMileageCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *UpdateVehicleMileageCommandHandler {
	return &UpdateVehicleMileageCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

func (h *UpdateVehicleMileageCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	mileageCmd, ok := cmd.(*command.UpdateVehicleMileageCommand)
	if !ok {
		return fmt.Errorf("invalid command type for UpdateVehicleMileageCommandHandler")
	}

	vehicleID, err := valueobject.NewVehicleID(mileageCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	newMileage, err := valueobject.NewMileage(mileageCmd.Mileage)
	if err != nil {
		return fmt.Errorf("invalid mileage: %w", err)
	}

	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	if err := vehicle.UpdateMileage(newMileage); err != nil {
		return fmt.Errorf("failed to update mileage: %w", err)
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
