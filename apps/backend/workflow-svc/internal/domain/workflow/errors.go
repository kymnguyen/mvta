package workflow

import "errors"

var (
	ErrWorkflowNotFound       = errors.New("workflow not found")
	ErrInstanceNotFound       = errors.New("instance not found")
	ErrInvalidTransition      = errors.New("invalid transition")
	ErrInvalidState           = errors.New("invalid state")
	ErrInvalidAction          = errors.New("invalid action")
	ErrInvalidWorkflow        = errors.New("invalid workflow definition")
	ErrDuplicateState         = errors.New("duplicate state found")
	ErrMissingInitialState    = errors.New("initial state not defined")
	ErrMissingTerminalState   = errors.New("no terminal states defined")
	ErrOrphanedNode           = errors.New("orphaned node found")
	ErrConcurrentModification = errors.New("concurrent modification detected")
	ErrDuplicateEvent         = errors.New("event already processed")
)
