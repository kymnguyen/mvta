package vehicle

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/api/handler"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/dto"
)

func (h *VehicleHandler) UpdateMileage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vehicleID := r.PathValue("id")

	if vehicleID == "" {
		handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_ID", "Vehicle ID is required")
		return
	}

	var req dto.UpdateVehicleMileageRequest
	if err := handler.DecodeJSON(r, &req); err != nil {
		h.logger.Error("failed to decode update mileage request", zap.Error(err))
		handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_REQUEST", "Invalid request body")
		return
	}

	cmd := &command.UpdateVehicleMileageCommand{
		VehicleID: vehicleID,
		Mileage:   req.Mileage,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to update vehicle mileage",
			zap.String("vehicleId", vehicleID),
			zap.Error(err))
		handler.RespondError(w, http.StatusInternalServerError, "ERR_UPDATE_FAILED", err.Error())
		return
	}

	h.logger.Info("vehicle mileage updated", zap.String("vehicleId", vehicleID), zap.Float64("mileage", req.Mileage))
	handler.RespondSuccess(w, http.StatusOK, map[string]string{
		"message": "mileage updated successfully",
	})
}
