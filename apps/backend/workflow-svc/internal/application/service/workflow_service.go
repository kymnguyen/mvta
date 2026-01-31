package service

import (
	"context"
	"time"

	"workflow-svc/internal/application/registry"
	"workflow-svc/internal/domain/repository"
	"workflow-svc/internal/domain/workflow"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TransitionHandler interface {
	OnTransition(ctx context.Context, instance *workflow.WorkflowInstance, transition workflow.StateTransition) error
}

type WorkflowService struct {
	registry          *registry.DefinitionRegistry
	repo              repository.InstanceRepository
	transitionHandler TransitionHandler
	logger            *zap.Logger
}

func NewWorkflowService(
	registry *registry.DefinitionRegistry,
	repo repository.InstanceRepository,
	transitionHandler TransitionHandler,
	logger *zap.Logger,
) *WorkflowService {
	return &WorkflowService{
		registry:          registry,
		repo:              repo,
		transitionHandler: transitionHandler,
		logger:            logger,
	}
}

func (s *WorkflowService) Start(ctx context.Context, workflowName, correlationID string, context map[string]interface{}) (*workflow.WorkflowInstance, error) {
	def, err := s.registry.Get(workflowName)
	if err != nil {
		return nil, err
	}

	initialState := def.GetInitialState()
	if initialState == nil {
		return nil, workflow.ErrMissingInitialState
	}

	instance := &workflow.WorkflowInstance{
		ID:            uuid.New().String(),
		WorkflowName:  workflowName,
		CorrelationID: correlationID,
		CurrentState:  initialState.Name,
		Context:       context,
		History:       []workflow.StateTransition{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Set timeout if state has one
	if initialState.Timeout != nil {
		timeoutAt := time.Now().Add(*initialState.Timeout)
		instance.TimeoutAt = &timeoutAt
	}

	if err := s.repo.Create(ctx, instance); err != nil {
		return nil, err
	}

	s.logger.Info("Started workflow instance",
		zap.String("instance_id", instance.ID),
		zap.String("workflow", workflowName),
		zap.String("correlation_id", correlationID),
		zap.String("initial_state", initialState.Name))

	return instance, nil
}

func (s *WorkflowService) ProcessEvent(ctx context.Context, correlationID, eventName string, context map[string]interface{}) (*workflow.WorkflowInstance, error) {
	// Find instance by correlation_id
	instance, err := s.repo.FindByCorrelationID(ctx, correlationID)
	if err != nil {
		return nil, err
	}

	return s.processTransition(ctx, instance, workflow.Trigger{Type: "event", Name: eventName}, context)
}

func (s *WorkflowService) ProcessAction(ctx context.Context, instanceID, actionName string, context map[string]interface{}) (*workflow.WorkflowInstance, error) {
	instance, err := s.repo.FindByID(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	return s.processTransition(ctx, instance, workflow.Trigger{Type: "action", Name: actionName}, context)
}

func (s *WorkflowService) processTransition(ctx context.Context, instance *workflow.WorkflowInstance, trigger workflow.Trigger, context map[string]interface{}) (*workflow.WorkflowInstance, error) {
	def, err := s.registry.Get(instance.WorkflowName)
	if err != nil {
		return nil, err
	}

	// Find valid transition
	var validTransition *workflow.Transition
	for _, t := range def.Transitions {
		if t.From == instance.CurrentState {
			if (trigger.Type == "event" && t.Event == trigger.Name) ||
				(trigger.Type == "action" && t.Action == trigger.Name) ||
				(trigger.Type == "timeout" && t.Event == "timeout") {
				validTransition = t
				break
			}
		}
	}

	if validTransition == nil {
		return nil, workflow.ErrInvalidTransition
	}

	// Merge context
	for k, v := range context {
		instance.Context[k] = v
	}

	// Record transition
	transition := workflow.StateTransition{
		FromState: instance.CurrentState,
		ToState:   validTransition.To,
		Trigger:   trigger,
		Context:   context,
		Timestamp: time.Now(),
	}
	instance.History = append(instance.History, transition)
	instance.CurrentState = validTransition.To

	// Update timeout
	if newState, ok := def.States[validTransition.To]; ok && newState.Timeout != nil {
		timeoutAt := time.Now().Add(*newState.Timeout)
		instance.TimeoutAt = &timeoutAt
	} else {
		instance.TimeoutAt = nil
	}

	// Save with optimistic locking
	if err := s.repo.Update(ctx, instance); err != nil {
		return nil, err
	}

	// Post-transition handler
	if s.transitionHandler != nil {
		if err := s.transitionHandler.OnTransition(ctx, instance, transition); err != nil {
			s.logger.Error("Transition handler failed",
				zap.String("instance_id", instance.ID),
				zap.Error(err))
		}
	}

	s.logger.Info("Processed transition",
		zap.String("instance_id", instance.ID),
		zap.String("workflow", instance.WorkflowName),
		zap.String("from", transition.FromState),
		zap.String("to", transition.ToState),
		zap.String("trigger_type", trigger.Type),
		zap.String("trigger_name", trigger.Name))

	return instance, nil
}

func (s *WorkflowService) GetInstance(ctx context.Context, instanceID string) (*workflow.WorkflowInstance, error) {
	return s.repo.FindByID(ctx, instanceID)
}

func (s *WorkflowService) ListInstances(ctx context.Context, filter workflow.InstanceFilter) ([]*workflow.WorkflowInstance, error) {
	return s.repo.List(ctx, filter)
}

func (s *WorkflowService) GetWorkflow(name string) (*workflow.WorkflowDefinition, error) {
	return s.registry.Get(name)
}

func (s *WorkflowService) ListWorkflows() []*workflow.WorkflowDefinition {
	return s.registry.List()
}

func (s *WorkflowService) ReloadWorkflows() error {
	return s.registry.Reload()
}
