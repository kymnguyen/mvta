package workflow

import "context"

type Engine interface {
	Start(ctx context.Context, workflowName, correlationID string, context map[string]interface{}) (*WorkflowInstance, error)
	ProcessEvent(ctx context.Context, correlationID, eventName string, context map[string]interface{}) (*WorkflowInstance, error)
	ProcessAction(ctx context.Context, instanceID, actionName string, context map[string]interface{}) (*WorkflowInstance, error)
	GetInstance(ctx context.Context, instanceID string) (*WorkflowInstance, error)
	ListInstances(ctx context.Context, filter InstanceFilter) ([]*WorkflowInstance, error)
}

func ValidateTransition(def *WorkflowDefinition, from, to, trigger string) error {
	for _, t := range def.Transitions {
		if t.From == from && t.To == to {
			if t.Event == trigger || t.Action == trigger {
				return nil
			}
		}
	}
	return ErrInvalidTransition
}
