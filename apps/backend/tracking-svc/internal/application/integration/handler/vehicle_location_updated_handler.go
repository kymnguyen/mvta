package handler

import (
	"context"
	"encoding/json"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/integration/event"
	"go.uber.org/zap"
)

type VehicleLocationUpdatedEventHandler struct {
	commandBus command.CommandBus
	logger     *zap.Logger
}

func NewVehicleLocationUpdatedEventHandler(commandBus command.CommandBus, logger *zap.Logger) *VehicleLocationUpdatedEventHandler {
	return &VehicleLocationUpdatedEventHandler{commandBus: commandBus, logger: logger}
}

func (h *VehicleLocationUpdatedEventHandler) Handle(ctx context.Context, payload []byte) error {
	var evt event.VehicleLocationUpdatedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		h.logger.Error("failed to unmarshal vehicle location updated event", zap.Error(err))
		return err
	}

	cmd := &command.UpdateVehicleLocationCommand{
		VehicleID: evt.VehicleID,
		Latitude:  evt.Latitude,
		Longitude: evt.Longitude,
		Altitude:  evt.Altitude,
		Timestamp: evt.Timestamp,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to dispatch update vehicle location command", zap.Error(err))
		return err
	}

	return nil
}
