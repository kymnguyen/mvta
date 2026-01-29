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

	cmd := &command.CreateVehicleCommand{
		VIN:           evt.VIN,
		VehicleName:   evt.VehicleName,
		VehicleModel:  evt.VehicleModel,
		LicenseNumber: evt.LicenseNumber,
		Status:        evt.Status,
		Latitude:      evt.Latitude,
		Longitude:     evt.Longitude,
		Altitude:      0,
		Mileage:       evt.Mileage,
		FuelLevel:     evt.FuelLevel,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to dispatch create vehicle command", zap.Error(err))
		return err
	}

	return nil
}
