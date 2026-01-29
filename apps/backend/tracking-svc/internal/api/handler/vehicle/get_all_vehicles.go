package vehicle

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/api/handler"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/query"
)

func (h *VehicleHandler) GetAllVehicles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if _, err := handler.ScanInt(l, &limit); err != nil {
			handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_LIMIT", "Invalid limit parameter")
			return
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if _, err := handler.ScanInt(o, &offset); err != nil {
			handler.RespondError(w, http.StatusBadRequest, "ERR_INVALID_OFFSET", "Invalid offset parameter")
			return
		}
	}

	q := &query.GetAllVehiclesQuery{
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.queryBus.Dispatch(ctx, q)
	if err != nil {
		h.logger.Error("failed to get all vehicles", zap.Error(err))
		handler.RespondError(w, http.StatusInternalServerError, "ERR_QUERY_FAILED", err.Error())
		return
	}

	handler.RespondSuccess(w, http.StatusOK, map[string]interface{}{
		"vehicles": result,
		"limit":    limit,
		"offset":   offset,
	})
}
