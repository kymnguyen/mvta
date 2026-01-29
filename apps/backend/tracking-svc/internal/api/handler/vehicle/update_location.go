package vehicle

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/api/handler"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/dto"
)

func (h *VehicleHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vehicleID := r.PathValue("id")

	if vehicleID == "" {
		handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_ID", "Vehicle ID is required")
		return
	}

	var req dto.UpdateVehicleLocationRequest
	if err := handler.DecodeJSON(r, &req); err != nil {
		h.logger.Error("failed to decode update location request", zap.Error(err))
		handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_REQUEST", "Invalid request body")
		return
	}

	cmd := &command.UpdateVehicleLocationCommand{
		VehicleID: vehicleID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Altitude:  req.Altitude,
		Timestamp: req.Timestamp,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to update vehicle location",
			zap.String("vehicleId", vehicleID),
			zap.Error(err))
		handler.RespondError(w, http.StatusInternalServerError, "ERR_UPDATE_FAILED", err.Error())
		return
	}

	h.logger.Info("vehicle location updated", zap.String("vehicleId", vehicleID))
	handler.RespondSuccess(w, http.StatusOK, map[string]string{
		"message": "location updated successfully",
	})
}
