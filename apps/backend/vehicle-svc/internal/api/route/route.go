package route

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/api/handler/vehicle"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/vehicle-svc/internal/application/query"
)

func RegisterRoutes(
	mux *http.ServeMux,
	commandBus command.CommandBus,
	queryBus query.QueryBus,
	logger *zap.Logger,
) {
	h := vehicle.InitVehicleHandler(commandBus, queryBus, logger)

	mux.HandleFunc("GET /health", healthCheck)

	mux.HandleFunc("POST /api/v1/vehicles", h.CreateVehicle)
	mux.HandleFunc("GET /api/v1/vehicles", h.GetAllVehicles)
	mux.HandleFunc("GET /api/v1/vehicles/{id}", h.GetVehicle)
	mux.HandleFunc("PATCH /api/v1/vehicles/{id}/location", h.UpdateLocation)
	mux.HandleFunc("PATCH /api/v1/vehicles/{id}/status", h.ChangeStatus)
	mux.HandleFunc("PATCH /api/v1/vehicles/{id}/mileage", h.UpdateMileage)
	mux.HandleFunc("PATCH /api/v1/vehicles/{id}/fuel", h.UpdateFuelLevel)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
