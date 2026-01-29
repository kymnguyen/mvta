package vehicle

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/api/handler"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/query"
)

func (h *VehicleHandler) GetVehicle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vehicleID := r.PathValue("id")

	if vehicleID == "" {
		handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_ID", "Vehicle ID is required")
		return
	}

	q := &query.GetVehicleQuery{
		VehicleID: vehicleID,
	}

	result, err := h.queryBus.Dispatch(ctx, q)
	if err != nil {
		h.logger.Error("failed to get vehicle", zap.String("vehicleId", vehicleID), zap.Error(err))
		handler.RespondError(w, http.StatusNotFound, "VEHICLE_NOT_FOUND", "Vehicle not found")
		return
	}

	handler.RespondSuccess(w, http.StatusOK, result)
}
