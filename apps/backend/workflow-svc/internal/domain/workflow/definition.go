package workflow

import "time"

type WorkflowDefinition struct {
	Name        string               `yaml:"name" json:"name"`
	Version     string               `yaml:"version" json:"version"`
	Description string               `yaml:"description" json:"description"`
	States      map[string]*StateDef `yaml:"states" json:"states"`
	Transitions []*Transition        `yaml:"transitions" json:"transitions"`
}

type StateDef struct {
	Name     string                 `yaml:"name" json:"name"`
	Type     string                 `yaml:"type" json:"type"` // initial, intermediate, terminal
	Timeout  *time.Duration         `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	OnEntry  []string               `yaml:"on_entry,omitempty" json:"on_entry,omitempty"`
	OnExit   []string               `yaml:"on_exit,omitempty" json:"on_exit,omitempty"`
	Metadata map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

type Transition struct {
	From      string                 `yaml:"from" json:"from"`
	To        string                 `yaml:"to" json:"to"`
	Event     string                 `yaml:"event,omitempty" json:"event,omitempty"`
	Action    string                 `yaml:"action,omitempty" json:"action,omitempty"`
	Condition string                 `yaml:"condition,omitempty" json:"condition,omitempty"`
	Metadata  map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

func (w *WorkflowDefinition) Validate() error {
	if w.Name == "" {
		return ErrInvalidWorkflow
	}

	// Check for duplicate states
	seen := make(map[string]bool)
	for name := range w.States {
		if seen[name] {
			return ErrDuplicateState
		}
		seen[name] = true
	}

	// Check for initial and terminal states
	hasInitial := false
	hasTerminal := false
	for _, state := range w.States {
		if state.Type == "initial" {
			hasInitial = true
		}
		if state.Type == "terminal" {
			hasTerminal = true
		}
	}
	if !hasInitial {
		return ErrMissingInitialState
	}
	if !hasTerminal {
		return ErrMissingTerminalState
	}

	// Validate transitions reference valid states
	for _, t := range w.Transitions {
		if _, ok := w.States[t.From]; !ok {
			return ErrInvalidState
		}
		if _, ok := w.States[t.To]; !ok {
			return ErrInvalidState
		}
	}

	return nil
}

func (w *WorkflowDefinition) GetInitialState() *StateDef {
	for _, state := range w.States {
		if state.Type == "initial" {
			return state
		}
	}
	return nil
}
