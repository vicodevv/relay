package storage

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/vicodevv/relay/pkg/workflow"
)

type WorkflowDefinitionModel struct {
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	Version    int       `db:"version"`
	Definition JSONBMap  `db:"definition"`
	IsActive   bool      `db:"is_active"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type WorkflowInstanceModel struct {
	ID           string                  `db:"id"`
	DefinitionID string                  `db:"definition_id"`
	Status       workflow.WorkflowStatus `db:"status"`
	CurrentStep  *string                 `db:"current_step"`
	Input        JSONBMap                `db:"input"`
	Output       JSONBMap                `db:"output"`
	StartedAt    *time.Time              `db:"started_at"`
	CompletedAt  *time.Time              `db:"completed_at"`
	CreatedAt    time.Time               `db:"created_at"`
	UpdatedAt    time.Time               `db:"updated_at"`
}

type WorkflowEventModel struct {
	ID         int64     `db:"id"`
	InstanceID string    `db:"instance_id"`
	EventType  string    `db:"event_type"`
	StepID     *string   `db:"step_id"`
	EventData  JSONBMap  `db:"event_data"`
	CreatedAt  time.Time `db:"created_at"`
}

type StepExecutionModel struct {
	ID          string              `db:"id"`
	InstanceID  string              `db:"instance_id"`
	StepID      string              `db:"step_id"`
	Status      workflow.StepStatus `db:"status"`
	Attempt     int                 `db:"attempt"`
	Input       JSONBMap            `db:"input"`
	Output      JSONBMap            `db:"output"`
	Error       *string             `db:"error"`
	StartedAt   *time.Time          `db:"started_at"`
	CompletedAt *time.Time          `db:"completed_at"`
	CreatedAt   time.Time           `db:"created_at"`
}

type ApprovalModel struct {
	ID            string     `db:"id"`
	InstanceID    string     `db:"instance_id"`
	StepID        string     `db:"step_id"`
	Status        string     `db:"status"`
	ApproverEmail *string    `db:"approver_email"`
	ApprovedBy    *string    `db:"approved_by"`
	ApprovedAt    *time.Time `db:"approved_at"`
	Comments      *string    `db:"comments"`
	CreatedAt     time.Time  `db:"created_at"`
}

type CompensationModel struct {
	ID             string     `db:"id"`
	InstanceID     string     `db:"instance_id"`
	StepID         string     `db:"step_id"`
	CompensationID string     `db:"compensation_id"`
	Status         string     `db:"status"`
	Result         JSONBMap   `db:"result"`
	CreatedAt      time.Time  `db:"created_at"`
	ExecutedAt     *time.Time `db:"executed_at"`
}

// JSONBMap handles PostgreSQL JSONB columns
type JSONBMap map[string]interface{}

func (j JSONBMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONBMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}
