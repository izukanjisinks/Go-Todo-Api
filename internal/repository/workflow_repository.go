package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"todo-api/internal/database"
	"todo-api/internal/models"
)

type WorkflowRepository struct {
	db *sql.DB
}

func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{
		db: database.DB,
	}
}

// CreateWorkflow creates a new workflow template
func (r *WorkflowRepository) CreateWorkflow(workflow *models.Workflow) error {
	_, err := r.db.Exec(`INSERT INTO workflows (id, name, description, is_active, created_by, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		workflow.ID, workflow.Name, workflow.Description, workflow.IsActive, workflow.CreatedBy, workflow.CreatedAt, workflow.UpdatedAt)
	return err
}

// GetWorkflow retrieves a workflow by ID
func (r *WorkflowRepository) GetWorkflow(id string) (*models.Workflow, error) {
	workflow := &models.Workflow{}
	err := r.db.QueryRow(`SELECT id, name, description, is_active, created_by, created_at, updated_at 
		FROM workflows WHERE id = $1`, id).Scan(
		&workflow.ID, &workflow.Name, &workflow.Description, &workflow.IsActive,
		&workflow.CreatedBy, &workflow.CreatedAt, &workflow.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workflow not found")
	}
	return workflow, err
}

// GetAllWorkflows retrieves all active workflows
func (r *WorkflowRepository) GetAllWorkflows() ([]*models.Workflow, error) {
	rows, err := r.db.Query(`SELECT id, name, description, is_active, created_by, created_at, updated_at 
		FROM workflows WHERE is_active = 1 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []*models.Workflow
	for rows.Next() {
		workflow := &models.Workflow{}
		err := rows.Scan(&workflow.ID, &workflow.Name, &workflow.Description, &workflow.IsActive,
			&workflow.CreatedBy, &workflow.CreatedAt, &workflow.UpdatedAt)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, workflow)
	}
	return workflows, nil
}

// CreateStep creates a new workflow step
func (r *WorkflowRepository) CreateStep(step *models.WorkflowStep) error {
	allowedRolesJSON, err := json.Marshal(step.AllowedRoles)
	if err != nil {
		return fmt.Errorf("failed to marshal allowed_roles: %w", err)
	}

	_, err = r.db.Exec(`INSERT INTO workflow_steps (id, workflow_id, step_name, step_order, initial, final, allowed_roles, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, @p8)`,
		step.ID, step.WorkflowID, step.StepName, step.StepOrder, step.Initial, step.Final, string(allowedRolesJSON), step.CreatedAt)
	return err
}

