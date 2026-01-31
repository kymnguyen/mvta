package workflow

type TransitionValidator interface {
	Validate(instance *WorkflowInstance, transition *Transition) error
}
