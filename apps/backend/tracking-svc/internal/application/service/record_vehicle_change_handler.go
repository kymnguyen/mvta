package service

import (
	"context"
	"fmt"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/entity"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/repository"
)

type RecordVehicleChangeCommandHandler struct {
	changeHistoryRepo repository.VehicleChangeHistoryRepository
}

func NewRecordVehicleChangeCommandHandler(changeHistoryRepo repository.VehicleChangeHistoryRepository) *RecordVehicleChangeCommandHandler {
	return &RecordVehicleChangeCommandHandler{changeHistoryRepo: changeHistoryRepo}
}

func (h *RecordVehicleChangeCommandHandler) Handle(ctx context.Context, cmd command.Command) error {
	recordCmd, ok := cmd.(*command.RecordVehicleChangeCommand)
	if !ok {
		return fmt.Errorf("invalid command type for RecordVehicleChangeCommandHandler")
	}

	history := entity.NewVehicleChangeHistory(
		recordCmd.VehicleID,
		recordCmd.VIN,
		recordCmd.ChangeType,
		recordCmd.OldValue,
		recordCmd.NewValue,
		recordCmd.Version,
	)

	if err := h.changeHistoryRepo.Save(ctx, history); err != nil {
		return fmt.Errorf("failed to save vehicle change history: %w", err)
	}

	return nil
}
