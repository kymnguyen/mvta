package repository

import (
	"context"
	"workflow-svc/internal/domain/workflow"
)

type InstanceRepository interface {
	Create(ctx context.Context, instance *workflow.WorkflowInstance) error
	Update(ctx context.Context, instance *workflow.WorkflowInstance) error
	FindByID(ctx context.Context, id string) (*workflow.WorkflowInstance, error)
	FindByCorrelationID(ctx context.Context, correlationID string) (*workflow.WorkflowInstance, error)
	List(ctx context.Context, filter workflow.InstanceFilter) ([]*workflow.WorkflowInstance, error)
	FindPendingTimeouts(ctx context.Context, limit int) ([]*workflow.WorkflowInstance, error)
}
