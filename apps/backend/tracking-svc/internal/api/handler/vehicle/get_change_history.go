package vehicle

import (
	"net/http"
	"strconv"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/api/handler"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/query"
	"go.uber.org/zap"
)

func (v *VehicleHandler) GetChangeHistory(w http.ResponseWriter, r *http.Request) {
	vehicleID := r.PathValue("id")
	if vehicleID == "" {
		handler.RespondError(w, http.StatusBadRequest, "INVALID_REQUEST", "vehicle id is required")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	q := &query.GetVehicleChangeHistoryQuery{
		VehicleID: vehicleID,
		Limit:     limit,
		Offset:    offset,
	}

	result, err := v.queryBus.Dispatch(r.Context(), q)
	if err != nil {
		v.logger.Error("failed to get vehicle change history",
			zap.String("vehicleId", vehicleID),
			zap.Error(err),
		)
		handler.RespondError(w, http.StatusInternalServerError, "QUERY_FAILED", "failed to get vehicle change history")
		return
	}

	handler.RespondSuccess(w, http.StatusOK, result)
}
