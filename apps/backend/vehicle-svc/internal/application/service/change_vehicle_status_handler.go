package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/repository"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/domain/valueobject"
)

type ChangeVehicleStatusCommandHandler struct {
	vehicleRepo repository.VehicleRepository
	outboxRepo  repository.OutboxRepository
}

func NewChangeVehicleStatusCommandHandler(
	vehicleRepo repository.VehicleRepository,
	outboxRepo repository.OutboxRepository,
) *ChangeVehicleStatusCommandHandler {
	return &ChangeVehicleStatusCommandHandler{
		vehicleRepo: vehicleRepo,
		outboxRepo:  outboxRepo,
	}
}

func (h *ChangeVehicleStatusCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	statusCmd, ok := cmd.(*command.ChangeVehicleStatusCommand)
	if !ok {
		return fmt.Errorf("invalid command type for ChangeVehicleStatusCommandHandler")
	}

	vehicleID, err := valueobject.NewVehicleID(statusCmd.VehicleID)
	if err != nil {
		return fmt.Errorf("invalid vehicle id: %w", err)
	}

	newStatus, err := valueobject.NewVehicleStatus(statusCmd.NewStatus)
	if err != nil {
		return fmt.Errorf("invalid status: %w", err)
	}

	vehicle, err := h.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find vehicle: %w", err)
	}

	if err := vehicle.ChangeStatus(newStatus); err != nil {
		return fmt.Errorf("failed to change status: %w", err)
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
