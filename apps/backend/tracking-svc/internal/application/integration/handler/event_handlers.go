package handler

import (
	"context"
	"encoding/json"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/integration/event"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/domain/repository"
	"go.uber.org/zap"
)

type VehicleCreatedEventHandler struct {
	logger *zap.Logger
}

func NewVehicleCreatedEventHandler(logger *zap.Logger) *VehicleCreatedEventHandler {
	return &VehicleCreatedEventHandler{logger: logger}
}

func (h *VehicleCreatedEventHandler) Handle(ctx context.Context, payload []byte) error {
	var evt event.VehicleCreatedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		h.logger.Error("failed to unmarshal vehicle created event", zap.Error(err))
		return err
	}

	h.logger.Info("vehicle created event received",
		zap.String("user_id", evt.UserID),
		zap.String("role", evt.Role),
	)

	// TODO: Update vehicle service authorization cache or state
	// e.g., update user permissions in Redis, mark user as authorized

	return nil
}

type TrackingCorrectionEventHandler struct {
	vehicleRepo repository.VehicleRepository
	logger      *zap.Logger
}

func NewTrackingCorrectionEventHandler(vehicleRepo repository.VehicleRepository, logger *zap.Logger) *TrackingCorrectionEventHandler {
	return &TrackingCorrectionEventHandler{vehicleRepo: vehicleRepo, logger: logger}
}

func (h *TrackingCorrectionEventHandler) Handle(ctx context.Context, payload []byte) error {
	var evt event.TrackingCorrectionAppliedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		h.logger.Error("failed to unmarshal tracking correction event", zap.Error(err))
		return err
	}

	h.logger.Info("tracking correction event received",
		zap.String("vehicle_id", evt.VehicleID),
		zap.String("field", evt.Field),
		zap.String("old_value", evt.OldValue),
		zap.String("new_value", evt.NewValue),
	)

	// TODO: Update vehicle state based on correction
	// e.g., if mileage corrected, update vehicle record
	// if fuel level corrected, update vehicle fuel state

	return nil
}

type TrackingAlertEventHandler struct {
	logger *zap.Logger
}

func NewTrackingAlertEventHandler(logger *zap.Logger) *TrackingAlertEventHandler {
	return &TrackingAlertEventHandler{logger: logger}
}

func (h *TrackingAlertEventHandler) Handle(ctx context.Context, payload []byte) error {
	var evt event.TrackingAlertEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		h.logger.Error("failed to unmarshal tracking alert event", zap.Error(err))
		return err
	}

	h.logger.Warn("tracking alert event received",
		zap.String("vehicle_id", evt.VehicleID),
		zap.String("alert_type", evt.AlertType),
		zap.String("message", evt.Message),
	)

	// TODO: Process alert
	// e.g., update vehicle status, trigger notifications, record alert in audit log

	return nil
}
