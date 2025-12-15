package repository

import (
	"database/sql"
	"fmt"
	"todo-api/internal/database"
	"todo-api/internal/models"
)

type WorkflowInstanceRepository struct {
	db *sql.DB
}

func NewWorkflowInstanceRepository() *WorkflowInstanceRepository {
	return &WorkflowInstanceRepository{
		db: database.DB,
	}
}

// CreateInstance creates a new workflow instance
func (r *WorkflowInstanceRepository) CreateInstance(instance *models.AssignedTodo) error {
	_, err := r.db.Exec(`INSERT INTO workflow_instances (id, workflow_id, current_step_id, assigned_to, created_at, updated_at) 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6 )`,
		instance.ID, instance.WorkflowId, instance.CurrentStepId, instance.AssignedTo, instance.CreatedAt, instance.UpdatedAt)
	return err
}

// GetInstance retrieves a workflow instance by ID
func (r *WorkflowInstanceRepository) GetInstance(id string) (*models.AssignedTodo, error) {
	instance := &models.AssignedTodo{}
	var todoData sql.NullString
	err := r.db.QueryRow(`SELECT id, workflow_id, current_step_id, title, description, task_data, assigned_to, created_by, created_at, updated_at 
		FROM workflow_instances WHERE id = @p1`, id).Scan(
		&instance.ID, &instance.WorkflowId, &instance.CurrentStepId,
		&todoData, &instance.AssignedTo, &instance.CreatedAt, &instance.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("instance not found")
	}
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// GetInstancesByWorkflow retrieves all instances for a workflow
func (r *WorkflowInstanceRepository) GetInstancesByWorkflow(workflowID string) ([]*models.AssignedTodo, error) {
	rows, err := r.db.Query(`SELECT id, workflow_id, current_step_id, title, description, task_data, assigned_to, created_by, created_at, updated_at 
		FROM workflow_instances WHERE workflow_id = @p1 ORDER BY created_at DESC`, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []*models.AssignedTodo
	for rows.Next() {
		instance := &models.AssignedTodo{}

		err := rows.Scan(&instance.ID, &instance.WorkflowId, &instance.CurrentStepId, &instance.AssignedTo, &instance.CreatedAt, &instance.UpdatedAt)
		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}
	return instances, nil
}

// GetInstancesByUser retrieves all instances assigned to a user
func (r *WorkflowInstanceRepository) GetInstancesByUser(userID string) ([]*models.AssignedTodo, error) {
	rows, err := r.db.Query(`SELECT id, workflow_id, current_step_id, title, description, task_data, assigned_to, created_by, created_at, updated_at 
		FROM workflow_instances WHERE assigned_to = @p1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []*models.AssignedTodo
	for rows.Next() {
		instance := &models.AssignedTodo{}

		err := rows.Scan(&instance.ID, &instance.WorkflowId, &instance.CurrentStepId, &instance.AssignedTo, &instance.CreatedAt, &instance.UpdatedAt)
		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}
	return instances, nil
}

// GetInstancesByStep retrieves all instances at a specific step
func (r *WorkflowInstanceRepository) GetInstancesByStep(stepID string) ([]*models.AssignedTodo, error) {
	rows, err := r.db.Query(`SELECT id, workflow_id, current_step_id, assigned_to, created_at, updated_at 
		FROM workflow_instances WHERE current_step_id = @p1 ORDER BY created_at DESC`, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []*models.AssignedTodo
	for rows.Next() {
		instance := &models.AssignedTodo{}

		err := rows.Scan(&instance.ID, &instance.WorkflowId, &instance.CurrentStepId, &instance.AssignedTo, &instance.CreatedAt, &instance.UpdatedAt)
		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}
	return instances, nil
}

// UpdateInstanceStep updates the current step of an instance
func (r *WorkflowInstanceRepository) UpdateInstanceStep(instanceID, newStepID string) error {
	_, err := r.db.Exec(`UPDATE workflow_instances SET current_step_id = @p1, updated_at = GETDATE() WHERE id = @p2`,
		newStepID, instanceID)
	return err
}

// UpdateInstance updates an instance
func (r *WorkflowInstanceRepository) UpdateInstance(instance *models.AssignedTodo) error {
	_, err := r.db.Exec(`UPDATE workflow_instances 
		SET assigned_to = @p1, updated_at = @p2 
		WHERE id = @p3`,
		instance.AssignedTo, instance.UpdatedAt, instance.ID)
	return err
}
