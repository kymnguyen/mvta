package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"workflow-svc/internal/application/service"
	"workflow-svc/internal/domain/workflow"
)

type InstanceHandler struct {
	svc *service.WorkflowService
}

func NewInstanceHandler(svc *service.WorkflowService) *InstanceHandler {
	return &InstanceHandler{svc: svc}
}

type StartWorkflowRequest struct {
	CorrelationID string                 `json:"correlation_id"`
	Context       map[string]interface{} `json:"context"`
}

func (h *InstanceHandler) StartWorkflow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract workflow name from /api/workflows/:name/start
	path := strings.TrimPrefix(r.URL.Path, "/api/workflows/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "start" {
		http.NotFound(w, r)
		return
	}
	workflowName := parts[0]

	var req StartWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.CorrelationID == "" {
		http.Error(w, "correlation_id is required", http.StatusBadRequest)
		return
	}

	instance, err := h.svc.Start(r.Context(), workflowName, req.CorrelationID, req.Context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(instance)
}

func (h *InstanceHandler) GetOrListInstances(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/instances")

	// GET /api/instances/:id
	if path != "" && path != "/" {
		instanceID := strings.Trim(path, "/")
		if strings.Contains(instanceID, "/") {
			http.NotFound(w, r)
			return
		}
		h.getInstance(w, r, instanceID)
		return
	}

	// GET /api/instances
	h.listInstances(w, r)
}

func (h *InstanceHandler) getInstance(w http.ResponseWriter, r *http.Request, id string) {
	instance, err := h.svc.GetInstance(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instance)
}

func (h *InstanceHandler) listInstances(w http.ResponseWriter, r *http.Request) {
	filter := workflow.InstanceFilter{
		WorkflowName:  r.URL.Query().Get("workflow_name"),
		State:         r.URL.Query().Get("state"),
		CorrelationID: r.URL.Query().Get("correlation_id"),
	}

	instances, err := h.svc.ListInstances(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"instances": instances})
}
