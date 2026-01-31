package route

import (
	"net/http"
	"strings"
	"workflow-svc/internal/api/handler"
	"workflow-svc/internal/application/service"
)

func SetupRoutes(workflowSvc *service.WorkflowService) http.Handler {
	mux := http.NewServeMux()

	// Workflow definition endpoints
	workflowHandler := handler.NewWorkflowHandler(workflowSvc)
	mux.HandleFunc("/api/workflows", workflowHandler.ListWorkflows)
	mux.HandleFunc("/api/workflows/", func(w http.ResponseWriter, r *http.Request) {
		// Route to different handlers based on path
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/start") {
			instanceHandler := handler.NewInstanceHandler(workflowSvc)
			instanceHandler.StartWorkflow(w, r)
		} else {
			workflowHandler.GetWorkflow(w, r)
		}
	})

	// Workflow instance endpoints
	instanceHandler := handler.NewInstanceHandler(workflowSvc)
	mux.HandleFunc("/api/instances", instanceHandler.GetOrListInstances)
	mux.HandleFunc("/api/instances/", func(w http.ResponseWriter, r *http.Request) {
		// Route to action handler if path contains /actions/
		if strings.Contains(r.URL.Path, "/actions/") {
			actionHandler := handler.NewActionHandler(workflowSvc)
			actionHandler.ProcessAction(w, r)
		} else {
			instanceHandler.GetOrListInstances(w, r)
		}
	})

	// Admin endpoints
	mux.HandleFunc("/admin/reload", workflowHandler.ReloadWorkflows)

	return mux
}