// GetWorkflowSteps retrieves all steps for a workflow
func (r *WorkflowRepository) GetWorkflowSteps(workflowID string) ([]*models.WorkflowStep, error) {
	rows, err := r.db.Query(`SELECT id, workflow_id, step_name, step_order, initial, final, allowed_roles, created_at 
		FROM workflow_steps WHERE workflow_id = $1 ORDER BY step_order`, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []*models.WorkflowStep
	for rows.Next() {
		step := &models.WorkflowStep{}
		var allowedRolesJSON string
		err := rows.Scan(&step.ID, &step.WorkflowID, &step.StepName, &step.StepOrder,
			&step.Initial, &step.Final, &allowedRolesJSON, &step.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Parse allowed_roles JSON
		if allowedRolesJSON != "" {
			json.Unmarshal([]byte(allowedRolesJSON), &step.AllowedRoles)
		}
		steps = append(steps, step)
	}
	return steps, nil
}

// GetStep retrieves a single step by ID
func (r *WorkflowRepository) GetStep(stepID string) (*models.WorkflowStep, error) {
	step := &models.WorkflowStep{}
	var allowedRolesJSON string
	err := r.db.QueryRow(`SELECT id, workflow_id, step_name, step_order, initial, final, allowed_roles, created_at 
		FROM workflow_steps WHERE id = $1`, stepID).Scan(
		&step.ID, &step.WorkflowID, &step.StepName, &step.StepOrder,
		&step.Initial, &step.Final, &allowedRolesJSON, &step.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("step not found")
	}
	if err != nil {
		return nil, err
	}

	// Parse allowed_roles JSON
	if allowedRolesJSON != "" {
		json.Unmarshal([]byte(allowedRolesJSON), &step.AllowedRoles)
	}
	return step, nil
}

// GetStartStep retrieves the start step for a workflow
func (r *WorkflowRepository) GetStartStep(workflowID string) (*models.WorkflowStep, error) {
	step := &models.WorkflowStep{}
	var allowedRolesJSON string
	err := r.db.QueryRow(`SELECT id, workflow_id, step_name, step_order, initial, final, allowed_roles, created_at 
		FROM workflow_steps WHERE workflow_id = $1 AND initial = 1`, workflowID).Scan(
		&step.ID, &step.WorkflowID, &step.StepName, &step.StepOrder,
		&step.Initial, &step.Final, &allowedRolesJSON, &step.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("start step not found for workflow")
	}
	if err != nil {
		return nil, err
	}

	// Parse allowed_roles JSON
	if allowedRolesJSON != "" {
		json.Unmarshal([]byte(allowedRolesJSON), &step.AllowedRoles)
	}
	return step, nil
}

// CreateTransition creates a new workflow transition
func (r *WorkflowRepository) CreateTransition(transition *models.WorkflowTransition) error {
	_, err := r.db.Exec(`INSERT INTO workflow_transitions (id, workflow_id, from_step_id, to_step_id, action_name, condition_type, condition_value, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, @p8)`,
		transition.ID, transition.WorkflowID, transition.FromStepID, transition.ToStepID,
		transition.ActionName, transition.ConditionType, transition.ConditionValue, transition.CreatedAt)
	return err
}

// GetTransitions retrieves all transitions for a workflow
func (r *WorkflowRepository) GetTransitions(workflowID string) ([]*models.WorkflowTransition, error) {
	rows, err := r.db.Query(`SELECT id, workflow_id, from_step_id, to_step_id, action_name, condition_type, condition_value, created_at 
		FROM workflow_transitions WHERE workflow_id = $1`, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transitions []*models.WorkflowTransition
	for rows.Next() {
		transition := &models.WorkflowTransition{}
		var conditionType, conditionValue sql.NullString
		err := rows.Scan(&transition.ID, &transition.WorkflowID, &transition.FromStepID, &transition.ToStepID,
			&transition.ActionName, &conditionType, &conditionValue, &transition.CreatedAt)
		if err != nil {
			return nil, err
		}
		if conditionType.Valid {
			transition.ConditionType = conditionType.String
		}
		if conditionValue.Valid {
			transition.ConditionValue = conditionValue.String
		}
		transitions = append(transitions, transition)
	}
	return transitions, nil
}

// GetAvailableTransitions retrieves possible transitions from a specific step
func (r *WorkflowRepository) GetAvailableTransitions(workflowID, fromStepID string) ([]*models.WorkflowTransition, error) {
	rows, err := r.db.Query(`SELECT id, workflow_id, from_step_id, to_step_id, action_name, condition_type, condition_value, created_at 
		FROM workflow_transitions WHERE workflow_id = $1 AND from_step_id = $2`, workflowID, fromStepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transitions []*models.WorkflowTransition
	for rows.Next() {
		transition := &models.WorkflowTransition{}
		var conditionType, conditionValue sql.NullString
		err := rows.Scan(&transition.ID, &transition.WorkflowID, &transition.FromStepID, &transition.ToStepID,
			&transition.ActionName, &conditionType, &conditionValue, &transition.CreatedAt)
		if err != nil {
			return nil, err
		}
		if conditionType.Valid {
			transition.ConditionType = conditionType.String
		}
		if conditionValue.Valid {
			transition.ConditionValue = conditionValue.String
		}
		transitions = append(transitions, transition)
	}
	return transitions, nil
}

// FindTransition finds a specific transition by action name
func (r *WorkflowRepository) FindTransition(workflowID, fromStepID, actionName string) (*models.WorkflowTransition, error) {
	transition := &models.WorkflowTransition{}
	var conditionType, conditionValue sql.NullString
	err := r.db.QueryRow(`SELECT id, workflow_id, from_step_id, to_step_id, action_name, condition_type, condition_value, created_at 
		FROM workflow_transitions WHERE workflow_id = $1 AND from_step_id = $2 AND action_name = $3`,
		workflowID, fromStepID, actionName).Scan(
		&transition.ID, &transition.WorkflowID, &transition.FromStepID, &transition.ToStepID,
		&transition.ActionName, &conditionType, &conditionValue, &transition.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transition not found")
	}
	if err != nil {
		return nil, err
	}

	if conditionType.Valid {
		transition.ConditionType = conditionType.String
	}
	if conditionValue.Valid {
		transition.ConditionValue = conditionValue.String
	}
	return transition, nil
}
