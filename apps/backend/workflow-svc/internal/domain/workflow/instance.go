package workflow

import "time"

type WorkflowInstance struct {
	ID            string                 `bson:"_id,omitempty" json:"id"`
	WorkflowName  string                 `bson:"workflow_name" json:"workflow_name"`
	CorrelationID string                 `bson:"correlation_id" json:"correlation_id"`
	CurrentState  string                 `bson:"current_state" json:"current_state"`
	Context       map[string]interface{} `bson:"context" json:"context"`
	History       []StateTransition      `bson:"history" json:"history"`
	Version       int                    `bson:"version" json:"version"`
	CreatedAt     time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time              `bson:"updated_at" json:"updated_at"`
	TimeoutAt     *time.Time             `bson:"timeout_at,omitempty" json:"timeout_at,omitempty"`
}

type StateTransition struct {
	FromState string                 `bson:"from_state" json:"from_state"`
	ToState   string                 `bson:"to_state" json:"to_state"`
	Trigger   Trigger                `bson:"trigger" json:"trigger"`
	Context   map[string]interface{} `bson:"context" json:"context"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
}

type Trigger struct {
	Type string `bson:"type" json:"type"` // event, action, timeout
	Name string `bson:"name" json:"name"`
}

type InstanceFilter struct {
	WorkflowName  string
	State         string
	CorrelationID string
}
