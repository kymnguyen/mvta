package handler

import (
	"context"
	"encoding/json"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/integration/event"
	"go.uber.org/zap"
)

type VehicleCreatedEventHandler struct {
	commandBus command.CommandBus
	logger     *zap.Logger
}

func NewVehicleCreatedEventHandler(commandBus command.CommandBus, logger *zap.Logger) *VehicleCreatedEventHandler {
	return &VehicleCreatedEventHandler{commandBus: commandBus, logger: logger}
}

func (h *VehicleCreatedEventHandler) Handle(ctx context.Context, payload []byte) error {
	var evt event.VehicleCreatedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		h.logger.Error("failed to unmarshal vehicle created event", zap.Error(err))
		return err
	}

	h.logger.Info("vehicle created event received",
		zap.String("vehicle_id", evt.VehicleID),
		zap.String("vin", evt.VIN),
	)

	changeCmd := &command.RecordVehicleChangeCommand{
		VehicleID:  evt.VehicleID,
		VIN:        evt.VIN,
		ChangeType: "created",
		OldValue:   map[string]interface{}{},
		NewValue: map[string]interface{}{
			"vin":          evt.VIN,
			"vehicleName":  evt.VehicleName,
			"vehicleModel": evt.VehicleModel,
			"status":       evt.Status,
			"latitude":     evt.Latitude,
			"longitude":    evt.Longitude,
			"mileage":      evt.Mileage,
			"fuelLevel":    evt.FuelLevel,
		},
		Version: 1,
	}

	if err := h.commandBus.Dispatch(ctx, changeCmd); err != nil {
		h.logger.Error("failed to record vehicle change history", zap.Error(err))
	}

	return nil
}
