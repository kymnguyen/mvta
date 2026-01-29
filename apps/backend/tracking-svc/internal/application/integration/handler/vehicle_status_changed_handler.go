package handler

import (
	"context"
	"encoding/json"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/integration/event"
	"go.uber.org/zap"
)

type VehicleStatusChangedEventHandler struct {
	commandBus command.CommandBus
	logger     *zap.Logger
}

func NewVehicleStatusChangedEventHandler(commandBus command.CommandBus, logger *zap.Logger) *VehicleStatusChangedEventHandler {
	return &VehicleStatusChangedEventHandler{commandBus: commandBus, logger: logger}
}

func (h *VehicleStatusChangedEventHandler) Handle(ctx context.Context, payload []byte) error {
	var evt event.VehicleStatusChangedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		h.logger.Error("failed to unmarshal vehicle status changed event", zap.Error(err))
		return err
	}

	changeCmd := &command.RecordVehicleChangeCommand{
		VehicleID:  evt.VehicleID,
		VIN:        "",
		ChangeType: "status_changed",
		OldValue: map[string]interface{}{
			"status": evt.OldStatus,
		},
		NewValue: map[string]interface{}{
			"status": evt.NewStatus,
		},
		Version: evt.Version,
	}

	if err := h.commandBus.Dispatch(ctx, changeCmd); err != nil {
		h.logger.Error("failed to record vehicle change history", zap.Error(err))
	}

	return nil
}
