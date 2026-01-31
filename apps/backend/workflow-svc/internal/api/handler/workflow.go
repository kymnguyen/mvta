package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"workflow-svc/internal/application/service"
)

type WorkflowHandler struct {
	svc *service.WorkflowService
}

func NewWorkflowHandler(svc *service.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{svc: svc}
}

func (h *WorkflowHandler) ListWorkflows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || r.URL.Path != "/api/workflows" {
		http.NotFound(w, r)
		return
	}

	workflows := h.svc.ListWorkflows()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"workflows": workflows})
}

func (h *WorkflowHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract workflow name from /api/workflows/:name
	path := strings.TrimPrefix(r.URL.Path, "/api/workflows/")
	if path == "" || strings.Contains(path, "/") {
		http.NotFound(w, r)
		return
	}

	workflow, err := h.svc.GetWorkflow(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

func (h *WorkflowHandler) ReloadWorkflows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.svc.ReloadWorkflows(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Workflows reloaded successfully"})
}
