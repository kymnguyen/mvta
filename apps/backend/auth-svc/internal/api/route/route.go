package route

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/api/handler"
	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/application/service"
)

func RegisterRoutes(
	mux *http.ServeMux,
	loginHandler *service.LoginHandler,
	registerHandler *service.RegisterUserHandler,
	logger *zap.Logger,
) {
	authHandler := handler.NewAuthHandler(loginHandler, registerHandler, logger)

	mux.HandleFunc("GET /health", healthCheck)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
