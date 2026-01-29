package handler

import (
	"context"
	"encoding/json"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/integration/event"
	"go.uber.org/zap"
)

type VehicleMileageUpdatedEventHandler struct {
	commandBus command.CommandBus
	logger     *zap.Logger
}

func NewVehicleMileageUpdatedEventHandler(commandBus command.CommandBus, logger *zap.Logger) *VehicleMileageUpdatedEventHandler {
	return &VehicleMileageUpdatedEventHandler{commandBus: commandBus, logger: logger}
}

func (h *VehicleMileageUpdatedEventHandler) Handle(ctx context.Context, payload []byte) error {
	var evt event.VehicleMileageUpdatedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		h.logger.Error("failed to unmarshal vehicle mileage updated event", zap.Error(err))
		return err
	}

	cmd := &command.UpdateVehicleMileageCommand{
		VehicleID: evt.VehicleID,
		Mileage:   evt.Mileage,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to dispatch update vehicle mileage command", zap.Error(err))
		return err
	}

	// Record change history
	changeCmd := &command.RecordVehicleChangeCommand{
		VehicleID:  evt.VehicleID,
		VIN:        "",
		ChangeType: "mileage_updated",
		OldValue:   map[string]interface{}{},
		NewValue: map[string]interface{}{
			"mileage": evt.Mileage,
		},
		Version: evt.Version,
	}

	if err := h.commandBus.Dispatch(ctx, changeCmd); err != nil {
		h.logger.Error("failed to record vehicle change history", zap.Error(err))
		// Don't return error - change history is non-critical
	}

	return nil
}
