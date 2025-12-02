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
func (r *WorkflowInstanceRepository) CreateInstance(instance *models.WorkflowInstance) error {
	_, err := r.db.Exec(`INSERT INTO workflow_instances (id, workflow_id, current_step_id, title, description, task_data, assigned_to, created_by, created_at, updated_at) 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10)`,
		instance.ID, instance.WorkflowID, instance.CurrentStepID, instance.Title, instance.Description,
		instance.TaskData, instance.AssignedTo, instance.CreatedBy, instance.CreatedAt, instance.UpdatedAt)
	return err
}

// GetInstance retrieves a workflow instance by ID
func (r *WorkflowInstanceRepository) GetInstance(id string) (*models.WorkflowInstance, error) {
	instance := &models.WorkflowInstance{}
	var taskData sql.NullString
	err := r.db.QueryRow(`SELECT id, workflow_id, current_step_id, title, description, task_data, assigned_to, created_by, created_at, updated_at 
		FROM workflow_instances WHERE id = @p1`, id).Scan(
		&instance.ID, &instance.WorkflowID, &instance.CurrentStepID, &instance.Title, &instance.Description,
		&taskData, &instance.AssignedTo, &instance.CreatedBy, &instance.CreatedAt, &instance.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("instance not found")
	}
	if err != nil {
		return nil, err
	}

	if taskData.Valid {
		instance.TaskData = taskData.String
	}
	return instance, nil
}

// GetInstancesByWorkflow retrieves all instances for a workflow
func (r *WorkflowInstanceRepository) GetInstancesByWorkflow(workflowID string) ([]*models.WorkflowInstance, error) {
	rows, err := r.db.Query(`SELECT id, workflow_id, current_step_id, title, description, task_data, assigned_to, created_by, created_at, updated_at 
		FROM workflow_instances WHERE workflow_id = @p1 ORDER BY created_at DESC`, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []*models.WorkflowInstance
	for rows.Next() {
		instance := &models.WorkflowInstance{}
		var taskData sql.NullString
		err := rows.Scan(&instance.ID, &instance.WorkflowID, &instance.CurrentStepID, &instance.Title, &instance.Description,
			&taskData, &instance.AssignedTo, &instance.CreatedBy, &instance.CreatedAt, &instance.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if taskData.Valid {
			instance.TaskData = taskData.String
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

// GetInstancesByUser retrieves all instances assigned to a user
func (r *WorkflowInstanceRepository) GetInstancesByUser(userID string) ([]*models.WorkflowInstance, error) {
	rows, err := r.db.Query(`SELECT id, workflow_id, current_step_id, title, description, task_data, assigned_to, created_by, created_at, updated_at 
		FROM workflow_instances WHERE assigned_to = @p1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []*models.WorkflowInstance
	for rows.Next() {
		instance := &models.WorkflowInstance{}
		var taskData sql.NullString
		err := rows.Scan(&instance.ID, &instance.WorkflowID, &instance.CurrentStepID, &instance.Title, &instance.Description,
			&taskData, &instance.AssignedTo, &instance.CreatedBy, &instance.CreatedAt, &instance.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if taskData.Valid {
			instance.TaskData = taskData.String
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

// GetInstancesByStep retrieves all instances at a specific step
func (r *WorkflowInstanceRepository) GetInstancesByStep(stepID string) ([]*models.WorkflowInstance, error) {
	rows, err := r.db.Query(`SELECT id, workflow_id, current_step_id, title, description, task_data, assigned_to, created_by, created_at, updated_at 
		FROM workflow_instances WHERE current_step_id = @p1 ORDER BY created_at DESC`, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []*models.WorkflowInstance
	for rows.Next() {
		instance := &models.WorkflowInstance{}
		var taskData sql.NullString
		err := rows.Scan(&instance.ID, &instance.WorkflowID, &instance.CurrentStepID, &instance.Title, &instance.Description,
			&taskData, &instance.AssignedTo, &instance.CreatedBy, &instance.CreatedAt, &instance.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if taskData.Valid {
			instance.TaskData = taskData.String
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
func (r *WorkflowInstanceRepository) UpdateInstance(instance *models.WorkflowInstance) error {
	_, err := r.db.Exec(`UPDATE workflow_instances 
		SET title = @p1, description = @p2, task_data = @p3, assigned_to = @p4, updated_at = @p5 
		WHERE id = @p6`,
		instance.Title, instance.Description, instance.TaskData, instance.AssignedTo, instance.UpdatedAt, instance.ID)
	return err
}

// CreateHistory creates a history record
func (r *WorkflowInstanceRepository) CreateHistory(history *models.WorkflowHistory) error {
	_, err := r.db.Exec(`INSERT INTO workflow_history (id, instance_id, from_step_id, to_step_id, action_taken, performed_by, comments, timestamp) 
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)`,
		history.ID, history.InstanceID, history.FromStepID, history.ToStepID,
		history.ActionTaken, history.PerformedBy, history.Comments, history.Timestamp)
	return err
}

// GetInstanceHistory retrieves the history for an instance
func (r *WorkflowInstanceRepository) GetInstanceHistory(instanceID string) ([]*models.WorkflowHistory, error) {
	rows, err := r.db.Query(`SELECT id, instance_id, from_step_id, to_step_id, action_taken, performed_by, comments, timestamp 
		FROM workflow_history WHERE instance_id = @p1 ORDER BY timestamp DESC`, instanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*models.WorkflowHistory
	for rows.Next() {
		h := &models.WorkflowHistory{}
		var fromStepID sql.NullString
		var comments sql.NullString
		err := rows.Scan(&h.ID, &h.InstanceID, &fromStepID, &h.ToStepID,
			&h.ActionTaken, &h.PerformedBy, &comments, &h.Timestamp)
		if err != nil {
			return nil, err
		}
		if fromStepID.Valid {
			h.FromStepID = &fromStepID.String
		}
		if comments.Valid {
			h.Comments = comments.String
		}
		history = append(history, h)
	}
	return history, nil
}
