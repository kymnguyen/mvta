package vehicle

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/api/handler"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/dto"
)

func (h *VehicleHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vehicleID := r.PathValue("id")

	if vehicleID == "" {
		handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_ID", "Vehicle ID is required")
		return
	}

	var req dto.ChangeVehicleStatusRequest
	if err := handler.DecodeJSON(r, &req); err != nil {
		h.logger.Error("failed to decode change status request", zap.Error(err))
		handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_REQUEST", "Invalid request body")
		return
	}

	cmd := &command.ChangeVehicleStatusCommand{
		VehicleID: vehicleID,
		NewStatus: req.Status,
	}

	if err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("failed to change vehicle status",
			zap.String("vehicleId", vehicleID),
			zap.Error(err))
		handler.RespondError(w, http.StatusInternalServerError, "ERR_UPDATE_FAILED", err.Error())
		return
	}

	h.logger.Info("vehicle status changed", zap.String("vehicleId", vehicleID))
	handler.RespondSuccess(w, http.StatusOK, map[string]string{
		"message": "status changed successfully",
	})
}
