package interfaces

import "todo-api/internal/models"

// WorkflowInterface defines the business logic contract for workflow operations
type WorkflowInterface interface {
	CreateWorkflow(workflow *models.Workflow) error
	GetWorkflow(id string) (*models.Workflow, error)
	GetAllWorkflows() ([]models.Workflow, error)
	UpdateWorkflow(workflow *models.Workflow) error
	DeleteWorkflow(id string) error

	CreateStep(step *models.WorkflowStep) error
	GetWorkflowSteps(workflowID string) ([]models.WorkflowStep, error)

	CreateTransition(transition *models.WorkflowTransition) error
	GetWorkflowTransitions(workflowID string) ([]models.WorkflowTransition, error)
	GetAvailableTransitions(stepID string) ([]models.WorkflowTransition, error)
}
