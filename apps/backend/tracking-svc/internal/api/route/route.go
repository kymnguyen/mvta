package route

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/api/handler/vehicle"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/api/middleware"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/tracking-svc/internal/application/query"
)

func RegisterRoutes(
	mux *http.ServeMux,
	commandBus command.CommandBus,
	queryBus query.QueryBus,
	logger *zap.Logger,
) {
	h := vehicle.InitVehicleHandler(commandBus, queryBus, logger)
	authMiddleware := middleware.AuthMiddleware("")

	mux.HandleFunc("GET /health", healthCheck)

	mux.HandleFunc("GET /api/v1/vehicles", h.GetAllVehicles)
	mux.HandleFunc("GET /api/v1/vehicles/{id}", h.GetVehicle)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("GET /api/v1/admin/vehicles", h.GetAllVehicles)
	middlewareHandler := authMiddleware(adminMux)

	mux.Handle("/api/v1/admin/", middlewareHandler)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
