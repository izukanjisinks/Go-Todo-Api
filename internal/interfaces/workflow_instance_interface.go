package interfaces

import "todo-api/internal/models"

// WorkflowInstanceInterface defines the business logic contract for workflow instance operations
type WorkflowInstanceInterface interface {
	StartTask(workflowID string, todoID string, assignedTo string) (*models.AssignedTodo, error)
	GetTask(instanceID string) (*models.WorkflowInstanceWithDetails, error)
	GetTasksByUser(userID string) ([]models.AssignedTodo, error)
	GetTasksByWorkflow(workflowID string) ([]models.AssignedTodo, error)
	ExecuteAction(instanceID string, transitionID string, performedBy string, comments string) error
	GetAvailableActions(instanceID string) ([]models.AvailableAction, error)
	GetTaskHistory(instanceID string) ([]models.WorkflowHistory, error)
}
