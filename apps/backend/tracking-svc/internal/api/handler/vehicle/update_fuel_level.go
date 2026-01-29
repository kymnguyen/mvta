package vehicle

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/api/handler"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/dto"
)

func (h *VehicleHandler) UpdateFuelLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vehicleID := r.PathValue("id")

	if vehicleID == "" {
		handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_ID", "Vehicle ID is required")
		return
	}

	var req dto.UpdateVehicleFuelLevelRequest
	if err := handler.DecodeJSON(r, &req); err != nil {
		h.logger.Error("failed to decode update fuel level request", zap.Error(err))
		handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_REQUEST", "Invalid request body")
		return
	}

	cmd := &command.UpdateVehicleFuelLevelCommand{
		VehicleID: vehicleID,
		FuelLevel: req.FuelLevel,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to update vehicle fuel level",
			zap.String("vehicleId", vehicleID),
			zap.Error(err))
		handler.RespondError(w, http.StatusInternalServerError, "ERR_UPDATE_FAILED", err.Error())
		return
	}

	h.logger.Info("vehicle fuel level updated", zap.String("vehicleId", vehicleID), zap.Float64("fuelLevel", req.FuelLevel))
	handler.RespondSuccess(w, http.StatusOK, map[string]string{
		"message": "fuel level updated successfully",
	})
}
