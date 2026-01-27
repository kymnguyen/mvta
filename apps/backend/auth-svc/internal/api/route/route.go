package route

import (
	"net/http"

	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/api/handler"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/api/v1/auth/register", handler.Register)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
