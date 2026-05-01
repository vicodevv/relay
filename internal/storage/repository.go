package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vicodevv/relay/pkg/workflow"
)

type WorkflowRepository struct {
	db *DB
}

func NewWorkflowRepository(db *DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

// Workflow Definitions
func (r *WorkflowRepository) CreateDefinition(ctx context.Context, def *workflow.WorkflowDefinition) (string, error) {
	defJSON, err := json.Marshal(def)
	if err != nil {
		return "", fmt.Errorf("failed to marshal definition: %w", err)
	}

	id := uuid.New().String()
	query := `
		INSERT INTO workflow_definitions (id, name, version, definition)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err = r.db.QueryRowContext(ctx, query, id, def.Name, def.Version, defJSON).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create definition: %w", err)
	}

	return id, nil
}

func (r *WorkflowRepository) GetDefinitionByName(ctx context.Context, name string) (*WorkflowDefinitionModel, error) {
	var model WorkflowDefinitionModel
	query := `
		SELECT id, name, version, definition, is_active, created_at, updated_at
		FROM workflow_definitions
		WHERE name = $1 AND is_active = true
		ORDER BY version DESC
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &model, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("definition not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get definition: %w", err)
	}

	return &model, nil
}

// Workflow Instances
func (r *WorkflowRepository) CreateInstance(ctx context.Context, defID string, input map[string]interface{}) (string, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO workflow_instances (id, definition_id, status, input, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		defID,
		workflow.StatusPending,
		JSONBMap(input),
		now,
		now,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to create instance: %w", err)
	}

	return id, nil
}

func (r *WorkflowRepository) GetInstance(ctx context.Context, id string) (*WorkflowInstanceModel, error) {
	var model WorkflowInstanceModel
	query := `
		SELECT id, definition_id, status, current_step, input, output, 
		       started_at, completed_at, created_at, updated_at
		FROM workflow_instances
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("instance not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	return &model, nil
}

func (r *WorkflowRepository) UpdateInstanceStatus(ctx context.Context, id string, status workflow.WorkflowStatus, currentStep *string) error {
	query := `
		UPDATE workflow_instances
		SET status = $1, current_step = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, status, currentStep, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update instance status: %w", err)
	}

	return nil
}

func (r *WorkflowRepository) CompleteInstance(ctx context.Context, id string, output map[string]interface{}) error {
	query := `
		UPDATE workflow_instances
		SET status = $1, output = $2, completed_at = $3, updated_at = $4
		WHERE id = $5
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, workflow.StatusCompleted, JSONBMap(output), now, now, id)
	if err != nil {
		return fmt.Errorf("failed to complete instance: %w", err)
	}

	return nil
}

// Workflow Events (Event Sourcing)
func (r *WorkflowRepository) CreateEvent(ctx context.Context, instanceID, eventType string, stepID *string, eventData map[string]interface{}) error {
	query := `
		INSERT INTO workflow_events (instance_id, event_type, step_id, event_data)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, instanceID, eventType, stepID, JSONBMap(eventData))
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (r *WorkflowRepository) GetInstanceEvents(ctx context.Context, instanceID string) ([]WorkflowEventModel, error) {
	var events []WorkflowEventModel
	query := `
		SELECT id, instance_id, event_type, step_id, event_data, created_at
		FROM workflow_events
		WHERE instance_id = $1
		ORDER BY created_at ASC
	`

	err := r.db.SelectContext(ctx, &events, query, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}

// Step Executions
func (r *WorkflowRepository) CreateStepExecution(ctx context.Context, instanceID, stepID string, input map[string]interface{}) (string, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO step_executions (id, instance_id, step_id, status, attempt, input, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		instanceID,
		stepID,
		workflow.StepStatusPending,
		1,
		JSONBMap(input),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to create step execution: %w", err)
	}

	return id, nil
}

func (r *WorkflowRepository) UpdateStepExecution(ctx context.Context, id string, status workflow.StepStatus, output map[string]interface{}, errorMsg *string) error {
	query := `
		UPDATE step_executions
		SET status = $1, output = $2, error = $3, completed_at = $4
		WHERE id = $5
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, status, JSONBMap(output), errorMsg, now, id)
	if err != nil {
		return fmt.Errorf("failed to update step execution: %w", err)
	}

	return nil
}

// List Instances
func (r *WorkflowRepository) ListInstances(ctx context.Context, status *workflow.WorkflowStatus, limit int) ([]WorkflowInstanceModel, error) {
	var instances []WorkflowInstanceModel
	query := `
		SELECT id, definition_id, status, current_step, input, output,
		       started_at, completed_at, created_at, updated_at
		FROM workflow_instances
	`

	args := []interface{}{}
	if status != nil {
		query += " WHERE status = $1"
		args = append(args, *status)
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit)

	err := r.db.SelectContext(ctx, &instances, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	return instances, nil
}
