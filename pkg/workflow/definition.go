package workflow

import "time"

type WorkflowDefinition struct {
	Name          string                 `yaml:"name" json:"name"`
	Description   string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Version       int                    `yaml:"version,omitempty" json:"version"`
	Steps         []Step                 `yaml:"steps" json:"steps"`
	Compensations []Compensation         `yaml:"compensations,omitempty" json:"compensations,omitempty"`
	Metadata      map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

type Step struct {
	ID            string                 `yaml:"id" json:"id"`
	Type          string                 `yaml:"type,omitempty" json:"type"`
	Service       string                 `yaml:"service,omitempty" json:"service"`
	Endpoint      string                 `yaml:"endpoint,omitempty" json:"endpoint"`
	Method        string                 `yaml:"method,omitempty" json:"method"`
	DependsOn     []string               `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	ParallelGroup string                 `yaml:"parallel_group,omitempty" json:"parallel_group,omitempty"`
	Retry         int                    `yaml:"retry,omitempty" json:"retry"`
	Timeout       string                 `yaml:"timeout,omitempty" json:"timeout"`
	Condition     string                 `yaml:"condition,omitempty" json:"condition,omitempty"`
	Compensate    string                 `yaml:"compensate,omitempty" json:"compensate,omitempty"`
	Approvers     []string               `yaml:"approvers,omitempty" json:"approvers,omitempty"`
	Input         map[string]interface{} `yaml:"input,omitempty" json:"input,omitempty"`
	Metadata      map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

type Compensation struct {
	ID       string                 `yaml:"id" json:"id"`
	Service  string                 `yaml:"service" json:"service"`
	Endpoint string                 `yaml:"endpoint" json:"endpoint"`
	Method   string                 `yaml:"method,omitempty" json:"method"`
	Input    map[string]interface{} `yaml:"input,omitempty" json:"input,omitempty"`
}

type WorkflowStatus string

const (
	StatusPending      WorkflowStatus = "PENDING"
	StatusRunning      WorkflowStatus = "RUNNING"
	StatusCompleted    WorkflowStatus = "COMPLETED"
	StatusFailed       WorkflowStatus = "FAILED"
	StatusCompensating WorkflowStatus = "COMPENSATING"
	StatusRolledBack   WorkflowStatus = "ROLLED_BACK"
	StatusCancelled    WorkflowStatus = "CANCELLED"
)

type StepStatus string

const (
	StepStatusPending   StepStatus = "PENDING"
	StepStatusRunning   StepStatus = "RUNNING"
	StepStatusCompleted StepStatus = "COMPLETED"
	StepStatusFailed    StepStatus = "FAILED"
	StepStatusSkipped   StepStatus = "SKIPPED"
)

type WorkflowInstance struct {
	ID           string                 `json:"id"`
	DefinitionID string                 `json:"definition_id"`
	Status       WorkflowStatus         `json:"status"`
	CurrentStep  string                 `json:"current_step,omitempty"`
	Input        map[string]interface{} `json:"input,omitempty"`
	Output       map[string]interface{} `json:"output,omitempty"`
	StartedAt    *time.Time             `json:"started_at,omitempty"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type StepExecution struct {
	ID          string                 `json:"id"`
	InstanceID  string                 `json:"instance_id"`
	StepID      string                 `json:"step_id"`
	Status      StepStatus             `json:"status"`
	Attempt     int                    `json:"attempt"`
	Input       map[string]interface{} `json:"input,omitempty"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Error       string                 `json:"error,omitempty"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}
