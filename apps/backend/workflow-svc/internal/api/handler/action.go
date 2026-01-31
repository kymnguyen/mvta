package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"workflow-svc/internal/application/service"
)

type ActionHandler struct {
	svc *service.WorkflowService
}

func NewActionHandler(svc *service.WorkflowService) *ActionHandler {
	return &ActionHandler{svc: svc}
}

type ProcessActionRequest struct {
	Context map[string]interface{} `json:"context"`
}

func (h *ActionHandler) ProcessAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract from /api/instances/:id/actions/:action
	path := strings.TrimPrefix(r.URL.Path, "/api/instances/")
	parts := strings.Split(path, "/")
	if len(parts) != 3 || parts[1] != "actions" {
		http.NotFound(w, r)
		return
	}
	instanceID := parts[0]
	action := parts[2]

	var req ProcessActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	instance, err := h.svc.ProcessAction(r.Context(), instanceID, action, req.Context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instance)
}
