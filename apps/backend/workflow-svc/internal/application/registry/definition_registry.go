package registry

import (
	"sync"

	"workflow-svc/internal/domain/workflow"
)

type DefinitionRegistry struct {
	mu         sync.RWMutex
	workflows  map[string]*workflow.WorkflowDefinition
	yamlLoader interface {
		LoadAll() ([]*workflow.WorkflowDefinition, error)
	}
}

func NewDefinitionRegistry(loader interface {
	LoadAll() ([]*workflow.WorkflowDefinition, error)
}) *DefinitionRegistry {
	return &DefinitionRegistry{
		workflows:  make(map[string]*workflow.WorkflowDefinition),
		yamlLoader: loader,
	}
}

func (r *DefinitionRegistry) Initialize() error {
	workflows, err := r.yamlLoader.LoadAll()
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, wf := range workflows {
		r.workflows[wf.Name] = wf
	}

	return nil
}

func (r *DefinitionRegistry) Reload() error {
	workflows, err := r.yamlLoader.LoadAll()
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.workflows = make(map[string]*workflow.WorkflowDefinition)
	for _, wf := range workflows {
		r.workflows[wf.Name] = wf
	}

	return nil
}

func (r *DefinitionRegistry) Get(name string) (*workflow.WorkflowDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	wf, ok := r.workflows[name]
	if !ok {
		return nil, workflow.ErrWorkflowNotFound
	}
	return wf, nil
}

func (r *DefinitionRegistry) List() []*workflow.WorkflowDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*workflow.WorkflowDefinition, 0, len(r.workflows))
	for _, wf := range r.workflows {
		result = append(result, wf)
	}
	return result
}
