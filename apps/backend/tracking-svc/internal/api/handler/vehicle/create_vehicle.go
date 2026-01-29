package vehicle

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/api/handler"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/dto"
)

func (h *VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateVehicleRequest
	if err := handler.DecodeJSON(r, &req); err != nil {
		h.logger.Error("failed to decode create vehicle request", zap.Error(err))
		handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_REQUEST", "Invalid request body")
		return
	}

	cmd := &command.CreateVehicleCommand{
		VIN:           req.VIN,
		VehicleName:   req.VehicleName,
		VehicleModel:  req.VehicleModel,
		LicenseNumber: req.LicenseNumber,
		Status:        req.Status,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		Altitude:      req.Altitude,
		Mileage:       req.Mileage,
		FuelLevel:     req.FuelLevel,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to create vehicle", zap.Error(err))
		handler.RespondError(w, http.StatusInternalServerError, "ERR_CREATE_FAILED", err.Error())
		return
	}

	h.logger.Info("vehicle created successfully", zap.String("vin", req.VIN))
	handler.RespondSuccess(w, http.StatusCreated, map[string]string{
		"message": "vehicle created successfully",
	})
}
