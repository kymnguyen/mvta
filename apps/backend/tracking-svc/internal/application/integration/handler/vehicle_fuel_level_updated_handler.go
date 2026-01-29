package handler

import (
	"context"
	"encoding/json"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/integration/event"
	"go.uber.org/zap"
)

type VehicleFuelLevelUpdatedEventHandler struct {
	commandBus command.CommandBus
	logger     *zap.Logger
}

func NewVehicleFuelLevelUpdatedEventHandler(commandBus command.CommandBus, logger *zap.Logger) *VehicleFuelLevelUpdatedEventHandler {
	return &VehicleFuelLevelUpdatedEventHandler{commandBus: commandBus, logger: logger}
}

func (h *VehicleFuelLevelUpdatedEventHandler) Handle(ctx context.Context, payload []byte) error {
	var evt event.VehicleFuelLevelUpdatedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		h.logger.Error("failed to unmarshal vehicle fuel level updated event", zap.Error(err))
		return err
	}

	cmd := &command.UpdateVehicleFuelLevelCommand{
		VehicleID: evt.VehicleID,
		FuelLevel: evt.FuelLevel,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to dispatch update vehicle fuel level command", zap.Error(err))
		return err
	}

	return nil
}
